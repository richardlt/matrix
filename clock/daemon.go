package clock

import (
	"context"
	"strconv"
	"time"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
)

// Start the clock software.
func Start(uri string) error {
	logrus.Infof("Start clock for uri %s\n", uri)

	c := &clock{}

	return software.Connect(uri, c, true)
}

type clock struct {
	api                         software.API
	layer                       software.Layer
	caracterBig, caracterMedium *software.CaracterDriver
	cancel                      func()
	green1, green2              common.Color
	blink                       bool
}

func (c *clock) Init(a software.API) (err error) {
	logrus.Debug("Init clock")

	c.api = a

	l := a.GetImageFromLocal("clock")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &l,
		MinPlayerCount: 1,
		MaxPlayerCount: 1,
	})

	c.layer, err = c.api.NewLayer()
	if err != nil {
		return err
	}

	c.caracterBig, err = c.layer.NewCaracterDriver(
		a.GetFontFromLocal("SevenByFour"))
	if err != nil {
		return err
	}

	c.caracterMedium, err = c.layer.NewCaracterDriver(
		a.GetFontFromLocal("FiveByThree"))
	if err != nil {
		return err
	}

	c.green1 = c.api.GetColorFromLocalThemeByName("flat", "green_1")
	c.green2 = c.api.GetColorFromLocalThemeByName("flat", "green_2")

	return a.Ready()
}

func (c *clock) Start(playerCount uint64) {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	c.print()

	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				c.print()
			}
		}
	}()
}

func (c *clock) Close() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c clock) ActionReceived(slot int, cmd common.Command) {}

func (c *clock) print() {
	c.layer.Clean()

	now := time.Now()
	h, m := now.Hour(), now.Minute()

	if h >= 10 {
		c.layer.SetWithCoord(common.Coord{X: 5, Y: 1}, c.green1)
	}
	if h >= 20 {
		c.layer.SetWithCoord(common.Coord{X: 6, Y: 1}, c.green1)
	}

	c.caracterBig.Render([]rune(strconv.Itoa(h % 10))[0],
		common.Coord{X: 3, Y: 4}, c.green2, common.Color{})

	c.caracterMedium.Render([]rune(strconv.Itoa(m / 10))[0],
		common.Coord{X: 9, Y: 5}, c.green2, common.Color{})

	c.caracterMedium.Render([]rune(strconv.Itoa(m % 10))[0],
		common.Coord{X: 13, Y: 5}, c.green2, common.Color{})

	if c.blink {
		c.layer.SetWithCoord(common.Coord{X: 6, Y: 4}, c.green1)
		c.layer.SetWithCoord(common.Coord{X: 6, Y: 6}, c.green1)
	}

	c.blink = !c.blink

	c.api.Print()
}
