package blocks

import (
	"math/rand"
	"time"

	"github.com/richardlt/matrix/sdk-go/common"
)

func newEngine(w, h int) *engine {
	return &engine{gridHeight: h, gridWidth: w, Stack: map[common.Coord]pieceType{}}
}

type block struct {
	Coord common.Coord
	Type  pieceType
}

type engine struct {
	gridHeight, gridWidth int
	Piece                 *piece
	Stack                 map[common.Coord]pieceType
	Score                 int
}

func (e *engine) ChangePieceDirection(direction string) {}

func (e *engine) MovePiece() {
	if e.Piece == nil {
		rand := rand.New(rand.NewSource(time.Now().Unix()))
		e.Piece = newRandomPiece(rand)
		e.Piece.Coord = common.Coord{X: -2, Y: 2 + rand.Int63n(4)}
		return
	}

	if e.isPieceStopped(*e.Piece) {
		for _, c := range e.Piece.ToCoords() {
			e.Stack[c] = e.Piece.Type
		}
		e.Piece = nil
		e.removeColumns()
	} else {
		e.Piece.Coord.X++
	}
}

func fibo(n int) int {
	if n <= 1 {
		return n
	}
	return fibo(n-1) + fibo(n-2)
}

func (e engine) getScoreFromColumn(cs []int) (score int) {
	count := 0
	for i, c := range cs {
		count++
		if len(cs) <= i+1 || cs[i+1] != c+1 || count == 4 {
			if count > 1 {
				score += fibo(count + 2)
			} else {
				score++
			}
			count = 0
		}
	}
	return
}

func (e *engine) removeColumns() {
	var columnsFull []int
	for x := 0; x < e.gridWidth; x++ {
		full := true
		for y := 0; y < e.gridHeight; y++ {
			if _, ok := e.Stack[common.Coord{X: int64(x), Y: int64(y)}]; !ok {
				full = false
				break
			}
		}
		if full {
			columnsFull = append(columnsFull, x)
		}
	}

	e.Score += e.getScoreFromColumn(columnsFull)

	for _, y := range columnsFull {
		m := map[common.Coord]pieceType{}
		for c, t := range e.Stack {
			if int(c.X) < y {
				c.X++
				m[c] = t
			} else if int(c.X) > y {
				m[c] = t
			}
		}
		e.Stack = m
	}
}

func (e *engine) MovePieceUp() {
	if e.Piece != nil {
		copy := *e.Piece
		copy.Coord.Y--
		if !e.isPieceStopped(copy) {
			e.Piece.Coord.Y--
		}
	}
}

func (e *engine) MovePieceDown() {
	if e.Piece != nil {
		copy := *e.Piece
		copy.Coord.Y++
		if !e.isPieceStopped(copy) {
			e.Piece.Coord.Y++
		}
	}
}

func (e *engine) RotatePiece() {
	if e.Piece != nil {
		copy := *e.Piece
		copy.Rotate()
		if !e.isPieceStopped(copy) {
			e.Piece.Rotate()
		}
	}
}

func (e *engine) isPieceStopped(p piece) bool {
	p.Coord.X++
	cs := p.ToCoords()
	for _, c := range cs {
		if int(c.X) >= e.gridWidth || c.Y < 0 || int(c.Y) >= e.gridHeight {
			return true
		}

		if _, ok := e.Stack[c]; ok {
			return true
		}
	}

	return false
}

func (e *engine) IsGameOver() bool {
	for c := range e.Stack {
		if c.X < 0 {
			return true
		}
	}
	return false
}
