package controller

import (
	"context"

	controllermodel "sms-worker/models/controller"
)

type IMessageRepository interface {
	MarkMessageStatusSubmitted(ctx context.Context, messageID uint64) error
}

type MessageHandler func(ctx context.Context, message controllermodel.Message) error

type IMessageConsumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
}

type IMessageProvider interface {
	SendMessage(ctx context.Context, message controllermodel.Message) error
}
