package blocks

import (
	"testing"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/stretchr/testify/assert"
)

func TestToCoords(t *testing.T) {
	assert := assert.New(t)

	p := newPiece(l)
	p.Coord = common.Coord{X: 5, Y: 5}

	cs := p.ToCoords()

	assert.Equal(4, len(cs))
	assert.Equal(common.Coord{X: 5, Y: 4}, cs[0])
	assert.Equal(common.Coord{X: 5, Y: 5}, cs[1])
	assert.Equal(common.Coord{X: 5, Y: 6}, cs[2])
	assert.Equal(common.Coord{X: 6, Y: 6}, cs[3])

	p.Coord = common.Coord{X: 0, Y: 0}

	cs = p.ToCoords()

	assert.Equal(4, len(cs))
	assert.Equal(common.Coord{X: 0, Y: -1}, cs[0])
	assert.Equal(common.Coord{X: 0, Y: 0}, cs[1])
	assert.Equal(common.Coord{X: 0, Y: 1}, cs[2])
	assert.Equal(common.Coord{X: 1, Y: 1}, cs[3])

	p.Coord = common.Coord{X: -1, Y: -1}

	cs = p.ToCoords()

	assert.Equal(4, len(cs))
	assert.Equal(common.Coord{X: -1, Y: -2}, cs[0])
	assert.Equal(common.Coord{X: -1, Y: -1}, cs[1])
	assert.Equal(common.Coord{X: -1, Y: 0}, cs[2])
	assert.Equal(common.Coord{X: 0, Y: 0}, cs[3])
}
