package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

const (
	numberOfUsers  = 1000
	defaultBalance = 10000
)

func seedDatabase(ctx context.Context, db *sql.DB) (err error) {
	seedUsers := createUsers(numberOfUsers)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}

	defer tx.Rollback()

	ids, err := insertUsers(ctx, tx, seedUsers)
	if err != nil {
		return fmt.Errorf("create user %w", err)
	}

	err = insertWallets(ctx, tx, ids)
	if err != nil {
		return fmt.Errorf("create wallet %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	return nil
}

func insertWallets(ctx context.Context, tx *sql.Tx, ids []int) error {
	placeholders := strings.
		TrimSuffix(strings.
			Repeat(fmt.Sprintf("(?, %d),", defaultBalance), len(ids)), ",")

	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf(`INSERT IGNORE INTO wallets (user_id, balance)
	VALUES %s`, placeholders)

	_, err := tx.ExecContext(ctx, query, args...)
	return err

}

func insertUsers(ctx context.Context, tx *sql.Tx, users map[int]string) ([]int, error) {
	placeholders := strings.TrimSuffix(strings.Repeat("(?, ?),", len(users)), ",")
	query := fmt.Sprintf(`INSERT IGNORE INTO users (id,name) VALUES %s`, placeholders)

	args := make([]any, 0, len(users))
	ids := make([]int, 0, len(users))

	for id, name := range users {
		args = append(args, id, name)
		ids = append(ids, id)
	}

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return ids, err
	}
	return ids, nil
}
func createUsers(number int) map[int]string {
	var users = make(map[int]string)

	for i := 1; i <= number; i++ {
		users[i] = fmt.Sprintf("user%d", i)
	}

	return users
}
