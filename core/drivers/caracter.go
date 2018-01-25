package drivers

import (
	"image/color"

	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
)

// NewCaracter returns a new caracter driver.
func NewCaracter(fr *render.Frame, fo render.Font) *Caracter {
	return &Caracter{frame: fr, font: fo}
}

// Caracter driver allows to render a caracter with specific font.
type Caracter struct {
	frame       *render.Frame
	font        render.Font
	endCallback func()
}

// SetFont allows to set and change caracter font.
func (c *Caracter) SetFont(f render.Font) { c.font = f }

// Render prints the caracter in frame.
func (c *Caracter) Render(value rune, center common.Coord,
	color, background color.RGBA) {
	ca := c.font.GetCaracterByValue(value)

	beginX, beginY := center.X-ca.Width/2, center.Y-c.font.Height/2
	endX, endY := beginX+ca.Width, beginY+c.font.Height

	index := 0
	for i := beginY; i < endY; i++ {
		for j := beginX; j < endX; j++ {
			if i < c.frame.Height && j < c.frame.Width {
				selectedColor := background
				if ca.Mask[index] > 0 {
					selectedColor = color
				}
				c.frame.SetWithCoord(common.Coord{X: j, Y: i}, selectedColor)
			}
			index++
		}
	}

	if c.endCallback != nil {
		c.endCallback()
	}
}

// OnEnd allows to set end callback.
func (c *Caracter) OnEnd(f func()) { c.endCallback = f }
