package menus

import (
	"github.com/richardlt/matrix/core/drivers"
	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/core/system"
	"github.com/richardlt/matrix/sdk-go/common"
)

// NewSoftware returns a new sotware menu.
func NewSoftware(f *render.Frame) Software {
	return Software{frame: f, selected: -1}
}

// Software menu display software logos.
type Software struct {
	frame                  *render.Frame
	softwaresMeta          []system.SoftwareMeta
	selected               int
	printCallback          func()
	selectSoftwareCallback func(system.SoftwareMeta)
}

// LoadMeta updates menu's metadata.
func (s *Software) LoadMeta(sms []system.SoftwareMeta) {
	s.softwaresMeta = sms
	if len(s.softwaresMeta) == 0 {
		s.selected = -1
	} else if s.selected >= len(s.softwaresMeta) {
		s.selected = len(s.softwaresMeta) - 1
	} else if s.selected == -1 {
		s.selected = 0
	}
	s.Print()
}

// Print renders software's logos.
func (s Software) Print() {
	s.frame.Clean()

	id := drivers.NewImage(s.frame)

	var i render.Image
	if len(s.softwaresMeta) == 0 {
		i = render.GetImageByName("empty")
	} else {
		i = s.softwaresMeta[s.selected].Logo
	}

	id.Render(i, common.Coord{X: 8, Y: 4})

	if s.printCallback != nil {
		go s.printCallback()
	}
}

// Action catches player's actions.
func (s *Software) Action(a system.Action) {
	// ignore action if no software to display
	if len(s.softwaresMeta) == 0 {
		return
	}

	switch a.Command {
	case common.Command_A_UP:
		if s.selectSoftwareCallback != nil {
			s.selectSoftwareCallback(s.softwaresMeta[s.selected])
		}
	case common.Command_LEFT_UP:
		if s.selected == 0 {
			s.selected = len(s.softwaresMeta) - 1
		} else {
			s.selected--
		}
		s.Print()
	case common.Command_RIGHT_UP:
		if s.selected == len(s.softwaresMeta)-1 {
			s.selected = 0
		} else {
			s.selected++
		}
		s.Print()
	}
}

// OnPrint allows to set print callback.
func (s *Software) OnPrint(c func()) { s.printCallback = c }

// OnSelectSoftware allows to set print callback.
func (s *Software) OnSelectSoftware(c func(meta system.SoftwareMeta)) {
	s.selectSoftwareCallback = c
}
