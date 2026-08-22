package zigzag

import (
	"fmt"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

func newRenderer(a software.API) (*renderer, error) {
	l1, err := a.NewLayer()
	if err != nil {
		return nil, err
	}

	l2, err := a.NewLayer()
	if err != nil {
		return nil, err
	}

	td, err := l2.NewTextDriver(a.GetFontFromLocal("FiveByFive"))
	if err != nil {
		return nil, err
	}

	td.OnStep(func(total, current uint64) { _ = a.Print() })

	return &renderer{
		api: a, layerInfo: l2, layer: l1, textDriver: td,
		snakesColor: []*common.Color{
			a.GetColorFromLocalThemeByName("zigzag", "player1"),
			a.GetColorFromLocalThemeByName("zigzag", "player2"),
			a.GetColorFromLocalThemeByName("zigzag", "player3"),
			a.GetColorFromLocalThemeByName("zigzag", "player4"),
		},
		candiesColor: []*common.Color{
			a.GetColorFromLocalThemeByName("zigzag", "candy1"),
			a.GetColorFromLocalThemeByName("zigzag", "candy2"),
		},
	}, nil
}

type renderer struct {
	api                       software.API
	snakesColor, candiesColor []*common.Color
	layerInfo, layer          software.Layer
	textDriver                *software.TextDriver
}

func (r *renderer) Clean() {
	_ = r.layerInfo.Clean()
	_ = r.layer.Clean()
}

func (r *renderer) Print(ss []snake, cs []candy) {
	_ = r.layer.Clean()
	for i, s := range ss {
		for j, b := range s.Body {
			color := r.snakesColor[i]
			if j == 0 {
				color.R -= 40
				color.G -= 40
				color.B -= 40
			}
			_ = r.layer.SetWithCoord(b.Convert(), color)
		}
	}
	for _, c := range cs {
		if c.State {
			_ = r.layer.SetWithCoord(c.Coord.Convert(), r.candiesColor[c.Points])
		}
	}
	_ = r.api.Print()
}

func (r *renderer) StartPrintWinners(winners []int) {
	_ = r.layerInfo.Clean()

	var text string
	if len(winners) > 1 {
		var list string
		for _, w := range winners {
			list = fmt.Sprintf("%s %d", list, w+1)
		}
		text = fmt.Sprintf("PLAYERS%s WON", list)
	} else {
		text = fmt.Sprintf("PLAYER %d WON", winners[0]+1)
	}

	_ = r.textDriver.Render(text, &common.Coord{X: 10, Y: 4},
		r.api.GetColorFromLocalThemeByName("flat", "red_2"),
		&common.Color{}, true)
}

func (r *renderer) StopPrintWinners() { _ = r.textDriver.Stop() }
