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
	GetImageFromRemote(name string) (Image, error)
	GetColorFromLocalThemeByName(themeName string, colorName string) common.Color
	GetColorFromRemoteThemeByName(themeName string, colorName string) (common.Color, error)
	GetFontFromLocal(name string) Font
	GetFontFromRemote(name string) (Font, error)
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

// GetImageFromLocal retrieve an image from local file, if not found returns
// a zero image.
func (a *api) GetImageFromLocal(name string) Image {
	for _, i := range is {
		if i.Name == name {
			return i
		}
	}
	return Image{}
}

// GetImageFromRemote retrieve an image from core, if not found returns
// a zero image.
func (a *api) GetImageFromRemote(name string) (Image, error) {
	res, err := a.ctx.SendLoadRequest(LoadRequest{
		Type:      LoadRequest_IMAGE,
		ImageData: &LoadRequest_ImageData{Name: name},
	})
	if err != nil {
		return Image{}, errors.WithStack(err)
	}
	return *res.Image, nil
}

// GetFontFromLocal retrieve a font from local file, if not found returns a
// zero font.
func (a *api) GetFontFromLocal(name string) Font {
	for _, f := range fs {
		if f.Name == name {
			return f
		}
	}
	return Font{}
}

// GetFontFromRemote retrieve a font from core, if not found returns
// a zero font.
func (a *api) GetFontFromRemote(name string) (Font, error) {
	res, err := a.ctx.SendLoadRequest(LoadRequest{
		Type:     LoadRequest_FONT,
		FontData: &LoadRequest_FontData{Name: name},
	})
	if err != nil {
		return Font{}, errors.WithStack(err)
	}
	return *res.Font, nil
}

// GetColorFromLocalThemeByName retrieve a loaded theme's color in memory, if
// theme or color not found returns a zero color.
func (a *api) GetColorFromLocalThemeByName(themeName, name string) common.Color {
	for _, t := range ts {
		if t.Name == themeName {
			for k, c := range t.Colors {
				if k == name {
					return *c
				}
			}
		}
	}
	return common.Color{}
}

// GetColorFromRemoteThemeByName retrieve a loaded theme's color from core, if
// theme or color not found returns a zero color.
func (a *api) GetColorFromRemoteThemeByName(themeName, name string) (common.Color, error) {
	res, err := a.ctx.SendLoadRequest(LoadRequest{
		Type:      LoadRequest_COLOR,
		ColorData: &LoadRequest_ColorData{Name: name, ThemeName: themeName},
	})
	if err != nil {
		return common.Color{}, errors.WithStack(err)
	}
	return *res.Color, nil
}

// NewLayer ask for a layer creation and returns its uuid.
func (a *api) NewLayer() (Layer, error) {
	res, err := a.ctx.SendCreateRequest(CreateRequest{
		Type:      CreateRequest_LAYER,
		LayerData: &CreateRequest_LayerData{SoftwareUUID: a.softwareUUID},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	l := &layer{
		ctx:          a.ctx,
		softwareUUID: a.softwareUUID,
		uuid:         res.UUID,
	}
	a.ctx.AddLayer(l.uuid, l)

	return l, nil
}
