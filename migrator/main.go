package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 1 {
		log.Fatal("specify an action: up, down")
	}

	databaseDSN := os.Getenv("DATABASE_URL")
	if databaseDSN == "" {
		log.Fatal("DATABASE_URL is required")
	}

	var (
		driverName = "mysql"
		action     = strings.ToLower(os.Args[1])
	)

	db, err := sql.Open(driverName, databaseDSN)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		_ = db.Close()
		log.Fatalf("create MySQL migration driver: %v", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://migrations", driverName, driver)
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}

	defer func() {
		_, _ = migrator.Close()
	}()

	switch action {
	case "up":
		err := migrator.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("run migration up: %v", err)
		}
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("database is already up to date")
		}

		log.Println("seeding users and wallets")
		if err := seedDatabase(context.Background(), db); err != nil {
			log.Fatalf("seed database: %v", err)
		}

	case "down":
		if err := migrator.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("database is already up to date")
				return
			}

			log.Fatalf("run migration down: %v", err)
		}
	default:
		log.Fatalf("unsupported action %q: use up, down", action)
	}

	log.Println("database migrations completed")
}
