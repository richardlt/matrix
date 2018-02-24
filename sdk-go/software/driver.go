package software

import common "github.com/richardlt/matrix/sdk-go/common"

type driver interface {
	End()
	Step(total, current uint64)
}

type RandomDriver struct {
	ctx          *ctx
	uuid         string
	endCallback  func()
	stepCallback func(total, current uint64)
}

func (r *RandomDriver) OnEnd(f func()) { r.endCallback = f }

func (r *RandomDriver) End() {
	if r.endCallback != nil {
		r.endCallback()
	}
}

func (r *RandomDriver) OnStep(f func(total, current uint64)) { r.stepCallback = f }

func (r *RandomDriver) Step(total, current uint64) {
	if r.stepCallback != nil {
		r.stepCallback(total, current)
	}
}

func (r *RandomDriver) Render() error {
	return r.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_DRIVER,
		DriverData: &ConnectRequest_DriverData{
			Action: ConnectRequest_DriverData_RENDER,
			UUID:   r.uuid,
		},
	})
}

type CaracterDriver struct {
	ctx            *ctx
	uuid           string
	endCallback    func()
	renderCallback func() error
	stepCallback   func(total, current uint64)
}

func (c *CaracterDriver) OnEnd(f func()) { c.endCallback = f }

func (c *CaracterDriver) End() {
	if c.endCallback != nil {
		c.endCallback()
	}
}

func (c *CaracterDriver) OnStep(f func(total, current uint64)) { c.stepCallback = f }

func (c *CaracterDriver) Step(total, current uint64) {
	if c.stepCallback != nil {
		c.stepCallback(total, current)
	}
}

func (c *CaracterDriver) Render(caracter rune, coord common.Coord, color, background common.Color) error {
	return c.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_DRIVER,
		DriverData: &ConnectRequest_DriverData{
			Action:     ConnectRequest_DriverData_RENDER,
			UUID:       c.uuid,
			Caracter:   string(caracter),
			Coord:      &coord,
			Color:      &color,
			Background: &background,
		},
	})
}

type TextDriver struct {
	ctx            *ctx
	uuid           string
	endCallback    func()
	renderCallback func() error
	stepCallback   func(total, current uint64)
}

func (t *TextDriver) OnEnd(f func()) { t.endCallback = f }

func (t *TextDriver) End() {
	if t.endCallback != nil {
		t.endCallback()
	}
}

func (t *TextDriver) OnStep(f func(total, current uint64)) { t.stepCallback = f }

func (t *TextDriver) Step(total, current uint64) {
	if t.stepCallback != nil {
		t.stepCallback(total, current)
	}
}

func (t *TextDriver) Render(text string, coord common.Coord,
	color, background common.Color, repeat bool) error {
	return t.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_DRIVER,
		DriverData: &ConnectRequest_DriverData{
			Action:     ConnectRequest_DriverData_RENDER,
			UUID:       t.uuid,
			Text:       text,
			Coord:      &coord,
			Color:      &color,
			Background: &background,
			Repeat:     repeat,
		},
	})
}

func (t *TextDriver) Stop() error {
	return t.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_DRIVER,
		DriverData: &ConnectRequest_DriverData{
			Action: ConnectRequest_DriverData_STOP,
			UUID:   t.uuid,
		},
	})
}

type ImageDriver struct {
	ctx            *ctx
	uuid           string
	endCallback    func()
	renderCallback func() error
	stepCallback   func(total, current uint64)
}

func (i *ImageDriver) OnEnd(f func()) { i.endCallback = f }

func (i *ImageDriver) End() {
	if i.endCallback != nil {
		i.endCallback()
	}
}

func (i *ImageDriver) OnStep(f func(total, current uint64)) { i.stepCallback = f }

func (i *ImageDriver) Step(total, current uint64) {
	if i.stepCallback != nil {
		i.stepCallback(total, current)
	}
}

func (c *ImageDriver) Render(image Image, coord common.Coord) error {
	return c.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_DRIVER,
		DriverData: &ConnectRequest_DriverData{
			Action: ConnectRequest_DriverData_RENDER,
			UUID:   c.uuid,
			Image:  &image,
			Coord:  &coord,
		},
	})
}
