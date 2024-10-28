package main

import (
	"context"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmcszk/effective-monorepo/pkg/queue"
)

func main() {
	subscriber, err := queue.NewSubscriber()
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	messages, err := subscriber.Subscribe(ctx, "test1")
	if err != nil {
		panic(err)
	}
	for message := range messages {
		if err := handleMessage(message); err != nil {
			slog.Error("handle message", "error", err.Error())
			message.Nack()
			continue
		}
		message.Ack()
	}
}

func handleMessage(message *message.Message) error {
	chessMessage, err := queue.FromQueueMessage(message)
	if err != nil {
		return err
	}
	slog.Info("chess message", "data", chessMessage)
	return nil
}
