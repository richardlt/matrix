package device

import (
	"context"
	
	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/sdk-go/display"
	"github.com/richardlt/matrix/sdk-go/player"
	"github.com/richardlt/matrix/sdk-go/software"
)

// Start device deamon.
func Start(uri string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := newMatrix()
	go func() {
		if err := m.OpenPorts(ctx); err != nil {
			logrus.Errorf("%+v", err)
		}
	}()

	g := newGamepad()
	go func() {
		if err := g.OpenDevices(ctx); err != nil {
			logrus.Errorf("%+v", err)
		}
	}()

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
