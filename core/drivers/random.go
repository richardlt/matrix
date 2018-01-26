package drivers

import (
	"math/rand"

	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
)

// NewRandom returns a new random driver.
func NewRandom(f *render.Frame) *Random { return &Random{frame: f} }

// Random driver allows to render a random colored frame.
type Random struct {
	frame       *render.Frame
	endCallback func()
}

// Render print colors in frame.
func (r *Random) Render() {
	for y := int64(0); y < int64(r.frame.Height); y++ {
		for x := int64(0); x < int64(r.frame.Width); x++ {
			r.frame.SetWithCoord(common.Coord{X: x, Y: y}, common.Color{
				R: uint64(rand.Intn(255)),
				G: uint64(rand.Intn(255)),
				B: uint64(rand.Intn(255)),
				A: 1,
			})
		}
	}

	if r.endCallback != nil {
		r.endCallback()
	}
}

// OnEnd allows to set end callback.
func (r *Random) OnEnd(f func()) { r.endCallback = f }
