package blocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCoords(t *testing.T) {
	assert := assert.New(t)

	p := newPiece(l)
	p.Coord = coord{x: 5, y: 5}

	cs := p.ToCoords()

	assert.Equal(4, len(cs))
	assert.Equal(coord{x: 5, y: 4}, cs[0])
	assert.Equal(coord{x: 5, y: 5}, cs[1])
	assert.Equal(coord{x: 5, y: 6}, cs[2])
	assert.Equal(coord{x: 6, y: 6}, cs[3])

	p.Coord = coord{x: 0, y: 0}

	cs = p.ToCoords()

	assert.Equal(4, len(cs))
	assert.Equal(coord{x: 0, y: -1}, cs[0])
	assert.Equal(coord{x: 0, y: 0}, cs[1])
	assert.Equal(coord{x: 0, y: 1}, cs[2])
	assert.Equal(coord{x: 1, y: 1}, cs[3])

	p.Coord = coord{x: -1, y: -1}

	cs = p.ToCoords()

	assert.Equal(4, len(cs))
	assert.Equal(coord{x: -1, y: -2}, cs[0])
	assert.Equal(coord{x: -1, y: -1}, cs[1])
	assert.Equal(coord{x: -1, y: 0}, cs[2])
	assert.Equal(coord{x: 0, y: 0}, cs[3])
}
