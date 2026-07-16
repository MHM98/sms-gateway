package controller

import (
	"context"
	controllermodel "sms-gateway/models/controller"
)

// in production we should read this value from env
const messageChargeAmount uint64 = 1

type messageController struct {
	messageRepository MessageRepository
}

func NewMessageController(messageRepository MessageRepository) *messageController {
	return &messageController{messageRepository: messageRepository}
}

func (m *messageController) CreateAndCharge(ctx context.Context, DTO controllermodel.CreateMessage) error {
	DTO.ChargeAmount = messageChargeAmount
	DTO.Status = "pending" // should read the status from env

	return m.messageRepository.CreateAndCharge(ctx, DTO)
}
