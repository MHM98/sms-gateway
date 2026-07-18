package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Open(ctx context.Context) (*sql.DB, error) {
	databaseDSN := os.Getenv("DATABASE_URL")
	if databaseDSN == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	conn, err := sql.Open("mysql", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	conn.SetMaxOpenConns(80)
	conn.SetMaxIdleConns(40)
	conn.SetConnMaxLifetime(2 * time.Minute)

	if err := conn.PingContext(ctx); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return conn, nil
}
