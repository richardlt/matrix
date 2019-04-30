package device

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	serial "go.bug.st/serial.v1"

	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

const refreshDelay = time.Millisecond * 30
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

func (m *matrix) ActionReceived(slot uint64, cmd common.Command) {
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
	connected := map[string]struct{}{}
	invalid := map[string]struct{}{}
	mutex := new(sync.Mutex)

	go func() {
		for {
			if ctx.Err() != nil {
				return
			}

			mutex.Lock()
			defered := func() {
				mutex.Unlock()
				time.Sleep(time.Second)
			}

			paths, err := serial.GetPortsList()
			if err != nil {
				logrus.Errorf("%+v", errors.WithStack(err))
				defered()
				continue
			}

			newInvalid := map[string]struct{}{}

			for _, path := range paths {
				if !(strings.Contains(strings.ToLower(path), "usb") || strings.Contains(path, "COM")) {
					continue
				}
				if _, ok := invalid[path]; ok {
					newInvalid[path] = struct{}{}
					continue
				}

				if _, ok := connected[path]; !ok {
					logrus.Debugf("Try to open port at %s", path)

					port, err := serial.Open(path, &serial.Mode{BaudRate: 115200})
					if err != nil {
						if err.Error() == "Serial port busy" {
							logrus.Debugf("Port at %s is not available", path)
							newInvalid[path] = struct{}{}
						} else {
							logrus.Errorf("%+v", errors.WithStack(err))
						}
						continue
					}

					logrus.Debugf("Port opened at %s", path)

					_ = port.ResetInputBuffer()  // ignore error, always occured on darwin
					_ = port.ResetOutputBuffer() // ignore error, always occured on darwin

					logrus.Debugf("Search for matrix signature and size at %s", path)
					buf := make([]byte, 2)

					ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
					go func() {
						defer cancel()
						_, _ = port.Read(buf) // ignore error, always occured on darwin
					}()

					<-ctxTimeout.Done()
					if ctxTimeout.Err() != nil && ctxTimeout.Err().Error() != "context canceled" {
						_ = port.Close()
						logrus.Debugf("Port at %s didn't answer to connect", path)
						newInvalid[path] = struct{}{}
						continue
					}

					if buf[0] != 0x15 {
						logrus.Debugf("Port at %s is not a matrix device", path)
						newInvalid[path] = struct{}{}
						continue
					}

					size := int(buf[1])
					logrus.Debugf("Receive %d size for port at %s", size, path)

					connected[path] = struct{}{}
					logrus.Infof("Serial %s connected", path)

					go func(path string) {
						t := time.NewTicker(refreshDelay)
						defer func() {
							t.Stop()

							mutex.Lock()
							delete(connected, path)
							mutex.Unlock()
							logrus.Infof("Serial %s disconnected", path)
						}()

						var lastBuffer []byte
						for {
							select {
							case <-ctx.Done():
								return
							case <-t.C:
								if !bytes.Equal(lastBuffer, m.buffer) {
									lastBuffer = m.buffer

									buffer := make([]byte, size*3+1)
									for i := 0; i < len(m.buffer) && i < len(buffer); i++ {
										buffer[i] = m.buffer[i]
									}

									if _, err := port.Write(buffer); err != nil {
										logrus.Errorf("%+v", errors.WithStack(err))
										return
									}

									// read the ack
									ack := make([]byte, 1)
									_, err = port.Read(ack)
									if err != nil {
										logrus.Errorf("%+v", errors.WithStack(err))
										return
									}
								}
							}
						}
					}(path)
				}
			}

			invalid = newInvalid

			defered()
		}
	}()

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
