package handler

import (
	"errors"

	"sms-gateway/controller"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func mapControllerError(err error) *fiber.Error {
	switch {
	case errors.Is(err, controller.ErrWalletNotFound):
		return fiber.NewError(fiber.StatusNotFound, "wallet not found")
	case errors.Is(err, controller.ErrInsufficientBalance):
		return fiber.NewError(fiber.StatusConflict, "insufficient wallet balance")
	default:
		log.Errorf("controller request failed: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
}
