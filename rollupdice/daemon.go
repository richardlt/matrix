package rollupdice

import (
	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// Start the rollup dice software.
func Start(uri string) error {
	logrus.Infof("Start rollup dice for uri %s\n", uri)

	r := &rollupDice{}

	return software.Connect(uri, r, true)
}

type rollupDice struct {
	engine   *engine
	renderer *renderer
}

func (r *rollupDice) Init(a software.API) (err error) {
	logrus.Debug("Init rollup dice")

	r.renderer, err = newRenderer(a)
	if err != nil {
		return err
	}

	if err := a.SetConfig(&software.ConnectRequest_SoftwareData_Config{
		Logo:           a.GetImageFromLocal("rollupdice"),
		MinPlayerCount: 1,
		MaxPlayerCount: 1,
	}); err != nil {
		return err
	}

	return a.Ready()
}

func (r *rollupDice) Start(uint64) {
	r.engine = newEngine()
	r.engine.roll()
	r.renderer.print(r.engine.values)
}

func (r *rollupDice) Close() { r.renderer.clean() }

func (r *rollupDice) ActionReceived(slot uint64, cmd common.Command) {
	if cmd != common.Command_A_UP {
		return
	}
	r.engine.roll()
	r.renderer.print(r.engine.values)
}
