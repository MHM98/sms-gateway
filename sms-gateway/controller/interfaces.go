package controller

import (
	"context"
	"errors"
)

var ErrWalletNotFound = errors.New("wallet not found")

type WalletRepository interface {
	TopUp(ctx context.Context, userID uint64, amount uint64) error
	GetUserBalance(ctx context.Context, userID uint64) (uint64, error)
}
