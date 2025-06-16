package cleanup

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bmcszk/effective-monorepo/services/cleanup-sidecar/internal/config"
)

type Cleaner struct {
	config          *config.Config
	rabbitmqCleaner *RabbitMQCleaner
	etcdCleaner     *EtcdCleaner
}

func New(cfg *config.Config) (*Cleaner, error) {
	cleaner := &Cleaner{
		config: cfg,
	}

	// Initialize RabbitMQ cleaner if enabled
	if cfg.CleanQueue {
		rabbitmqConfig := RabbitMQConfig{
			URI:           cfg.AMQPURI,
			Topic:         cfg.QueueTopic,
			TopicDLQ:      cfg.QueueTopicDLQ,
			RetryAttempts: cfg.RetryAttempts,
			RetryBackoff:  cfg.RetryBackoff,
		}

		rabbitmqCleaner, err := NewRabbitMQCleaner(rabbitmqConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create RabbitMQ cleaner: %w", err)
		}
		cleaner.rabbitmqCleaner = rabbitmqCleaner
	}

	// Initialize etcd cleaner if enabled
	if cfg.CleanDatabase {
		etcdConfig := EtcdConfig{
			Endpoints:     cfg.EtcdURIs,
			Prefix:        cfg.EtcdPrefix,
			RetryAttempts: cfg.RetryAttempts,
			RetryBackoff:  cfg.RetryBackoff,
		}

		etcdCleaner, err := NewEtcdCleaner(etcdConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create etcd cleaner: %w", err)
		}
		cleaner.etcdCleaner = etcdCleaner
	}

	return cleaner, nil
}

func (c *Cleaner) Run(ctx context.Context) error {
	slog.Info("starting cleanup process",
		"clean_queue", c.config.CleanQueue,
		"clean_database", c.config.CleanDatabase)

	if err := c.cleanQueue(ctx); err != nil {
		return err
	}

	if err := c.cleanDatabase(ctx); err != nil {
		return err
	}

	slog.Info("cleanup process completed successfully")
	return nil
}

func (c *Cleaner) cleanQueue(ctx context.Context) error {
	if !c.config.CleanQueue {
		slog.Info("skipping RabbitMQ cleanup (disabled)")
		return nil
	}

	if c.rabbitmqCleaner == nil {
		return nil
	}

	if err := c.rabbitmqCleaner.Cleanup(ctx); err != nil {
		return fmt.Errorf("RabbitMQ cleanup failed: %w", err)
	}

	return nil
}

func (c *Cleaner) cleanDatabase(ctx context.Context) error {
	if !c.config.CleanDatabase {
		slog.Info("skipping etcd cleanup (disabled)")
		return nil
	}

	if c.etcdCleaner == nil {
		return nil
	}

	if err := c.etcdCleaner.Cleanup(ctx); err != nil {
		return fmt.Errorf("etcd cleanup failed: %w", err)
	}

	return nil
}

func (c *Cleaner) Close() error {
	var errs []error

	if c.rabbitmqCleaner != nil {
		if err := c.rabbitmqCleaner.Close(); err != nil {
			errs = append(errs, fmt.Errorf("RabbitMQ cleaner close error: %w", err))
		}
	}

	if c.etcdCleaner != nil {
		if err := c.etcdCleaner.Close(); err != nil {
			errs = append(errs, fmt.Errorf("etcd cleaner close error: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cleanup close errors: %v", errs)
	}

	return nil
}
