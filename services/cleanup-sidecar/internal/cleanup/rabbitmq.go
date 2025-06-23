package cleanup

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQCleaner struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	config  RabbitMQConfig
}

type RabbitMQConfig struct {
	URI           string
	Topic         string
	TopicDLQ      string
	RetryAttempts int
	RetryBackoff  time.Duration
}

func NewRabbitMQCleaner(config RabbitMQConfig) (*RabbitMQCleaner, error) {
	conn, err := amqp091.Dial(config.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		if closeErr := conn.Close(); closeErr != nil {
			slog.Error("failed to close connection", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQCleaner{
		conn:    conn,
		channel: channel,
		config:  config,
	}, nil
}

func (r *RabbitMQCleaner) Cleanup(ctx context.Context) error {
	slog.Info("starting RabbitMQ cleanup")

	// Clean main topic queue
	if err := r.cleanQueue(ctx, r.config.Topic); err != nil {
		return fmt.Errorf("failed to clean main queue: %w", err)
	}

	// Clean DLQ
	if err := r.cleanQueue(ctx, r.config.TopicDLQ); err != nil {
		return fmt.Errorf("failed to clean DLQ: %w", err)
	}

	slog.Info("RabbitMQ cleanup completed")
	return nil
}

func (r *RabbitMQCleaner) cleanQueue(ctx context.Context, queueName string) error {
	slog.Info("cleaning queue", "queue", queueName)
	return r.retryQueueCleanup(ctx, queueName)
}

func (r *RabbitMQCleaner) retryQueueCleanup(ctx context.Context, queueName string) error {
	var lastErr error
	for attempt := 0; attempt < r.config.RetryAttempts; attempt++ {
		if err := checkContextCancellationRabbit(ctx); err != nil {
			return err
		}

		if err := r.attemptQueueCleanup(queueName); err != nil {
			lastErr = err
			if r.shouldContinueRetrying(queueName, attempt, err) {
				continue
			}
			return lastErr
		}
		return nil
	}
	return lastErr
}

func checkContextCancellationRabbit(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (r *RabbitMQCleaner) shouldContinueRetrying(queueName string, attempt int, err error) bool {
	if attempt >= r.config.RetryAttempts-1 {
		return false
	}
	slog.Warn("retrying queue cleanup", "queue", queueName, "attempt", attempt+1, "error", err)
	time.Sleep(r.config.RetryBackoff)
	return true
}

func (r *RabbitMQCleaner) attemptQueueCleanup(queueName string) error {
	queue, err := r.declareQueue(queueName)
	if err != nil {
		return err
	}

	purged, err := r.purgeQueue(queueName)
	if err != nil {
		return err
	}

	slog.Info("queue cleaned", "queue", queueName, "messages_purged", purged, "total_messages", queue.Messages)
	return nil
}

func (r *RabbitMQCleaner) declareQueue(queueName string) (*amqp091.Queue, error) {
	queue, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}
	return &queue, nil
}

func (r *RabbitMQCleaner) purgeQueue(queueName string) (int, error) {
	purged, err := r.channel.QueuePurge(queueName, false)
	if err != nil {
		return 0, fmt.Errorf("failed to purge queue %s: %w", queueName, err)
	}
	return purged, nil
}

func (r *RabbitMQCleaner) Close() error {
	var errs []error

	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close channel: %w", err))
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %v", errs)
	}

	return nil
}
