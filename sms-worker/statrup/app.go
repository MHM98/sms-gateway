package statrup

import (
	"context"
	"database/sql"

	"sms-worker/controller"
	"sms-worker/pkg/rabbitmq"
)

type application struct {
	db                *sql.DB
	rabbitClient      *rabbitmq.Client
	messageController *controller.MessageController
}

func Run(ctx context.Context) error {
	app, err := newApplication(ctx)
	if err != nil {
		return err
	}

	defer app.Close()

	return app.Run(ctx)
}

func newApplication(ctx context.Context) (*application, error) {
	databasePool, err := openDatabase(ctx)
	if err != nil {
		return nil, err
	}

	resources := &application{db: databasePool}

	rabbit, err := openRabbitMQ(ctx)
	if err != nil {
		return nil, err
	}
	resources.rabbitClient = rabbit.client

	controllers := buildDependencies(databasePool, rabbit.consumer)
	resources.messageController = controllers.message

	return resources, nil
}
