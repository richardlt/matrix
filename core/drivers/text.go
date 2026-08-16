package drivers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// NewText returns a new text driver.
func NewText(fr *render.Frame, fo *software.Font) *Text {
	return &Text{
		frame: fr,
		font:  fo,
	}
}

// Text allows to render a given text in frame.
type Text struct {
	frame        *render.Frame
	font         *software.Font
	lock         sync.Mutex
	done         chan struct{}
	endCallback  func()
	stepCallback func(total, current uint64)
}

// Render displays given text from left to right with scroll effect if too long.
func (t *Text) Render(text string, center *common.Coord, color, background *common.Color, repeat bool) {
	t.Stop()

	text = fmt.Sprintf(" %s ", strings.Join(strings.Split(text, ""), " "))
	if repeat && 0 < center.X {
		text += strings.Repeat(" ", int(center.X))
	}

	len := uint64(0)
	for _, c := range text {
		len += render.GetFontCaracterByValue(t.font, c).Width
	}

	stepCount := int(center.X) + int(len-t.frame.Width)
	if stepCount < 0 {
		stepCount = 0
	}

	f := render.NewFrame(len, t.frame.Height)
	cd := NewCaracter(&f, t.font)

	offset := int64(0)
	for _, c := range text {
		caracterWidth := render.GetFontCaracterByValue(t.font, c).Width
		cd.Render(c, &common.Coord{
			X: offset + int64(caracterWidth-caracterWidth/2) - 1,
			Y: center.Y,
		}, color, background)
		offset += int64(caracterWidth)
	}

	// The ticker and its stop signal are set up here rather than in the goroutine below:
	// a Stop landing before the goroutine got to run would otherwise leave a scroll
	// running that nothing can reach any more.
	ticker := time.NewTicker(100 * time.Millisecond)
	done := make(chan struct{})

	t.lock.Lock()
	t.done = done
	t.lock.Unlock()

	step := 0
	go func() {
		defer ticker.Stop()

		xStart := int(center.X)
		offset := 0
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
			}

			if step > stepCount && !repeat {
				break
			}

			for i, j := xStart, offset; i < int(t.frame.Width); i++ {
				t.frame.SetColumn(i, f.GetColumn(j))
				if j < int(f.Width)-1 {
					j++
				} else {
					j = 0
				}
			}
			if xStart > 0 {
				xStart--
			} else {
				if offset < int(f.Width)-1 {
					offset++
				} else {
					offset = 0
				}
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
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.done != nil {
		close(t.done)
		t.done = nil
	}
}
