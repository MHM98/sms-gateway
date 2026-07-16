package handler

import (
	"strconv"

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
		return mapControllerError(err)
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

	if err := h.controller.TopUp(c.Context(), request.UserID, request.Amount); err != nil {
		return mapControllerError(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
