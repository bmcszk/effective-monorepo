package main

import (
	"encoding/json"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmcszk/effective-monorepo/pkg/queue"
)

func main() {
	publisher, err := queue.NewPublisher()
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	router := http.NewServeMux()
	router.HandleFunc("POST /moves", createMove(publisher))
	http.ListenAndServe(":8080", router)
}

func createMove(publisher message.Publisher) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var chMvReq ChessMoveRequest
		if err := json.NewDecoder(r.Body).Decode(&chMvReq); err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		chMvMsg := chMvReq.ToMessage()
		qMsg, err := chMvMsg.ToQueueMessage()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		if err := publisher.Publish("test1", qMsg); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(201)
	}
}
