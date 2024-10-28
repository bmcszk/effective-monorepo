package main

import (
	"time"

	"github.com/bmcszk/effective-monorepo/pkg/model"
	"github.com/bmcszk/effective-monorepo/pkg/queue"
)

type ChessMoveRequest struct {
	model.ChessMove
}

func (r ChessMoveRequest) ToMessage() queue.ChessMoveMessage {
	return queue.ChessMoveMessage{
		ChessMove:   r.ChessMove,
		PublishedAt: time.Now(),
	}
}
