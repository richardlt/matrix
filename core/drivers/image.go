package drivers

import (
	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// NewImage returns a new image driver.
func NewImage(f *render.Frame) *Image { return &Image{frame: f} }

// Image driver allows to render a pixel image.
type Image struct {
	frame       *render.Frame
	endCallback func()
}

// Render prints the image in frame.
func (i *Image) Render(im software.Image, c common.Coord) {
	beginX, beginY := c.X-int64(im.Width)/2, c.Y-int64(im.Height)/2
	endX, endY := beginX+int64(im.Width), beginY+int64(im.Height)

	index := 0
	for y := beginY; y < endY; y++ {
		for x := beginX; x < endX; x++ {
			i.frame.SetWithCoord(common.Coord{X: x, Y: y}, render.GetImagePixelWithIndex(im, index))
			index++
		}
	}

	if i.endCallback != nil {
		i.endCallback()
	}
}

// OnEnd allows to set end callback.
func (i *Image) OnEnd(f func()) { i.endCallback = f }
