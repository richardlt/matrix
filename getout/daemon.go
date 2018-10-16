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
	g.print()
}

func (g *getout) Close() {
	g.renderer.Clean()
}

func (g *getout) ActionReceived(slot uint64, cmd common.Command) {}

func (g *getout) print() { g.renderer.Print(g.engine.grid, g.engine.start, g.engine.end) }
