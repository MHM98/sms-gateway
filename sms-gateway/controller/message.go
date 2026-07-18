package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	controllermodel "sms-gateway/models/controller"
)

// in production we should read this value from env
const messageChargeAmount uint64 = 1

type MessageController struct {
	repository IMessageRepository
	publisher  IMessagePublisher
}

func NewMessageController(repository IMessageRepository, publisher IMessagePublisher) *MessageController {
	return &MessageController{
		repository: repository,
		publisher:  publisher,
	}
}

func (c *MessageController) CreateAndCharge(ctx context.Context, message controllermodel.Message) error {
	message.ChargeAmount = messageChargeAmount

	return c.repository.CreateAndCharge(ctx, message)
}

func (c *MessageController) GetUserReport(ctx context.Context, userID uint64, from, to time.Time) (controllermodel.Messages, error) {
	messages, err := c.repository.GetUserReport(ctx, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("get report for user %d: %w", userID, err)
	}

	return messages, nil
}

func (c *MessageController) DispatchPendingMessages(ctx context.Context, serviceType controllermodel.ServiceType, limit int) error {

	messages, err := c.repository.ClaimPendingMessages(ctx, serviceType, limit)
	if err != nil {
		return fmt.Errorf("claim pending %s messages: %w", serviceType, err)
	}

	var dispatchErrors []error
	for _, message := range messages {
		if err := c.publisher.Publish(ctx, message); err != nil {
			publishErr := fmt.Errorf("publish message %d: %w", message.ID, err)

			//change message status for furture processsing
			releaseErr := c.repository.ReleaseMessage(ctx, message.ID, message.CreatedAt)
			if releaseErr != nil {
				dispatchErrors = append(dispatchErrors, errors.Join(
					publishErr,
					fmt.Errorf("release message %d: %w", message.ID, releaseErr),
				))
				continue
			}

			dispatchErrors = append(dispatchErrors, publishErr)
		}
	}

	return errors.Join(dispatchErrors...)
}
