package db

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	messageStatusSubmitted = "submitted"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *messageRepository {
	return &messageRepository{db: db}
}

func (m *messageRepository) MarkMessageStatusSubmitted(ctx context.Context, messageID uint64) error {
	query := `UPDATE messages 
	SET status=?,
	submission_latency_ms = TIMESTAMPDIFF(SECOND, created_at, CURRENT_TIMESTAMP())
	WHERE id=?`

	_, err := m.db.ExecContext(ctx, query, messageStatusSubmitted, messageID)
	if err != nil {
		return fmt.Errorf("failed to update messge status. err: %w", err)
	}

	return nil
}
