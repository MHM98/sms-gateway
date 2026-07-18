package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Delivery struct {
	ID          string
	Body        []byte
	Redelivered bool
}

// requeue=true immediately returns the failed message to the queue.
type Handler func(ctx context.Context, message Delivery) (requeue bool, err error)

type Consumer struct {
	client   *Client
	queue    string
	prefetch int
}

type consumerSession struct {
	channel    *amqp.Channel
	deliveries <-chan amqp.Delivery
}

func NewConsumer(client *Client, queue string, prefetch int) (*Consumer, error) {
	if client == nil {
		return nil, fmt.Errorf("RabbitMQ client is required")
	}

	if queue == "" {
		return nil, fmt.Errorf("consumer queue is required")
	}

	if prefetch <= 0 {
		prefetch = 1
	}

	return &Consumer{
		client:   client,
		queue:    queue,
		prefetch: prefetch,
	}, nil
}

// When the RabbitMQ connection or channel closes, it retries opening a new
// consumer session with exponential backoff for up to 60 seconds.
func (c *Consumer) TryConsume(ctx context.Context, handler Handler) error {
	if handler == nil {
		return fmt.Errorf("consumer handler is required")
	}

	for ctx.Err() == nil {
		session, err := c.tryOpenSession(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			return fmt.Errorf(
				"open RabbitMQ consumer session: %w",
				err,
			)
		}

		err = c.consume(ctx, session, handler)

		_ = session.channel.Close()

		if ctx.Err() != nil {
			return nil
		}

		log.Printf(
			"RabbitMQ consumer disconnected: queue=%s error=%v",
			c.queue,
			err,
		)
	}

	return nil
}

func (c *Consumer) tryOpenSession(ctx context.Context) (*consumerSession, error) {
	return backoff.Retry(
		ctx,
		func() (*consumerSession, error) {
			session, err := c.openSession(ctx)
			if errors.Is(err, ErrClientClosed) {
				err = backoff.Permanent(err)
			}

			return session, err
		},
		backoff.WithBackOff(
			backoff.NewExponentialBackOff(),
		),
		backoff.WithNotify(
			func(err error, next time.Duration) {
				log.Printf(
					"RabbitMQ consumer connection failed: queue=%s retry_in=%s error=%v",
					c.queue,
					next,
					err,
				)
			},
		),
		backoff.WithMaxElapsedTime(retryMaxElapsedTime),
	)
}

func (c *Consumer) openSession(ctx context.Context) (*consumerSession, error) {
	conn, err := c.client.connection()
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf(
			"open consumer channel: %w",
			err,
		)
	}

	if err := ch.Qos(
		c.prefetch,
		0,
		false,
	); err != nil {
		_ = ch.Close()

		return nil, fmt.Errorf(
			"set consumer prefetch: %w",
			err,
		)
	}

	deliveries, err := ch.ConsumeWithContext(
		ctx,
		c.queue,
		"",
		false, // manual ACK
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		_ = ch.Close()

		return nil, fmt.Errorf(
			"start consumer: %w",
			err,
		)
	}

	return &consumerSession{
		channel:    ch,
		deliveries: deliveries,
	}, nil
}

func (c *Consumer) consume(ctx context.Context, session *consumerSession, handler Handler) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case delivery, ok := <-session.deliveries:
			if !ok {
				return fmt.Errorf(
					"RabbitMQ delivery channel closed",
				)
			}

			requeue, handlerErr := handler(
				ctx,
				Delivery{
					ID:          delivery.MessageId,
					Body:        delivery.Body,
					Redelivered: delivery.Redelivered,
				},
			)

			if handlerErr == nil {
				if err := delivery.Ack(false); err != nil {
					return fmt.Errorf(
						"ack message: %w",
						err,
					)
				}

				continue
			}

			if err := delivery.Nack(
				false,
				requeue,
			); err != nil {
				return fmt.Errorf(
					"nack message: %w",
					err,
				)
			}

			log.Printf(
				"RabbitMQ message failed: queue=%s message_id=%s requeue=%t error=%v",
				c.queue,
				delivery.MessageId,
				requeue,
				handlerErr,
			)
		}
	}
}
