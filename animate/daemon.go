package animate

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/internal/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// Start the animate software.
func Start(uri string) error {
	logrus.Infof("Start animate for uri %s\n", uri)

	a := &animate{}

	// list animation headers and set index
	if err := filepath.Walk("./animations", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".json") {
			buf, err := os.ReadFile(path)
			if err != nil {
				return errors.Errorf("reading animation %s: %w", path, err)
			}

			var h header
			if err := json.Unmarshal(buf, &h); err != nil {
				return errors.Errorf("unmarshaling animation %s: %w", path, err)
			}

			a.headers = append(a.headers, h)
		}
		return nil
	}); err != nil {
		return errors.Errorf("walking the animations directory: %w", err)
	}

	return software.Connect(uri, a, true)
}

type header struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	FPS    int    `json:"fps"`
}

type animation []byte

func (a animation) readFrame(width, height, index int) *software.Image {
	pixels := 3 * width * height
	start := index * pixels
	end := start + pixels
	buf := a[start:end]

	// The palette key packs the channels into one integer. Formatting it as a string
	// instead costs a reflection-based format and an allocation for every pixel of every
	// frame, which on this hardware is the bulk of the work of decoding one.
	// Frames tend to reuse a handful of colours, so the palette is sized for that rather
	// than grown from empty a pixel at a time.
	colors := make([]*common.Color, 0, 32)
	mapColors := make(map[uint64]uint64, 32)
	mask := make([]uint64, width*height)
	var cursor int
	for i := range mask {
		r, g, b := uint64(buf[cursor]), uint64(buf[cursor+1]), uint64(buf[cursor+2])
		cursor += 3

		key := r<<16 | g<<8 | b
		if v, ok := mapColors[key]; ok {
			mask[i] = v
			continue
		}

		// Only a colour the frame has not used yet needs its own value.
		colors = append(colors, &common.Color{R: r, G: g, B: b, A: 1})
		mask[i] = uint64(len(colors) - 1)
		mapColors[key] = mask[i]
	}

	return &software.Image{
		Width:  uint64(width),
		Height: uint64(height),
		Colors: colors,
		Mask:   mask,
	}
}

type animate struct {
	api         software.API
	layer       software.Layer
	imageDriver *software.ImageDriver
	cancel      context.CancelFunc
	headers     []header
	index       int
}

func (a *animate) Init(api software.API) (err error) {
	logrus.Debug("Init animate")

	a.api = api

	i := api.GetImageFromLocal("animate")

	if err := api.SetConfig(&software.ConnectRequest_SoftwareData_Config{
		Logo:           i,
		MinPlayerCount: 1,
		MaxPlayerCount: 1,
	}); err != nil {
		return err
	}

	a.layer, err = api.NewLayer()
	if err != nil {
		return err
	}

	a.imageDriver, err = a.layer.NewImageDriver()
	if err != nil {
		return err
	}
	a.imageDriver.OnEnd(func() { _ = a.api.Print() })

	return api.Ready()
}

func (a *animate) Start(playerCount uint64) { a.play() }

func (a *animate) Close() { a.reset() }

func (a *animate) ActionReceived(slot uint64, cmd common.Command) {
	switch cmd {
	case common.Command_LEFT_UP:
		if a.index < 1 {
			a.index = len(a.headers) - 1
		} else {
			a.index--
		}
		a.play()
	case common.Command_RIGHT_UP:
		if a.index+1 == len(a.headers) {
			a.index = 0
		} else {
			a.index++
		}
		a.play()
	}
}

func (a *animate) reset() {
	if a.cancel != nil {
		a.cancel()
	}
}

func (a *animate) play() {
	a.reset()

	_ = a.layer.Clean()
	_ = a.api.Print()

	if len(a.headers) == 0 {
		return
	}

	var anim animation
	var err error
	anim, err = os.ReadFile(fmt.Sprintf("./animations/%s", a.headers[a.index].Name))
	if err != nil {
		logrus.Errorf("%+v", errors.Errorf("reading animation %s: %w", a.headers[a.index].Name, err))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	h := a.headers[a.index]

	ticker := time.NewTicker(time.Second / time.Duration(h.FPS))
	defer ticker.Stop()
	maxIndex := len(anim) / (h.Width * h.Height * 3)
	index := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = a.imageDriver.Render(
				anim.readFrame(h.Width, h.Height, index),
				&common.Coord{X: 8, Y: 4}, // middle of the screen
			)
			if index+1 == maxIndex {
				index = 0
			} else {
				index++
			}
		}
	}
}
