package yumyum

import (
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

func newRenderer(a software.API) (*renderer, error) {
	l, err := a.NewLayer()
	if err != nil {
		return nil, err
	}

	return &renderer{
		api:   a,
		layer: l,
		playersColor: []common.Color{
			a.GetColorFromThemeByName("yumyum", "player1"),
			a.GetColorFromThemeByName("yumyum", "player2"),
			a.GetColorFromThemeByName("yumyum", "player3"),
			a.GetColorFromThemeByName("yumyum", "player4"),
		},
		candiesColor: []common.Color{
			a.GetColorFromThemeByName("yumyum", "candy1"),
			a.GetColorFromThemeByName("yumyum", "candy2"),
		},
	}, nil
}

type renderer struct {
	api          software.API
	playersColor []common.Color
	candiesColor []common.Color
	layer        software.Layer
}

func (r renderer) Print(ps []player, cs []candy) {
	r.layer.Clean()
	for i, p := range ps {
		r.layer.SetWithCoord(p.Coord.Convert(), r.playersColor[i])
	}
	for _, c := range cs {
		if c.State {
			r.layer.SetWithCoord(c.Coord.Convert(), r.candiesColor[c.Points])
		}
	}
	r.api.Print()
}
