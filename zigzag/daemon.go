package zigzag

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
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

	return a.Ready()
}

func (z *zigzag) Start(playerCount uint64) {
	z.engine = newEngine(playerCount, 16, 9)
	z.print()

	ctx, cancel := context.WithCancel(context.Background())
	z.cancel = cancel

	go func() {
		t := time.NewTicker(time.Millisecond * 300)
		defer t.Stop()

		var gameOver bool
		for !gameOver {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				z.engine.MovePlayers()
				z.print()
				gameOver = z.engine.IsGameOver()
			}
		}

		z.renderer.StartPrintWinners(z.engine.GetWinners())
	}()
}

func (z *zigzag) Close() {
	if z.cancel != nil {
		z.cancel()
	}
	z.renderer.Clean()
	z.renderer.StopPrintWinners()
}

func (z *zigzag) ActionReceived(slot uint64, cmd common.Command) {
	pSlot := int(slot)
	switch cmd {
	case common.Command_LEFT_UP:
		z.engine.ChangePlayerDirection(pSlot, "left")
	case common.Command_UP_UP:
		z.engine.ChangePlayerDirection(pSlot, "up")
	case common.Command_RIGHT_UP:
		z.engine.ChangePlayerDirection(pSlot, "right")
	case common.Command_DOWN_UP:
		z.engine.ChangePlayerDirection(pSlot, "down")
	}
	z.print()
}

func (z *zigzag) print() {
	z.renderer.Print(z.engine.GetSnakes(), z.engine.GetCandies())
}
