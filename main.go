package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/richardlt/matrix/animate"
	"github.com/richardlt/matrix/blocks"
	"github.com/richardlt/matrix/clock"
	"github.com/richardlt/matrix/core"
	"github.com/richardlt/matrix/demo"
	"github.com/richardlt/matrix/device"
	"github.com/richardlt/matrix/draw"
	"github.com/richardlt/matrix/emulator"
	"github.com/richardlt/matrix/gamepad"
	"github.com/richardlt/matrix/getout"
	"github.com/richardlt/matrix/internal/errors"
	"github.com/richardlt/matrix/light"
	"github.com/richardlt/matrix/rollupdice"
	"github.com/richardlt/matrix/yumyum"
	"github.com/richardlt/matrix/zigzag"
)

// version is stamped in at build time from the git tag, by the Makefile. A binary built
// with a bare `go build` carries the default instead.
var version = "dev"

func main() {
	cmd := &cli.Command{
		Name:    "matrix",
		Usage:   "video game console operating system for a 16x9 RGB LED matrix",
		Version: version,
		Commands: []*cli.Command{{
			Name:  "start",
			Usage: "start the matrix components",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "core-uri",
					Value:   "localhost:8080",
					Sources: cli.EnvVars("MATRIX_CORE_URI"),
					Usage:   "Core URI is used by softwares, players and displays.",
				},
				&cli.IntFlag{Name: "core-port", Value: 8080, Sources: cli.EnvVars("MATRIX_CORE_PORT")},
				&cli.IntFlag{Name: "emulator-port", Value: 3000, Sources: cli.EnvVars("MATRIX_EMULATOR_PORT")},
				&cli.IntFlag{Name: "gamepad-port", Value: 4000, Sources: cli.EnvVars("MATRIX_GAMEPAD_PORT")},
				&cli.StringFlag{
					Name:  "log-level",
					Value: "warning",
					Usage: "[panic fatal error warning info debug]",
				},
			},
			ArgsUsage: "[core emulator gamepad device zigzag yumyum demo clock draw blocks getout animate light rollupdice]",
			Action:    startAction,
		}},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
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

func startAction(ctx context.Context, cmd *cli.Command) error {
	level, err := logrus.ParseLevel(cmd.String("log-level"))
	if err != nil {
		return errors.Errorf("invalid given log level: %w", err)
	}
	logrus.SetLevel(level)

	args := cmd.Args()

	if args.Len() < 1 {
		return errors.Errorf("missing component name, expected at least one of %s", cmd.ArgsUsage)
	}

	var cs []component
	for _, arg := range args.Slice() {
		switch arg {
		case "core":
			cs = append(cs, component(func() error { return core.Start(cmd.Int("core-port")) }))
		case "emulator":
			cs = append(cs, component(func() error { return emulator.Start(cmd.Int("emulator-port"), cmd.String("core-uri")) }))
		case "gamepad":
			cs = append(cs, component(func() error { return gamepad.Start(cmd.Int("gamepad-port"), cmd.String("core-uri")) }))
		case "device":
			cs = append(cs, component(func() error { return device.Start(cmd.String("core-uri")) }))
		case "zigzag":
			cs = append(cs, component(func() error { return zigzag.Start(cmd.String("core-uri")) }))
		case "yumyum":
			cs = append(cs, component(func() error { return yumyum.Start(cmd.String("core-uri")) }))
		case "demo":
			cs = append(cs, component(func() error { return demo.Start(cmd.String("core-uri")) }))
		case "clock":
			cs = append(cs, component(func() error { return clock.Start(cmd.String("core-uri")) }))
		case "draw":
			cs = append(cs, component(func() error { return draw.Start(cmd.String("core-uri")) }))
		case "blocks":
			cs = append(cs, component(func() error { return blocks.Start(cmd.String("core-uri")) }))
		case "getout":
			cs = append(cs, component(func() error { return getout.Start(cmd.String("core-uri")) }))
		case "animate":
			cs = append(cs, component(func() error { return animate.Start(cmd.String("core-uri")) }))
		case "light":
			cs = append(cs, component(func() error { return light.Start(cmd.String("core-uri")) }))
		case "rollupdice":
			cs = append(cs, component(func() error { return rollupdice.Start(cmd.String("core-uri")) }))
		default:
			return errors.Errorf("invalid component name %q, expected one of %s", arg, cmd.ArgsUsage)
		}
	}

	if len(cs) == 1 {
		return cs[0]()
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, c := range cs {
		go c.run(cancel)
	}

	<-ctx.Done()

	return nil
}
