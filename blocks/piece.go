package blocks

import (
	"math/rand"

	"github.com/richardlt/matrix/sdk-go/common"
)

type pieceType int

const (
	i pieceType = 0
	l pieceType = 1
	j pieceType = 2
	o pieceType = 3
	s pieceType = 4
	z pieceType = 5
	t pieceType = 6
)

func newRandomPiece() *piece { return newPiece(pieceType(rand.Intn(7))) }

func newPiece(t pieceType) *piece {
	switch t {
	case i:
		return &piece{Type: i, views: []view{{
			Width: 1, Height: 4, Center: common.Coord{X: 0, Y: 1},
			Mask: []bool{true, true, true, true},
		}, {
			Width: 4, Height: 1, Center: common.Coord{X: 2, Y: 0},
			Mask: []bool{true, true, true, true},
		}}}
	case l:
		return &piece{Type: l, views: []view{{
			Width: 2, Height: 3, Center: common.Coord{X: 0, Y: 1},
			Mask: []bool{true, false, true, false, true, true},
		}, {
			Width: 3, Height: 2, Center: common.Coord{X: 1, Y: 0},
			Mask: []bool{true, true, true, true, false, false},
		}, {
			Width: 2, Height: 3, Center: common.Coord{X: 1, Y: 1},
			Mask: []bool{true, true, false, true, false, true},
		}, {
			Width: 3, Height: 2, Center: common.Coord{X: 1, Y: 1},
			Mask: []bool{false, false, true, true, true, true},
		}}}
	case j:
		return &piece{Type: j, views: []view{{
			Width: 2, Height: 3, Center: common.Coord{X: 1, Y: 1},
			Mask: []bool{false, true, false, true, true, true},
		}, {
			Width: 3, Height: 2, Center: common.Coord{X: 1, Y: 1},
			Mask: []bool{true, false, false, true, true, true},
		}, {
			Width: 2, Height: 3, Center: common.Coord{X: 0, Y: 1},
			Mask: []bool{true, true, true, false, true, false},
		}, {
			Width: 3, Height: 2, Center: common.Coord{X: 1, Y: 0},
			Mask: []bool{true, true, true, false, false, true},
		}}}
	case o:
		return &piece{Type: o, views: []view{{
			Width: 2, Height: 2, Center: common.Coord{X: 0, Y: 1},
			Mask: []bool{true, true, true, true},
		}}}
	case s:
		return &piece{Type: s, views: []view{{
			Width: 3, Height: 2, Center: common.Coord{X: 1, Y: 1},
			Mask: []bool{false, true, true, true, true, false},
		}, {
			Width: 2, Height: 3, Center: common.Coord{X: 0, Y: 1},
			Mask: []bool{true, false, true, true, false, true},
		}}}
	case z:
		return &piece{Type: z, views: []view{{
			Width: 3, Height: 2, Center: common.Coord{X: 1, Y: 1},
			Mask: []bool{true, true, false, false, true, true},
		}, {
			Width: 2, Height: 3, Center: common.Coord{X: 1, Y: 1},
			Mask: []bool{false, true, true, true, true, false},
		}}}
	}
	return &piece{Type: t, views: []view{{
		Width: 3, Height: 2, Center: common.Coord{X: 1, Y: 0},
		Mask: []bool{true, true, true, false, true, false},
	}, {
		Width: 2, Height: 3, Center: common.Coord{X: 1, Y: 1},
		Mask: []bool{false, true, true, true, false, true},
	}, {
		Width: 3, Height: 2, Center: common.Coord{X: 1, Y: 1},
		Mask: []bool{false, true, false, true, true, true},
	}, {
		Width: 2, Height: 3, Center: common.Coord{X: 0, Y: 1},
		Mask: []bool{true, false, true, true, true, false},
	}}}
}

type view struct {
	Width, Height int64
	Mask          []bool
	Center        common.Coord
}

type piece struct {
	Type  pieceType
	views []view
	angle int
	Coord common.Coord
}

func (p *piece) Rotate() {
	if p.angle < len(p.views)-1 {
		p.angle++
	} else {
		p.angle = 0
	}
}

func (p *piece) ToCoords() (coords []common.Coord) {
	v := p.views[p.angle]
	for y := int64(0); y < v.Height; y++ {
		diffY := y - v.Center.Y
		for x := int64(0); x < v.Width; x++ {
			diffX := x - v.Center.X
			if v.Mask[x+y*v.Width] {
				coords = append(coords, common.Coord{X: p.Coord.X + diffX, Y: p.Coord.Y + diffY})
			}
		}
	}
	return
}
