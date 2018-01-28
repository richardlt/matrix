package draw

import (
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
)

// Start the draw software.
func Start(uri string) error {
	logrus.Infof("Start draw for uri %s\n", uri)

	d := &draw{}

	return software.Connect(uri, d, true)
}

type draw struct {
	api                    software.API
	layerDraw, layerPlayer software.Layer
	colors                 []common.Color
	players                []*player
	playerCount            uint64
}

func (d *draw) Init(a software.API) (err error) {
	logrus.Debug("Init draw")

	d.api = a

	l := a.GetImageFromLocal("draw")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &l,
		MinPlayerCount: 1,
		MaxPlayerCount: 4,
	})

	d.layerDraw, err = d.api.NewLayer()
	if err != nil {
		return err
	}
	d.layerPlayer, err = d.api.NewLayer()
	if err != nil {
		return err
	}

	d.colors = []common.Color{
		d.api.GetColorFromThemeByName("flat", "turquoise_1"),
		d.api.GetColorFromThemeByName("flat", "green_1"),
		d.api.GetColorFromThemeByName("flat", "blue_1"),
		d.api.GetColorFromThemeByName("flat", "violet_1"),
		d.api.GetColorFromThemeByName("flat", "dark_grey_1"),
		d.api.GetColorFromThemeByName("flat", "grey_1"),
		d.api.GetColorFromThemeByName("flat", "white_1"),
		d.api.GetColorFromThemeByName("flat", "red_1"),
		d.api.GetColorFromThemeByName("flat", "orange_1"),
		d.api.GetColorFromThemeByName("flat", "yellow_1"),
	}

	a.Ready()
	return nil
}

func (d *draw) Start(playerCount uint64) {
	d.playerCount = playerCount

	d.players = make([]*player, d.playerCount)
	for i := 0; i < int(d.playerCount); i++ {
		var x, y int64
		if i == 0 || i == 3 {
			x = 0
		} else {
			x = 15
		}
		if i == 0 || i == 2 {
			y = 0
		} else {
			y = 8
		}
		d.players[i] = &player{
			Color: 6,
			Coord: common.Coord{X: x, Y: y},
		}
		d.layerPlayer.SetWithCoord(d.players[i].Coord, d.colors[d.players[i].Color])
	}

	d.api.Print()
}

func (d draw) Close() {}

func (d *draw) ActionReceived(slot int, cmd common.Command) {
	switch cmd {
	case common.Command_A_UP:
		d.layerDraw.SetWithCoord(d.players[slot].Coord, d.colors[d.players[slot].Color])
		d.print()
	case common.Command_B_UP:
		d.layerDraw.SetWithCoord(d.players[slot].Coord, common.Color{})
		d.print()
	case common.Command_X_UP:
		if d.players[slot].Color < len(d.colors)-1 {
			d.players[slot].Color++
		} else {
			d.players[slot].Color = 0
		}
		d.print()
	case common.Command_LEFT_UP:
		if d.players[slot].Coord.X > 0 {
			d.players[slot].Coord.X--
		}
		d.print()
	case common.Command_UP_UP:
		if d.players[slot].Coord.Y > 0 {
			d.players[slot].Coord.Y--
		}
		d.print()
	case common.Command_RIGHT_UP:
		if d.players[slot].Coord.X < 15 {
			d.players[slot].Coord.X++
		}
		d.print()
	case common.Command_DOWN_UP:
		if d.players[slot].Coord.Y < 8 {
			d.players[slot].Coord.Y++
		}
		d.print()
	}

}

func (d *draw) print() {
	d.layerPlayer.Clean()

	for _, p := range d.players {
		d.layerPlayer.SetWithCoord(p.Coord, d.colors[p.Color])
	}

	d.api.Print()
}

type player struct {
	Color int
	Coord common.Coord
}
