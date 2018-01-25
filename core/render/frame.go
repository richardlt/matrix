package render

import (
	"image/color"
	"math/rand"

	"github.com/richardlt/matrix/sdk-go/common"
)

// NewFrame returns a clean transparent frame.
func NewFrame(w, h uint64) Frame {
	return Frame{Width: w, Height: h, Pixels: make([]color.RGBA, w*h)}
}

// NewFrameRandom returns a frame with random colors.
func NewFrameRandom(w, h uint64) Frame {
	ps := make([]color.RGBA, w*h)
	for i := uint64(0); i < w*h; i++ {
		ps[i].R = uint8(rand.Intn(255))
		ps[i].G = uint8(rand.Intn(255))
		ps[i].B = uint8(rand.Intn(255))
		ps[i].A = 1
	}
	return Frame{Width: w, Height: h, Pixels: ps}
}

// Frame is a rectangle with given width
// and height that contains pixels.
type Frame struct {
	Width, Height uint64
	Pixels        []color.RGBA
}

// SetWithCoord allows to set a pixel by coord.
func (f *Frame) SetWithCoord(coo common.Coord, col color.RGBA) {
	i := int(coo.X + coo.Y*f.Width)
	if i < len(f.Pixels) {
		f.Pixels[i] = col
	}
}

// Clean set all pixels to transparent.
func (f *Frame) Clean() { f.Pixels = make([]color.RGBA, f.Width*f.Height) }
