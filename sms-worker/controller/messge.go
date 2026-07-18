package controller

import (
	"context"
	"fmt"

	controllermodel "sms-worker/models/controller"
)

type MessageController struct {
	repo            IMessageRepository
	consumer        IMessageConsumer
	messageProvider IMessageProvider
}

func NewMessageController(messageRepo IMessageRepository,
	consumer IMessageConsumer, provider IMessageProvider) *MessageController {
	return &MessageController{
		repo:            messageRepo,
		consumer:        consumer,
		messageProvider: provider,
	}
}

func (c *MessageController) Consume(ctx context.Context) error {
	return c.consumer.Consume(ctx, c.handleMessage)
}

// TODO check for concureency here
func (c *MessageController) handleMessage(ctx context.Context, message controllermodel.Message) error {
	if err := c.messageProvider.SendMessage(ctx, message); err != nil {
		return fmt.Errorf("provider faile to submit message. %w", err)
	}
	if err := c.repo.MarkMessageStatusSubmitted(ctx, message.ID); err != nil {
		return fmt.Errorf("mark message %d as submitted: %w", message.ID, err)
	}

	// we log here for test purpose , in production it should not be the case
	fmt.Printf("message with ID %d successfully sent to %s\n", message.ID, message.Recipient)

	return nil
}
