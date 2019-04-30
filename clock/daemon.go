package clock

import (
	"context"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
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
	blink                       bool
	colors                      []common.Color
	model, color                int
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

	c.colors = []common.Color{
		a.GetColorFromLocalThemeByName("flat", "green_2"),
		a.GetColorFromLocalThemeByName("flat", "blue_2"),
		a.GetColorFromLocalThemeByName("flat", "violet_2"),
		a.GetColorFromLocalThemeByName("flat", "white_2"),
		a.GetColorFromLocalThemeByName("flat", "red_2"),
		a.GetColorFromLocalThemeByName("flat", "orange_2"),
		a.GetColorFromLocalThemeByName("flat", "yellow_2"),
	}

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
				c.blink = !c.blink
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

func (c *clock) ActionReceived(slot uint64, cmd common.Command) {
	switch cmd {
	case common.Command_A_UP:
		if c.color < len(c.colors)-1 {
			c.color++
		} else {
			c.color = 0
		}
		c.print()
	case common.Command_LEFT_UP:
		if c.model > 0 {
			c.model--
		} else {
			c.model = 2
		}
		c.print()
	case common.Command_RIGHT_UP:
		if c.model < 2 {
			c.model++
		} else {
			c.model = 0
		}
		c.print()
	}
}

func (c *clock) print() {
	c.layer.Clean()

	now := time.Now()
	h, m := now.Hour(), now.Minute()
	color := c.colors[c.color]

	switch c.model {
	case 0:
		c.printModelBig(h, m, color)
	case 1:
		c.printModelMedium(h, m, color, false)
	case 2:
		c.printModelMedium(h, m, color, true)
	}

	c.api.Print()
}

func (c *clock) printModelBig(h, m int, color common.Color) {
	if h >= 10 {
		c.layer.SetWithCoord(common.Coord{X: 0, Y: 7}, color)
	}
	if h >= 20 {
		c.layer.SetWithCoord(common.Coord{X: 0, Y: 6}, color)
	}

	c.caracterBig.Render([]rune(strconv.Itoa(h % 10))[0],
		common.Coord{X: 4, Y: 4}, color, common.Color{})

	c.caracterMedium.Render([]rune(strconv.Itoa(m / 10))[0],
		common.Coord{X: 10, Y: 5}, color, common.Color{})

	c.caracterMedium.Render([]rune(strconv.Itoa(m % 10))[0],
		common.Coord{X: 14, Y: 5}, color, common.Color{})

	if c.blink {
		c.layer.SetWithCoord(common.Coord{X: 7, Y: 4}, color)
		c.layer.SetWithCoord(common.Coord{X: 7, Y: 6}, color)
	}
}

func (c *clock) printModelMedium(h, m int, color common.Color, static bool) {
	var offset int64
	if static || c.blink {
		offset = 1
	}

	c.caracterMedium.Render([]rune(strconv.Itoa(h / 10))[0],
		common.Coord{X: 1, Y: 3 + offset}, color, common.Color{})

	c.caracterMedium.Render([]rune(strconv.Itoa(h % 10))[0],
		common.Coord{X: 5, Y: 3 + offset}, color, common.Color{})

	c.caracterMedium.Render([]rune(strconv.Itoa(m / 10))[0],
		common.Coord{X: 10, Y: 5 - offset}, color, common.Color{})

	c.caracterMedium.Render([]rune(strconv.Itoa(m % 10))[0],
		common.Coord{X: 14, Y: 5 - offset}, color, common.Color{})
}
