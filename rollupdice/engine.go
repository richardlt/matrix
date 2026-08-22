package rollupdice

import "math/rand"

// diceCount dice are rolled together, one per half of the display.
const diceCount = 2

func newEngine() *engine { return &engine{} }

type engine struct {
	values [diceCount]int
}

// roll draws a new face, one to six, for every die. The global source is used rather than
// one seeded from the clock: the board this runs on has no clock to seed from, and the
// runtime seeds the global source itself.
func (e *engine) roll() {
	for i := range e.values {
		e.values[i] = rand.Intn(6) + 1
	}
}
