package gamepad

import (
	"fmt"
	"log"
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

func Start(port, corePort int) error {
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
		if err := display.Connect(fmt.Sprintf("localhost:%d", corePort), gs, true); err != nil {
			logrus.Errorf("%+v", err)
		}
	}()

	go func() {
		if err := player.Connect(fmt.Sprintf("localhost:%d", corePort), gs, true); err != nil {
			logrus.Errorf("%+v", err)
		}
	}()

	e := echo.New()
	e.Any("/socket.io/", echo.WrapHandler(s))
	e.Static("/", "./gamepad/client/dist")

	log.Printf("Start gamepad on port %d\n", port)
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

func (g *gamepadServer) AskSlot(gp *gamepad, slot int) {}

func (g *gamepadServer) GetSlots(gp *gamepad) {}

func newSocketIOServer(gs *gamepadServer) (*socketio.Server, error) {
	s, err := socketio.NewServer(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := s.On("connection", func(so socketio.Socket) {
		g := newGamepad(so)
		so.Join("display")
		so.On("ask_slot", func(slot int) {
			gs.AskSlot(g, slot)
		})
		so.On("get_slots", func() {
			gs.GetSlots(g)
		})
		so.On("command", func(cmd string) {
			gs.Command(g, commandFromString(cmd))
		})
		so.On("disconnection", func() {
			gs.RemoveGamepad(g)
		})
		gs.AddGamepad(g)
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := s.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	}); err != nil {
		return nil, err
	}

	return s, nil
}

func newGamepad(so socketio.Socket) *gamepad {
	return &gamepad{so, 2}
}

type gamepad struct {
	so   socketio.Socket
	Slot int
}
