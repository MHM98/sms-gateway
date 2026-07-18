package statrup

import (
	"context"
	"database/sql"
	"os"
	"sms-worker/controller"
	infraDB "sms-worker/infrastructure/db"
	"sms-worker/pkg/rabbitmq"
)

type application struct {
	db           *sql.DB
	rabbitClient *rabbitmq.Client
}

func Run(ctx context.Context) error {

	db, err := openDatabase(ctx)
	if err != nil {
		return err
	}

	messageRepo := infraDB.NewMessageRepository(db)
	messageCntl := controller.NewMessageController(messageRepo)

	rabbitmq.NewClient(rabbitmq.Config{
		URL: os.Getenv("RABBITMQ_URL"),
		
	})
}
