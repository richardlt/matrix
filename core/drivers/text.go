package drivers

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
)

// NewText returns a new text driver.
func NewText(fr *render.Frame, fo render.Font) *Text {
	return &Text{
		frame:          fr,
		font:           fo,
		caracterDriver: NewCaracter(fr, fo),
	}
}

// Text allows to render a given text in frame.
type Text struct {
	frame          *render.Frame
	font           render.Font
	caracterDriver *Caracter
	ticker         *time.Ticker
	endCallback    func()
	stepCallback   func(total, current uint64)
}

// Render displays given text from left to right with scroll effect if too long.
func (t *Text) Render(text string, center common.Coord, color, background color.RGBA) {
	if t.ticker != nil {
		t.ticker.Stop()
	}

	t.ticker = time.NewTicker(100 * time.Millisecond)

	spacesText := fmt.Sprintf("   %s  ", strings.Join(strings.Split(text, ""), " "))

	textLength := uint64(0)
	for _, c := range spacesText {
		textLength += t.font.GetCaracterByValue(c).Width
	}

	stepCount := textLength - t.frame.Width
	if stepCount < 0 {
		stepCount = 0
	} else {
		stepCount++
	}

	step := uint64(0)

	go func() {
		for _ = range t.ticker.C {
			if step > stepCount {
				t.Stop()
				if t.endCallback != nil {
					t.endCallback()
				}
			} else {
				if t.stepCallback != nil {
					go t.stepCallback(stepCount, step)
				}
				beginX := uint64(0) - step
				for _, c := range spacesText {
					caracterWidth := t.font.GetCaracterByValue(c).Width
					t.caracterDriver.Render(c, common.Coord{
						X: beginX + caracterWidth - uint64(caracterWidth/2) - 1,
						Y: center.Y,
					}, color, background)
					beginX += caracterWidth
				}
				step++
			}
		}
	}()
}

// OnEnd allows to set end callback.
func (t *Text) OnEnd(c func()) { t.endCallback = c }

// OnStep allows to set step callback.
func (t *Text) OnStep(c func(total, current uint64)) { t.stepCallback = c }

// Stop the rendering process if not ended.
func (t *Text) Stop() { t.ticker.Stop() }
