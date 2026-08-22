package software

import (
	"testing"
	"time"

	common "github.com/richardlt/matrix/sdk-go/common"
)

type plainSoftware struct{}

func (plainSoftware) Init(API) error                        { return nil }
func (plainSoftware) Start(uint64)                          {}
func (plainSoftware) Close()                                {}
func (plainSoftware) ActionReceived(uint64, common.Command) {}

type pausableSoftware struct {
	plainSoftware
	paused chan bool
}

func (p *pausableSoftware) Paused(paused bool) { p.paused <- paused }

func pauseResponse(paused bool) *ConnectResponse {
	return &ConnectResponse{
		Type: ConnectResponse_SOFTWARE,
		SoftwareData: &ConnectResponse_SoftwareData{
			Action: ConnectResponse_SoftwareData_PAUSE,
			Paused: paused,
		},
	}
}

// A software that keeps a clock of its own has to hear about the pause, with the right
// value: the core has already stopped drawing it by then.
func TestPauseReachesPausableSoftware(t *testing.T) {
	s := &pausableSoftware{paused: make(chan bool, 2)}

	processResponse(s, pauseResponse(true))
	processResponse(s, pauseResponse(false))

	for _, want := range []bool{true, false} {
		select {
		case got := <-s.paused:
			if got != want {
				t.Errorf("got %t, want %t", got, want)
			}
		case <-time.After(time.Second):
			t.Fatal("software was never told about the pause")
		}
	}
}

// Pause is optional: a software driven only by player commands declares itself pausable
// without implementing anything, and must not be troubled by the signal.
func TestPauseIgnoredByPlainSoftware(t *testing.T) {
	processResponse(plainSoftware{}, pauseResponse(true))
}
