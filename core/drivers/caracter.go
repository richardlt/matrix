package drivers

import (
	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// NewCaracter returns a new caracter driver.
func NewCaracter(fr *render.Frame, fo software.Font) *Caracter {
	return &Caracter{frame: fr, font: fo}
}

// Caracter driver allows to render a caracter with specific font.
type Caracter struct {
	frame       *render.Frame
	font        software.Font
	endCallback func()
}

// SetFont allows to set and change caracter font.
func (c *Caracter) SetFont(f software.Font) { c.font = f }

// Render prints the caracter in frame.
func (c *Caracter) Render(value rune, center common.Coord,
	color, background common.Color) {
	ca := render.GetFontCaracterByValue(c.font, value)

	beginX, beginY := center.X-int64(ca.Width)/2, center.Y-int64(c.font.Height)/2
	endX, endY := beginX+int64(ca.Width), beginY+int64(c.font.Height)

	index := 0
	for i := beginY; i < endY; i++ {
		for j := beginX; j < endX; j++ {
			if i < int64(c.frame.Height) && j < int64(c.frame.Width) {
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
