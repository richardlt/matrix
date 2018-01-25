package core

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/richardlt/matrix/core/os"
	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/core/system"
	"github.com/richardlt/matrix/sdk-go/display"
	"github.com/richardlt/matrix/sdk-go/player"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Start the matrix core.
func Start(port int) error {
	logrus.Infof("Start core on port %d\n", port)

	if err := render.Init(); err != nil {
		return err
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.WithStack(err)
	}

	s := grpc.NewServer()

	ps := system.NewPlayerServer()
	ds := system.NewDisplayServer()
	ss := system.NewSoftwareServer()

	c := os.NewContext(ps, ds, ss)

	player.RegisterPlayerServer(s, ps)
	display.RegisterDisplayServer(s, ds)
	software.RegisterSoftwareServer(s, ss)

	os.StartContext(c)

	return errors.WithStack(s.Serve(ln))
}
