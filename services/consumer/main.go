package main

import (
	"context"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmcszk/effective-monorepo/pkg/queue"
	"github.com/bmcszk/effective-monorepo/services/consumer/internal/game"
	"github.com/bmcszk/effective-monorepo/services/consumer/repo"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	repository, err := repo.NewRepo()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := repository.Close(); err != nil {
			slog.Error("failed to close repo", "error", err)
		}
	}()

	provider := game.NewGameProvider(repository)

	subscriber, err := queue.NewSubscriber(handleMessage(provider))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := subscriber.Close(); err != nil {
			slog.Error("failed to close subscriber", "error", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := subscriber.Run(ctx); err != nil {
		panic(err)
	}
}

func handleMessage(provider *game.GameProvider) func(*message.Message) error {
	return func(message *message.Message) error {
		ctx := message.Context()
		chessMessage, err := queue.FromQueueMessage(message)
		if err != nil {
			slog.Error("map from queue message", "error", err.Error())
			return err
		}
		slog.Info("chess message", "data", chessMessage)
		if err := provider.HandleMove(ctx, chessMessage.GameID, chessMessage.Move); err != nil {
			slog.Error("handle move", "error", err.Error())
			return err
		}
		return nil
	}
}
