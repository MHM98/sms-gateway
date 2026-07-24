package controller

import (
	"context"
	"time"

	controllermodel "sms-gateway/models/controller"
)

type IWalletRepository interface {
	TopUp(ctx context.Context, userID uint64, amount uint64) error
	GetUserBalance(ctx context.Context, userID uint64) (uint64, error)
}

type IMessageRepository interface {
	CreateAndCharge(ctx context.Context, message controllermodel.Message) error
	GetUserReport(ctx context.Context, query controllermodel.UserMessageReportQuery) (controllermodel.Messages, error)
	ClaimPendingMessages(ctx context.Context, serviceType controllermodel.ServiceType, limit int) (controllermodel.Messages, error)
	ReleaseMessage(ctx context.Context, messageID uint64, createdAt time.Time) error
}

type IMessagePublisher interface {
	Publish(ctx context.Context, message controllermodel.Message) error
}
