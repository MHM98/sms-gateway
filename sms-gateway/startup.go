package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"sms-gateway/controller"
	"sms-gateway/handler"
	infrastructuredb "sms-gateway/infrastructure/db"
	database "sms-gateway/pkg/db"

	"github.com/gofiber/fiber/v3"
)

const (
	httpAddress     = ":3000"
	shutdownTimeout = 10 * time.Second
)

type application struct {
	http *fiber.App
	db   *sql.DB
}

func newApplication(ctx context.Context) (_ *application, err error) {
	databasePool, err := database.Open(ctx)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	walletRepository := infrastructuredb.NewWalletDB(databasePool)
	walletController := controller.NewWalletController(walletRepository)
	walletHandler := handler.NewWalletHandler(walletController)

	httpApp := fiber.New()
	initRoutes(httpApp, handlers{
		wallet: walletHandler,
	})

	return &application{
		http: httpApp,
		db:   databasePool,
	}, nil
}

func (a *application) Run(ctx context.Context) error {
	if err := a.http.Listen(httpAddress, fiber.ListenConfig{
		GracefulContext: ctx,
		ShutdownTimeout: shutdownTimeout,
	}); err != nil {
		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
}

func (a *application) Close() error {
	var closeErrors []error

	if a.http != nil {
		if err := a.http.ShutdownWithTimeout(shutdownTimeout); err != nil && !errors.Is(err, fiber.ErrNotRunning) {
			closeErrors = append(closeErrors, fmt.Errorf("shut down HTTP server: %w", err))
		}
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("close database: %w", err))
		}
	}

	return errors.Join(closeErrors...)
}
