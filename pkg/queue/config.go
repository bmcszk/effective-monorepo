package queue

import "os"

type Config struct {
	AmqpURI  string
	Topic    string
	TopicDLQ string
}

func NewConfig() *Config {
	return &Config{
		AmqpURI:  os.Getenv("AMQP_URI"),
		Topic:    os.Getenv("QUEUE_TOPIC"),
		TopicDLQ: os.Getenv("QUEUE_TOPIC_DLQ"),
	}
}
