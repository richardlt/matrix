package light

import (
	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// Start the light software.
func Start(uri string) error {
	logrus.Infof("Start light for uri %s\n", uri)

	d := &light{}

	return software.Connect(uri, d, true)
}

type light struct {
	api      software.API
	layer    software.Layer
	colors   []common.Color
	selected int
}

func (l *light) Init(a software.API) (err error) {
	logrus.Debug("Init light")

	l.api = a

	logo := a.GetImageFromLocal("light")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &logo,
		MinPlayerCount: 1,
		MaxPlayerCount: 1,
	})

	l.layer, err = l.api.NewLayer()
	if err != nil {
		return err
	}

	l.colors = []common.Color{
		l.api.GetColorFromLocalThemeByName("flat", "white_1"),
		l.api.GetColorFromLocalThemeByName("flat", "turquoise_1"),
		l.api.GetColorFromLocalThemeByName("flat", "green_1"),
		l.api.GetColorFromLocalThemeByName("flat", "blue_1"),
		l.api.GetColorFromLocalThemeByName("flat", "violet_1"),
		l.api.GetColorFromLocalThemeByName("flat", "dark_grey_1"),
		l.api.GetColorFromLocalThemeByName("flat", "red_1"),
		l.api.GetColorFromLocalThemeByName("flat", "orange_1"),
		l.api.GetColorFromLocalThemeByName("flat", "yellow_1"),
	}

	return a.Ready()
}

func (l *light) Start(playerCount uint64) {
	l.print()
}

func (l light) Close() {}

func (l *light) ActionReceived(slot uint64, cmd common.Command) {
	switch cmd {
	case common.Command_LEFT_UP:
		if 0 < l.selected {
			l.selected--
		} else {
			l.selected = len(l.colors) - 1
		}
		l.print()
	case common.Command_RIGHT_UP:
		if l.selected < len(l.colors)-1 {
			l.selected++
		} else {
			l.selected = 0
		}
		l.print()
	}
}

func (l *light) print() {
	for x := 0; x < 16; x++ {
		for y := 0; y < 9; y++ {
			l.layer.SetWithCoord(common.Coord{
				X: int64(x),
				Y: int64(y),
			}, l.colors[l.selected])
		}
	}
	l.api.Print()
}
