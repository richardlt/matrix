package gamepad

import (
	"fmt"
	"sync"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/display"
	"github.com/richardlt/matrix/sdk-go/player"
	"github.com/richardlt/matrix/websocket"
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

	s := newWebSocketServer(gs)

	go func() {
		for f := range frameChannel {
			s.Broadcast("frame", f)
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
	e.HideBanner = true

	e.Any("/websocket", echo.WrapHandler(s))
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

func newWebSocketServer(gs *gamepadServer) *websocket.Server {
	s := websocket.NewServer()

	s.OnConnect(func(c *websocket.Client) {
		g := newGamepad(c)
		gs.AddGamepad(g)
		c.OnEvent(func(eventType string, data interface{}) {
			switch eventType {
			case "select-slot":
				g.SelectSlot(int(data.(float64)))
			case "command":
				gs.Command(g, commandFromString(data.(string)))
			}
		})
		c.OnDisconnect(func() {
			gs.RemoveGamepad(g)
		})
	})

	return s
}

func newGamepad(so *websocket.Client) *gamepad { return &gamepad{so, -1} }

type gamepad struct {
	so   *websocket.Client
	Slot int
}

func (g *gamepad) SelectSlot(slot int) {
	g.Slot = slot
	g.so.Send("slot", slot)
}
