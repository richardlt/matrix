package os

import (
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
	st := &softwareState{
		meta:                  meta,
		playerCount:           count,
		multiPressSelectStart: mp,
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

type softwareState struct {
	meta                  system.SoftwareMeta
	playerCount           uint64
	multiPressSelectStart software.ActionGenerator
	ctx                   *Context
}

func (s *softwareState) Init(ctx *Context) {
	s.ctx = ctx

	ctx.softwareServer.OnPrint(func(f render.Frame) {
		ctx.displayServer.Print([]render.Frame{f})
	})

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
			ctx.SetState(newSoftMenuState(ctx.GetSoftwareMeta()))
		}
	})

	ctx.softwareServer.StartSoftware(s.meta, s.playerCount)
}

func (s *softwareState) catchAction(a system.Action) {
	switch a.Command {
	case commandSelectStart:
		s.ctx.softwareServer.CloseSoftware()
		s.ctx.SetState(newSoftMenuState(s.ctx.GetSoftwareMeta()))
	default:
		s.ctx.softwareServer.Command(a.Slot, a.Command)
	}
}
