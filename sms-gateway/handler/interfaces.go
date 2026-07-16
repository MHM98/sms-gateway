package handler

import (
	"context"
)

type WalletController interface {
	TopUp(ctx context.Context, userID uint64, amount uint64) error
	GetUserBalance(ctx context.Context, userID uint64) (uint64, error)
}

type WalletHandler struct {
	controller WalletController
}

func NewWalletHandler(controller WalletController) *WalletHandler {
	return &WalletHandler{
		controller: controller,
	}
}
