package model

import (
	"time"

	"github.com/google/uuid"
)

type ChessMove struct {
	ID     uuid.UUID `json:"id" bson:"id"`
	SentAt time.Time `json:"sent_at" bson:"sent_at"`
	GameID uuid.UUID `json:"game_id" bson:"game_id"`
	Move   string    `json:"move" bson:"move"`
}
