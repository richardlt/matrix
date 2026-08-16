package zigzag

import (
	"context"
	"sync/atomic"
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
	paused   atomic.Bool
}

func (z *zigzag) Init(a software.API) (err error) {
	logrus.Debug("Init zigzag")

	z.renderer, err = newRenderer(a)
	if err != nil {
		return err
	}

	l := a.GetImageFromLocal("zigzag")

	if err := a.SetConfig(&software.ConnectRequest_SoftwareData_Config{
		Logo:           l,
		MinPlayerCount: 1,
		MaxPlayerCount: 4,
		Pausable:       true,
	}); err != nil {
		return err
	}

	return a.Ready()
}

func (z *zigzag) Start(playerCount uint64) {
	z.paused.Store(false)
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
				if z.paused.Load() {
					continue
				}
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

// Paused satisfies software.Pausable: the core suspends play, and the snakes have to
// stop moving for as long as it lasts.
func (z *zigzag) Paused(paused bool) { z.paused.Store(paused) }

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
