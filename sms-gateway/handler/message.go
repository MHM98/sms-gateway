package handler

import (
	apimodel "sms-gateway/models/api"
	controllermodel "sms-gateway/models/controller"

	"github.com/gofiber/fiber/v3"
)

func (m *MessageHandler) SendMessage(c fiber.Ctx) error {
	var request apimodel.MessageRequest
	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	input := controllermodel.Message{
		UserID:      request.UserID,
		Recipient:   request.Recipient,
		Body:        request.Body,
		ServiceType: controllermodel.ServiceType(request.ServiceType),
	}

	if err := m.controller.CreateAndCharge(c.Context(), input); err != nil {
		return mapControllerError(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
