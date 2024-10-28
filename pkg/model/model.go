package model

import (
	"time"

	"github.com/google/uuid"
)

type ChessMove struct {
	ID     uuid.UUID `json:"id" bson:"id"`
	SentAt time.Time `json:"sent_at" bson:"sent_at"`
	Piece  string    `json:"piece" bson:"piece"`
	Number int       `json:"number" bson:"number"`
	Move   string    `json:"move" bson:"move"`
}
