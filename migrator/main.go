package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("specify a migration action: up or down")
	}

	databaseDSN := os.Getenv("DATABASE_URL")
	if databaseDSN == "" {
		log.Fatal("DATABASE_URL is required")
	}

	var (
		driverName = "mysql"
		action     = strings.ToLower(os.Args[1])
	)

	driver, err := createDriverInstance(databaseDSN, driverName)
	if err != nil {
		log.Fatalf("create database driver: %v", err)
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
		if err := migrator.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("database is already up to date")
				return
			}

			log.Fatalf("run migration up: %v", err)
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
		log.Fatalf("unsupported migration action %q: use up or down", action)
	}

	log.Println("database migrations completed")
}

func createDriverInstance(dsn, driverName string) (database.Driver, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create MySQL migration driver: %w", err)
	}

	return driver, nil
}
