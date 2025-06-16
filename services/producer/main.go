package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bmcszk/effective-monorepo/pkg/queue"
	"github.com/bmcszk/effective-monorepo/services/producer/producer"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("failed to load .env")
	}
	publisher, err := queue.NewPublisher()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := publisher.Close(); err != nil {
			slog.Error("failed to close publisher", "error", err)
		}
	}()

	router := http.NewServeMux()
	router.HandleFunc("POST /moves", createMove(publisher))
	if err := http.ListenAndServe(":8080", router); err != nil {
		slog.Error("server failed", "error", err)
	}
}

func createMove(publisher *queue.Publisher) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var chMvReq producer.ChessMoveRequest
		if err := json.NewDecoder(r.Body).Decode(&chMvReq); err != nil {
			w.WriteHeader(400)
			if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
				slog.Error("failed to write response", "error", writeErr)
			}
			return
		}
		chMvMsg := chMvReq.ToMessage()
		qMsg, err := chMvMsg.ToQueueMessage()
		if err != nil {
			w.WriteHeader(500)
			if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
				slog.Error("failed to write response", "error", writeErr)
			}
			return
		}
		if err := publisher.Publish(qMsg); err != nil {
			w.WriteHeader(500)
			if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
				slog.Error("failed to write response", "error", writeErr)
			}
			return
		}
		w.WriteHeader(201)
	}
}
