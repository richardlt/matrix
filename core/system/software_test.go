package system

import (
	"testing"

	softwareSDK "github.com/richardlt/matrix/sdk-go/software"
)

// The core drives the pause, so the running software has to be told, and only it.
func TestSetPausedReachesRunningSoftware(t *testing.T) {
	s := NewSoftwareServer()
	so := &software{
		UUID:                   "a",
		connectResponseChannel: make(chan *softwareSDK.ConnectResponse, 2),
	}
	s.current = so

	s.SetPaused(true)
	if !so.Paused {
		t.Error("software should be marked paused")
	}

	res := <-so.connectResponseChannel
	if res.SoftwareData.Action != softwareSDK.ConnectResponse_SoftwareData_PAUSE ||
		!res.SoftwareData.Paused {
		t.Fatalf("got %v, want a PAUSE response carrying true", res.SoftwareData)
	}

	s.SetPaused(false)
	res = <-so.connectResponseChannel
	if res.SoftwareData.Paused {
		t.Error("resuming should carry false")
	}
	if so.Paused {
		t.Error("software should no longer be marked paused")
	}
}

// Nothing on screen means nothing to pause, and no send on a channel nobody reads.
func TestSetPausedWithoutRunningSoftware(t *testing.T) {
	NewSoftwareServer().SetPaused(true)
}

// Whether the pause button is caught by the core or handed to the software depends on
// this flag, so it has to survive into the metadata the menus work from.
func TestConfigCarriesPausable(t *testing.T) {
	so := &software{UUID: "a"}
	so.SetConfig(&softwareSDK.ConnectRequest_SoftwareData_Config{
		MinPlayerCount: 1,
		MaxPlayerCount: 2,
		Pausable:       true,
	})
	if !so.GetMeta().Pausable {
		t.Error("meta should report the software as pausable")
	}
}

// Starting a software that was left paused must not leave it looking paused.
func TestStartClearsPaused(t *testing.T) {
	so := &software{
		UUID:                   "a",
		Paused:                 true,
		connectResponseChannel: make(chan *softwareSDK.ConnectResponse, 1),
	}
	so.Start(1)
	if so.Paused {
		t.Error("Start must clear the paused flag")
	}
}
