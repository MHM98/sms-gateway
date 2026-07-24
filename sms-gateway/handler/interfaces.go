package handler

import (
	"context"

	controllermodel "sms-gateway/models/controller"
)

// wallet
type IWalletController interface {
	TopUp(ctx context.Context, userID uint64, amount uint64) error
	GetUserBalance(ctx context.Context, userID uint64) (uint64, error)
}

type WalletHandler struct {
	controller IWalletController
}

func NewWalletHandler(controller IWalletController) *WalletHandler {
	return &WalletHandler{
		controller: controller,
	}
}

// message
type IMessageController interface {
	CreateAndCharge(ctx context.Context, input controllermodel.Message) error
	GetUserReport(ctx context.Context, query controllermodel.UserMessageReportQuery) (controllermodel.Messages, error)
}

type MessageHandler struct {
	controller IMessageController
}

func NewMessageHandler(controller IMessageController) *MessageHandler {
	return &MessageHandler{controller: controller}
}
