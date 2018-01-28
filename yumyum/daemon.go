package yumyum

import (
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
)

// Start the zigzag software.
func Start(uri string) error {
	logrus.Infof("Start yumyum for uri %s\n", uri)

	y := &yumyum{}

	return software.Connect(uri, y, true)
}

type yumyum struct{}

func (y yumyum) Init(a software.API) (err error) {
	logrus.Debug("Init yumyum")

	l := a.GetImageFromLocal("yumyum")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &l,
		MinPlayerCount: 1,
		MaxPlayerCount: 4,
	})

	a.Ready()
	return nil
}

func (y yumyum) Start() {}

func (y yumyum) Close() {}

func (y yumyum) ActionReceived(slot int, cmd common.Command) {}
