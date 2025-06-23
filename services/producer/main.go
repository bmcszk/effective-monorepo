package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
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
		chMvReq, err := parseRequest(r)
		if err != nil {
			writeErrorResponse(w, 400, err)
			return
		}

		qMsg, err := convertToQueueMessage(chMvReq)
		if err != nil {
			writeErrorResponse(w, 500, err)
			return
		}

		if err := publisher.Publish(qMsg); err != nil {
			writeErrorResponse(w, 500, err)
			return
		}

		w.WriteHeader(201)
	}
}

func parseRequest(r *http.Request) (producer.ChessMoveRequest, error) {
	var chMvReq producer.ChessMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&chMvReq); err != nil {
		return chMvReq, err
	}
	return chMvReq, nil
}

func convertToQueueMessage(chMvReq producer.ChessMoveRequest) (*message.Message, error) {
	chMvMsg := chMvReq.ToMessage()
	return chMvMsg.ToQueueMessage()
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
		slog.Error("failed to write response", "error", writeErr)
	}
}
