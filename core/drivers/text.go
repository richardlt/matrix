package drivers

import (
	"fmt"
	"strings"
	"time"

	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// NewText returns a new text driver.
func NewText(fr *render.Frame, fo software.Font) *Text {
	return &Text{
		frame:          fr,
		font:           fo,
		caracterDriver: NewCaracter(fr, fo),
	}
}

// Text allows to render a given text in frame.
type Text struct {
	frame          *render.Frame
	font           software.Font
	caracterDriver *Caracter
	ticker         *time.Ticker
	endCallback    func()
	stepCallback   func(total, current uint64)
}

// Render displays given text from left to right with scroll effect if too long.
func (t *Text) Render(text string, center common.Coord, color, background common.Color, repeat bool) {
	if t.ticker != nil {
		t.ticker.Stop()
	}

	text = fmt.Sprintf(" %s ", strings.Join(strings.Split(text, ""), " "))

	len := uint64(0)
	for _, c := range text {
		len += render.GetFontCaracterByValue(t.font, c).Width
	}

	stepCount := int(center.X) + int(len-t.frame.Width)
	if stepCount < 0 {
		stepCount = 0
	}

	step := 0

	go func() {
		t.ticker = time.NewTicker(100 * time.Millisecond)
		for _ = range t.ticker.C {
			if step > stepCount && !repeat {
				t.Stop()
				break
			}

			beginX := int(center.X) - step
			for _, c := range text {
				caracterWidth := render.GetFontCaracterByValue(t.font, c).Width
				t.caracterDriver.Render(c, common.Coord{
					X: int64(beginX) + int64(caracterWidth) - int64(caracterWidth/2) - 1,
					Y: center.Y,
				}, color, background)
				beginX += int(caracterWidth)
			}

			if t.stepCallback != nil {
				go t.stepCallback(uint64(stepCount), uint64(step))
			}
			step++
		}

		if t.endCallback != nil {
			t.endCallback()
		}
	}()
}

// OnEnd allows to set end callback.
func (t *Text) OnEnd(c func()) { t.endCallback = c }

// OnStep allows to set step callback.
func (t *Text) OnStep(c func(total, current uint64)) { t.stepCallback = c }

// Stop the rendering process if not ended.
func (t *Text) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}
