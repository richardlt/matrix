package animate

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ovh/cds/sdk"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return sdk.WithStack(err)
			}

			var h header
			if err := json.Unmarshal(buf, &h); err != nil {
				return errors.WithStack(err)
			}

			a.headers = append(a.headers, h)
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return software.Connect(uri, a, true)
}

type header struct {
	Name        string `json:"name"`
	Orientation string `json:"orientation"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

type animation []byte

func (a animation) readFrame(width, height, index int) software.Image {
	pixels := 3 * width * height
	start := index * pixels
	end := start + pixels
	buf := a[start:end]

	var colors []*common.Color
	mapColors := map[string]uint64{}
	mask := make([]uint64, width*height)
	var cursor int
	for i := range mask {
		c := common.Color{
			R: uint64(buf[cursor]),
			G: uint64(buf[cursor+1]),
			B: uint64(buf[cursor+2]),
			A: 1,
		}
		key := fmt.Sprintf("%d%d%d", c.R, c.G, c.B)
		if v, ok := mapColors[key]; !ok {
			colors = append(colors, &c)
			mask[i] = uint64(len(colors) - 1)
			mapColors[key] = mask[i]
		} else {
			mask[i] = v
		}
		cursor += 3
	}

	return software.Image{
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

	api.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &i,
		MinPlayerCount: 1,
		MaxPlayerCount: 1,
	})

	a.layer, err = api.NewLayer()
	if err != nil {
		return err
	}

	a.imageDriver, err = a.layer.NewImageDriver()
	if err != nil {
		return err
	}
	a.imageDriver.OnEnd(func() { a.api.Print() })

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

	a.layer.Clean()
	a.api.Print()

	if len(a.headers) == 0 {
		return
	}

	var anim animation
	var err error
	anim, err = ioutil.ReadFile(fmt.Sprintf("./animations/%s", a.headers[a.index].Name))
	if err != nil {
		logrus.Error(errors.WithStack(err))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()
	maxIndex := len(anim) / (a.headers[a.index].Width * a.headers[a.index].Height * 3)
	index := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.imageDriver.Render(
				anim.readFrame(a.headers[a.index].Width, a.headers[a.index].Height, index),
				common.Coord{
					X: int64(a.headers[a.index].Width / 2),
					Y: int64(a.headers[a.index].Height / 2),
				},
			)
			if index+1 == maxIndex {
				index = 0
			} else {
				index++
			}
		}
	}
}
