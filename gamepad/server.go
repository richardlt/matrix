package gamepad

import (
	"fmt"
	"sync"

	"github.com/googollee/go-socket.io"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/display"
	"github.com/richardlt/matrix/sdk-go/player"
	"github.com/sirupsen/logrus"
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

// Start the gamepad server.
func Start(port int, uri string) error {
	frameChannel := make(chan frame)
	defer close(frameChannel)

	gs := &gamepadServer{frameChannel: frameChannel}

	s, err := newSocketIOServer(gs)
	if err != nil {
		return err
	}

	go func() {
		for f := range frameChannel {
			s.BroadcastTo("display", "frame", f)
		}
	}()

	go func() {
		if err := display.Connect(uri, gs, true); err != nil {
			logrus.Errorf("%+v", err)
		}
	}()

	go func() {
		if err := player.Connect(uri, gs, true); err != nil {
			logrus.Errorf("%+v", err)
		}
	}()

	e := echo.New()
	e.Any("/socket.io/", echo.WrapHandler(s))
	e.Static("/", "./gamepad/public")

	logrus.Infof("Start gamepad on port %d\n", port)
	return errors.WithStack(e.Start(fmt.Sprintf(":%d", port)))
}

type gamepadServer struct {
	frameChannel chan frame
	gamepads     []*gamepad
	gamepadLock  sync.RWMutex
	playerAPI    *player.API
}

func (g *gamepadServer) FramesReceived(fs []*common.Frame) {
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
		g.frameChannel <- frame
	}
}

func (g *gamepadServer) Init(api *player.API) error {
	g.playerAPI = api
	return nil
}

func (g *gamepadServer) AddGamepad(gp *gamepad) {
	g.gamepadLock.Lock()
	g.gamepads = append(g.gamepads, gp)
	g.gamepadLock.Unlock()
}

func (g *gamepadServer) RemoveGamepad(gp *gamepad) {
	g.gamepadLock.Lock()
	gs := []*gamepad{}
	for _, gamepad := range g.gamepads {
		if gamepad != gp {
			gs = append(gs, gamepad)
		}
	}
	g.gamepads = gs
	g.gamepadLock.Unlock()
}

func (g *gamepadServer) Command(gp *gamepad, cmd common.Command) {
	if 0 <= gp.Slot {
		g.playerAPI.Command(uint64(gp.Slot), cmd)
	}
}

func newSocketIOServer(gs *gamepadServer) (*socketio.Server, error) {
	s, err := socketio.NewServer(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := s.On("connection", func(so socketio.Socket) {
		g := newGamepad(so)
		so.Join("display")
		so.On("select-slot", func(slot int) { g.SelectSlot(slot) })
		so.On("command", func(cmd string) { gs.Command(g, commandFromString(cmd)) })
		so.On("disconnection", func() { gs.RemoveGamepad(g) })
		gs.AddGamepad(g)
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

func newGamepad(so socketio.Socket) *gamepad { return &gamepad{so, -1} }

type gamepad struct {
	so   socketio.Socket
	Slot int
}

func (g *gamepad) SelectSlot(slot int) {
	g.Slot = slot
	g.so.Emit("slot", slot)
}
