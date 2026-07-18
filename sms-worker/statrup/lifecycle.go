package statrup

import (
	"context"
	"errors"
	"fmt"
)

func (a *application) Run(ctx context.Context) error {
	return a.messageController.Consume(ctx)
}

func (a *application) Close() error {
	var closeErrors []error

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
