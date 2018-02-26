package blocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetScoreFromColumn(t *testing.T) {
	assert := assert.New(t)

	e := newEngine(16, 9)
	assert.Equal(0, e.getScoreFromColumn([]int{}))
	assert.Equal(1, e.getScoreFromColumn([]int{0}))
	assert.Equal(9, e.getScoreFromColumn([]int{0, 5, 6, 10, 11, 12}))
	assert.Equal(16, e.getScoreFromColumn([]int{5, 6, 7, 8, 9, 10, 11, 12}))
}
