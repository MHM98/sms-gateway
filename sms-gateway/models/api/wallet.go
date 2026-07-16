package apimodel

type WalletBalanceResponse struct {
	UserID  uint64 `json:"user_id"`
	Balance uint64 `json:"balance"`
}

type TopUpWalletRequest struct {
	UserID uint64 `json:"user_id" validate:"gt=0"`
	Amount uint64 `json:"amount" validate:"gt=0"`
}
