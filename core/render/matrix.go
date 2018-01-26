package render

import (
	"github.com/richardlt/matrix/sdk-go/common"
)

// NewMatrix returns a new matrix with given size.
func NewMatrix(w, h uint64) *Matrix {
	return &Matrix{
		topFrame: NewFrame(w, h),
		width:    w,
		height:   h,
	}
}

// Matrix is a multi layer frames container.
type Matrix struct {
	frames        []*Frame
	topFrame      Frame
	width, height uint64
}

func (m *Matrix) getTopPixelAtIndex(idx uint64) common.Color {
	if m.width*m.height <= idx {
		return common.Color{}
	}

	for i := len(m.frames) - 1; i >= 0; i-- {
		c := m.frames[i].Pixels[idx]
		if c.A > 0 {
			return *c
		}
	}

	return common.Color{}
}

// PrintFrame renders the top frame of the matrix.
func (m *Matrix) PrintFrame() {
	f := NewFrame(m.width, m.height)
	for i := uint64(0); i < m.width*m.height; i++ {
		top := m.getTopPixelAtIndex(i)
		f.Pixels[i] = &top
	}
	m.topFrame = f
}

// NewFrame returns a frame with matrix size under existing frames.
func (m *Matrix) NewFrame() *Frame {
	f := NewFrame(m.width, m.height)
	m.frames = append(m.frames, &f)
	return &f
}

// RemoveLayer allows to remove a layer in the matrix.
func (m *Matrix) RemoveLayer(f *Frame) {
	for i, frame := range m.frames {
		if frame == f {
			if len(m.frames) > 1 {
				m.frames = append(m.frames[:i], m.frames[i+1:]...)
			} else {
				m.frames = nil
			}
			break
		}
	}
}

// GetTopFrame returns matrix's top frame.
func (m *Matrix) GetTopFrame() Frame { return m.topFrame }
