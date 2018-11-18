package device

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/karalabe/hid"
	"github.com/pkg/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/player"
	"github.com/sirupsen/logrus"
)

func newGamepad() *gamepad {
	m := map[string]config{}

	for _, c := range configs {
		// sort buttons by values to detect multi press
		sort.Slice(c.Buttons, func(i, j int) bool {
			return c.Buttons[j].Value < c.Buttons[i].Value
		})
		m[fmt.Sprintf("%d_%d", c.VendorID, c.ProductID)] = c
	}

	return &gamepad{configs: m}
}

type gamepad struct {
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

func (g *gamepad) OpenDevices(ctx context.Context) error {
	cAction := make(chan action)
	defer close(cAction)

	vid, pid := uint16(0x0079), uint16(0x0011)

	go func() {
		connected := map[string]*device{}
		mutex := new(sync.Mutex)

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

			if len(freeSlots) > 0 {
				devs := hid.Enumerate(vid, pid)
				for i := 0; i < len(devs) && i < len(freeSlots); i++ {
					if _, ok := connected[devs[i].Path]; !ok {
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
			}

			mutex.Unlock()
			time.Sleep(time.Second)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case a := <-cAction:
			if g.api != nil {
				g.api.Command(uint64(a.Slot), a.Command)
			}
		}
	}
}

func (g *gamepad) listenDevice(cAction chan action, dev *device) {
	defer logrus.Debugf("Stop listening from %s controller", dev.HID.Path)

	d, err := dev.HID.Open()
	if err != nil {
		logrus.Errorf("%+v", errors.WithStack(err))
		return
	}
	defer func() {
		if err := d.Close(); err != nil {
			logrus.Errorf("%+v", errors.WithStack(err))
		}
	}()

	logrus.Debugf("Listening from %s controller", dev.HID.Path)

	handler := g.handleData(dev)

	buf := make([]byte, 7)
	for {
		if _, err := d.Read(buf); err != nil {
			logrus.Errorf("%+v", errors.WithStack(err))
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
