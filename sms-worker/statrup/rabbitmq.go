package statrup

import (
	"context"
	"fmt"
	"os"

	infraRabbitMQ "sms-worker/infrastructure/rabbitmq"
	controllermodel "sms-worker/models/controller"
	"sms-worker/pkg/rabbitmq"
)

const (
	rabbitMQExchange = "sms.messages"

	normalQueueName   = "sms.messages.normal"
	normalRoutingKey  = "message.normal"
	expressQueueName  = "sms.messages.express"
	expressRoutingKey = "message.express"

	consumerPrefetch = 4
	consumerWorkers  = 4
)

type rabbitResources struct {
	client   *rabbitmq.Client
	consumer *infraRabbitMQ.MessageConsumer
}

var consumerQueues = map[controllermodel.ServiceType]rabbitmq.Queue{
	controllermodel.ServiceTypeNormal:  {Name: normalQueueName, RoutingKey: normalRoutingKey},
	controllermodel.ServiceTypeExpress: {Name: expressQueueName, RoutingKey: expressRoutingKey},
}

func openRabbitMQ(ctx context.Context) (result *rabbitResources, err error) {
	resources := &rabbitResources{}

	serviceType := controllermodel.ServiceType(os.Getenv("WORKER_SERVICE_TYPE"))
	queue, ok := consumerQueues[serviceType]
	if !ok {
		return nil, fmt.Errorf(
			"service type must be %q or %q",
			controllermodel.ServiceTypeNormal,
			controllermodel.ServiceTypeExpress,
		)
	}

	client, err := rabbitmq.NewClient(rabbitmq.Config{
		URL:      os.Getenv("RABBITMQ_URL"),
		Exchange: rabbitMQExchange,
		Queues:   []rabbitmq.Queue{queue},
	})
	if err != nil {
		return nil, fmt.Errorf("create RabbitMQ client: %w", err)
	}

	resources.client = client

	if err := client.TryToConnect(ctx); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	consumer, err := infraRabbitMQ.NewMessageConsumer(infraRabbitMQ.ConsumerConfig{
		Client:   client,
		Queue:    queue.Name,
		Prefetch: consumerPrefetch,
		Workers:  consumerWorkers,
	})
	if err != nil {
		return nil, fmt.Errorf("create %s message consumer: %w", serviceType, err)
	}
	resources.consumer = consumer

	return resources, nil
}
