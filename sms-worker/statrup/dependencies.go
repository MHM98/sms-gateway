package statrup

import (
	"database/sql"

	"sms-worker/controller"
	infrastructuredb "sms-worker/infrastructure/db"
)

type controllers struct {
	message *controller.MessageController
}

func buildDependencies(db *sql.DB, consumer controller.IMessageConsumer) controllers {

	messageProvider := controller.NewMessageProviderController()

	messageRepository := infrastructuredb.NewMessageRepository(db)
	messageController := controller.NewMessageController(messageRepository, consumer, messageProvider)

	return controllers{
		message: messageController,
	}
}
