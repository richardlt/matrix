package blocks

import (
	"context"
	"time"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
)

// Start the blocks software.
func Start(uri string) error {
	logrus.Infof("Start blocks for uri %s\n", uri)

	t := &blocks{}

	return software.Connect(uri, t, true)
}

type blocks struct {
	engine      *engine
	renderer    *renderer
	cancel      func()
	commandChan chan common.Command
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
		ti := time.NewTicker(time.Millisecond * 600)
		defer ti.Stop()

		var gameOver bool
		for !gameOver {
			select {
			case <-ctx.Done():
				close(b.commandChan)
				b.commandChan = nil
				return
			case <-ti.C:
				b.engine.MovePiece()
				b.print()
				gameOver = b.engine.IsGameOver()
			case cmd := <-b.commandChan:
				switch cmd {
				case common.Command_UP_UP:
					b.engine.MovePieceUp()
				case common.Command_DOWN_UP:
					b.engine.MovePieceDown()
				case common.Command_A_UP:
					b.engine.RotatePiece()
				case common.Command_RIGHT_UP:
					b.engine.MovePiece()
				}
				b.print()
			}
		}
	}()
}

func (b *blocks) Close() {
	if b.cancel != nil {
		b.cancel()
	}
	b.renderer.Clean()
}

func (b *blocks) ActionReceived(slot int, cmd common.Command) {
	if b.commandChan != nil {
		b.commandChan <- cmd
	}
}

func (b *blocks) print() {
	p := b.engine.GetPiece()
	if p != nil {
		b.renderer.Print(b.engine.GetBlocks(), *p)
	}
}
