package e2e

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bmcszk/effective-monorepo/pkg/model"
	"github.com/bmcszk/effective-monorepo/services/consumer/repo"
	"github.com/bmcszk/effective-monorepo/services/producer/producer"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Block struct {
	*testing.T
	ctx           context.Context
	client        *http.Client
	producerUri   string
	repo          *repo.Repo
	gameID        uuid.UUID
	move          string
	respReturned  *http.Response
	errReturned   error
	boardReturned string
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		slog.Warn("failed to load .env")
	}
}

func NewBlocks(t *testing.T) (*Block, *Block, *Block) {
	ctx, cancel := context.WithCancel(context.Background())
	repo, err := repo.NewRepo()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cancel()
		repo.Close()
	})
	b := &Block{
		T:           t,
		ctx:         ctx,
		client:      http.DefaultClient,
		producerUri: os.Getenv("PRODUCER_URI"),
		repo:        repo,
	}
	return b, b, b
}

func (b *Block) and() *Block {
	return b
}

func (b *Block) aGame() *Block {
	b.gameID = uuid.New()
	return b
}

func (b *Block) aWhiteOpeningMove() *Block {
	b.move = "e2e4"
	return b
}

func (b *Block) dispatchingMove() *Block {
	request := producer.ChessMoveRequest{
		ChessMove: model.ChessMove{
			ID:     uuid.New(),
			GameID: b.gameID,
			SentAt: time.Now(),
			Move:   b.move,
		},
	}
	requestBody, err := request.ToJSON()
	if err != nil {
		b.T.Fatal(err)
	}
	b.respReturned, b.errReturned = b.client.Post(b.producerUri+"/moves", "application/json", bytes.NewReader(requestBody))
	return b
}

func (b *Block) moveIsDispatched() *Block {
	if b.errReturned != nil {
		b.T.Fatal(b.errReturned)
	}
	if b.respReturned.StatusCode != http.StatusCreated {
		b.T.Fatal("expected status code 201")
	}
	return b
}

func (b *Block) fetchingBoard() *Block {
	b.boardReturned, b.errReturned = b.repo.Get(b.ctx, b.gameID.String())
	return b
}

func (b *Block) boardIsFetched() *Block {
	if b.errReturned != nil {
		b.T.Fatal(b.errReturned)
	}
	if b.boardReturned == "" {
		b.T.Fatal("expected board to be not empty")
	}
	return b
}
