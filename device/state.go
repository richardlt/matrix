package device

import "sync/atomic"

// coreState mirrors what the core reports it is doing, shared between the panel and the
// controller discovery in this package.
//
// Both of those poll, and both are expensive on the USB bus: hidapi opens each matching
// device to read its string descriptors, and the panel's serial adapter shares that bus,
// so a poll lands as a visible stutter in whatever is being drawn. Holding the polls
// while a software is running keeps them out of the way; they run in the menus, where
// the picture is still and where a player joins anyway.
//
// The trade is deliberate: hardware plugged in mid-game is not picked up until you step
// back to the menu.
type coreState struct{ softwareRunning atomic.Bool }

// setSoftwareRunning records the core's state, as reported over the display stream.
func (c *coreState) setSoftwareRunning(running bool) { c.softwareRunning.Store(running) }

// idle reports whether housekeeping that competes with the panel may run.
func (c *coreState) idle() bool { return !c.softwareRunning.Load() }
