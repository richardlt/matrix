package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/clock"
	"github.com/richardlt/matrix/core"
	"github.com/richardlt/matrix/demo"
	"github.com/richardlt/matrix/device"
	"github.com/richardlt/matrix/draw"
	"github.com/richardlt/matrix/emulator"
	"github.com/richardlt/matrix/gamepad"
	"github.com/richardlt/matrix/yumyum"
	"github.com/richardlt/matrix/zigzag"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	logrus.SetLevel(logrus.DebugLevel)

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "core-uri", Value: "localhost:8080",
			EnvVar: "MATRIX_CORE_URI"},
	}

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
		Action: func(c *cli.Context) error { return device.Start(c.Parent().String("core-uri")) },
	}, {
		Name:   "zigzag",
		Usage:  "start zigzag game",
		Action: func(c *cli.Context) error { return zigzag.Start(c.Parent().String("core-uri")) },
	}, {
		Name:   "yumyum",
		Usage:  "start yumyum game",
		Action: func(c *cli.Context) error { return yumyum.Start(c.Parent().String("core-uri")) },
	}, {
		Name:   "demo",
		Usage:  "start demo software",
		Action: func(c *cli.Context) error { return demo.Start(c.Parent().String("core-uri")) },
	}, {
		Name:   "clock",
		Usage:  "start clock software",
		Action: func(c *cli.Context) error { return clock.Start(c.Parent().String("core-uri")) },
	}, {
		Name:   "draw",
		Usage:  "start draw software",
		Action: func(c *cli.Context) error { return draw.Start(c.Parent().String("core-uri")) },
	}, {
		Name:  "all",
		Usage: "start all",
		Action: func(c *cli.Context) error {
			go func() {
				if err := device.Start(c.Parent().String("core-uri")); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			go func() {
				if err := gamepad.Start(4000, 8080); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			go func() {
				if err := emulator.Start(3000, 8080); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			go func() {
				if err := zigzag.Start(c.Parent().String("core-uri")); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			go func() {
				if err := yumyum.Start(c.Parent().String("core-uri")); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			go func() {
				if err := demo.Start(c.Parent().String("core-uri")); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			go func() {
				if err := clock.Start(c.Parent().String("core-uri")); err != nil {
					logrus.Errorf("%+v", err)
				}
			}()
			go func() {
				if err := draw.Start(c.Parent().String("core-uri")); err != nil {
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
