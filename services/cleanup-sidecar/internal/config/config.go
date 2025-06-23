package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// Target resources to clean
	CleanQueue    bool
	CleanDatabase bool

	// RabbitMQ configuration
	AMQPURI       string
	QueueTopic    string
	QueueTopicDLQ string

	// etcd configuration
	EtcdURIs   []string
	EtcdPrefix string

	// Retry configuration
	RetryAttempts int
	RetryBackoff  time.Duration

	// Timeout configuration
	Timeout time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		CleanQueue:    getEnvBool("CLEANUP_QUEUE", true),
		CleanDatabase: getEnvBool("CLEANUP_DATABASE", true),

		AMQPURI:       getEnv("AMQP_URI", "amqp://guest:guest@rabbitmq:5672/"),
		QueueTopic:    getEnv("QUEUE_TOPIC", "monorepo"),
		QueueTopicDLQ: getEnv("QUEUE_TOPIC_DLQ", "monorepo-dlq"),

		EtcdURIs:   getEnvStringSlice("ETCD_URIS", []string{"etcd:2379"}),
		EtcdPrefix: getEnv("ETCD_PREFIX", "monorepo"),

		RetryAttempts: getEnvInt("RETRY_ATTEMPTS", 3),
		RetryBackoff:  getEnvDuration("RETRY_BACKOFF", 5*time.Second),

		Timeout: getEnvDuration("CLEANUP_TIMEOUT", 30*time.Second),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
