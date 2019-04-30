package demo

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// Start the demo software.
func Start(uri string) error {
	logrus.Infof("Start demo for uri %s\n", uri)

	d := &demo{}

	return software.Connect(uri, d, true)
}

type demo struct {
	api            software.API
	layer          software.Layer
	randomDriver   *software.RandomDriver
	caracterDriver *software.CaracterDriver
	textDriver     *software.TextDriver
	imageDriver    *software.ImageDriver
	step           int
	cancel         context.CancelFunc
}

func (d *demo) Init(a software.API) (err error) {
	logrus.Debug("Init demo")

	d.api = a

	i := a.GetImageFromLocal("demo")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &i,
		MinPlayerCount: 1,
		MaxPlayerCount: 4,
	})

	d.layer, err = a.NewLayer()
	if err != nil {
		return err
	}

	d.randomDriver, err = d.layer.NewRandomDriver()
	if err != nil {
		return err
	}
	d.randomDriver.OnEnd(func() { d.api.Print() })

	d.caracterDriver, err = d.layer.NewCaracterDriver(a.GetFontFromLocal("FiveByFive"))
	if err != nil {
		return err
	}
	d.caracterDriver.OnEnd(func() { d.api.Print() })

	d.textDriver, err = d.layer.NewTextDriver(a.GetFontFromLocal("FiveByFive"))
	if err != nil {
		return err
	}
	d.textDriver.OnStep(func(total, current uint64) { d.api.Print() })

	d.imageDriver, err = d.layer.NewImageDriver()
	if err != nil {
		return err
	}
	d.imageDriver.OnEnd(func() { d.api.Print() })

	return a.Ready()
}

func (d *demo) Start(playerCount uint64) { d.play() }

func (d *demo) Close() { d.reset() }

func (d *demo) ActionReceived(slot uint64, cmd common.Command) {
	switch cmd {
	case common.Command_LEFT_UP:
		if d.step < 1 {
			d.step = 4
		} else {
			d.step--
		}
		d.play()
	case common.Command_RIGHT_UP:
		if 3 < d.step {
			d.step = 0
		} else {
			d.step++
		}
		d.play()
	}
}

func (d *demo) reset() {
	if d.cancel != nil {
		d.cancel()
	}
	if d.textDriver != nil {
		d.textDriver.Stop()
	}
}

func (d *demo) play() {
	d.reset()

	d.layer.Clean()
	d.api.Print()

	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel

	switch d.step {
	case 0:
		d.playRandom(ctx)
	case 1:
		d.playCaracter()
	case 2:
		d.playText()
	case 3:
		d.playTextRepeat()
	case 4:
		d.playImage(ctx)
	}
}

func (d *demo) playRandom(ctx context.Context) {
	ticker := time.NewTicker(time.Millisecond * 25)
	defer ticker.Stop()
	d.randomDriver.Render()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.randomDriver.Render()
		}
	}
}

func (d *demo) playCaracter() {
	d.caracterDriver.Render('A', common.Coord{X: 5, Y: 3},
		d.api.GetColorFromLocalThemeByName("flat", "red_2"), common.Color{})
	d.caracterDriver.Render('B', common.Coord{X: 8, Y: 4},
		d.api.GetColorFromLocalThemeByName("flat", "orange_2"), common.Color{})
	d.caracterDriver.Render('C', common.Coord{X: 10, Y: 5},
		d.api.GetColorFromLocalThemeByName("flat", "green_2"), common.Color{})
}

func (d *demo) playText() {
	d.textDriver.OnEnd(func() {
		time.Sleep(500 * time.Millisecond)
		d.layer.Clean()
		d.textDriver.OnEnd(func() {
			time.Sleep(500 * time.Millisecond)
			d.layer.Clean()
			d.playText()
		})
		d.textDriver.Render("SOFTWARE", common.Coord{X: 0, Y: 6},
			d.api.GetColorFromLocalThemeByName("flat", "green_2"),
			common.Color{}, false)
	})
	d.textDriver.Render("EXAMPLE", common.Coord{X: 0, Y: 2},
		d.api.GetColorFromLocalThemeByName("flat", "red_2"),
		common.Color{}, false)
}

func (d *demo) playTextRepeat() {
	d.textDriver.OnEnd(func() {})
	d.textDriver.Render("REPEAT", common.Coord{X: 10, Y: 4},
		d.api.GetColorFromLocalThemeByName("flat", "red_2"),
		common.Color{}, true)
}

func (d *demo) playImage(ctx context.Context) {
	exec := func(nb int) {
		d.layer.Clean()
		if nb == 0 {
			d.imageDriver.Render(d.api.GetImageFromLocal("monster-one"),
				common.Coord{X: 6, Y: 4})
		} else {
			d.imageDriver.Render(d.api.GetImageFromLocal("monster-two"),
				common.Coord{X: 11, Y: 5})
		}
	}

	var nb int
	exec(nb)

	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			nb++
			if 1 < nb {
				nb = 0
			}
			exec(nb)
		}
	}
}
