package queue

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

var amqpURI = "amqp://guest:guest@localhost:5672/"

func NewSubscriber() (message.Subscriber, error) {
	config := NewConfig()
	amqpConfig := amqp.NewDurableQueueConfig(config.AmqpURI)

	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return nil, err
	}
	return subscriber, err
}
