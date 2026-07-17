package startup

import (
	"context"
	"database/sql"
	"fmt"

	database "sms-gateway/pkg/db"
)

func openDatabase(ctx context.Context) (*sql.DB, error) {
	databasePool, err := database.Open(ctx)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return databasePool, nil
}
