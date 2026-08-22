package os

import (
	"github.com/richardlt/matrix/core/system"
)

type state interface {
	Init(*Context)
	// SoftwareRunning distinguishes a software drawing from a menu being shown. Displays
	// are told, so one driving hardware can keep its own housekeeping out of the way of
	// the frames it is about to receive.
	SoftwareRunning() bool
}

// NewContext return a init os.
func NewContext(ps *system.PlayerServer, ds *system.DisplayServer,
	ss *system.SoftwareServer) *Context {
	return &Context{
		playerServer:   ps,
		displayServer:  ds,
		softwareServer: ss,
	}
}

// Context manage softwares, displays, players.
type Context struct {
	playerServer   *system.PlayerServer
	displayServer  *system.DisplayServer
	softwareServer *system.SoftwareServer
}

// StartContext set the initial state of the context.
func StartContext(c *Context) { c.SetState(newSoftMenuState(c.GetSoftwareMeta())) }

// SetState allows to change context's state.
func (c *Context) SetState(s state) {
	c.playerServer.ResetCallback()
	c.softwareServer.ResetCallback()
	c.displayServer.SetSoftwareRunning(s.SoftwareRunning())
	s.Init(c)
}

// GetSoftwareMeta returns metadata from software server.
func (c *Context) GetSoftwareMeta() []system.SoftwareMeta {
	return c.softwareServer.GetSoftwaresMeta()
}
