package rollupdice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Only the six faces of the DiceValues font can be rendered, so a value outside one to
// six would print an empty face rather than a die.
func TestRollStaysOnTheFaces(t *testing.T) {
	assert := assert.New(t)

	e := newEngine()
	seen := map[int]bool{}

	for i := 0; i < 1000; i++ {
		e.roll()
		for _, v := range e.values {
			assert.GreaterOrEqual(v, 1)
			assert.LessOrEqual(v, 6)
			seen[v] = true
		}
	}

	assert.Len(seen, 6)
}
