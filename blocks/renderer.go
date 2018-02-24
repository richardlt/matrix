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

	td, err := l1.NewTextDriver(a.GetFontFromLocal("FiveByFive"))
	if err != nil {
		return nil, err
	}

	td.OnStep(func(total, current uint64) { a.Print() })

	return &renderer{
		api: a, layerInfo: l3, layerPiece: l2, layerStack: l1, textDriver: td,
		pieceColors: []common.Color{
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
	pieceColors                       []common.Color
	layerInfo, layerPiece, layerStack software.Layer
	textDriver                        *software.TextDriver
}

func (r *renderer) Clean() {
	r.layerPiece.Clean()
	r.layerStack.Clean()
	r.layerInfo.Clean()
}

func (r *renderer) Print(stack map[common.Coord]pieceType, p *piece) {
	r.layerStack.Clean()
	for c, t := range stack {
		r.layerStack.SetWithCoord(c, r.pieceColors[int(t)])
	}

	r.layerPiece.Clean()
	if p != nil {
		for _, c := range p.ToCoords() {
			r.layerPiece.SetWithCoord(c, r.pieceColors[int(p.Type)])
		}
	}

	r.api.Print()
}

func (r *renderer) StartPrintScore(score int) {
	r.layerInfo.Clean()
	r.textDriver.Render(fmt.Sprintf("%d PTS", score), common.Coord{X: 10, Y: 4},
		r.api.GetColorFromLocalThemeByName("flat", "red_2"), common.Color{}, true)
}

func (r *renderer) StopPrintScore() { r.textDriver.Stop() }
