package queue

import "os"

type Config struct {
	AmqpURI string
}

func NewConfig() *Config {
	return &Config{
		AmqpURI: os.Getenv("AMQP_URI"),
	}
}
