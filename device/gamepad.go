package device

import (
	"context"
	"fmt"
	"sort"

	"github.com/richardlt/matrix/sdk-go/common"

	"github.com/google/gousb"
	"github.com/pkg/errors"
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
	HID    *gousb.Device
	Conf   config
	States map[string]bool
	Index  int
	Key    string
}

type action struct {
	Slot    int
	Command common.Command
}

func (g *gamepad) OpenDevices(ctx context.Context) error {
	uctx := gousb.NewContext()

	cAction := make(chan action)
	defer close(cAction)

	vid, pid := gousb.ID(0x0079), gousb.ID(0x0011)
	devs, err := uctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == vid && desc.Product == pid
	})
	if err != nil {
		return errors.WithStack(err)
	}

	logrus.Debugf("Found %d controller", len(devs))

	for i := 0; i < len(devs); i++ {
		key := fmt.Sprintf("%d_%d", devs[i].Desc.Vendor, devs[i].Desc.Product)
		go g.listenDevice(cAction, &device{
			HID:    devs[i],
			Conf:   g.configs[key],
			States: map[string]bool{},
			Index:  i,
			Key:    key,
		})
	}

	for {
		select {
		case <-ctx.Done():
			uctx.Close()
			return nil
		case a := <-cAction:
			if g.api != nil {
				g.api.Command(uint64(a.Slot), a.Command)
			}
		}
	}
}

func (g *gamepad) listenDevice(cAction chan action, dev *device) {
	defer dev.HID.Close()
	defer logrus.Debugf("Stop listening from %s controller", dev.Key)

	if err := dev.HID.SetAutoDetach(true); err != nil {
		logrus.Errorf("%+v", errors.WithStack(err))
		return
	}

	cfg, err := dev.HID.Config(1)
	if err != nil {
		logrus.Errorf("%+v", errors.WithStack(err))
		return
	}
	defer cfg.Close()

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		logrus.Errorf("%+v", errors.WithStack(err))
		return
	}
	defer intf.Close()

	epIn, err := intf.InEndpoint(1)
	if err != nil {
		logrus.Errorf("%+v", errors.WithStack(err))
		return
	}

	logrus.Debugf("Listening from %s controller", dev.Key)

	handler := g.handleData(dev)

	buf := make([]byte, epIn.Desc.MaxPacketSize)
	for {
		if _, err := epIn.Read(buf); err != nil {
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
				a = &action{Slot: dev.Index, Command: commandFromString(b.Name + ":press")}
			} else if !isPressed && currentState {
				a = &action{Slot: dev.Index, Command: commandFromString(b.Name + ":release")}
			}

			dev.States[b.Name] = isPressed

			if a != nil {
				return a
			}
		}

		return nil
	}
}
