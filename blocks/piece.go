package blocks

import (
	"math/rand"
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

func newRandomPiece(rand *rand.Rand) *piece {
	return newPiece(pieceType(rand.Intn(7)))
}

func newPiece(t pieceType) *piece {
	switch t {
	case i:
		return &piece{Type: i, views: []view{{
			Width: 1, Height: 4, Center: coord{x: 0, y: 1},
			Mask: []bool{true, true, true, true},
		}, {
			Width: 4, Height: 1, Center: coord{x: 2, y: 0},
			Mask: []bool{true, true, true, true},
		}}}
	case l:
		return &piece{Type: l, views: []view{{
			Width: 2, Height: 3, Center: coord{x: 0, y: 1},
			Mask: []bool{true, false, true, false, true, true},
		}, {
			Width: 3, Height: 2, Center: coord{x: 1, y: 0},
			Mask: []bool{true, true, true, true, false, false},
		}, {
			Width: 2, Height: 3, Center: coord{x: 1, y: 1},
			Mask: []bool{true, true, false, true, false, true},
		}, {
			Width: 3, Height: 2, Center: coord{x: 1, y: 1},
			Mask: []bool{false, false, true, true, true, true},
		}}}
	case j:
		return &piece{Type: j, views: []view{{
			Width: 2, Height: 3, Center: coord{x: 1, y: 1},
			Mask: []bool{false, true, false, true, true, true},
		}, {
			Width: 3, Height: 2, Center: coord{x: 1, y: 1},
			Mask: []bool{true, false, false, true, true, true},
		}, {
			Width: 2, Height: 3, Center: coord{x: 0, y: 1},
			Mask: []bool{true, true, true, false, true, false},
		}, {
			Width: 3, Height: 2, Center: coord{x: 1, y: 0},
			Mask: []bool{true, true, true, false, false, true},
		}}}
	case o:
		return &piece{Type: o, views: []view{{
			Width: 2, Height: 2, Center: coord{x: 0, y: 1},
			Mask: []bool{true, true, true, true},
		}}}
	case s:
		return &piece{Type: s, views: []view{{
			Width: 3, Height: 2, Center: coord{x: 1, y: 1},
			Mask: []bool{false, true, true, true, true, false},
		}, {
			Width: 2, Height: 3, Center: coord{x: 0, y: 1},
			Mask: []bool{true, false, true, true, false, true},
		}}}
	case z:
		return &piece{Type: z, views: []view{{
			Width: 3, Height: 2, Center: coord{x: 1, y: 1},
			Mask: []bool{true, true, false, false, true, true},
		}, {
			Width: 2, Height: 3, Center: coord{x: 1, y: 1},
			Mask: []bool{false, true, true, true, true, false},
		}}}
	}
	return &piece{Type: t, views: []view{{
		Width: 3, Height: 2, Center: coord{x: 1, y: 0},
		Mask: []bool{true, true, true, false, true, false},
	}, {
		Width: 2, Height: 3, Center: coord{x: 1, y: 1},
		Mask: []bool{false, true, true, true, false, true},
	}, {
		Width: 3, Height: 2, Center: coord{x: 1, y: 1},
		Mask: []bool{false, true, false, true, true, true},
	}, {
		Width: 2, Height: 3, Center: coord{x: 0, y: 1},
		Mask: []bool{true, false, true, true, true, false},
	}}}
}

type view struct {
	Width, Height int
	Mask          []bool
	Center        coord
}

type piece struct {
	Type  pieceType
	views []view
	angle int
	Coord coord
}

func (p *piece) Rotate() {
	if p.angle < len(p.views)-1 {
		p.angle++
	} else {
		p.angle = 0
	}
}

func (p *piece) ToCoords() (coords []coord) {
	v := p.views[p.angle]
	for y := 0; y < v.Height; y++ {
		diffY := y - v.Center.y
		for x := 0; x < v.Width; x++ {
			diffX := x - v.Center.x
			if v.Mask[x+y*v.Width] {
				coords = append(coords, coord{x: p.Coord.x + diffX, y: p.Coord.y + diffY})
			}
		}
	}
	return
}
