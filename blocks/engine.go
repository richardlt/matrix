package blocks

import (
	"math/rand"

	"github.com/richardlt/matrix/sdk-go/common"
)

func newEngine(w, h uint64) *engine {
	return &engine{gridHeight: h, gridWidth: w, blocks: map[int]block{}}
}

type block struct {
	Coord common.Coord
	Type  pieceType
}

type engine struct {
	gridHeight, gridWidth uint64
	piece                 *piece
	blocks                map[int]block
}

func (e *engine) ChangePieceDirection(direction string) {}

func (e *engine) MovePiece() {
	if e.piece == nil {
		e.piece = newRandomPiece()
		e.piece.Coord = common.Coord{X: 0, Y: int64(2 + rand.Intn(4))}
		return
	}

	if e.isPieceStopped(*e.piece) {
		for _, c := range e.piece.ToCoords() {
			i := c.X + c.Y*int64(e.gridWidth)
			e.blocks[int(i)] = block{c, e.piece.Type}
		}
		e.piece = nil

		e.removeColumns()
	} else {
		e.piece.Coord.X++
	}
}

func (e *engine) removeColumns() {
	var columnsFull []int
	for x := 0; x < int(e.gridWidth); x++ {
		full := true
		for y := 0; y < int(e.gridHeight); y++ {
			if _, ok := e.blocks[x+y*int(e.gridWidth)]; !ok {
				full = false
				break
			}
		}
		if full {
			columnsFull = append(columnsFull, x)
		}
	}

	for _, c := range columnsFull {
		newMap := map[int]block{}
		for index, b := range e.blocks {
			if int(b.Coord.X) < c {
				b.Coord.X++
				i := b.Coord.X + b.Coord.Y*int64(e.gridWidth)
				newMap[int(i)] = b
			} else if int(b.Coord.X) > c {
				newMap[index] = b
			}
		}
		e.blocks = newMap
	}
}

func (e *engine) MovePieceUp() {
	if e.piece != nil {
		copy := *e.piece
		copy.Coord.Y--
		if !e.isPieceStopped(copy) {
			e.piece.Coord.Y--
		}
	}
}

func (e *engine) MovePieceDown() {
	if e.piece != nil {
		copy := *e.piece
		copy.Coord.Y++
		if !e.isPieceStopped(copy) {
			e.piece.Coord.Y++
		}
	}
}

func (e *engine) RotatePiece() {
	if e.piece != nil {
		copy := *e.piece
		copy.Rotate()
		if !e.isPieceStopped(copy) {
			e.piece.Rotate()
		}
	}
}

func (e *engine) isPieceStopped(p piece) bool {
	p.Coord.X++
	cs := p.ToCoords()
	for _, c := range cs {
		if c.X >= int64(e.gridWidth) || c.Y < 0 || c.Y >= int64(e.gridHeight) {
			return true
		}

		i := c.X + c.Y*int64(e.gridWidth)
		if _, ok := e.blocks[int(i)]; ok {
			return true
		}
	}

	return false
}

func (e *engine) IsGameOver() bool {
	for i := range e.blocks {
		if i < 0 {
			return true
		}
	}
	return false
}

func (e engine) GetBlocks() []block {
	var blocks []block
	for _, b := range e.blocks {
		blocks = append(blocks, b)
	}
	return blocks
}

func (e engine) GetPiece() *piece { return e.piece }
