package system

import (
	"testing"
	"time"

	displaySDK "github.com/richardlt/matrix/sdk-go/display"
)

// A display driving hardware relies on this to know when to stay out of the way, so the
// notification has to reach every connected display and carry the right state.
func TestSetSoftwareRunningNotifiesDisplays(t *testing.T) {
	ds := NewDisplayServer()
	ch := make(chan *displaySDK.Response, 1)
	ds.AddDisplay(display{UUID: "test", responseChannel: ch})

	for _, tc := range []struct {
		running bool
		want    displaySDK.Response_StateData_State
	}{
		{true, displaySDK.Response_StateData_SOFTWARE},
		{false, displaySDK.Response_StateData_MENU},
	} {
		go ds.SetSoftwareRunning(tc.running)

		select {
		case r := <-ch:
			if r.Type != displaySDK.Response_STATE {
				t.Fatalf("running=%t: got type %v", tc.running, r.Type)
			}
			if r.StateData.GetState() != tc.want {
				t.Errorf("running=%t: got state %v, want %v", tc.running, r.StateData.GetState(), tc.want)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("running=%t: no state sent", tc.running)
		}
	}
}
