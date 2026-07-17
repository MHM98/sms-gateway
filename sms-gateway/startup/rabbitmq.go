package startup

import (
	"context"
	"fmt"
	"os"

	rabbitmqinfra "sms-gateway/infrastructure/rabbitmq"
	controllermodel "sms-gateway/models/controller"
	rabbitmq "sms-gateway/pkg/rabbitmq"
)

const (
	rabbitMQExchange = "sms.messages"

	normalQueueName   = "sms.messages.normal"
	normalRoutingKey  = "message.normal"
	expressQueueName  = "sms.messages.express"
	expressRoutingKey = "message.express"
)

type rabbitResources struct {
	client           *rabbitmq.Client
	normalPublisher  *rabbitmq.Publisher
	expressPublisher *rabbitmq.Publisher
	adapter          *rabbitmqinfra.MessagePublisher
}

func openRabbitMQ(ctx context.Context) (result *rabbitResources, err error) {
	resources := &rabbitResources{}

	client, err := rabbitmq.NewClient(rabbitmq.Config{
		URL:      os.Getenv("RABBITMQ_URL"),
		Exchange: rabbitMQExchange,
		Queues: []rabbitmq.Queue{
			{Name: normalQueueName, RoutingKey: normalRoutingKey},
			{Name: expressQueueName, RoutingKey: expressRoutingKey},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create RabbitMQ client: %w", err)
	}
	resources.client = client

	if err := client.TryToConnect(ctx); err != nil {
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	normalPublisher, err := rabbitmq.NewPublisher(client)
	if err != nil {
		return nil, fmt.Errorf("create normal RabbitMQ publisher: %w", err)
	}
	resources.normalPublisher = normalPublisher

	expressPublisher, err := rabbitmq.NewPublisher(client)
	if err != nil {
		return nil, fmt.Errorf("create express RabbitMQ publisher: %w", err)
	}
	resources.expressPublisher = expressPublisher

	adapter, err := rabbitmqinfra.NewMessagePublisher(rabbitmqinfra.MessagePublisherConfig{
		controllermodel.ServiceTypeNormal: {
			Publisher:  normalPublisher,
			RoutingKey: normalRoutingKey,
		},
		controllermodel.ServiceTypeExpress: {
			Publisher:  expressPublisher,
			RoutingKey: expressRoutingKey,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create message publisher: %w", err)
	}
	resources.adapter = adapter

	return resources, nil
}
