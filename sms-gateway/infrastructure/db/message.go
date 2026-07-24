package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"sms-gateway/controller"
	controllermodel "sms-gateway/models/controller"
)

const (
	messageStatusPending             = "pending"
	messageStatusProcessing          = "processing"
	messagePreviousDayOverlapMinutes = 15
	userMessageReportPageSize        = 500
)

type messageDB struct {
	db *sql.DB
}

func NewMessageDB(db *sql.DB) *messageDB {
	return &messageDB{db: db}
}

func (m *messageDB) CreateAndCharge(ctx context.Context, data controllermodel.Message) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin message transaction: %w", err)
	}
	defer tx.Rollback()

	if err := m.debitWallet(ctx, tx, data); err != nil {
		return err
	}

	if err := m.insertMessage(ctx, tx, data); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit message transaction: %w", err)
	}

	return nil
}

func (m *messageDB) GetUserReport(ctx context.Context, filter controllermodel.UserMessageReportQuery) (controllermodel.Messages, error) {
	query := `SELECT id, user_id, recipient, body, service_type, status, created_at,
			IFNULL(submission_latency_seconds,0)
			FROM messages
			WHERE user_id = ? AND created_at >= ? AND created_at < ? AND id > ?
			ORDER BY id ASC
			LIMIT ?`

	rows, err := m.db.QueryContext(ctx, query, filter.UserID, filter.From,
		filter.To, filter.LastSeen, userMessageReportPageSize)
	if err != nil {
		return nil, fmt.Errorf("select user messages report: %w", err)
	}
	defer rows.Close()

	messages := make(controllermodel.Messages, 0, userMessageReportPageSize)
	for rows.Next() {
		var message controllermodel.Message

		if err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Recipient,
			&message.Body,
			&message.ServiceType,
			&message.Status,
			&message.CreatedAt,
			&message.SubmissionLatencySeconds,
		); err != nil {
			return nil, fmt.Errorf("scan user messages report: %w", err)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user messages report: %w", err)
	}

	return messages, nil
}

func (m *messageDB) ClaimPendingMessages(ctx context.Context, serviceType controllermodel.ServiceType, limit int) (controllermodel.Messages, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin claim messages transaction: %w", err)
	}
	defer tx.Rollback()

	// double check table index here
	query := `SELECT id, user_id, recipient, body, service_type, status, created_at
		FROM messages
		WHERE status = ? AND service_type = ?
		AND created_at >= CURRENT_DATE - INTERVAL ? MINUTE
		AND created_at < CURRENT_DATE + INTERVAL 1 DAY
		ORDER BY created_at, id
		LIMIT ?
		FOR UPDATE SKIP LOCKED`

	rows, err := tx.QueryContext(ctx, query,
		messageStatusPending, serviceType,
		messagePreviousDayOverlapMinutes, limit)
	if err != nil {
		return nil, fmt.Errorf("select pending messages for claim: %w", err)
	}

	defer rows.Close()

	messages := make(controllermodel.Messages, 0, limit)
	for rows.Next() {
		var message controllermodel.Message
		if err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Recipient,
			&message.Body,
			&message.ServiceType,
			&message.Status,
			&message.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan pending message for claim: %w", err)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pending messages for claim: %w", err)
	}

	if len(messages) < 1 {
		return messages, nil
	}

	if err = m.updateMessagesStatus(ctx, tx, messages); err != nil {
		return nil, fmt.Errorf("failed to update message status : %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit claim messages transaction: %w", err)
	}

	return messages, nil
}

func (m *messageDB) updateMessagesStatus(ctx context.Context, tx *sql.Tx, messages controllermodel.Messages) error {
	placeholders := strings.TrimSuffix(strings.Repeat("(?, ?),", len(messages)), ",")

	query := fmt.Sprintf(`
	UPDATE messages
			SET status = ?
			WHERE status = ? AND (id, created_at) IN (%s)`, placeholders)

	args := make([]any, 0, 2+len(messages)*2)
	args = append(args, messageStatusProcessing, messageStatusPending)

	for _, message := range messages {
		args = append(args, message.ID, message.CreatedAt)
	}

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("mark claimed messages as processing: %w", err)
	}

	// update new message status to processing
	for i := range messages {
		messages[i].Status = messageStatusProcessing
	}

	return nil
}

func (m *messageDB) ReleaseMessage(ctx context.Context, messageID uint64, createdAt time.Time) error {

	query := `UPDATE messages
		SET status = ?
		WHERE id = ? AND created_at = ? AND status = ?`
	_, err := m.db.ExecContext(
		ctx,
		query,
		messageStatusPending,
		messageID,
		createdAt,
		messageStatusProcessing,
	)
	if err != nil {
		return fmt.Errorf("release message to pending: %w", err)
	}

	return nil
}

func (m *messageDB) debitWallet(ctx context.Context, tx *sql.Tx, data controllermodel.Message) error {
	query := `UPDATE wallets
	SET balance = balance - ?,
	updated_at = CURRENT_TIMESTAMP()
	WHERE user_id = ? AND balance >= ?`

	result, err := tx.ExecContext(ctx, query, data.ChargeAmount,
		data.UserID, data.ChargeAmount)
	if err != nil {
		return fmt.Errorf("debit wallet: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("debit wallet rows affected: %w", err)
	}

	if affectedRows == 0 {
		return controller.ErrInsufficientBalance
	}
	return nil
}

func (m *messageDB) insertMessage(ctx context.Context, tx *sql.Tx, data controllermodel.Message) error {
	query := `INSERT INTO messages (user_id, recipient, body, service_type, status, created_at)
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP())`

	_, err := tx.ExecContext(ctx, query, data.UserID, data.Recipient,
		data.Body, data.ServiceType, messageStatusPending)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}

	return nil
}
