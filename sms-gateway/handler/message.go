package handler

import (
	"time"

	apimodel "sms-gateway/models/api"
	controllermodel "sms-gateway/models/controller"

	"github.com/gofiber/fiber/v3"
)

const reportDateLayout = "2006-01-02"

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

func (m *MessageHandler) GetUserReport(c fiber.Ctx) error {
	var request apimodel.UserMessageReportRequest
	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user_id, from, or to")
	}

	from, err := time.Parse(reportDateLayout, request.From)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid from date")
	}

	to, err := time.Parse(reportDateLayout, request.To)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid to date")
	}
	if !to.After(from) {
		return fiber.NewError(fiber.StatusBadRequest, "to must be after from")
	}

	messages, err := m.controller.GetUserReport(
		c.Context(),
		request.UserID,
		from,
		to,
	)
	if err != nil {
		return mapControllerError(err)
	}

	response := apimodel.UserMessageReportResponse{
		UserID:   request.UserID,
		From:     request.From,
		To:       request.To,
		Messages: make(apimodel.MessagesReportResponse, 0, len(messages)),
	}
	for _, message := range messages {
		response.Messages = append(response.Messages, apimodel.MessageReportResponse{
			ID:                       message.ID,
			Recipient:                message.Recipient,
			Body:                     message.Body,
			ServiceType:              string(message.ServiceType),
			Status:                   message.Status,
			CreatedAt:                message.CreatedAt,
			SubmissionLatencySeconds: message.SubmissionLatencySeconds,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
