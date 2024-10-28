package queue

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmcszk/effective-monorepo/pkg/model"
)

type ChessMoveMessage struct {
	model.ChessMove
	PublishedAt time.Time `json:"published_at" bson:"published_at"`
}

func (m ChessMoveMessage) ToQueueMessage() (*message.Message, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(m); err != nil {
		return nil, err
	}
	return message.NewMessage(m.ID.String(), buf.Bytes()), nil
}

func FromQueueMessage(msg *message.Message) (ChessMoveMessage, error) {
	var m ChessMoveMessage
	if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&m); err != nil {
		return m, err
	}
	return m, nil
}
