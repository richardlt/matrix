package device

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
	"github.com/sirupsen/logrus"
	serial "go.bug.st/serial.v1"
)

const refreshDelay = time.Millisecond * 65 // ~15hz
const defaultBrightness = 204

func newMatrix() *matrix {
	return &matrix{
		brightness: defaultBrightness,
		buffer:     []byte{defaultBrightness},
	}
}

type matrix struct {
	api         software.API
	frame       common.Frame
	brightness  uint8
	buffer      []byte
	layer       software.Layer
	imageDriver *software.ImageDriver
}

func (m *matrix) FramesReceived(fs []*common.Frame) {
	if len(fs) > 0 {
		m.frame = *fs[0]
		m.updateBuffer()
	}
}

func (m *matrix) Init(a software.API) (err error) {
	m.api = a

	i := a.GetImageFromLocal("device")

	a.SetConfig(software.ConnectRequest_SoftwareData_Config{
		Logo:           &i,
		MinPlayerCount: 1,
		MaxPlayerCount: 1,
	})

	m.layer, err = a.NewLayer()
	if err != nil {
		return err
	}

	m.imageDriver, err = m.layer.NewImageDriver()
	if err != nil {
		return err
	}

	a.Ready()
	return nil
}

func (m *matrix) Start(uint64) { m.print() }

func (m *matrix) print() {
	m.imageDriver.Render(m.api.GetImageFromLocal("arrow-left"), common.Coord{X: 2, Y: 4})
	m.imageDriver.Render(m.api.GetImageFromLocal("arrow-right"), common.Coord{X: 13, Y: 4})

	grey := m.api.GetColorFromLocalThemeByName("flat", "grey_2")
	c := m.api.GetColorFromLocalThemeByName("flat", "yellow_2")

	if m.brightness < 51 {
		c = grey
	}
	m.layer.SetWithCoord(common.Coord{X: 5, Y: 6}, c)

	if m.brightness < 102 {
		c = grey
	}
	m.layer.SetWithCoord(common.Coord{X: 6, Y: 6}, c)
	m.layer.SetWithCoord(common.Coord{X: 6, Y: 5}, c)

	if m.brightness < 153 {
		c = grey
	}
	m.layer.SetWithCoord(common.Coord{X: 7, Y: 6}, c)
	m.layer.SetWithCoord(common.Coord{X: 7, Y: 5}, c)
	m.layer.SetWithCoord(common.Coord{X: 7, Y: 4}, c)

	if m.brightness < 204 {
		c = grey
	}
	m.layer.SetWithCoord(common.Coord{X: 8, Y: 6}, c)
	m.layer.SetWithCoord(common.Coord{X: 8, Y: 5}, c)
	m.layer.SetWithCoord(common.Coord{X: 8, Y: 4}, c)
	m.layer.SetWithCoord(common.Coord{X: 8, Y: 3}, c)
	m.layer.SetWithCoord(common.Coord{X: 9, Y: 6}, c)
	m.layer.SetWithCoord(common.Coord{X: 9, Y: 5}, c)
	m.layer.SetWithCoord(common.Coord{X: 9, Y: 4}, c)
	m.layer.SetWithCoord(common.Coord{X: 9, Y: 3}, c)

	if m.brightness < 255 {
		c = grey
	}
	m.layer.SetWithCoord(common.Coord{X: 10, Y: 6}, c)
	m.layer.SetWithCoord(common.Coord{X: 10, Y: 5}, c)
	m.layer.SetWithCoord(common.Coord{X: 10, Y: 4}, c)
	m.layer.SetWithCoord(common.Coord{X: 10, Y: 3}, c)
	m.layer.SetWithCoord(common.Coord{X: 10, Y: 2}, c)

	m.api.Print()
}

func (m *matrix) Close() {}

func (m *matrix) ActionReceived(slot int, cmd common.Command) {
	switch cmd {
	case common.Command_LEFT_UP:
		if m.brightness > 0 {
			m.brightness -= 51
			m.updateBuffer()
			m.print()
		}
	case common.Command_RIGHT_UP:
		if m.brightness < 255 {
			m.brightness += 51
			m.updateBuffer()
			m.print()
		}
	}
}

func (m *matrix) OpenPorts(ctx context.Context) error {
	matches, err := filepath.Glob("/dev/tty*")
	if err != nil {
		return errors.WithStack(err)
	}

	paths := []string{}
	for _, ma := range matches {
		if strings.Contains(strings.ToLower(ma), "usb") {
			paths = append(paths, ma)
		}
	}

	for _, path := range paths {
		logrus.Debugf("Try to open port at %s", path)

		port, err := serial.Open(path, &serial.Mode{BaudRate: 115200})
		if err != nil {
			return errors.WithStack(err)
		}

		logrus.Debugf("Port opened at %s", path)

		if err := port.ResetInputBuffer(); err != nil {
			return errors.WithStack(err)
		}
		if err := port.ResetOutputBuffer(); err != nil {
			return errors.WithStack(err)
		}

		// read the matrix size
		buf := make([]byte, 1)
		_, err = port.Read(buf)
		if err != nil {
			return errors.WithStack(err)
		}

		size := int(buf[0])

		logrus.Debugf("Receive %d size for port at %s", size, path)

		go func(path string) {
			t := time.NewTicker(refreshDelay)
			defer t.Stop()

			var lastBuffer []byte
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					if !bytes.Equal(lastBuffer, m.buffer) {
						buffer := make([]byte, size*3+1)

						for i := 0; i < len(m.buffer) && i < len(buffer); i++ {
							buffer[i] = m.buffer[i]
						}

						if _, err := port.Write(buffer); err != nil {
							logrus.Errorf("%+v", errors.WithStack(err))
							return
						}

						lastBuffer = m.buffer
					}
				}
			}
		}(path)
	}

	return nil
}

func (m *matrix) updateBuffer() {
	buffer := []byte{m.brightness}

	for _, p := range m.frame.Pixels {
		if p.A > 0 {
			buffer = append(buffer, byte(p.R), byte(p.G), byte(p.B))
		} else {
			buffer = append(buffer, 0, 0, 0)
		}
	}

	m.buffer = buffer
}
