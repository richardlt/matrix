package render

import (
	"github.com/richardlt/matrix/sdk-go/common"
)

// NewFrame returns a clean transparent frame.
func NewFrame(w, h uint64) Frame {
	f := Frame{Width: w, Height: h}
	f.Clean()
	return f
}

// Frame is a rectangle with given width
// and height that contains pixels.
type Frame struct {
	common.Frame
	Width, Height uint64
}

// SetWithCoord allows to set a pixel by coord.
func (f *Frame) SetWithCoord(coo common.Coord, col common.Color) {
	if coo.X >= 0 && coo.Y >= 0 {
		i := int(coo.X) + int(coo.Y)*int(f.Width)
		if 0 <= i && i < len(f.Pixels) {
			f.Pixels[i] = &col
		}
	}
}

// Clean set all pixels to transparent.
func (f *Frame) Clean() {
	f.Pixels = make([]*common.Color, f.Width*f.Height)
	for i := 0; i < len(f.Pixels); i++ {
		f.Pixels[i] = &common.Color{}
	}
}
