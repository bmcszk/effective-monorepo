package queue

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Publisher struct {
	message.Publisher
	config *Config
}

func NewPublisher() (*Publisher, error) {
	logger := watermill.NewSlogLogger(nil)
	config := NewConfig()
	amqpConfig := amqp.NewDurableQueueConfig(config.AmqpURI)

	publisher, err := amqp.NewPublisher(
		amqpConfig,
		logger)
	if err != nil {
		return nil, err
	}
	return &Publisher{Publisher: publisher, config: config}, err
}

func (p *Publisher) Publish(msg *message.Message) error {
	return p.Publisher.Publish(p.config.Topic, msg)
}
