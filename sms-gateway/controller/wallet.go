package controller

import (
	"context"
)

type WalletController struct {
	walletRepository WalletRepository
}

func NewWalletController(walletRepository WalletRepository) *WalletController {
	return &WalletController{
		walletRepository: walletRepository,
	}
}

func (w *WalletController) TopUp(ctx context.Context, userID uint64, amount uint64) error {
	return w.walletRepository.TopUp(ctx, userID, amount)
}

func (w *WalletController) GetUserBalance(ctx context.Context, userID uint64) (uint64, error) {
	return w.walletRepository.GetUserBalance(ctx, userID)
}
