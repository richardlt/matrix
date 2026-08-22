package blocks

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

	l3, err := a.NewLayer()
	if err != nil {
		return nil, err
	}

	td, err := l3.NewTextDriver(a.GetFontFromLocal("FiveByFive"))
	if err != nil {
		return nil, err
	}

	td.OnStep(func(total, current uint64) { _ = a.Print() })

	return &renderer{
		api: a, layerInfo: l3, layerPiece: l2, layerStack: l1, textDriver: td,
		pieceColors: []*common.Color{
			a.GetColorFromLocalThemeByName("flat", "green_2"),
			a.GetColorFromLocalThemeByName("flat", "blue_2"),
			a.GetColorFromLocalThemeByName("flat", "violet_2"),
			a.GetColorFromLocalThemeByName("flat", "white_2"),
			a.GetColorFromLocalThemeByName("flat", "red_2"),
			a.GetColorFromLocalThemeByName("flat", "orange_2"),
			a.GetColorFromLocalThemeByName("flat", "yellow_2"),
		},
	}, nil
}

type renderer struct {
	api                               software.API
	pieceColors                       []*common.Color
	layerInfo, layerPiece, layerStack software.Layer
	textDriver                        *software.TextDriver
}

func (r *renderer) Clean() {
	_ = r.layerPiece.Clean()
	_ = r.layerStack.Clean()
	_ = r.layerInfo.Clean()
}

func (r *renderer) Print(stack map[coord]pieceType, p *piece) {
	_ = r.layerStack.Clean()
	for c, t := range stack {
		_ = r.layerStack.SetWithCoord(&common.Coord{X: int64(c.x), Y: int64(c.y)},
			r.pieceColors[int(t)])
	}

	_ = r.layerPiece.Clean()
	if p != nil {
		for _, c := range p.ToCoords() {
			_ = r.layerPiece.SetWithCoord(&common.Coord{X: int64(c.x), Y: int64(c.y)},
				r.pieceColors[int(p.Type)])
		}
	}

	_ = r.api.Print()
}

func (r *renderer) StartPrintScore(score int) {
	_ = r.layerInfo.Clean()
	_ = r.textDriver.Render(fmt.Sprintf("%d PTS", score), &common.Coord{X: 10, Y: 4},
		r.api.GetColorFromLocalThemeByName("flat", "red_2"), &common.Color{}, true)
}

func (r *renderer) StopPrintInfo() {
	_ = r.layerInfo.Clean()
	_ = r.textDriver.Stop()
}
