package controller

import (
	"context"
)

type WalletController struct {
	walletDB WalletRepository
}

func NewWalletController(walletDB WalletRepository) *WalletController {
	return &WalletController{
		walletDB: walletDB,
	}
}

func (w *WalletController) TopUp(ctx context.Context, userID uint64, amount uint64) error {
	return w.walletDB.TopUp(ctx, userID, amount)
}

func (w *WalletController) GetUserBalance(ctx context.Context, userID uint64) (uint64, error) {
	return w.walletDB.GetUserBalance(ctx, userID)
}
