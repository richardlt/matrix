package software

import (
	"github.com/pkg/errors"
	common "github.com/richardlt/matrix/sdk-go/common"
)

type API interface {
	NewLayer() (Layer, error)
	Ready() error
	Print() error
	SetConfig(ConnectRequest_SoftwareData_Config) error
	GetImageFromLocal(name string) Image
	GetColorFromThemeByName(themeName string, colorName string) common.Color
	GetFontFromLocal(name string) Font
}

// API allow the software to send events to the matrix core.
type api struct {
	ctx          *ctx
	softwareUUID string
}

// Ready should be used when the software initialization finished.
func (a *api) Ready() error {
	return a.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_SOFTWARE,
		SoftwareData: &ConnectRequest_SoftwareData{
			Action: ConnectRequest_SoftwareData_READY,
		},
	})
}

// Print compute the final frame from all software's layers.
func (a *api) Print() error {
	return a.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_SOFTWARE,
		SoftwareData: &ConnectRequest_SoftwareData{
			Action: ConnectRequest_SoftwareData_PRINT,
		},
	})
}

// SetConfig send software config to core.
func (a *api) SetConfig(c ConnectRequest_SoftwareData_Config) error {
	return a.ctx.SendConnectRequest(ConnectRequest{
		Type: ConnectRequest_SOFTWARE,
		SoftwareData: &ConnectRequest_SoftwareData{
			Action: ConnectRequest_SoftwareData_SET_CONFIG,
			Config: &c,
		},
	})
}

// GetImageFromLocal retrieve an image from local file.
func (a *api) GetImageFromLocal(name string) Image {
	for _, i := range is {
		if i.Name == name {
			return i
		}
	}
	return Image{}
}

// GetFontFromLocal retrieve a font from local file.
func (a *api) GetFontFromLocal(name string) Font {
	for _, f := range fs {
		if f.Name == name {
			return f
		}
	}
	return Font{}
}

// GetColorFromThemeByName retrieve an loaded theme's color in memory.
func (a *api) GetColorFromThemeByName(themeName, colorName string) common.Color {
	for _, t := range ts {
		if t.Name == themeName {
			for k, c := range t.Colors {
				if k == colorName {
					return *c
				}
			}
		}
	}

	return common.Color{}
}

// NewLayer ask for a layer creation and return its uuid.
func (a *api) NewLayer() (Layer, error) {
	res, err := a.ctx.SendCreateRequest(CreateRequest{
		Type: CreateRequest_LAYER,
		LayerData: &CreateRequest_LayerData{
			SoftwareUUID: a.softwareUUID,
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	l := &layer{
		ctx:          a.ctx,
		softwareUUID: a.softwareUUID,
		uuid:         res.LayerData.UUID,
	}
	a.ctx.AddLayer(l.uuid, l)

	return l, nil
}
