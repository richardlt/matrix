package blocks

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

	return &renderer{
		api: a, layerPiece: l2, layerStack: l1,
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
	api                    software.API
	pieceColors            []common.Color
	layerPiece, layerStack software.Layer
}

func (r *renderer) Clean() {
	r.layerPiece.Clean()
	r.layerStack.Clean()
}

func (r *renderer) Print(blocks []block, piece piece) {
	r.printBlocks(blocks)
	r.printPiece(piece)
	r.api.Print()
}

func (r *renderer) printBlocks(blocks []block) {
	r.layerStack.Clean()
	for _, b := range blocks {
		r.layerStack.SetWithCoord(b.Coord, r.pieceColors[int(b.Type)])
	}
}

func (r *renderer) printPiece(piece piece) {
	r.layerPiece.Clean()
	for _, c := range piece.ToCoords() {
		r.layerPiece.SetWithCoord(c, r.pieceColors[int(piece.Type)])
	}
}
