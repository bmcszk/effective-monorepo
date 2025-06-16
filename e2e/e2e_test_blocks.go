package e2e

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bmcszk/effective-monorepo/pkg/model"
	"github.com/bmcszk/effective-monorepo/services/consumer/repo"
	"github.com/bmcszk/effective-monorepo/services/producer/producer"
	"github.com/bmcszk/xmlsurf"
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

func NewBlocks(t *testing.T) (firstBlock, secondBlock, thirdBlock *Block) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	repository, err := repo.NewRepo()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cancel()
		repository.Close()
	})
	b := &Block{
		T:           t,
		ctx:         ctx,
		client:      http.DefaultClient,
		producerUri: os.Getenv("PRODUCER_URI"),
		repo:        repository,
	}
	return b, b, b
}

func (b *Block) And() *Block {
	return b
}

func (b *Block) AGame() *Block {
	b.gameID = uuid.New()
	return b
}

func (b *Block) AWhiteOpeningMove() *Block {
	b.move = "e2e4"
	return b
}

func (b *Block) DispatchingMove() *Block {
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
		b.Fatal(err)
	}
	resp, err := b.client.Post(
		b.producerUri+"/moves",
		"application/json",
		bytes.NewReader(requestBody),
	)
	if err == nil && resp != nil {
		defer resp.Body.Close()
	}
	b.respReturned = resp
	b.errReturned = err
	return b
}

func (b *Block) MoveIsDispatched() *Block {
	if b.errReturned != nil {
		b.Fatal(b.errReturned)
	}
	if b.respReturned != nil && b.respReturned.StatusCode != http.StatusCreated {
		b.Fatal("expected status code 201")
	}
	return b
}

func (b *Block) FetchingBoard() *Block {
	b.boardReturned, b.errReturned = b.repo.Get(b.ctx, b.gameID.String())
	return b
}

func (b *Block) BoardIsFetched() *Block {
	if b.errReturned != nil {
		b.Fatal(b.errReturned)
	}
	_, _ = xmlsurf.ParseToMap(strings.NewReader("asdf"))
	if b.boardReturned == "" {
		b.Fatal("expected board to be not empty")
	}
	return b
}
