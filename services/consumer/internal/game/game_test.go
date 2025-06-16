package game_test

import (
	"testing"

	"github.com/bmcszk/effective-monorepo/services/consumer/internal/game"
	"github.com/google/uuid"
)

func TestGame_Move(t *testing.T) {
	gameUUID := uuid.New()

	tests := []struct {
		name          string
		board         game.Board
		expectedBoard game.Board
		move          string
		wantErr       bool
		expectedErr   string
	}{
		{
			name:  "move",
			board: game.NewBoard(game.StartingBoardStr),
			expectedBoard: game.NewBoard(`
				♜♞♝♛♚♝♞♜
				♟♟♟♟♟♟♟♟


				♙

				 ♙♙♙♙♙♙♙
				♖♘♗♕♔♗♘♖
				`),
			move:        "a2a4",
			wantErr:     false,
			expectedErr: "",
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
			move:        "g1f3",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "almost empty board",
			board: game.NewBoard(`
				



				
				     ♘
				

				`),
			expectedBoard: game.NewBoard(`
				


				    ♘
				    
				     
				

				`),
			move:        "f3e5",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "pieces are of the same color",
			board:       game.NewBoard(game.StartingBoardStr),
			move:        "a2b2",
			wantErr:     true,
			expectedErr: "pieces are of the same color",
		},
		{
			name:        "square is empty",
			board:       game.NewBoard(game.StartingBoardStr),
			move:        "a3a4",
			wantErr:     true,
			expectedErr: "square is empty",
		},
		{
			name:        "invalid square",
			board:       game.NewBoard(game.StartingBoardStr),
			move:        "i8i9",
			wantErr:     true,
			expectedErr: "invalid square",
		},
		{
			name:        "invalid move",
			board:       game.NewBoard(game.StartingBoardStr),
			move:        "",
			wantErr:     true,
			expectedErr: "invalid move",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := game.NewGame(gameUUID, tt.board)
			err := g.Move(tt.move)
			if (err != nil) != tt.wantErr {
				t.Errorf("Game.Move() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.expectedErr {
				t.Errorf("Game.Move() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			if tt.expectedBoard.String() != "" && g.Board != tt.expectedBoard {
				t.Errorf("Game.Move() board = %v, expectedBoard %v", g.Board, tt.expectedBoard)
			}
		})
	}
}
