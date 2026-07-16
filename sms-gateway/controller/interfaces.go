package controller

import (
	"context"
	controllermodel "sms-gateway/models/controller"
)

type WalletRepository interface {
	TopUp(ctx context.Context, userID uint64, amount uint64) error
	GetUserBalance(ctx context.Context, userID uint64) (uint64, error)
}

type MessageRepository interface {
	CreateAndCharge(ctx context.Context, data controllermodel.CreateMessage) error
}
