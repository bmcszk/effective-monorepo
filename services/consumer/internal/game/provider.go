package game

import (
	"context"

	"github.com/bmcszk/effective-monorepo/services/consumer/repo"
	"github.com/google/uuid"
)

type GameProvider struct {
	repo *repo.Repo
}

func NewGameProvider(repo *repo.Repo) *GameProvider {
	return &GameProvider{
		repo: repo,
	}
}

func (p *GameProvider) HandleMove(ctx context.Context, gameID uuid.UUID, move string) error {
	bstr, err := p.repo.Get(ctx, gameID.String())
	if err != nil {
		return err
	}
	var board Board
	if bstr != "" {
		board = NewBoard(bstr)
	} else {
		board = NewBoard(StartingBoardStr)
	}
	game := NewGame(gameID, board)
	if err := game.Move(move); err != nil {
		return err
	}
	if err := p.repo.Put(ctx, gameID.String(), game.Board.String()); err != nil {
		return err
	}

	return nil
}
