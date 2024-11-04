package queue

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewPublisher() (message.Publisher, error) {
	config := NewConfig()
	amqpConfig := amqp.NewDurableQueueConfig(config.AmqpURI)

	publisher, err := amqp.NewPublisher(
		amqpConfig,
		watermill.NewStdLogger(false, false))
	if err != nil {
		return nil, err
	}
	return publisher, err
}
