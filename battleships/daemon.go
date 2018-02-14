package battleships

import (
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
)

// Start the battleships software.
func Start(uri string) error {
	logrus.Infof("Start battleships for uri %s\n", uri)

	b := &battleships{}

	return software.Connect(uri, b, true)
}

type battleships struct {
	api software.API
}

func (b *battleships) Init(a software.API) (err error) {
	logrus.Debug("Init battleships")

	b.api = a

	l := a.GetImageFromLocal("battleships")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &l,
		MinPlayerCount: 2,
		MaxPlayerCount: 2,
	})

	return a.Ready()
}

func (b battleships) Start(playerCount uint64) {}

func (b battleships) Close() {}

func (b battleships) ActionReceived(slot int, cmd common.Command) {}
