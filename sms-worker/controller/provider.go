package controller

import (
	"context"
	"fmt"
	controllermodel "sms-worker/models/controller"
)

type messageProviderController struct{}

func NewMessageProviderController() *messageProviderController {
	return &messageProviderController{}
}

func (p *messageProviderController) SendMessage(ctx context.Context, message controllermodel.Message) error {

	fmt.Printf("message with ID %d, User ID %d, Body %q delivered to %s",
		message.ID, message.UserID, message.Body, message.Recipient)
	return nil
}
