package startup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
)

const (
	httpAddress     = ":3000"
	shutdownTimeout = 10 * time.Second
)

func (a *application) Run(ctx context.Context) error {
	if a.cron != nil {
		a.cron.Start()
	}

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

	if a.cron != nil {
		<-a.cron.Stop().Done()
	}

	if a.http != nil {
		if err := a.http.ShutdownWithTimeout(shutdownTimeout); err != nil && !errors.Is(err, fiber.ErrNotRunning) {
			closeErrors = append(closeErrors, fmt.Errorf("shut down HTTP server: %w", err))
		}
	}

	if a.rabbitExpressPublisher != nil {
		if err := a.rabbitExpressPublisher.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("close express RabbitMQ publisher: %w", err))
		}
	}

	if a.rabbitNormalPublisher != nil {
		if err := a.rabbitNormalPublisher.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("close normal RabbitMQ publisher: %w", err))
		}
	}

	if a.rabbitClient != nil {
		if err := a.rabbitClient.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("close RabbitMQ client: %w", err))
		}
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("close database: %w", err))
		}
	}

	return errors.Join(closeErrors...)
}
