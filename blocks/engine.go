package blocks

import (
	"math/rand"
	"time"
)

type coord struct{ x, y int }

func newEngine(w, h int) *engine {
	return &engine{
		gridHeight: h, gridWidth: w,
		Stack: map[coord]pieceType{},
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type block struct {
	Coord coord
	Type  pieceType
}

type engine struct {
	gridHeight, gridWidth int
	Piece                 *piece
	Stack                 map[coord]pieceType
	Score                 int
	rand                  *rand.Rand
}

func (e *engine) ChangePieceDirection(direction string) {}

func (e *engine) MovePiece() {
	if e.Piece == nil {
		e.Piece = newRandomPiece(e.rand)
		e.Piece.Coord = coord{x: -2, y: 2 + e.rand.Intn(4)}
		for i := e.rand.Int63n(4); i > 0; i-- {
			e.Piece.Rotate()
		}
		return
	}

	if e.isPieceStopped(*e.Piece) {
		for _, c := range e.Piece.ToCoords() {
			e.Stack[c] = e.Piece.Type
		}
		e.Piece = nil
		e.removeColumns()
	} else {
		e.Piece.Coord.x++
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
			if _, ok := e.Stack[coord{x: x, y: y}]; !ok {
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
		m := map[coord]pieceType{}
		for c, t := range e.Stack {
			if int(c.x) < y {
				c.x++
				m[c] = t
			} else if int(c.x) > y {
				m[c] = t
			}
		}
		e.Stack = m
	}
}

func (e *engine) MovePieceUp() {
	if e.Piece != nil {
		copy := *e.Piece
		copy.Coord.y--
		if !e.isPieceStopped(copy) {
			e.Piece.Coord.y--
		}
	}
}

func (e *engine) MovePieceDown() {
	if e.Piece != nil {
		copy := *e.Piece
		copy.Coord.y++
		if !e.isPieceStopped(copy) {
			e.Piece.Coord.y++
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
	p.Coord.x++
	cs := p.ToCoords()
	for _, c := range cs {
		if c.x >= e.gridWidth || c.y < 0 || c.y >= e.gridHeight {
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
		if c.x < 0 {
			return true
		}
	}
	return false
}
