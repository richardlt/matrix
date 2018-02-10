package device

import (
	"github.com/richardlt/matrix/sdk-go/player"
)

func newGamepad() *gamepad {
	return &gamepad{}
}

type gamepad struct {
	api *player.API
}

func (g *gamepad) Init(api *player.API) error {
	g.api = api
	return nil
}
