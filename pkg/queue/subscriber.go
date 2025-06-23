package queue

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type Subscriber struct {
	logger     watermill.LoggerAdapter
	config     *Config
	amqpConfig amqp.Config
	subscriber *amqp.Subscriber
	router     *message.Router
	handler    *message.Handler
}

func NewSubscriber(handlerFunc message.NoPublishHandlerFunc) (*Subscriber, error) {
	logger := watermill.NewSlogLogger(nil)
	config := NewConfig()
	amqpConfig := amqp.NewDurableQueueConfig(config.AmqpURI)

	subscriber, err := amqp.NewSubscriber(amqpConfig, logger)
	if err != nil {
		return nil, err
	}
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}
	handlerName := fmt.Sprintf("%s-handler", config.Topic)
	handler := router.AddNoPublisherHandler(handlerName, config.Topic, subscriber, handlerFunc)
	result := &Subscriber{
		logger:     logger,
		config:     config,
		amqpConfig: amqpConfig,
		subscriber: subscriber,
		router:     router,
		handler:    handler,
	}
	if err := result.addMiddlewares(); err != nil {
		return nil, err
	}
	return result, err
}

func (s *Subscriber) Run(ctx context.Context) error {
	return s.router.Run(ctx)
}

func (s *Subscriber) Close() error {
	return s.router.Close()
}

func (s *Subscriber) addMiddlewares() error {
	if s.config.TopicDLQ != "" {
		if err := s.subscriber.SubscribeInitialize(s.config.TopicDLQ); err != nil {
			return err
		}
		dlqPublisher, err := amqp.NewPublisher(s.amqpConfig, s.logger)
		if err != nil {
			return err
		}
		dlqMiddleware, err := middleware.PoisonQueue(dlqPublisher, s.config.TopicDLQ)
		if err != nil {
			return err
		}
		s.router.AddMiddleware(dlqMiddleware)
	}
	return nil
}
