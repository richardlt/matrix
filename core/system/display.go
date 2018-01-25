package system

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
	displaySDK "github.com/richardlt/matrix/sdk-go/display"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// NewDisplayServer return new display server.
func NewDisplayServer() *DisplayServer { return &DisplayServer{} }

// DisplayServer expose RPC server for displays.
type DisplayServer struct {
	displays    []display
	displayLock sync.RWMutex
	lastFrames  []render.Frame
}

// Connect display action.
func (d *DisplayServer) Connect(stream displaySDK.Display_ConnectServer) error {
	chRes := make(chan displaySDK.Response)
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
				chRes <- displaySDK.Response{Type: displaySDK.Response_PING}
			}
		}
	}()

	go func() {
		for r := range chRes {
			if err := stream.Send(&r); err != nil {
				logrus.Errorf("%+v", errors.WithStack(err))
			}
		}
	}()

	// print last frames on connect
	di.Print(d.lastFrames)

	for {
		_, err := stream.Recv()
		if err != nil {
			return errors.WithStack(err)
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

func newDisplay(chRes chan displaySDK.Response) display {
	return display{uuid.NewV4().String(), chRes}
}

type display struct {
	UUID            string
	responseChannel chan displaySDK.Response
}

func (d *display) Print(fs []render.Frame) {
	r := displaySDK.Response{
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
