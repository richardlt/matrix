package menus

import (
	"github.com/richardlt/matrix/core/drivers"
	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
)

// NewPause returns a new pause screen.
func NewPause(f *render.Frame) Pause { return Pause{frame: f} }

// Pause is the screen the core holds up while play is suspended. It belongs here rather
// than to a software so that every software gets the same one, and so the screen keeps
// moving while the software behind it is stopped.
type Pause struct {
	frame         *render.Frame
	text          *drivers.Text
	printCallback func()
}

// Start scrolls the paused label until Stop is called.
func (p *Pause) Start() {
	p.frame.Clean()

	p.text = drivers.NewText(p.frame, render.GetFontByName("FiveByFive"))
	p.text.OnStep(func(total, current uint64) { p.print() })
	p.text.Render("PAUSE", &common.Coord{X: 4, Y: 4},
		render.GetColorFromLocalThemeByName("flat", "dark_grey_2"),
		&common.Color{}, true)
}

// Stop ends the animation.
func (p *Pause) Stop() {
	if p.text != nil {
		p.text.Stop()
		p.text = nil
	}
	p.frame.Clean()
}

func (p *Pause) print() {
	if p.printCallback != nil {
		p.printCallback()
	}
}

// OnPrint allows to set print callback.
func (p *Pause) OnPrint(c func()) { p.printCallback = c }
