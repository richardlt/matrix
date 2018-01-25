package zigzag

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
		snakesColor: []common.Color{
			a.GetColorFromThemeByName("zigzag", "player1"),
			a.GetColorFromThemeByName("zigzag", "player2"),
			a.GetColorFromThemeByName("zigzag", "player3"),
			a.GetColorFromThemeByName("zigzag", "player4"),
		},
		candiesColor: []common.Color{
			a.GetColorFromThemeByName("zigzag", "candy1"),
			a.GetColorFromThemeByName("zigzag", "candy2"),
		},
	}, nil
}

type renderer struct {
	api          software.API
	snakesColor  []common.Color
	candiesColor []common.Color
	layer        software.Layer
}

func (r renderer) Print(ss []snake, cs []candy) {
	r.layer.Clean()
	for i, s := range ss {
		for j, b := range s.Body {
			color := r.snakesColor[i]
			if j == 0 {
				color.R -= 40
				color.G -= 40
				color.B -= 40
			}
			r.layer.SetWithCoord(b.Convert(), color)
		}
	}
	for _, c := range cs {
		if c.State {
			r.layer.SetWithCoord(c.Coord.Convert(), r.candiesColor[c.Points])
		}
	}
	r.api.Print()
}
