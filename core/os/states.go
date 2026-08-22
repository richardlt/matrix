package os

import (
	"sync/atomic"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/core/menus"
	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/core/system"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

func newSoftMenuState(sms []system.SoftwareMeta) *softMenuState {
	f := render.NewFrame(16, 9)
	sm := menus.NewSoftware(&f)
	sm.LoadMeta(sms)
	return &softMenuState{sm: &sm, f: &f}
}

type softMenuState struct {
	sm *menus.Software
	f  *render.Frame
}

func (s *softMenuState) Init(ctx *Context) {
	ctx.playerServer.OnAction(func(a system.Action) { s.sm.Action(a) })

	ctx.softwareServer.OnSoftwareChange(func(sms []system.SoftwareMeta) {
		s.sm.LoadMeta(sms)
	})

	s.sm.OnPrint(func() { ctx.displayServer.Print([]render.Frame{*s.f}) })

	s.sm.OnSelectSoftware(func(meta system.SoftwareMeta) {
		if meta.MinPlayerCount == meta.MaxPlayerCount {
			ctx.SetState(newSoftwareState(meta, meta.MinPlayerCount))
		} else {
			ctx.SetState(newPlayerMenuState(meta))
		}
	})

	s.sm.Print()
}

func newPlayerMenuState(meta system.SoftwareMeta) *playerMenuState {
	f := render.NewFrame(16, 9)
	pm := menus.NewPlayer(&f, meta.MinPlayerCount, meta.MaxPlayerCount)
	return &playerMenuState{pm: &pm, f: &f, meta: meta}
}

type playerMenuState struct {
	pm   *menus.Player
	f    *render.Frame
	meta system.SoftwareMeta
}

func (p *playerMenuState) Init(ctx *Context) {
	ctx.playerServer.OnAction(func(a system.Action) { p.pm.Action(a) })

	ctx.softwareServer.OnSoftwareChange(func(sms []system.SoftwareMeta) {
		found := false
		for _, sm := range sms {
			if sm.UUID == p.meta.UUID {
				found = true
				break
			}
		}
		if !found {
			ctx.SetState(newSoftMenuState(ctx.GetSoftwareMeta()))
		}
	})

	p.pm.OnPrint(func() { ctx.displayServer.Print([]render.Frame{*p.f}) })

	p.pm.OnSelectCount(func(count uint64) {
		ctx.SetState(newSoftwareState(p.meta, count))
	})

	p.pm.OnGoBack(func() {
		ctx.SetState(newSoftMenuState(ctx.GetSoftwareMeta()))
	})

	p.pm.Print()
}

func newSoftwareState(meta system.SoftwareMeta, count uint64) *softwareState {
	mp := software.NewMultiPress(common.Button_SELECT, common.Button_START)
	f := render.NewFrame(16, 9)
	pa := menus.NewPause(&f)
	st := &softwareState{
		meta:                  meta,
		playerCount:           count,
		multiPressSelectStart: mp,
		pause:                 &pa,
		pauseFrame:            &f,
	}

	mp.OnAction(func(slot uint64) {
		st.catchAction(system.Action{
			Slot:    slot,
			Command: commandSelectStart,
		})
	})

	return st
}

const commandSelectStart = 100

// paused and left are read and written from different goroutines: player commands arrive
// on the player stream, frames on the software stream, and SetState asks about them from
// wherever the handover happened.
type softwareState struct {
	meta                  system.SoftwareMeta
	playerCount           uint64
	multiPressSelectStart software.ActionGenerator
	pause                 *menus.Pause
	pauseFrame            *render.Frame
	paused                atomic.Bool
	left                  atomic.Bool
	ctx                   *Context
}

func (s *softwareState) Init(ctx *Context) {
	s.ctx = ctx

	ctx.softwareServer.OnPrint(func(f render.Frame) {
		// A software told to pause may still have a frame in flight, and it must not
		// paint over the paused screen the core is holding up.
		if s.paused.Load() {
			return
		}
		ctx.displayServer.Print([]render.Frame{f})
	})

	s.pause.OnPrint(func() { ctx.displayServer.Print([]render.Frame{*s.pauseFrame}) })

	ctx.playerServer.OnAction(func(a system.Action) {
		s.multiPressSelectStart.SendAction(a.Slot, a.Command)
		s.catchAction(a)
	})

	ctx.softwareServer.OnSoftwareChange(func(sms []system.SoftwareMeta) {
		found := false
		for _, sm := range sms {
			if sm.UUID == s.meta.UUID {
				found = true
				break
			}
		}
		if !found {
			s.leave()
			ctx.SetState(newSoftMenuState(ctx.GetSoftwareMeta()))
		}
	})

	if err := ctx.softwareServer.StartSoftware(s.meta, s.playerCount); err != nil {
		logrus.Errorf("%+v", err)
	}
}

func (s *softwareState) catchAction(a system.Action) {
	// An action already on its way when the state handed over must not restart anything
	// here, or the paused screen would end up drawing over the menu.
	if s.left.Load() {
		return
	}

	switch a.Command {
	case commandSelectStart:
		s.leave()
		s.ctx.softwareServer.CloseSoftware()
		s.ctx.SetState(newSoftMenuState(s.ctx.GetSoftwareMeta()))
	case common.Command_START_UP:
		// Pause belongs to the core: it owns the button, the screen it holds up and the
		// state displays are told about. A software only declares that it can be paused.
		if s.meta.Pausable {
			s.setPaused(!s.paused.Load())
		} else {
			s.ctx.softwareServer.Command(a.Slot, a.Command)
		}
	default:
		// A paused software receives nothing, so play cannot advance behind the screen
		// the core is holding up.
		if !s.paused.Load() {
			s.ctx.softwareServer.Command(a.Slot, a.Command)
		}
	}
}

func (s *softwareState) setPaused(paused bool) {
	// Two players can hit the button at once, and only the press that flips the state does
	// the work that goes with it.
	if !s.paused.CompareAndSwap(!paused, paused) {
		return
	}

	// A paused software is not drawing, so for displays it counts as a menu: work they
	// hold back during play resumes, which is what allows a controller that dropped out
	// to be picked up again without leaving the game.
	s.ctx.softwareServer.SetPaused(paused)
	s.ctx.displayServer.SetSoftwareRunning(!paused)

	if paused {
		s.pause.Start()
	} else {
		s.pause.Stop()
		s.ctx.softwareServer.PrintCurrent()
	}
}

// leave stops what this state owns before it hands over, so no animation is left
// drawing over whatever comes next.
func (s *softwareState) leave() {
	// Marked first, so an action arriving alongside this one turns back at the top of
	// catchAction rather than acting on a state that is halfway out.
	s.left.Store(true)

	// A software on its way out is told that play resumed, so it is never left holding a
	// pause that nothing is going to lift.
	if s.paused.Swap(false) {
		s.ctx.softwareServer.SetPaused(false)
	}

	s.pause.Stop()
}

func (s *softMenuState) SoftwareRunning() bool   { return false }
func (p *playerMenuState) SoftwareRunning() bool { return false }
func (s *softwareState) SoftwareRunning() bool   { return !s.paused.Load() }
