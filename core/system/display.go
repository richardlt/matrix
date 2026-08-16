package system

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/internal/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	displaySDK "github.com/richardlt/matrix/sdk-go/display"
)

// NewDisplayServer return new display server.
func NewDisplayServer() *DisplayServer { return &DisplayServer{} }

// DisplayServer expose RPC server for displays.
type DisplayServer struct {
	displaySDK.UnimplementedDisplayServer
	displays    []display
	displayLock sync.RWMutex
	lastFrames  []render.Frame
	// softwareRunning is replayed to each display as it connects, so one that joins
	// mid-game is not left assuming a menu is showing.
	softwareRunning bool
}

// Connect display action.
func (d *DisplayServer) Connect(stream displaySDK.Display_ConnectServer) error {
	chRes := make(chan *displaySDK.Response)
	defer close(chRes)

	di := newDisplay(chRes)

	logrus.Debugf("Display %s connect", di.UUID)
	defer logrus.Debugf("Display %s disconnect", di.UUID)

	d.AddDisplay(di)
	defer d.RemoveDisplay(di)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// every 3 seconds ping the display to test the conn
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				chRes <- &displaySDK.Response{Type: displaySDK.Response_PING}
			}
		}
	}()

	go func() {
		for r := range chRes {
			if err := stream.Send(r); err != nil {
				logrus.Errorf("%+v", errors.Errorf("sending response to display %s: %w", di.UUID, err))
			}
		}
	}()

	// bring the new display up to date on connect
	di.Print(d.lastFrames)
	di.SoftwareRunning(d.softwareRunning)

	for {
		if _, err := stream.Recv(); err != nil {
			return errors.Errorf("receiving from display %s: %w", di.UUID, err)
		}
	}
}

// AddDisplay allows to push a new display in server.
func (d *DisplayServer) AddDisplay(di display) {
	d.displayLock.Lock()
	d.displays = append(d.displays, di)
	d.displayLock.Unlock()
}

// RemoveDisplay removes a existing display in server.
func (d *DisplayServer) RemoveDisplay(di display) {
	d.displayLock.Lock()
	var ds []display
	for _, ed := range d.displays {
		if ed.UUID != di.UUID {
			ds = append(ds, ed)
		}
	}
	d.displays = ds
	d.displayLock.Unlock()
}

// Print send frame to all displays.
func (d *DisplayServer) Print(fs []render.Frame) {
	d.lastFrames = fs
	d.displayLock.RLock()
	logrus.Debugf("Print %d frames to %d displays", len(fs), len(d.displays))
	for _, di := range d.displays {
		di.Print(fs)
	}
	d.displayLock.RUnlock()
}

// SetSoftwareRunning records whether a software is drawing rather than a menu being
// shown, and tells every connected display. A display that drives hardware uses it to
// keep housekeeping out of the way of the frames.
func (d *DisplayServer) SetSoftwareRunning(running bool) {
	d.displayLock.RLock()
	defer d.displayLock.RUnlock()

	d.softwareRunning = running
	logrus.Debugf("Software running %t, notifying %d displays", running, len(d.displays))
	for _, di := range d.displays {
		di.SoftwareRunning(running)
	}
}

func newDisplay(chRes chan *displaySDK.Response) display {
	return display{uuid.NewString(), chRes}
}

type display struct {
	UUID            string
	responseChannel chan *displaySDK.Response
}

func (d *display) Print(fs []render.Frame) {
	r := &displaySDK.Response{
		Type: displaySDK.Response_DISPLAY,
		DisplayData: &displaySDK.Response_DisplayData{
			Action: displaySDK.Response_DisplayData_FRAMES,
			Frames: []*common.Frame{},
		},
	}
	for _, f := range fs {
		frame := &common.Frame{}
		for _, p := range f.Pixels {
			frame.Pixels = append(frame.Pixels, &common.Color{
				R: uint64(p.R),
				G: uint64(p.G),
				B: uint64(p.B),
				A: uint64(p.A),
			})
		}
		r.DisplayData.Frames = append(r.DisplayData.Frames, frame)
	}
	d.responseChannel <- r
}

// SoftwareRunning forwards the core's state to this display.
func (d *display) SoftwareRunning(running bool) {
	state := displaySDK.Response_StateData_MENU
	if running {
		state = displaySDK.Response_StateData_SOFTWARE
	}
	d.responseChannel <- &displaySDK.Response{
		Type:      displaySDK.Response_STATE,
		StateData: &displaySDK.Response_StateData{State: state},
	}
}
