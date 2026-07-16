package handler

import (
	"errors"
	"strconv"

	"sms-gateway/controller"
	apimodel "sms-gateway/models/api"

	"github.com/gofiber/fiber/v3"
)

func (h *WalletHandler) GetWalletBalance(c fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("user_id"), 10, 64)
	if err != nil || userID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user_id")
	}

	balance, err := h.controller.GetUserBalance(c.Context(), userID)
	if err != nil {
		return walletError(err)
	}

	return c.Status(fiber.StatusOK).JSON(apimodel.WalletBalanceResponse{
		UserID:  userID,
		Balance: balance,
	})
}

func (h *WalletHandler) TopUpWallet(c fiber.Ctx) error {
	var request apimodel.TopUpWalletRequest
	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	// we should think of middleware layer in handler
	if request.UserID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user_id")
	}
	if request.Amount == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "amount must be greater than zero")
	}

	if err := h.controller.TopUp(c.Context(), request.UserID, request.Amount); err != nil {
		return walletError(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func walletError(err error) error {
	if errors.Is(err, controller.ErrWalletNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "wallet not found")
	}

	return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
}
