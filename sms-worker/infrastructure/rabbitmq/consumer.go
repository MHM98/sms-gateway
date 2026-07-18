package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"sms-worker/controller"
	controllermodel "sms-worker/models/controller"
	rabbitmqclient "sms-worker/pkg/rabbitmq"
)

type ConsumerConfig struct {
	Queue    string
	Prefetch int
	Client   *rabbitmqclient.Client
	Workers  int
}

type MessageConsumer struct {
	config ConsumerConfig
}

type messagePayload struct {
	ID          uint64                      `json:"id"`
	UserID      uint64                      `json:"user_id"`
	Recipient   string                      `json:"recipient"`
	Body        string                      `json:"body"`
	ServiceType controllermodel.ServiceType `json:"service_type"`
}

var _ controller.IMessageConsumer = (*MessageConsumer)(nil)

func NewMessageConsumer(config ConsumerConfig) (*MessageConsumer, error) {
	if config.Client == nil {
		return nil, errors.New("RabbitMQ client is required")
	}
	if config.Queue == "" {
		return nil, errors.New("consumer queue is required")
	}
	if config.Workers <= 0 {
		config.Workers = 1
	}

	return &MessageConsumer{
		config: config,
	}, nil
}

func (m *MessageConsumer) Consume(ctx context.Context, handler controller.MessageHandler) error {
	if handler == nil {
		return errors.New("message handler is required")
	}

	var (
		wg    sync.WaitGroup
		errCh = make(chan error, m.config.Workers)
	)

	for range m.config.Workers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			consumer, err := rabbitmqclient.NewConsumer(
				m.config.Client,
				m.config.Queue,
				m.config.Prefetch,
			)
			if err != nil {
				errCh <- fmt.Errorf("create consumer for queue %q: %w", m.config.Queue, err)
				return
			}

			if err := consumer.TryConsume(ctx, messageDeliveryHandler(handler)); err != nil {
				errCh <- fmt.Errorf(
					"try consume queue %s : %w",
					m.config.Queue,
					err,
				)
			}
		}()
	}

	wg.Wait()
	close(errCh)

	var consumeErrors []error
	for err := range errCh {
		consumeErrors = append(consumeErrors, err)
	}

	return errors.Join(consumeErrors...)
}

func messageDeliveryHandler(handler controller.MessageHandler) rabbitmqclient.Handler {
	return func(ctx context.Context, delivery rabbitmqclient.Delivery) (bool, error) {
		var data messagePayload
		if err := json.Unmarshal(delivery.Body, &data); err != nil {
			return false, fmt.Errorf("unmarshal message data: %w", err)
		}

		message := controllermodel.Message{
			ID:          data.ID,
			UserID:      data.UserID,
			Recipient:   data.Recipient,
			Body:        data.Body,
			ServiceType: data.ServiceType,
		}

		// if handler face error, we should requeue for furture processing
		if err := handler(ctx, message); err != nil {
			return true, err
		}

		return false, nil
	}
}
