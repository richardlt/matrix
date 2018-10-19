package getout

import (
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
)

// Start the getout software.
func Start(uri string) error {
	logrus.Infof("Start getout for uri %s\n", uri)

	g := &getout{}

	return software.Connect(uri, g, true)
}

type getout struct {
	engine   *engine
	renderer *renderer
}

func (g *getout) Init(a software.API) (err error) {
	logrus.Debug("Init getout")

	g.renderer, err = newRenderer(a)
	if err != nil {
		return err
	}

	l := a.GetImageFromLocal("getout")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &l,
		MinPlayerCount: 1,
		MaxPlayerCount: 1,
	})

	return a.Ready()
}

func (g *getout) Start(uint64) {
	g.engine = newEngine(16, 9)
	g.renderer.printGrid(g.engine.grid)
	g.renderer.printPlayer(g.engine.player, g.engine.end)
}

func (g *getout) Close() {
	g.renderer.clean()
	g.renderer.stopPrintGameOver()
}

func (g *getout) ActionReceived(slot uint64, cmd common.Command) {
	if g.engine.isGameOver() {
		return
	}

	switch cmd {
	case common.Command_LEFT_UP:
		g.engine.move("left")
	case common.Command_UP_UP:
		g.engine.move("up")
	case common.Command_RIGHT_UP:
		g.engine.move("right")
	case common.Command_DOWN_UP:
		g.engine.move("down")
	}
	g.renderer.printPlayer(g.engine.player, g.engine.end)

	if g.engine.isGameOver() {
		g.renderer.startPrintGameOver()
	}
}
