package render

import (
	"github.com/richardlt/matrix/sdk-go/common"
)

// NewFrame returns a clean transparent frame.
func NewFrame(w, h uint64) Frame {
	f := Frame{Frame: &common.Frame{}, Width: w, Height: h}
	f.Clean()
	return f
}

// Frame is a rectangle with given width
// and height that contains pixels.
//
// common.Frame is embedded by pointer: it is a protobuf message and carries a
// no-copy marker, so embedding it by value would make every Frame copy a vet
// error. Copies share the pixel slice either way.
type Frame struct {
	*common.Frame
	Width, Height uint64
}

// SetWithCoord allows to set a pixel by coord.
//
// The colour is copied rather than referenced, so a frame always owns its pixels. That
// matters because callers routinely pass a colour they reuse: the caracter driver hands
// the same one to every pixel of a glyph, and it comes from the loaded theme, which
// storing by reference would leave open to being overwritten through the frame. It also
// keeps the argument off the heap, and at one allocation per pixel per frame that is the
// bulk of what the collector has to chase while rendering.
func (f *Frame) SetWithCoord(coo *common.Coord, col *common.Color) {
	if col == nil || coo.X < 0 || coo.Y < 0 {
		return
	}

	i := int(coo.X) + int(coo.Y)*int(f.Width)
	if i < 0 || i >= len(f.Pixels) {
		return
	}

	p := f.Pixels[i]
	if p == nil {
		p = &common.Color{}
		f.Pixels[i] = p
	}
	p.R, p.G, p.B, p.A = col.R, col.G, col.B, col.A
}

// Clean set all pixels to transparent.
func (f *Frame) Clean() {
	f.Pixels = make([]*common.Color, f.Width*f.Height)
	for i := 0; i < len(f.Pixels); i++ {
		f.Pixels[i] = &common.Color{}
	}
}

// GetColumn returns pixels for given column index.
func (f *Frame) GetColumn(idx int) []*common.Color {
	c := make([]*common.Color, f.Height)
	if 0 <= idx && idx < int(f.Width) {
		for i := 0; i < int(f.Height); i++ {
			c[i] = f.Pixels[i*int(f.Width)+idx]
		}
	}
	return c
}

// SetColumn update pixels of a given column by index.
func (f *Frame) SetColumn(idx int, col []*common.Color) {
	if idx < 0 || idx >= int(f.Width) {
		return
	}
	// Copied for the same reason as SetWithCoord: the text driver feeds this the output of
	// GetColumn on another frame, and sharing the colours would tie the two together.
	for i := 0; i < int(f.Height) && i < len(col); i++ {
		f.SetWithCoord(&common.Coord{X: int64(idx), Y: int64(i)}, col[i])
	}
}
