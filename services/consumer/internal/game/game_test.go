package game_test

import (
	"testing"

	"github.com/bmcszk/effective-monorepo/services/consumer/internal/game"
	"github.com/google/uuid"
)

func TestGame_Move_ValidMoves(t *testing.T) {
	testValidMoves(t)
}

func TestGame_Move_InvalidMoves(t *testing.T) {
	testInvalidMoves(t)
}

func testValidMoves(t *testing.T) {
	t.Helper()
	gameUUID := uuid.New()

	tests := []struct {
		name          string
		board         game.Board
		expectedBoard game.Board
		move          string
	}{
		{
			name:  "pawn move",
			board: game.NewBoard(game.StartingBoardStr),
			expectedBoard: game.NewBoard(`
				♜♞♝♛♚♝♞♜
				♟♟♟♟♟♟♟♟


				♙

				 ♙♙♙♙♙♙♙
				♖♘♗♕♔♗♘♖
				`),
			move: "a2a4",
		},
		{
			name:  "knight move",
			board: game.NewBoard(game.StartingBoardStr),
			expectedBoard: game.NewBoard(`
				♜♞♝♛♚♝♞♜
				♟♟♟♟♟♟♟♟


				
				     ♘
				♙♙♙♙♙♙♙♙
				♖♘♗♕♔♗ ♖
				`),
			move: "g1f3",
		},
		{
			name: "knight move on empty board",
			board: game.NewBoard(`
				



				
				     ♘
				

				`),
			expectedBoard: game.NewBoard(`
				


				    ♘
				    
				     
				

				`),
			move: "f3e5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := game.NewGame(gameUUID, tt.board)
			err := g.Move(tt.move)
			if err != nil {
				t.Errorf("Game.Move() unexpected error = %v", err)
			}
			if g.Board != tt.expectedBoard {
				t.Errorf("Game.Move() board = %v, expectedBoard %v", g.Board, tt.expectedBoard)
			}
		})
	}
}

func testInvalidMoves(t *testing.T) {
	t.Helper()
	gameUUID := uuid.New()

	tests := []struct {
		name        string
		board       game.Board
		move        string
		expectedErr string
	}{
		{
			name:        "pieces are of the same color",
			board:       game.NewBoard(game.StartingBoardStr),
			move:        "a2b2",
			expectedErr: "pieces are of the same color",
		},
		{
			name:        "square is empty",
			board:       game.NewBoard(game.StartingBoardStr),
			move:        "a3a4",
			expectedErr: "square is empty",
		},
		{
			name:        "invalid square",
			board:       game.NewBoard(game.StartingBoardStr),
			move:        "i8i9",
			expectedErr: "invalid square",
		},
		{
			name:        "invalid move",
			board:       game.NewBoard(game.StartingBoardStr),
			move:        "",
			expectedErr: "invalid move",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := game.NewGame(gameUUID, tt.board)
			err := g.Move(tt.move)
			if err == nil {
				t.Errorf("Game.Move() expected error but got none")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Game.Move() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
