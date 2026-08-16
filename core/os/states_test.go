package os

import "testing"

// SetState reports this to the displays, so a menu must never claim a software is
// running and the software state must never claim otherwise.
func TestStatesReportSoftwareRunning(t *testing.T) {
	for _, tc := range []struct {
		name string
		s    state
		want bool
	}{
		{"software menu", &softMenuState{}, false},
		{"player menu", &playerMenuState{}, false},
		{"software", &softwareState{}, true},
		{"paused software", &softwareState{paused: true}, false},
	} {
		if got := tc.s.SoftwareRunning(); got != tc.want {
			t.Errorf("%s: got %t, want %t", tc.name, got, tc.want)
		}
	}
}
