package yumyum

import (
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
)

// Start the zigzag software.
func Start(uri string) error {
	logrus.Infof("Start yumyum for uri %s\n", uri)

	y := &yumyum{}

	return software.Connect(uri, y, true)
}

type yumyum struct {
	engine   *engine
	renderer *renderer
}

func (y *yumyum) Init(a software.API) (err error) {
	logrus.Debug("Init yumyum")

	y.renderer, err = newRenderer(a)
	if err != nil {
		return err
	}

	l := a.GetImageFromLocal("yumyum")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &l,
		MinPlayerCount: 1,
		MaxPlayerCount: 4,
	})

	a.Ready()
	return nil
}

func (y *yumyum) Start(playerCount uint64) {
	y.engine = newEngine(playerCount, 16, 9)
	y.print()
}

func (y yumyum) Close() {}

func (y *yumyum) ActionReceived(slot int, cmd common.Command) {
	switch cmd {
	case common.Command_LEFT_UP:
		y.engine.MovePlayer(slot, "left")
	case common.Command_UP_UP:
		y.engine.MovePlayer(slot, "up")
	case common.Command_RIGHT_UP:
		y.engine.MovePlayer(slot, "right")
	case common.Command_DOWN_UP:
		y.engine.MovePlayer(slot, "down")
	}
	y.print()
}

func (y *yumyum) print() {
	y.renderer.Print(y.engine.GetPlayers(), y.engine.GetCandies())
}
