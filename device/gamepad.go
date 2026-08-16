package device

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/karalabe/hid"
	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/internal/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/player"
)

func newGamepad(a *coreState) *gamepad {
	m := map[string]config{}

	for _, c := range configs {
		// sort buttons by values to detect multi press
		sort.Slice(c.Buttons, func(i, j int) bool {
			return c.Buttons[j].Value < c.Buttons[i].Value
		})
		m[fmt.Sprintf("%d_%d", c.VendorID, c.ProductID)] = c
	}

	return &gamepad{configs: m, state: a}
}

type gamepad struct {
	state   *coreState
	api     *player.API
	configs map[string]config
}

func (g *gamepad) Init(api *player.API) error {
	g.api = api
	return nil
}

type device struct {
	HID    hid.DeviceInfo
	Conf   config
	States map[string]bool
	Slot   int
}

type action struct {
	Slot    int
	Command common.Command
}

// Controller discovery polls, and each poll is more expensive than it looks: hidapi opens
// every matching device to read its string descriptors, which is a run of USB control
// transfers. On a board where the panel's serial adapter shares the bus, doing that every
// second is enough to stutter the display, so the interval backs off once the set of
// controllers has settled and snaps back the moment it changes.
const (
	devicePollMin = time.Second
	devicePollMax = 8 * time.Second
)

func (g *gamepad) OpenDevices(ctx context.Context) error {
	cAction := make(chan action)
	defer close(cAction)

	vid, pid := uint16(0x0079), uint16(0x0011)

	go func() {
		connected := map[string]*device{}
		mutex := new(sync.Mutex)
		poll := devicePollMin

		for {
			if ctx.Err() != nil {
				return
			}

			mutex.Lock()

			mFreeSlots := map[int]bool{0: true, 1: true, 2: true, 3: true}
			for _, d := range connected {
				mFreeSlots[d.Slot] = false
			}
			freeSlots := []int{}
			for i := 0; i < len(mFreeSlots); i++ {
				if mFreeSlots[i] {
					freeSlots = append(freeSlots, i)
				}
			}

			// With every slot taken there is nothing to discover, so the bus is left alone
			// entirely until a controller drops off. Nor is it worth looking while a
			// software is drawing: the poll competes with the frames for the same bus, and
			// a player joins from the menu, where the picture is still.
			switch {
			case len(freeSlots) == 0:
				poll = devicePollMax
			case !g.state.idle():
				poll = devicePollMin
			default:
				var found bool
				devs := hid.Enumerate(vid, pid)
				for i := 0; i < len(devs) && i < len(freeSlots); i++ {
					if _, ok := connected[devs[i].Path]; !ok {
						found = true
						key := fmt.Sprintf("%d_%d", devs[i].VendorID, devs[i].ProductID)
						d := &device{
							HID:    devs[i],
							Conf:   g.configs[key],
							States: map[string]bool{},
							Slot:   freeSlots[i],
						}
						connected[d.HID.Path] = d
						logrus.Infof("Controller %s connected on slot %d", d.HID.Path, d.Slot)
						go func(d *device) {
							g.listenDevice(cAction, d)
							mutex.Lock()
							delete(connected, d.HID.Path)
							mutex.Unlock()
							logrus.Infof("Controller %s at slot %d disconnected", d.HID.Path, d.Slot)
						}(d)
					}
				}

				// A change means more may be arriving, so look again promptly; otherwise
				// ease off towards the maximum.
				if found {
					poll = devicePollMin
				} else if poll < devicePollMax {
					poll *= 2
					if poll > devicePollMax {
						poll = devicePollMax
					}
				}
			}

			mutex.Unlock()
			time.Sleep(poll)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case a := <-cAction:
			if g.api != nil {
				if err := g.api.Command(uint64(a.Slot), a.Command); err != nil {
					logrus.Errorf("%+v", err)
				}
			}
		}
	}
}

func (g *gamepad) listenDevice(cAction chan action, dev *device) {
	defer logrus.Debugf("Stop listening from %s controller", dev.HID.Path)

	d, err := dev.HID.Open()
	if err != nil {
		logrus.Errorf("%+v", errors.Errorf("opening controller at %s: %w", dev.HID.Path, err))
		return
	}
	defer func() {
		if err := d.Close(); err != nil {
			logrus.Errorf("%+v", errors.Errorf("closing controller at %s: %w", dev.HID.Path, err))
		}
	}()

	logrus.Debugf("Listening from %s controller", dev.HID.Path)

	handler := g.handleData(dev)

	buf := make([]byte, 7)
	for {
		if _, err := d.Read(buf); err != nil {
			logrus.Errorf("%+v", errors.Errorf("reading from controller at %s: %w", dev.HID.Path, err))
			return
		}
		if a := handler(buf); a != nil {
			cAction <- *a
		}
	}
}

func (g *gamepad) handleData(dev *device) func([]byte) *action {
	mapPins := map[int]pin{}
	for _, p := range dev.Conf.Pins {
		mapPins[p.Number] = p
	}

	return func(d []byte) *action {
		data := make([]int, len(d))
		for i, b := range d {
			data[i] = int(b)
		}

		// substract button value from pin value to detect button state
		// for multitouch pins or compare value for monotouch pins
		for _, b := range dev.Conf.Buttons {
			var isPressed bool
			if mapPins[b.Pin].Multi {
				isPressed = data[b.Pin]-b.Value >= 0
			} else {
				isPressed = data[b.Pin] == b.Value
			}

			if isPressed && mapPins[b.Pin].Multi {
				data[b.Pin] -= b.Value
			}

			currentState := dev.States[b.Name]

			var a *action
			if isPressed && !currentState {
				a = &action{Slot: dev.Slot, Command: commandFromString(b.Name + ":press")}
			} else if !isPressed && currentState {
				a = &action{Slot: dev.Slot, Command: commandFromString(b.Name + ":release")}
			}

			dev.States[b.Name] = isPressed

			if a != nil {
				return a
			}
		}

		return nil
	}
}
