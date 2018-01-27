package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/core"
	"github.com/richardlt/matrix/demo"
	"github.com/richardlt/matrix/device"
	"github.com/richardlt/matrix/emulator"
	"github.com/richardlt/matrix/gamepad"
	"github.com/richardlt/matrix/yumyum"
	"github.com/richardlt/matrix/zigzag"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	logrus.SetLevel(logrus.DebugLevel)

	app.Commands = []cli.Command{{
		Name:   "core",
		Usage:  "start the matrix core",
		Action: func(c *cli.Context) error { return core.Start(8080) },
	}, {
		Name:   "emulator",
		Usage:  "start the matrix device emulator",
		Action: func(c *cli.Context) error { return emulator.Start(3000, 8080) },
	}, {
		Name:   "gamepad",
		Usage:  "start the matrix gamepad",
		Action: func(c *cli.Context) error { return gamepad.Start(4000, 8080) },
	}, {
		Name:   "device",
		Usage:  "start the matrix device",
		Action: func(c *cli.Context) error { return device.Start(8080) },
	}, {
		Name:   "zigzag",
		Usage:  "start zigzag game",
		Action: func(c *cli.Context) error { return zigzag.Start("localhost:8080") },
	}, {
		Name:   "yumyum",
		Usage:  "start yumyum game",
		Action: func(c *cli.Context) error { return yumyum.Start("localhost:8080") },
	}, {
		Name:   "demo",
		Usage:  "start demo game",
		Action: func(c *cli.Context) error { return demo.Start("localhost:8080") },
	}, {
		Name:  "all",
		Usage: "start all",
		Action: func(c *cli.Context) error {
			go func() {
				if err := gamepad.Start(4000, 8080); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			go func() {
				if err := zigzag.Start("localhost:8080"); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			go func() {
				if err := yumyum.Start("localhost:8080"); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			go func() {
				if err := demo.Start("localhost:8080"); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			return core.Start(8080)
		},
	}}

	if err := app.Run(os.Args); err != nil {
		logrus.Errorf("%+v", err)
	}
}
