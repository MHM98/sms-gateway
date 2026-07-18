package statrup

import (
	"context"
	"database/sql"
	"fmt"
	"sms-worker/pkg/db"
)

func openDatabase(ctx context.Context) (*sql.DB, error) {

	databasePool, err := db.Open(ctx)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return databasePool, nil
}
