package device

import (
	"context"

	"github.com/richardlt/matrix/sdk-go/display"
	"github.com/richardlt/matrix/sdk-go/player"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
)

type action struct {
	Slot    uint64 `json:"slot"`
	Command string `json:"command"`
}

// Start device deamon.
func Start(uri string) error {
	m := newMatrix()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := m.OpenPorts(ctx); err != nil {
			logrus.Errorf("%+v", err)
		}
	}()

	g := newGamepad()

	go func() {
		if err := display.Connect(uri, m, true); err != nil {
			logrus.Errorf("%+v", err)
		}
	}()

	go func() {
		if err := player.Connect(uri, g, true); err != nil {
			logrus.Errorf("%+v", err)
		}
	}()

	return software.Connect(uri, m, true)
}
