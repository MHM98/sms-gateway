package startup

import (
	"database/sql"

	"sms-gateway/controller"
	"sms-gateway/handler"
	infrastructuredb "sms-gateway/infrastructure/db"
)

type handlers struct {
	wallet  *handler.WalletHandler
	message *handler.MessageHandler
}

type controllers struct {
	message *controller.MessageController
}

func buildDependencies(db *sql.DB, publisher controller.IMessagePublisher) (handlers, controllers) {
	walletRepository := infrastructuredb.NewWalletDB(db)
	walletController := controller.NewWalletController(walletRepository)
	walletHandler := handler.NewWalletHandler(walletController)

	messageRepository := infrastructuredb.NewMessageDB(db)
	messageController := controller.NewMessageController(messageRepository, publisher)
	messageHandler := handler.NewMessageHandler(messageController)

	return handlers{
			wallet:  walletHandler,
			message: messageHandler,
		}, controllers{
			message: messageController,
		}
}
