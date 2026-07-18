package controller

type MessageController struct {
	repo IMessageRepository
	// consumer
}

func NewMessageController(messageRepo IMessageRepository) *MessageController {
	return &MessageController{repo: messageRepo}
}


