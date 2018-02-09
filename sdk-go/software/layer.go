package software

import (
	common "github.com/richardlt/matrix/sdk-go/common"
)

type Layer interface {
	Clean() error
	Remove() error
	SetWithCoord(common.Coord, common.Color) error
	NewRandomDriver() (*RandomDriver, error)
	NewCaracterDriver(Font) (*CaracterDriver, error)
	NewTextDriver(Font) (*TextDriver, error)
	NewImageDriver() (*ImageDriver, error)
}

type layer struct {
	ctx          *ctx
	softwareUUID string
	uuid         string
}

// Clean set all layer's pixels to default color.
func (l *layer) Clean() error {
	return l.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_LAYER,
		LayerData: &ConnectRequest_LayerData{
			Action: ConnectRequest_LayerData_CLEAN,
			UUID:   l.uuid,
		},
	})
}

// Remove an existing layer.
func (l *layer) Remove() error {
	return l.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_LAYER,
		LayerData: &ConnectRequest_LayerData{
			Action: ConnectRequest_LayerData_REMOVE,
			UUID:   l.uuid,
		},
	})
}

// SetWithCoord allows to change the color of a given pixel.
func (l *layer) SetWithCoord(coord common.Coord, color common.Color) error {
	return l.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_LAYER,
		LayerData: &ConnectRequest_LayerData{
			Action: ConnectRequest_LayerData_SET_WITH_COORD,
			Coord:  &coord,
			Color:  &color,
			UUID:   l.uuid,
		},
	})
}

func (l *layer) NewRandomDriver() (*RandomDriver, error) {
	res, err := l.ctx.SendCreateRequest(CreateRequest{
		Type: CreateRequest_DRIVER,
		DriverData: &CreateRequest_DriverData{
			Type:         CreateRequest_DriverData_RANDOM,
			LayerUUID:    l.uuid,
			SoftwareUUID: l.softwareUUID,
		},
	})
	if err != nil {
		return nil, err
	}

	rd := &RandomDriver{ctx: l.ctx, uuid: res.UUID}
	l.ctx.AddDriver(rd.uuid, rd)

	return rd, nil
}

func (l *layer) NewCaracterDriver(font Font) (*CaracterDriver, error) {
	res, err := l.ctx.SendCreateRequest(CreateRequest{
		Type: CreateRequest_DRIVER,
		DriverData: &CreateRequest_DriverData{
			Type:         CreateRequest_DriverData_CARACTER,
			LayerUUID:    l.uuid,
			SoftwareUUID: l.softwareUUID,
			Font:         &font,
		},
	})
	if err != nil {
		return nil, err
	}

	cd := &CaracterDriver{ctx: l.ctx, uuid: res.UUID}
	l.ctx.AddDriver(cd.uuid, cd)

	return cd, nil
}

func (l *layer) NewTextDriver(font Font) (*TextDriver, error) {
	res, err := l.ctx.SendCreateRequest(CreateRequest{
		Type: CreateRequest_DRIVER,
		DriverData: &CreateRequest_DriverData{
			Type:         CreateRequest_DriverData_TEXT,
			LayerUUID:    l.uuid,
			SoftwareUUID: l.softwareUUID,
			Font:         &font,
		},
	})
	if err != nil {
		return nil, err
	}

	td := &TextDriver{ctx: l.ctx, uuid: res.UUID}
	l.ctx.AddDriver(td.uuid, td)

	return td, nil
}

func (l *layer) NewImageDriver() (*ImageDriver, error) {
	res, err := l.ctx.SendCreateRequest(CreateRequest{
		Type: CreateRequest_DRIVER,
		DriverData: &CreateRequest_DriverData{
			Type:         CreateRequest_DriverData_IMAGE,
			LayerUUID:    l.uuid,
			SoftwareUUID: l.softwareUUID,
		},
	})
	if err != nil {
		return nil, err
	}

	id := &ImageDriver{ctx: l.ctx, uuid: res.UUID}
	l.ctx.AddDriver(id.uuid, id)

	return id, nil
}
