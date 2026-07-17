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

type Queue struct {
	Name       string
	RoutingKey string
}

type Config struct {
	URL      string
	Exchange string
	Queues   []Queue
}

type Client struct {
	config Config

	mu     sync.Mutex
	conn   *amqp.Connection
	closed bool
}

const retryMaxElapsedTime = 60 * time.Second

var ErrClientClosed = errors.New("RabbitMQ client is closed")

func NewClient(config Config) (*Client, error) {
	if config.URL == "" {
		return nil, errors.New("RabbitMQ URL is required")
	}

	if config.Exchange == "" {
		return nil, errors.New("RabbitMQ exchange is required")
	}

	if len(config.Queues) == 0 {
		return nil, errors.New("at least one RabbitMQ queue is required")
	}

	for _, queue := range config.Queues {
		if queue.Name == "" || queue.RoutingKey == "" {
			return nil, errors.New("queue name and routing key are required")
		}
	}

	return &Client{
		config: config,
	}, nil
}

// retries with ExponentialBackOff until 60s or ctx expires.
func (c *Client) TryToConnect(ctx context.Context) error {
	_, err := backoff.Retry(ctx, func() (string, error) {
		_, err := c.connection()
		if errors.Is(err, ErrClientClosed) {
			err = backoff.Permanent(err)
		}

		return "", err
	},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxElapsedTime(retryMaxElapsedTime),
	)

	return err
}

// connection returns the current connection or creates a new one.
func (c *Client) connection() (*amqp.Connection, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil, ErrClientClosed
	}

	if c.conn != nil && !c.conn.IsClosed() {
		return c.conn, nil
	}

	conn, err := amqp.Dial(c.config.URL)
	if err != nil {
		return nil, fmt.Errorf("dial RabbitMQ: %w", err)
	}

	if err := c.declareTopology(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	c.conn = conn

	// try to reconnect client here
	closeNotify := make(chan *amqp.Error, 1)
	c.conn.NotifyClose(closeNotify)
	go c.reconnectWatcher(closeNotify)

	return conn, nil
}

func (c *Client) reconnectWatcher(ch <-chan *amqp.Error) {
	<-ch
	if c.closed {
		return
	}
	c.TryToConnect(context.Background())

}

func (c *Client) declareTopology(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open topology channel: %w", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		c.config.Exchange,
		"direct",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	for _, queue := range c.config.Queues {
		if _, err := ch.QueueDeclare(
			queue.Name,
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			nil,
		); err != nil {
			return fmt.Errorf(
				"declare queue %q: %w",
				queue.Name,
				err,
			)
		}

		if err := ch.QueueBind(
			queue.Name,
			queue.RoutingKey,
			c.config.Exchange,
			false,
			nil,
		); err != nil {
			return fmt.Errorf(
				"bind queue %q: %w",
				queue.Name,
				err,
			)
		}
	}

	return nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	if c.conn == nil || c.conn.IsClosed() {
		return nil
	}

	return c.conn.Close()
}
