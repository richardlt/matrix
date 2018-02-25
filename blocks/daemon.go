package blocks

import (
	"context"
	"time"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
)

const (
	commandLR             common.Command = 100
	commandHoldLeft       common.Command = 101
	commandHoldUp         common.Command = 102
	commandHoldRight      common.Command = 103
	commandHoldDown       common.Command = 104
	longPressTriggerDelay time.Duration  = 200 * time.Millisecond
	longPressFireDelay    time.Duration  = 50 * time.Millisecond
)

// Start the blocks software.
func Start(uri string) error {
	logrus.Infof("Start blocks for uri %s\n", uri)

	mp := software.NewMultiPress(common.Button_L, common.Button_R)
	lpl := software.NewLongPress(common.Button_LEFT,
		longPressTriggerDelay, longPressFireDelay)
	lpu := software.NewLongPress(common.Button_UP,
		longPressTriggerDelay, longPressFireDelay)
	lpr := software.NewLongPress(common.Button_RIGHT,
		longPressTriggerDelay, longPressFireDelay)
	lpd := software.NewLongPress(common.Button_DOWN,
		longPressTriggerDelay, longPressFireDelay)
	b := &blocks{
		mutliPressLR:   mp,
		longPressLeft:  lpl,
		longPressUp:    lpu,
		longPressRight: lpr,
		longPressDown:  lpd,
	}

	mp.OnAction(func(slot uint64) { b.action(commandLR) })
	lpl.OnAction(func(slot uint64) { b.action(commandHoldLeft) })
	lpu.OnAction(func(slot uint64) { b.action(commandHoldUp) })
	lpr.OnAction(func(slot uint64) { b.action(commandHoldRight) })
	lpd.OnAction(func(slot uint64) { b.action(commandHoldDown) })

	return software.Connect(uri, b, true)
}

type blocks struct {
	engine      *engine
	renderer    *renderer
	cancel      func()
	commandChan chan common.Command
	mutliPressLR, longPressLeft, longPressUp,
	longPressRight, longPressDown software.ActionGenerator
	rotateCommand bool
}

func (b *blocks) Init(a software.API) (err error) {
	logrus.Debug("Init blocks")

	b.renderer, err = newRenderer(a)
	if err != nil {
		return err
	}

	l := a.GetImageFromLocal("blocks")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &l,
		MinPlayerCount: 1,
		MaxPlayerCount: 1,
	})

	return a.Ready()
}

func (b *blocks) Start(uint64) {
	b.engine = newEngine(16, 9)
	b.print()

	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel

	b.commandChan = make(chan common.Command)

	go func() {
		ti := time.NewTicker(time.Millisecond * 500)
		defer ti.Stop()

		var gameOver bool
		for !gameOver {
			select {
			case <-ctx.Done():
				return
			case <-ti.C:
				b.engine.MovePiece()
				b.print()
				gameOver = b.engine.IsGameOver()
			case cmd := <-b.commandChan:
				switch cmd {
				case common.Command_LEFT_UP, commandHoldLeft:
					if b.rotateCommand {
						b.engine.MovePieceDown()
					}
				case common.Command_UP_UP, commandHoldUp:
					if !b.rotateCommand {
						b.engine.MovePieceUp()
					}
				case common.Command_RIGHT_UP, commandHoldRight:
					if b.rotateCommand {
						b.engine.MovePieceUp()
					} else {
						b.engine.MovePiece()
					}
				case common.Command_DOWN_UP, commandHoldDown:
					if !b.rotateCommand {
						b.engine.MovePieceDown()
					} else {
						b.engine.MovePiece()
					}
				case common.Command_A_UP:
					b.engine.RotatePiece()
				case commandLR:
					b.rotateCommand = !b.rotateCommand
				}
				b.print()
			}
		}

		b.renderer.Clean()
		b.renderer.StartPrintScore(b.engine.Score)
	}()
}

func (b *blocks) Close() {
	if b.cancel != nil {
		b.cancel()
	}
	b.renderer.Clean()
	b.renderer.StopPrintScore()
}

func (b *blocks) ActionReceived(slot uint64, cmd common.Command) {
	b.mutliPressLR.SendAction(slot, cmd)
	b.longPressLeft.SendAction(slot, cmd)
	b.longPressUp.SendAction(slot, cmd)
	b.longPressRight.SendAction(slot, cmd)
	b.longPressDown.SendAction(slot, cmd)
	b.action(cmd)
}

func (b *blocks) action(cmd common.Command) {
	if b.commandChan != nil {
		b.commandChan <- cmd
	}
}

func (b *blocks) print() { b.renderer.Print(b.engine.Stack, b.engine.Piece) }
