package getout

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
	}, nil
}

type renderer struct {
	api   software.API
	layer software.Layer
}

func (r *renderer) Clean() {}

func (r *renderer) Print(grid [][]int, start, end coord) {
	r.layer.Clean()
	for x, column := range grid {
		for y, cell := range column {
			if start.x == x && start.y == y {
				r.layer.SetWithCoord(common.Coord{X: int64(x), Y: int64(y)},
					common.Color{B: 255, A: 1})
			} else if end.x == x && end.y == y {
				r.layer.SetWithCoord(common.Coord{X: int64(x), Y: int64(y)},
					common.Color{R: 255, A: 1})
			} else if cell == 0 || cell == 2 {
				r.layer.SetWithCoord(common.Coord{X: int64(x), Y: int64(y)},
					common.Color{R: 255, G: 255, B: 255, A: 1})
			}
		}
	}
	r.api.Print()
}
