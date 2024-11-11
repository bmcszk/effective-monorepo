package game

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Game struct {
	ID    uuid.UUID `json:"id" bson:"id"`
	Board Board     `json:"board" bson:"board"`
}

func NewGame(id uuid.UUID, board Board) *Game {
	return &Game{
		ID:    id,
		Board: board,
	}
}

func (g *Game) Move(move string) error {
	var err error
	g.Board, err = g.Board.move(move)
	return err
}

type Board [64]Piece

const (
	StartingBoardStr = `
		♜♞♝♛♚♝♞♜
		♟♟♟♟♟♟♟♟




		♙♙♙♙♙♙♙♙
		♖♘♗♕♔♗♘♖`
)

var StartingBoard = NewBoard(StartingBoardStr)

func NewBoard(bstr string) Board {
	b := Board{}
	lines := strings.Split(bstr, "\n")
	lines = trimLines(lines)
	lines = normalize(lines)
	for y := 0; y < 8; y++ {
		line := []rune(lines[7-y])
		for x := 0; x < 8 && x < len(line); x++ {
			b[y*8+x] = NewPiece(line[x])
		}
	}
	return b
}

func (b Board) Empty() bool {
	for _, p := range b {
		if p != 0 {
			return false
		}
	}
	return true
}

func (b Board) String() string {
	if b.Empty() {
		return ""
	}
	lines := [8]string{}
	for y := 0; y < 8; y++ {
		var line [8]rune
		for x := 0; x < 8; x++ {
			line[x] = b[y*8+x].Rune()
		}
		lines[7-y] = string(line[:])
	}
	return strings.Join(lines[:], "\n")
}

func (b Board) set(square string, piece Piece) Board {
	b[squareIdx(square)] = piece
	return b
}

func (b Board) get(square string) Piece {
	return b[squareIdx(square)]
}

// move https://en.wikipedia.org/wiki/Universal_Chess_Interface
func (b Board) move(move string) (Board, error) {
	move = strings.TrimSpace(move)
	move = strings.ToLower(move)
	if len(move) < 4 {
		return b, errors.New("invalid move")
	}
	ps := move[:2]
	ts := move[2:4]
	if err := validateSquare(ps); err != nil {
		return b, err
	}
	if err := validateSquare(ts); err != nil {
		return b, err
	}
	p := b.get(ps)
	t := b.get(ts)
	if p == 0 {
		return b, errors.New("square is empty")
	}
	if t != 0 && t.Color() == p.Color() {
		if (p == '♔' && t == '♖') || (p == '♚' && t == '♜') {
			// TODO proper castle
			b = b.set(ts, p)
			b = b.set(ps, t)
			return b, nil
		}
		return b, errors.New("pieces are of the same color")
	}
	// TODO piece move rules
	// TODO promotions

	b = b.set(ps, 0)
	b = b.set(ts, p)
	return b, nil
}

func trimLines(lines []string) []string {
	for i, line := range lines {
		line = strings.TrimLeft(line, "\t")
		line = strings.TrimRight(line, " ")
		lines[i] = line
	}
	return lines
}

func normalize(lines []string) []string {
	l := len(lines)
	if l == 8 {
		return lines
	}
	if l < 8 {
		return normalize(append(lines, ""))
	}
	if l > 9 && lines[0] == "" && lines[l-1] == "" {
		return normalize(lines[1 : l-1])
	}
	if lines[0] == "" {
		return normalize(lines[1:])
	}
	if lines[l-1] == "" {
		return normalize(lines[:l-1])
	}
	return lines
}

func validateSquare(square string) error {
	idx := squareIdx(square)
	if idx > 63 {
		return errors.New("invalid square")
	}
	// TODO proper validation
	return nil
}

func squareIdx(square string) byte {
	return 8*(square[1]-'1') + (square[0] - 'a')
}

type Piece rune
type Color byte

const (
	Zero Color = iota
	White
	Black
)

func NewPiece(r rune) Piece {
	if r == ' ' {
		return 0
	}
	return Piece(r)
}

func (p Piece) Rune() rune {
	if p == 0 {
		return ' '
	}
	return rune(p)
}

func (p Piece) Color() Color {
	if p >= '♔' && p <= '♙' {
		return White
	}
	if p >= '♚' && p <= '♟' {
		return Black
	}
	return Zero
}
