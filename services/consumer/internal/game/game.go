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
	parsedMove, err := parseMove(move)
	if err != nil {
		return b, err
	}

	piece := b.get(parsedMove.from)
	if piece == 0 {
		return b, errors.New("square is empty")
	}

	return b.executeMove(parsedMove, piece)
}

type parsedMove struct {
	from string
	to   string
}

func parseMove(move string) (parsedMove, error) {
	move = strings.TrimSpace(strings.ToLower(move))
	if len(move) < 4 {
		return parsedMove{}, errors.New("invalid move")
	}

	from := move[:2]
	to := move[2:4]

	if err := validateSquare(from); err != nil {
		return parsedMove{}, err
	}
	if err := validateSquare(to); err != nil {
		return parsedMove{}, err
	}

	return parsedMove{from: from, to: to}, nil
}

func (b Board) executeMove(move parsedMove, piece Piece) (Board, error) {
	target := b.get(move.to)

	if target != 0 && target.Color() == piece.Color() {
		if isCastleMove(piece, target) {
			return b.performCastle(move, piece, target), nil
		}
		return b, errors.New("pieces are of the same color")
	}

	// TODO piece move rules
	// TODO promotions
	return b.performRegularMove(move, piece), nil
}

func isCastleMove(piece, target Piece) bool {
	return (piece == '♔' && target == '♖') || (piece == '♚' && target == '♜')
}

func (b Board) performCastle(move parsedMove, piece, target Piece) Board {
	// TODO proper castle logic
	return b.set(move.to, piece).set(move.from, target)
}

func (b Board) performRegularMove(move parsedMove, piece Piece) Board {
	return b.set(move.from, 0).set(move.to, piece)
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

	lines = padLines(lines, l)
	lines = trimEmptyEdges(lines)

	return normalize(lines) // Recursive call to handle remaining cases
}

func padLines(lines []string, l int) []string {
	if l < 8 {
		return append(lines, "")
	}
	return lines
}

func trimEmptyEdges(lines []string) []string {
	l := len(lines)
	if l <= 9 {
		return lines
	}

	if lines[0] == "" && lines[l-1] == "" {
		return lines[1 : l-1]
	}
	if lines[0] == "" {
		return lines[1:]
	}
	if lines[l-1] == "" {
		return lines[:l-1]
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
