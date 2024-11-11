package game

import (
	"testing"

	"github.com/google/uuid"
)

func TestGame_Move(t *testing.T) {
	uuid := uuid.New()

	tests := []struct {
		name          string
		board         Board
		expectedBoard Board
		move          string
		wantErr       bool
		expectedErr   string
	}{
		{
			name:  "move",
			board: NewBoard(StartingBoardStr),
			expectedBoard: NewBoard(`
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
			board: NewBoard(StartingBoardStr),
			expectedBoard: NewBoard(`
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
			board: NewBoard(`
				



				
				     ♘
				

				`),
			expectedBoard: NewBoard(`
				


				    ♘
				    
				     
				

				`),
			move:        "f3e5",
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:        "pieces are of the same color",
			board:       NewBoard(StartingBoardStr),
			move:        "a2b2",
			wantErr:     true,
			expectedErr: "pieces are of the same color",
		},
		{
			name:        "square is empty",
			board:       NewBoard(StartingBoardStr),
			move:        "a3a4",
			wantErr:     true,
			expectedErr: "square is empty",
		},
		{
			name:        "invalid square",
			board:       NewBoard(StartingBoardStr),
			move:        "i8i9",
			wantErr:     true,
			expectedErr: "invalid square",
		},
		{
			name:        "invalid move",
			board:       NewBoard(StartingBoardStr),
			move:        "",
			wantErr:     true,
			expectedErr: "invalid move",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame(uuid, tt.board)
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
