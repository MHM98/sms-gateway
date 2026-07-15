package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// this slice used to feed dummy data to database
var seedUsers = []string{
	"user1",
	"user2",
	"user3",
	"user4",
	"user5",
}

func seedDatabase(ctx context.Context, db *sql.DB) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, name := range seedUsers {
		userID, seedErr := findOrCreateSeedUser(ctx, tx, name)
		if seedErr != nil {
			return seedErr
		}

		if _, seedErr = tx.ExecContext(ctx, `
			INSERT INTO wallets (user_id, balance)
			VALUES (?, 0)
			ON DUPLICATE KEY UPDATE user_id = wallets.user_id`, userID); seedErr != nil {
			return fmt.Errorf("create wallet for %q: %w", name, seedErr)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	return nil
}

func findOrCreateSeedUser(ctx context.Context, tx *sql.Tx, name string) (uint64, error) {
	var userID uint64
	err := tx.QueryRowContext(ctx, `
		SELECT id
		FROM users
		WHERE name = ?`,
		name).Scan(&userID)
	if err == nil {
		return userID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("find seed user %q: %w", name, err)
	}

	result, err := tx.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)", name)
	if err != nil {
		return 0, fmt.Errorf("create seed user %q: %w", name, err)
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get id for seed user %q: %w", name, err)
	}

	return uint64(insertedID), nil
}
