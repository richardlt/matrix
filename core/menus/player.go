package menus

import (
	"strconv"

	"github.com/richardlt/matrix/core/drivers"
	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/core/system"
	"github.com/richardlt/matrix/sdk-go/common"
)

// NewPlayer returns a new player menu.
func NewPlayer(f *render.Frame, min, max uint64) Player {
	return Player{frame: f, min: min, max: max, selected: min - 1}
}

// Player menu displays count of players.
type Player struct {
	frame               *render.Frame
	min, max, selected  uint64
	printCallback       func()
	selectCountCallback func(uint64)
	goBackCallback      func()
}

// Print renders a count of players from min to max.
func (p *Player) Print() {
	p.frame.Clean()

	green := render.GetColorFromThemeByName("flat", "green_2")

	cdTBT := drivers.NewCaracter(p.frame, render.GetFontByName("ThreeByThree"))
	cdTBT.Render('P', common.Coord{X: 11, Y: 5}, green, common.Color{})

	cdFBF := drivers.NewCaracter(p.frame, render.GetFontByName("FiveByFive"))
	cdFBF.Render([]rune(strconv.Itoa(int(p.selected + 1)))[0], common.Coord{X: 5, Y: 4}, green,
		common.Color{})

	if p.printCallback != nil {
		p.printCallback()
	}
}

// Action catches player's actions.
func (p *Player) Action(a system.Action) {
	switch a.Command {
	case common.Command_A_UP:
		if p.selectCountCallback != nil {
			p.selectCountCallback(p.selected + 1)
		}
	case common.Command_B_UP:
		if p.goBackCallback != nil {
			p.goBackCallback()
		}
	case common.Command_LEFT_UP:
		if p.selected == p.min-1 {
			p.selected = p.max - 1
		} else {
			p.selected--
		}
		p.Print()
	case common.Command_RIGHT_UP:
		if p.selected == p.max-1 {
			p.selected = p.min - 1
		} else {
			p.selected++
		}
		p.Print()
	}
}

// OnPrint allows to set print callback.
func (p *Player) OnPrint(c func()) { p.printCallback = c }

// OnSelectCount allows to set select count callback.
func (p *Player) OnSelectCount(c func(count uint64)) { p.selectCountCallback = c }

// OnGoBack allows to set go back callback.
func (p *Player) OnGoBack(c func()) { p.goBackCallback = c }
