package getout

import (
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

	l3, err := a.NewLayer()
	if err != nil {
		return nil, err
	}

	td, err := l3.NewTextDriver(a.GetFontFromLocal("FiveByFive"))
	if err != nil {
		return nil, err
	}

	td.OnStep(func(total, current uint64) { a.Print() })

	return &renderer{
		api:         a,
		layerGrid:   l1,
		layerPlayer: l2,
		layerInfo:   l3,
		textDriver:  td,
	}, nil
}

type renderer struct {
	api                               software.API
	layerGrid, layerPlayer, layerInfo software.Layer
	textDriver                        *software.TextDriver
}

func (r *renderer) clean() {
	r.layerGrid.Clean()
	r.layerPlayer.Clean()
	r.layerInfo.Clean()
}

func (r *renderer) printGrid(grid [][]int) {
	r.layerGrid.Clean()
	for x, column := range grid {
		for y, cell := range column {
			if cell == 0 || cell == 2 {
				r.layerGrid.SetWithCoord(common.Coord{X: int64(x), Y: int64(y)},
					common.Color{R: 255, G: 255, B: 255, A: 1})
			}
		}
	}
	r.api.Print()
}

func (r *renderer) printPlayer(player, end coord) {
	r.layerPlayer.Clean()
	r.layerPlayer.SetWithCoord(common.Coord{X: int64(end.x), Y: int64(end.y)},
		r.api.GetColorFromLocalThemeByName("flat", "red_2"))
	r.layerPlayer.SetWithCoord(common.Coord{X: int64(player.x), Y: int64(player.y)},
		r.api.GetColorFromLocalThemeByName("flat", "blue_2"))
	r.api.Print()
}

func (r *renderer) startPrintGameOver() {
	r.layerInfo.Clean()
	r.textDriver.Render("GAME OVER", common.Coord{X: 10, Y: 4},
		r.api.GetColorFromLocalThemeByName("flat", "red_2"),
		common.Color{A: 1}, true)
}

func (r *renderer) stopPrintGameOver() { r.textDriver.Stop() }
