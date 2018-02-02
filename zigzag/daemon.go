package zigzag

import (
	"context"
	"time"

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
	cancel   func()
}

func (z *zigzag) Init(a software.API) (err error) {
	logrus.Debug("Init zigzag")

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

func (z *zigzag) Start(playerCount uint64) {
	z.engine = newEngine(playerCount, 16, 9)
	z.print()

	ctx, cancel := context.WithCancel(context.Background())
	z.cancel = cancel

	go func() {
		t := time.NewTicker(time.Millisecond * 300)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				z.engine.MovePlayers()
				z.print()
			}
		}
	}()
}

func (z *zigzag) Close() {
	if z.cancel != nil {
		z.cancel()
	}
}

func (z *zigzag) ActionReceived(slot int, cmd common.Command) {
	switch cmd {
	case common.Command_LEFT_UP:
		z.engine.ChangePlayerDirection(slot, "left")
	case common.Command_UP_UP:
		z.engine.ChangePlayerDirection(slot, "up")
	case common.Command_RIGHT_UP:
		z.engine.ChangePlayerDirection(slot, "right")
	case common.Command_DOWN_UP:
		z.engine.ChangePlayerDirection(slot, "down")
	}
	z.print()
}

func (z *zigzag) print() {
	z.renderer.Print(z.engine.GetSnakes(), z.engine.GetCandies())
}
