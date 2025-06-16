package cleanup

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdCleaner struct {
	client *clientv3.Client
	config EtcdConfig
}

type EtcdConfig struct {
	Endpoints     []string
	Prefix        string
	RetryAttempts int
	RetryBackoff  time.Duration
}

func NewEtcdCleaner(config EtcdConfig) (*EtcdCleaner, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	return &EtcdCleaner{
		client: client,
		config: config,
	}, nil
}

func (e *EtcdCleaner) Cleanup(ctx context.Context) error {
	slog.Info("starting etcd cleanup", "prefix", e.config.Prefix)

	var lastErr error
	for attempt := 0; attempt < e.config.RetryAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Get all keys with the prefix first to count them
		resp, err := e.client.Get(ctx, e.config.Prefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
		if err != nil {
			lastErr = fmt.Errorf("failed to count keys with prefix %s: %w", e.config.Prefix, err)
			if attempt < e.config.RetryAttempts-1 {
				slog.Warn("retrying etcd key count", "prefix", e.config.Prefix, "attempt", attempt+1, "error", err)
				time.Sleep(e.config.RetryBackoff)
				continue
			}
			return lastErr
		}

		keyCount := resp.Count
		slog.Info("found keys to delete", "prefix", e.config.Prefix, "count", keyCount)

		if keyCount == 0 {
			slog.Info("no keys found to delete", "prefix", e.config.Prefix)
			return nil
		}

		// Delete all keys with the prefix
		delResp, err := e.client.Delete(ctx, e.config.Prefix, clientv3.WithPrefix())
		if err != nil {
			lastErr = fmt.Errorf("failed to delete keys with prefix %s: %w", e.config.Prefix, err)
			if attempt < e.config.RetryAttempts-1 {
				slog.Warn("retrying etcd key deletion", "prefix", e.config.Prefix, "attempt", attempt+1, "error", err)
				time.Sleep(e.config.RetryBackoff)
				continue
			}
			return lastErr
		}

		slog.Info("etcd cleanup completed", "prefix", e.config.Prefix, "keys_deleted", delResp.Deleted)
		return nil
	}

	return lastErr
}

func (e *EtcdCleaner) Close() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}
