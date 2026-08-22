package blocks

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

const (
	commandLR             common.Command = 100
	commandHoldLeft       common.Command = 101
	commandHoldUp         common.Command = 102
	commandHoldRight      common.Command = 103
	commandHoldDown       common.Command = 104
	commandPause          common.Command = 105
	commandResume         common.Command = 106
	longPressTriggerDelay time.Duration  = 200 * time.Millisecond
	longPressFireDelay    time.Duration  = 100 * time.Millisecond
	defaultMoveDelay      time.Duration  = 500 * time.Millisecond
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
	engine   *engine
	renderer *renderer
	// running is cancelled both when the software is closed and when the game ends, so
	// nothing is left waiting on a game loop that has already returned.
	running     context.Context
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

	if err := a.SetConfig(&software.ConnectRequest_SoftwareData_Config{
		Logo:           l,
		MinPlayerCount: 1,
		MaxPlayerCount: 1,
		Pausable:       true,
	}); err != nil {
		return err
	}

	return a.Ready()
}

func (b *blocks) Start(uint64) {
	b.engine = newEngine(16, 9)
	b.print()

	ctx, cancel := context.WithCancel(context.Background())
	b.running, b.cancel = ctx, cancel

	b.commandChan = make(chan common.Command)

	go func() {
		// Cancelling on the way out covers the game-over exit as well as Close, so a
		// player still pressing buttons afterwards is never left blocked on the send.
		defer cancel()

		ti := time.NewTicker(defaultMoveDelay)
		defer ti.Stop()

		var paused bool
		var gameOver bool
		for !gameOver {
			select {
			case <-ctx.Done():
				return
			case <-ti.C:
				if !paused {
					b.engine.MovePiece()
					b.print()
					gameOver = b.engine.IsGameOver()
				}
			case cmd := <-b.commandChan:
				switch cmd {
				case commandPause, commandResume:
					// The core holds the screen while paused, so there is nothing to
					// draw here, only the fall to stop.
					paused = cmd == commandPause
					continue
				}
				if paused {
					continue
				}

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
	b.renderer.StopPrintInfo()
}

// Paused satisfies software.Pausable: the core suspends play, and the piece has to stop
// falling for as long as it lasts.
func (b *blocks) Paused(paused bool) {
	if !paused {
		b.action(commandResume)
		return
	}

	b.action(commandPause)

	// No command reaches a paused software, so a button held when play stopped would
	// never be seen coming back up and would still be repeating on the way out.
	for _, g := range b.generators() {
		g.Reset()
	}
}

func (b *blocks) generators() []software.ActionGenerator {
	return []software.ActionGenerator{
		b.mutliPressLR, b.longPressLeft, b.longPressUp, b.longPressRight, b.longPressDown,
	}
}

func (b *blocks) ActionReceived(slot uint64, cmd common.Command) {
	for _, g := range b.generators() {
		g.SendAction(slot, cmd)
	}
	b.action(cmd)
}

// action hands a command to the running game. The channel is unbuffered and its only
// reader is the game loop, so the send has to give up once that loop is gone: every
// action arrives on its own goroutine, and a blocked send would strand one for good.
func (b *blocks) action(cmd common.Command) {
	if b.commandChan == nil || b.running == nil {
		return
	}
	select {
	case b.commandChan <- cmd:
	case <-b.running.Done():
	}
}

func (b *blocks) print() { b.renderer.Print(b.engine.Stack, b.engine.Piece) }
