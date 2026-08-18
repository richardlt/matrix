package device

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	serial "go.bug.st/serial.v1"

	"github.com/richardlt/matrix/internal/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// serialBaudRate must match Serial.begin() in device/firmware/firmware.ino. A mismatch
// is not silent: the signature handshake fails and the port is simply never recognised
// as a matrix, so no display appears.
//
// A frame is NUMPIXELS*3+1 = 433 bytes, so at 8N1 the wire time is 4330/baud seconds:
// 37.6ms at 115200, 8.7ms at 500000. The panel's own show() adds a further 4.3ms with
// interrupts disabled, and the write waits for an ack, so the two do not overlap.
const serialBaudRate = 500000

// refreshDelay paces the writes. It has to stay above the per-frame wire time plus
// show(), around 13ms at this baud rate, or the writer simply queues behind the link.
const refreshDelay = time.Millisecond * 15

const defaultBrightness = 204

func newMatrix(a *coreState) *matrix {
	return &matrix{
		state:      a,
		brightness: defaultBrightness,
		buffer:     []byte{defaultBrightness},
	}
}

type matrix struct {
	state       *coreState
	api         software.API
	frame       *common.Frame
	brightness  uint8
	buffer      []byte
	layer       software.Layer
	imageDriver *software.ImageDriver
}

// SoftwareRunning satisfies display.StateAware: the core tells this display when a
// software takes over from the menus, and both pollers in this package hold off while
// one is drawing.
func (m *matrix) SoftwareRunning(running bool) {
	logrus.Debugf("Core reports software running: %t", running)
	m.state.setSoftwareRunning(running)
}

func (m *matrix) FramesReceived(fs []*common.Frame) {
	if len(fs) > 0 {
		m.frame = fs[0]
		m.updateBuffer()
	}
}

func (m *matrix) Init(a software.API) (err error) {
	m.api = a

	i := a.GetImageFromLocal("device")

	if err := a.SetConfig(&software.ConnectRequest_SoftwareData_Config{
		Logo:           i,
		MinPlayerCount: 1,
		MaxPlayerCount: 1,
	}); err != nil {
		return err
	}

	m.layer, err = a.NewLayer()
	if err != nil {
		return err
	}

	m.imageDriver, err = m.layer.NewImageDriver()
	if err != nil {
		return err
	}

	return a.Ready()
}

func (m *matrix) Start(uint64) { m.print() }

func (m *matrix) print() {
	_ = m.imageDriver.Render(m.api.GetImageFromLocal("arrow-left"), &common.Coord{X: 2, Y: 4})
	_ = m.imageDriver.Render(m.api.GetImageFromLocal("arrow-right"), &common.Coord{X: 13, Y: 4})

	grey := m.api.GetColorFromLocalThemeByName("flat", "grey_2")
	c := m.api.GetColorFromLocalThemeByName("flat", "yellow_2")

	if m.brightness < 51 {
		c = grey
	}
	_ = m.layer.SetWithCoord(&common.Coord{X: 5, Y: 6}, c)

	if m.brightness < 102 {
		c = grey
	}
	_ = m.layer.SetWithCoord(&common.Coord{X: 6, Y: 6}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 6, Y: 5}, c)

	if m.brightness < 153 {
		c = grey
	}
	_ = m.layer.SetWithCoord(&common.Coord{X: 7, Y: 6}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 7, Y: 5}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 7, Y: 4}, c)

	if m.brightness < 204 {
		c = grey
	}
	_ = m.layer.SetWithCoord(&common.Coord{X: 8, Y: 6}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 8, Y: 5}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 8, Y: 4}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 8, Y: 3}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 9, Y: 6}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 9, Y: 5}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 9, Y: 4}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 9, Y: 3}, c)

	if m.brightness < 255 {
		c = grey
	}
	_ = m.layer.SetWithCoord(&common.Coord{X: 10, Y: 6}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 10, Y: 5}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 10, Y: 4}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 10, Y: 3}, c)
	_ = m.layer.SetWithCoord(&common.Coord{X: 10, Y: 2}, c)

	_ = m.api.Print()
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

// Port discovery polls, and every poll walks the ports the system knows about. Once a
// panel is attached there is rarely a second one coming, so the interval eases off and
// only returns to the minimum when the set of ports changes.
const (
	portPollMin = time.Second
	portPollMax = 8 * time.Second
)

func (m *matrix) OpenPorts(ctx context.Context) error {
	connected := map[string]struct{}{}
	invalid := map[string]struct{}{}
	mutex := new(sync.Mutex)

	go func() {
		poll := portPollMin

		for {
			if ctx.Err() != nil {
				return
			}

			// Scanning walks every port and opens any it has not seen, which competes with
			// the frames already going out over one of them. Nothing to gain while a
			// software is drawing.
			if !m.state.idle() {
				time.Sleep(portPollMin)
				continue
			}

			mutex.Lock()

			// Under the lock, since a port that drops off is deleted from this map by the
			// goroutine reading it.
			before := len(connected)

			defered := func() {
				// A newly attached panel means another may follow, so look again soon.
				if len(connected) != before {
					poll = portPollMin
				} else if poll < portPollMax {
					poll *= 2
					if poll > portPollMax {
						poll = portPollMax
					}
				}
				mutex.Unlock()
				time.Sleep(poll)
			}

			paths, err := serial.GetPortsList()
			if err != nil {
				logrus.Errorf("%+v", errors.Errorf("listing serial ports: %w", err))
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

					port, err := serial.Open(path, &serial.Mode{BaudRate: serialBaudRate})
					if err != nil {
						if err.Error() == "Serial port busy" {
							logrus.Debugf("Port at %s is not available", path)
							newInvalid[path] = struct{}{}
						} else {
							logrus.Errorf("%+v", errors.Errorf("opening serial port %s: %w", path, err))
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
										logrus.Errorf("%+v", errors.Errorf("writing frame to serial port %s: %w", path, err))
										return
									}

									// read the ack
									ack := make([]byte, 1)
									_, err = port.Read(ack)
									if err != nil {
										logrus.Errorf("%+v", errors.Errorf("reading ack from serial port %s: %w", path, err))
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
	// Sized up front: growing from one byte reallocates about ten times per frame, and at
	// this frame rate that is a meaningful share of what the collector has to chase. It
	// stays a fresh slice each call, because the writer compares it against the previous
	// one to decide whether to send.
	buffer := make([]byte, 0, len(m.frame.Pixels)*3+1)
	buffer = append(buffer, m.brightness)

	for _, p := range m.frame.Pixels {
		if p.A > 0 {
			buffer = append(buffer, byte(p.R), byte(p.G), byte(p.B))
		} else {
			buffer = append(buffer, 0, 0, 0)
		}
	}

	m.buffer = buffer
}
