package device

import (
	"testing"

	"github.com/richardlt/matrix/sdk-go/display"
)

// matrix has to satisfy the SDK's optional interface, or the core's state notifications
// are silently dropped and the pollers never hold off.
func TestMatrixIsStateAware(t *testing.T) {
	var _ display.StateAware = (*matrix)(nil)
}

func TestCoreStateGatesPolling(t *testing.T) {
	s := new(coreState)

	if !s.idle() {
		t.Error("menus are idle by default")
	}

	s.setSoftwareRunning(true)
	if s.idle() {
		t.Error("must not poll while a software is drawing")
	}

	s.setSoftwareRunning(false)
	if !s.idle() {
		t.Error("polling resumes back at the menu")
	}
}
