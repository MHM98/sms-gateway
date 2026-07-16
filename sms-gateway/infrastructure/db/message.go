package db

import (
	"context"
	"database/sql"
	"fmt"

	"sms-gateway/controller"
	controllermodel "sms-gateway/models/controller"
)

type messageDB struct {
	db *sql.DB
}

func NewMessageDB(db *sql.DB) *messageDB {
	return &messageDB{db: db}
}

func (m *messageDB) CreateAndCharge(ctx context.Context, data controllermodel.CreateMessage) error {
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

func (m *messageDB) debitWallet(ctx context.Context, tx *sql.Tx, data controllermodel.CreateMessage) error {
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

func (m *messageDB) insertMessage(ctx context.Context, tx *sql.Tx, data controllermodel.CreateMessage) error {
	query := `INSERT INTO messages (user_id, recipient, body, service_type, status, created_at)
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP())`

	_, err := tx.ExecContext(ctx, query, data.UserID, data.Recipient,
		data.Body, data.ServiceType, data.Status)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}

	return nil
}
