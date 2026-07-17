package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"sms-gateway/controller"
	controllermodel "sms-gateway/models/controller"
	rabbitmqclient "sms-gateway/pkg/rabbitmq"
)

type PublisherRoute struct {
	Publisher  *rabbitmqclient.Publisher
	RoutingKey string
}

type MessagePublisherConfig map[controllermodel.ServiceType]PublisherRoute

type MessagePublisher struct {
	routes MessagePublisherConfig
}

type messagePayload struct {
	ID          uint64                      `json:"id"`
	UserID      uint64                      `json:"user_id"`
	Recipient   string                      `json:"recipient"`
	Body        string                      `json:"body"`
	ServiceType controllermodel.ServiceType `json:"service_type"`
}

var _ controller.IMessagePublisher = (*MessagePublisher)(nil)

func NewMessagePublisher(config MessagePublisherConfig) (*MessagePublisher, error) {
	if len(config) == 0 {
		return nil, errors.New("at least one message publisher route is required")
	}

	routes := make(MessagePublisherConfig, len(config))
	for serviceType, route := range config {
		if serviceType == "" {
			return nil, errors.New("message publisher service type is required")
		}
		if route.Publisher == nil {
			return nil, fmt.Errorf("RabbitMQ publisher is required for service type %q", serviceType)
		}
		if route.RoutingKey == "" {
			return nil, fmt.Errorf("routing key is required for service type %q", serviceType)
		}

		routes[serviceType] = route
	}

	return &MessagePublisher{routes: routes}, nil
}

func (p *MessagePublisher) Publish(ctx context.Context, message controllermodel.Message) error {
	route, ok := p.routes[message.ServiceType]
	if !ok {
		return fmt.Errorf("publish message %d: unsupported service type %q", message.ID, message.ServiceType)
	}

	body, err := json.Marshal(messagePayload{
		ID:          message.ID,
		UserID:      message.UserID,
		Recipient:   message.Recipient,
		Body:        message.Body,
		ServiceType: message.ServiceType,
	})
	if err != nil {
		return fmt.Errorf("marshal message %d payload: %w", message.ID, err)
	}

	if err := route.Publisher.TryToPublish(ctx, rabbitmqclient.PublishMessage{
		ID:         strconv.FormatUint(message.ID, 10),
		RoutingKey: route.RoutingKey,
		Body:       body,
	}); err != nil {
		return fmt.Errorf("publish message %d to RabbitMQ: %w", message.ID, err)
	}

	return nil
}
