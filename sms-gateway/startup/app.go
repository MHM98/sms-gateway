package startup

import (
	"context"
	"database/sql"

	"sms-gateway/pkg/rabbitmq"

	"github.com/gofiber/fiber/v3"
	"github.com/robfig/cron/v3"
)

type application struct {
	http                   *fiber.App
	db                     *sql.DB
	rabbitClient           *rabbitmq.Client
	rabbitNormalPublisher  *rabbitmq.Publisher
	rabbitExpressPublisher *rabbitmq.Publisher
	cron                   *cron.Cron
}

func Run(ctx context.Context) (err error) {
	app, err := newApplication(ctx)
	if err != nil {
		return err
	}

	defer app.Close()

	return app.Run(ctx)
}

func newApplication(ctx context.Context) (result *application, err error) {
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
	resources.rabbitNormalPublisher = rabbit.normalPublisher
	resources.rabbitExpressPublisher = rabbit.expressPublisher

	handlers, controllers := buildDependencies(databasePool, rabbit.adapter)
	resources.http = newHTTPApplication(handlers)

	cronRunner, err := buildScheduler(controllers.message)
	if err != nil {
		return nil, err
	}
	resources.cron = cronRunner

	return resources, nil
}
