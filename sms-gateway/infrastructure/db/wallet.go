package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"sms-gateway/controller"
)

type WalletDB struct {
	db *sql.DB
}

func NewWalletDB(db *sql.DB) *WalletDB {
	return &WalletDB{db: db}
}

func (w *WalletDB) TopUp(ctx context.Context, userID uint64, amount uint64) error {
	query := `UPDATE wallets
	SET balance = balance + ?
	WHERE user_id = ?;`

	result, err := w.db.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("top up wallet: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("top up wallet rows affected: %w", err)
	}

	if affectedRows == 0 {
		return fmt.Errorf("top up wallet: %w", controller.ErrWalletNotFound)
	}

	return nil
}

func (w *WalletDB) GetUserBalance(ctx context.Context, userID uint64) (uint64, error) {
	query := `SELECT balance FROM wallets WHERE user_id=?`

	var balance uint64
	err := w.db.QueryRowContext(ctx, query, userID).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("get wallet balance: %w", controller.ErrWalletNotFound)
	}
	if err != nil {
		return 0, fmt.Errorf("get wallet balance: %w", err)
	}

	return balance, nil
}
