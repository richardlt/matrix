package device

import (
	"fmt"

	"github.com/googollee/go-socket.io"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/display"
	"github.com/richardlt/matrix/sdk-go/player"
	"github.com/sirupsen/logrus"
)

// Start device deamon.
func Start(port, corePort int) error {
	frameChannel := make(chan common.Frame)
	defer close(frameChannel)

	d := &device{frameChannel: frameChannel}

	s, err := newSocketIOServer(d)
	if err != nil {
		return err
	}

	go func() {
		for f := range frameChannel {
			s.BroadcastTo("display", "frame", f)
		}
	}()

	go func() {
		if err := display.Connect(fmt.Sprintf("localhost:%d", corePort), d, true); err != nil {
			logrus.Errorf("%+v", errors.WithStack(err))
		}
	}()

	go func() {
		if err := player.Connect(fmt.Sprintf("localhost:%d", corePort), d, true); err != nil {
			logrus.Errorf("%+v", err)
		}
	}()

	e := echo.New()
	e.Any("/socket.io/", echo.WrapHandler(s))

	logrus.Infof("Start device on port %d\n", port)
	return e.Start(fmt.Sprintf(":%d", port))
}

func newSocketIOServer(d *device) (*socketio.Server, error) {
	s, err := socketio.NewServer(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := s.On("connection", func(so socketio.Socket) {
		so.Join("display")
		so.On("command", func(cmd string) { d.Command(commandFromString(cmd)) })
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

type device struct {
	frameChannel chan common.Frame
	playerAPI    *player.API
}

func (d *device) FramesReceived(fs []*common.Frame) {
	if len(fs) > 0 {
		d.frameChannel <- *fs[0]
	}
}

func (d *device) Init(api *player.API) error {
	d.playerAPI = api
	return nil
}

func (d *device) Command(cmd common.Command) {
	if d.playerAPI != nil {
		d.playerAPI.Command(uint64(0), cmd)
	}
}
