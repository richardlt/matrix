package rollupdice

import (
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// A face is a square of this many pixels, which is what fits two dice on the 16x9
// display: a one pixel border all around and a two pixel gap down the middle.
const faceSize = 7

func newRenderer(a software.API) (*renderer, error) {
	faceLayer, err := a.NewLayer()
	if err != nil {
		return nil, err
	}

	pipLayer, err := a.NewLayer()
	if err != nil {
		return nil, err
	}

	cd, err := pipLayer.NewCaracterDriver(a.GetFontFromLocal("DiceValues"))
	if err != nil {
		return nil, err
	}

	return &renderer{
		api:            a,
		faceLayer:      faceLayer,
		pipLayer:       pipLayer,
		caracterDriver: cd,
		faceColor:      &common.Color{R: 255, G: 255, B: 255, A: 1},
		pipColor:       &common.Color{A: 1},
	}, nil
}

type renderer struct {
	api                 software.API
	faceLayer, pipLayer software.Layer
	caracterDriver      *software.CaracterDriver
	faceColor, pipColor *common.Color
}

func (r *renderer) clean() {
	_ = r.faceLayer.Clean()
	_ = r.pipLayer.Clean()
}

func (r *renderer) print(values [diceCount]int) {
	r.clean()

	for i, value := range values {
		originX := i * (faceSize + 2)

		for x := originX; x < originX+faceSize; x++ {
			for y := 1; y <= faceSize; y++ {
				_ = r.faceLayer.SetWithCoord(&common.Coord{X: int64(x), Y: int64(y)}, r.faceColor)
			}
		}

		// The pips of a value are one caracter of the DiceValues font, five by five, and
		// the driver places it around its center. The background is left transparent so
		// the face below shows through.
		_ = r.caracterDriver.Render(rune('0'+value),
			&common.Coord{X: int64(originX + faceSize/2), Y: int64(1 + faceSize/2)},
			r.pipColor, &common.Color{})
	}

	_ = r.api.Print()
}
