package main

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/richardlt/matrix/blocks"
	"github.com/richardlt/matrix/clock"
	"github.com/richardlt/matrix/core"
	"github.com/richardlt/matrix/demo"
	"github.com/richardlt/matrix/device"
	"github.com/richardlt/matrix/draw"
	"github.com/richardlt/matrix/emulator"
	"github.com/richardlt/matrix/gamepad"
	"github.com/richardlt/matrix/getout"
	"github.com/richardlt/matrix/yumyum"
	"github.com/richardlt/matrix/zigzag"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{{
		Name:  "start",
		Usage: "start the matrix components",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "core-uri",
				Value:  "localhost:8080",
				EnvVar: "MATRIX_CORE_URI",
				Usage:  "Core URI is used by softwares, players and displays.",
			},
			cli.IntFlag{Name: "core-port", Value: 8080, EnvVar: "MATRIX_CORE_PORT"},
			cli.IntFlag{Name: "emulator-port", Value: 3000, EnvVar: "MATRIX_EMULATOR_PORT"},
			cli.IntFlag{Name: "gamepad-port", Value: 4000, EnvVar: "MATRIX_GAMEPAD_PORT"},
			cli.StringFlag{
				Name:  "log-level",
				Value: "warning",
				Usage: "[panic fatal error warning info debug]",
			},
		},
		ArgsUsage: "[core emulator gamepad device zigzag yumyum demo clock draw blocks]",
		Action:    startAction,
	}}

	if err := app.Run(os.Args); err != nil {
		logrus.Errorf("%+v", err)
	}
}

type component func() error

func (c component) run(cancel func()) {
	if err := c(); err != nil {
		logrus.Errorf("%+v", err)
		cancel()
	}
}

func startAction(c *cli.Context) error {
	level, err := logrus.ParseLevel(c.String("log-level"))
	if err != nil {
		return errors.Wrap(err, "Invalid given log level")
	}
	logrus.SetLevel(level)

	args := c.Args()

	if len(args) < 1 {
		return errors.New("Missing component name")
	}

	var cs []component
	for _, arg := range args {
		switch arg {
		case "core":
			cs = append(cs, component(func() error { return core.Start(c.Int("core-port")) }))
		case "emulator":
			cs = append(cs, component(func() error { return emulator.Start(c.Int("emulator-port"), c.String("core-uri")) }))
		case "gamepad":
			cs = append(cs, component(func() error { return gamepad.Start(c.Int("gamepad-port"), c.String("core-uri")) }))
		case "device":
			cs = append(cs, component(func() error { return device.Start(c.String("core-uri")) }))
		case "zigzag":
			cs = append(cs, component(func() error { return zigzag.Start(c.String("core-uri")) }))
		case "yumyum":
			cs = append(cs, component(func() error { return yumyum.Start(c.String("core-uri")) }))
		case "demo":
			cs = append(cs, component(func() error { return demo.Start(c.String("core-uri")) }))
		case "clock":
			cs = append(cs, component(func() error { return clock.Start(c.String("core-uri")) }))
		case "draw":
			cs = append(cs, component(func() error { return draw.Start(c.String("core-uri")) }))
		case "blocks":
			cs = append(cs, component(func() error { return blocks.Start(c.String("core-uri")) }))
		case "getout":
			cs = append(cs, component(func() error { return getout.Start(c.String("core-uri")) }))
		default:
			return errors.New("Invalid given component name")
		}
	}

	if len(cs) == 1 {
		return cs[0]()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, c := range cs {
		go c.run(cancel)
	}

	<-ctx.Done()

	return nil
}
