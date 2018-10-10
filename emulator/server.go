package emulator

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/googollee/go-socket.io"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/display"
)

type frame struct {
	Number int     `json:"number"`
	Pixels []pixel `json:"pixels"`
}

type pixel struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

// Start the emulator server.
func Start(port int, uri string) error {
	frameChannel := make(chan frame)
	defer close(frameChannel)

	s, err := newSocketIOServer()
	if err != nil {
		return err
	}

	go func() {
		for f := range frameChannel {
			s.BroadcastTo("display", "frame", f)
		}
	}()

	go func() {
		if err := display.Connect(uri, emulator{frameChannel}, true); err != nil {
			logrus.Errorf("%+v", errors.WithStack(err))
		}
	}()

	e := echo.New()
	e.Any("/socket.io/", echo.WrapHandler(s))
	e.Static("/", "./emulator/public")

	logrus.Infof("Start emulator on port %d\n", port)
	return e.Start(fmt.Sprintf(":%d", port))
}

type emulator struct{ frameChannel chan frame }

func (e emulator) FramesReceived(fs []*common.Frame) {
	for i, f := range fs {
		frame := frame{
			Number: i,
			Pixels: make([]pixel, len(f.Pixels)),
		}
		for i, c := range f.Pixels {
			frame.Pixels[i].R = int(c.R)
			frame.Pixels[i].G = int(c.G)
			frame.Pixels[i].B = int(c.B)
		}
		e.frameChannel <- frame
	}
}

func newSocketIOServer() (*socketio.Server, error) {
	s, err := socketio.NewServer(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := s.On("connection", func(so socketio.Socket) {
		so.Join("display")
		so.On("disconnection", func() {})
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := s.On("error", func(so socketio.Socket, err error) {
		logrus.Errorf("%+v", errors.WithStack(err))
	}); err != nil {
		return nil, err
	}

	return s, nil
}
