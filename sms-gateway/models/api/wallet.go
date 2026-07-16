package apimodel

type WalletBalanceResponse struct {
	UserID  uint64 `json:"user_id"`
	Balance uint64 `json:"balance"`
}

type TopUpWalletRequest struct {
	UserID uint64 `json:"user_id"`
	Amount uint64 `json:"amount"`
}
