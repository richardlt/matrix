package zigzag

import (
	"fmt"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
)

// Start the zigzag software.
func Start(uri string) error {
	logrus.Infof("Start zigzag for uri %s\n", uri)

	z := &zigzag{}

	return software.Connect(uri, z, true)
}

type zigzag struct {
	engine   *engine
	renderer *renderer
}

func (z *zigzag) Init(a software.API) (err error) {
	logrus.Debug("Init zigzag")

	z.engine = newEngine(4, 16, 9)

	z.renderer, err = newRenderer(a)
	if err != nil {
		return err
	}

	l := a.GetImageFromLocal("zigzag")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &l,
		MinPlayerCount: 1,
		MaxPlayerCount: 4,
	})

	a.Ready()
	return nil
}

func (z *zigzag) Start() { z.print() }

func (z zigzag) Close() { fmt.Println("close") }

func (z *zigzag) ActionReceived(slot int, cmd common.Command) {
	switch cmd {
	case common.Command_LEFT_UP:
		z.engine.MovePlayer(slot, "left")
	case common.Command_UP_UP:
		z.engine.MovePlayer(slot, "up")
	case common.Command_RIGHT_UP:
		z.engine.MovePlayer(slot, "right")
	case common.Command_DOWN_UP:
		z.engine.MovePlayer(slot, "down")
	}
	z.print()
}

func (z *zigzag) print() {
	z.renderer.Print(z.engine.GetSnakes(), z.engine.GetCandies())
}
