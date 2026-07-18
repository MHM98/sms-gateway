package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

const publishRetryTimeout = 60 * time.Second

var ErrPublishUncertain = errors.New(
	"RabbitMQ publish result is uncertain",
)

type PublishMessage struct {
	ID         string
	RoutingKey string
	Body       []byte
}

type Publisher struct {
	client *Client

	mu sync.Mutex
	ch *amqp.Channel
}

func NewPublisher(client *Client) (*Publisher, error) {
	if client == nil {
		return nil, errors.New("RabbitMQ client is required")
	}

	return &Publisher{
		client: client,
	}, nil
}

// retries with ExponentialBackOff until 60s or ctx expires.
func (p *Publisher) TryToPublish(ctx context.Context, message PublishMessage) error {
	if err := p.validateMessage(message); err != nil {
		return err
	}

	_, err := backoff.Retry(ctx, func() (string, error) {
		err := p.Publish(ctx, message)
		if err == nil {
			return "", nil
		}

		// RabbitMQ may have received the message.
		// Retrying could publish a duplicate.
		if errors.Is(err, ErrPublishUncertain) {
			return "", backoff.Permanent(err)
		}

		return "", err
	},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxElapsedTime(publishRetryTimeout),
	)

	if retryErr := backoff.AsRetryError(err); retryErr != nil {
		if errors.Is(err, backoff.ErrPermanent) {
			return retryErr.LastErr
		}
	}

	return err
}

// publish message to queue. it should not used concurrently
func (p *Publisher) Publish(ctx context.Context, message PublishMessage) error {
	if message.RoutingKey == "" {
		return errors.New("routing key is required")
	}

	if len(message.Body) == 0 {
		return errors.New("message body is required")
	}

	// amqp.Channel must not be used concurrently.
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.ensureChannel(); err != nil {
		return err
	}

	confirmation, err :=
		p.ch.PublishWithDeferredConfirmWithContext(
			ctx,
			p.client.config.Exchange,
			message.RoutingKey,
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				MessageId:    message.ID,
				Timestamp:    time.Now().UTC(),
				Body:         message.Body,
			},
		)

	if err != nil {
		p.resetChannel()

		return fmt.Errorf(
			"%w: %v",
			ErrPublishUncertain,
			err,
		)
	}

	confirmed, err := confirmation.WaitContext(ctx)
	if err != nil {
		p.resetChannel()

		return fmt.Errorf(
			"%w: %v",
			ErrPublishUncertain,
			err,
		)
	}

	if !confirmed {
		return errors.New("RabbitMQ rejected the message")
	}

	return nil
}

func (p *Publisher) validateMessage(message PublishMessage) error {
	if message.RoutingKey == "" {
		return errors.New("routing key is required")
	}

	if len(message.Body) == 0 {
		return errors.New("message body is required")
	}
	return nil
}

func (p *Publisher) ensureChannel() error {
	if p.ch != nil && !p.ch.IsClosed() {
		return nil
	}

	conn, err := p.client.connection()
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open publisher channel: %w", err)
	}

	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		return fmt.Errorf("enable publisher confirms: %w", err)
	}

	p.ch = ch

	return nil
}

func (p *Publisher) resetChannel() {
	if p.ch != nil {
		_ = p.ch.Close()
		p.ch = nil
	}
}

func (p *Publisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.ch == nil || p.ch.IsClosed() {
		return nil
	}

	return p.ch.Close()
}
