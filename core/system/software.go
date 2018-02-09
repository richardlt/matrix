package system

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/richardlt/matrix/core/drivers"
	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
	softwareSDK "github.com/richardlt/matrix/sdk-go/software"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// NewSoftwareServer returns new software server.
func NewSoftwareServer() *SoftwareServer { return &SoftwareServer{} }

// SoftwareServer exposes RPC server for softwares.
type SoftwareServer struct {
	softwares              []*software
	softwareLock           sync.RWMutex
	current                *software
	printCallback          func(f render.Frame)
	softwareChangeCallback func(sms []SoftwareMeta)
}

// ResetCallback removes existing callbacks.
func (s *SoftwareServer) ResetCallback() {
	s.printCallback = nil
	s.softwareChangeCallback = nil
}

// OnPrint allows to set print callback.
func (s *SoftwareServer) OnPrint(f func(render.Frame)) { s.printCallback = f }

// OnSoftwareChange allows to set software change callback.
func (s *SoftwareServer) OnSoftwareChange(f func([]SoftwareMeta)) {
	s.softwareChangeCallback = f
}

// StartSoftware executes start for a given software meta and player count, it
// also close current software if exists.
func (s *SoftwareServer) StartSoftware(meta SoftwareMeta, playerCount uint64) error {
	var new *software

	for _, so := range s.softwares {
		if so.UUID == meta.UUID {
			new = so
			break
		}
	}

	if new == nil {
		return errors.New("Invalid given software meta")
	}

	if s.current != nil {
		s.CloseSoftware()
	}

	s.current = new
	s.current.Start(playerCount)

	return nil
}

// CloseSoftware allows to stop running software if exists.
func (s *SoftwareServer) CloseSoftware() {
	if s.current != nil {
		s.current.Close()
		s.current = nil
	}
}

// Command sends player command to current software.
func (s *SoftwareServer) Command(slot uint64, cmd common.Command) {
	if s.current != nil {
		s.current.Command(slot, cmd)
	}
}

func (s *SoftwareServer) notifyChanges() {
	if s.softwareChangeCallback != nil {
		go s.softwareChangeCallback(s.GetSoftwaresMeta())
	}
}

// AddSoftware allows to push a new software in server.
func (s *SoftwareServer) AddSoftware(so *software) {
	s.softwareLock.Lock()
	s.softwares = append(s.softwares, so)
	s.softwareLock.Unlock()

	s.notifyChanges()
}

// GetSoftwaresMeta returns metadata from softwares.
func (s *SoftwareServer) GetSoftwaresMeta() []SoftwareMeta {
	s.softwareLock.RLock()
	sms := []SoftwareMeta{}
	for _, so := range s.softwares {
		if so.Ready {
			sms = append(sms, so.GetMeta())
		}
	}
	s.softwareLock.RUnlock()
	return sms
}

// RemoveSoftware removes a existing software in server.
func (s *SoftwareServer) RemoveSoftware(so *software) {
	if s.current == so {
		s.current = nil
	}

	s.softwareLock.Lock()
	var ss []*software
	for _, es := range s.softwares {
		if es.UUID != so.UUID {
			ss = append(ss, es)
		}
	}
	s.softwares = ss
	s.softwareLock.Unlock()

	s.notifyChanges()
}

// Connect software action.
func (s *SoftwareServer) Connect(stream softwareSDK.Software_ConnectServer) error {
	chRes := make(chan softwareSDK.ConnectResponse)
	defer close(chRes)

	so := newSoftware(chRes)

	logrus.Debugf("Software %s connect", so.UUID)
	defer logrus.Debugf("Software %s disconnect", so.UUID)

	s.AddSoftware(so)
	defer s.RemoveSoftware(so)

	// wait for the register request
	req, err := stream.Recv()
	if err != nil {
		return errors.WithStack(err)
	}
	if req.Type != softwareSDK.ConnectRequest_SOFTWARE ||
		req.SoftwareData.Action != softwareSDK.ConnectRequest_SoftwareData_REGISTER {
		return errors.New("error register software")
	}

	// init the software with a random uuid
	if err := stream.Send(&softwareSDK.ConnectResponse{
		Type: softwareSDK.ConnectResponse_SOFTWARE,
		SoftwareData: &softwareSDK.ConnectResponse_SoftwareData{
			Action: softwareSDK.ConnectResponse_SoftwareData_INIT,
			UUID:   so.UUID,
		},
	}); err != nil {
		return errors.WithStack(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// every 3 seconds ping the spftware to test the conn
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				chRes <- softwareSDK.ConnectResponse{Type: softwareSDK.ConnectResponse_PING}
			}
		}
	}()

	go func() {
		for r := range chRes {
			if err := stream.Send(&r); err != nil {
				logrus.Errorf("%+v", errors.WithStack(err))
			}
		}
	}()

	for {
		req, err := stream.Recv()
		if err != nil {
			return errors.WithStack(err)
		}
		s.processRequest(so, req)
	}
}

func (s *SoftwareServer) processRequest(so *software, req *softwareSDK.ConnectRequest) {
	switch req.Type {
	case softwareSDK.ConnectRequest_SOFTWARE:
		switch req.SoftwareData.Action {
		case softwareSDK.ConnectRequest_SoftwareData_SET_CONFIG:
			so.SetConfig(*req.SoftwareData.Config)
		case softwareSDK.ConnectRequest_SoftwareData_READY:
			so.Ready = true
			s.notifyChanges()
		case softwareSDK.ConnectRequest_SoftwareData_PRINT:
			if s.current != nil && s.current.UUID == so.UUID && s.printCallback != nil {
				s.printCallback(so.GetTopFrame())
			}
		}
	case softwareSDK.ConnectRequest_LAYER:
		switch req.LayerData.Action {
		case softwareSDK.ConnectRequest_LayerData_SET_WITH_COORD:
			so.LayerSetWithCoord(req.LayerData.UUID, *req.LayerData.Coord, *req.LayerData.Color)
		case softwareSDK.ConnectRequest_LayerData_CLEAN:
			so.LayerClean(req.LayerData.UUID)
		case softwareSDK.ConnectRequest_LayerData_REMOVE:
			so.LayerRemove(req.LayerData.UUID)
		}
	case softwareSDK.ConnectRequest_DRIVER:
		switch req.DriverData.Action {
		case softwareSDK.ConnectRequest_DriverData_RENDER:
			if rd, ok := so.randomDrivers[req.DriverData.UUID]; ok {
				rd.Render()
			}
			if cd, ok := so.caracterDrivers[req.DriverData.UUID]; ok {
				cd.Render([]rune(req.DriverData.Caracter)[0], *req.DriverData.Coord,
					*req.DriverData.Color, *req.DriverData.Background)
			}
			if td, ok := so.textDrivers[req.DriverData.UUID]; ok {
				td.Render(req.DriverData.Text, *req.DriverData.Coord,
					*req.DriverData.Color, *req.DriverData.Background)
			}
			if id, ok := so.imageDrivers[req.DriverData.UUID]; ok {
				id.Render(*req.DriverData.Image, *req.DriverData.Coord)
			}
		case softwareSDK.ConnectRequest_DriverData_STOP:
			if td, ok := so.textDrivers[req.DriverData.UUID]; ok {
				td.Stop()
			}
		}
	}
}

// Create software action.
func (s *SoftwareServer) Create(ctx context.Context,
	req *softwareSDK.CreateRequest) (res *softwareSDK.CreateResponse, err error) {
	switch req.Type {
	case softwareSDK.CreateRequest_LAYER:
		if so := s.getSoftwareByUUID(req.LayerData.SoftwareUUID); so != nil {
			res = &softwareSDK.CreateResponse{
				Type: softwareSDK.CreateResponse_LAYER,
				UUID: so.CreateLayer(),
			}
		}
	case softwareSDK.CreateRequest_DRIVER:
		if so := s.getSoftwareByUUID(req.DriverData.SoftwareUUID); so != nil {
			if l := so.GetLayerByUUID(req.DriverData.LayerUUID); l != nil {
				res = &softwareSDK.CreateResponse{
					Type: softwareSDK.CreateResponse_DRIVER,
					UUID: so.CreateDriver(l, *req.DriverData),
				}
			}
		}
	}
	return
}

// Load software action.
func (s *SoftwareServer) Load(ctx context.Context,
	req *softwareSDK.LoadRequest) (res *softwareSDK.LoadResponse, err error) {
	switch req.Type {
	case softwareSDK.LoadRequest_IMAGE:
		i := render.GetImageByName(req.ImageData.Name)
		res = &softwareSDK.LoadResponse{
			Type:  softwareSDK.LoadResponse_IMAGE,
			Image: &i,
		}
	case softwareSDK.LoadRequest_COLOR:
		c := render.GetColorFromLocalThemeByName(req.ColorData.ThemeName,
			req.ColorData.Name)
		res = &softwareSDK.LoadResponse{
			Type:  softwareSDK.LoadResponse_COLOR,
			Color: &c,
		}
	case softwareSDK.LoadRequest_FONT:
		f := render.GetFontByName(req.FontData.Name)
		res = &softwareSDK.LoadResponse{
			Type: softwareSDK.LoadResponse_FONT,
			Font: &f,
		}
	}
	return
}

func (s *SoftwareServer) getSoftwareByUUID(uuid string) *software {
	s.softwareLock.RLock()
	defer s.softwareLock.RUnlock()

	for _, so := range s.softwares {
		if so.UUID == uuid {
			return so
		}
	}
	return nil
}

func newSoftware(connectResponseChannel chan softwareSDK.ConnectResponse) *software {
	return &software{
		UUID: uuid.NewV4().String(),
		connectResponseChannel: connectResponseChannel,
		matrix:                 render.NewMatrix(16, 9),
		layers:                 make(map[string]*render.Frame),
		randomDrivers:          make(map[string]*drivers.Random),
		caracterDrivers:        make(map[string]*drivers.Caracter),
		textDrivers:            make(map[string]*drivers.Text),
		imageDrivers:           make(map[string]*drivers.Image),
	}
}

type software struct {
	UUID                           string
	Logo                           softwareSDK.Image
	MinPlayerCount, MaxPlayerCount uint64
	connectResponseChannel         chan softwareSDK.ConnectResponse
	Ready                          bool
	matrix                         *render.Matrix
	layers                         map[string]*render.Frame
	randomDrivers                  map[string]*drivers.Random
	caracterDrivers                map[string]*drivers.Caracter
	textDrivers                    map[string]*drivers.Text
	imageDrivers                   map[string]*drivers.Image
}

func (s software) GetMeta() SoftwareMeta {
	return SoftwareMeta{
		UUID:           s.UUID,
		Logo:           s.Logo,
		MinPlayerCount: s.MinPlayerCount,
		MaxPlayerCount: s.MaxPlayerCount,
	}
}

func (s software) GetTopFrame() render.Frame {
	s.matrix.PrintFrame()
	return s.matrix.GetTopFrame()
}

func (s *software) LayerSetWithCoord(layerUUID string, coord common.Coord, col common.Color) {
	if l, ok := s.layers[layerUUID]; ok {
		l.SetWithCoord(coord, col)
	}
}

func (s *software) LayerClean(layerUUID string) {
	if l, ok := s.layers[layerUUID]; ok {
		l.Clean()
	}
}

func (s *software) LayerRemove(layerUUID string) {
	if l, ok := s.layers[layerUUID]; ok {
		s.matrix.RemoveLayer(l)
	}
}

func (s *software) SetConfig(c softwareSDK.ConnectRequest_SoftwareData_Config) {
	if c.Logo != nil {
		s.Logo = *c.Logo
		s.MinPlayerCount = c.MinPlayerCount
		s.MaxPlayerCount = c.MaxPlayerCount
	}
}

func (s *software) Start(playerCount uint64) {
	s.connectResponseChannel <- softwareSDK.ConnectResponse{
		Type: softwareSDK.ConnectResponse_SOFTWARE,
		SoftwareData: &softwareSDK.ConnectResponse_SoftwareData{
			Action:      softwareSDK.ConnectResponse_SoftwareData_START,
			PlayerCount: playerCount,
		},
	}
}

func (s *software) Close() {
	s.connectResponseChannel <- softwareSDK.ConnectResponse{
		Type: softwareSDK.ConnectResponse_SOFTWARE,
		SoftwareData: &softwareSDK.ConnectResponse_SoftwareData{
			Action: softwareSDK.ConnectResponse_SoftwareData_CLOSE,
		},
	}
}

func (s *software) Command(slot uint64, command common.Command) {
	s.connectResponseChannel <- softwareSDK.ConnectResponse{
		Type: softwareSDK.ConnectResponse_SOFTWARE,
		SoftwareData: &softwareSDK.ConnectResponse_SoftwareData{
			Action:  softwareSDK.ConnectResponse_SoftwareData_PLAYER_COMMAND,
			Slot:    slot,
			Command: command,
		},
	}
}

func (s *software) CreateLayer() string {
	uuid := uuid.NewV4().String()
	s.layers[uuid] = s.matrix.NewFrame()
	return uuid
}

func (s *software) CreateDriver(l *render.Frame, driverData softwareSDK.CreateRequest_DriverData) string {
	uuid := uuid.NewV4().String()

	switch driverData.Type {
	case softwareSDK.CreateRequest_DriverData_RANDOM:
		rd := drivers.NewRandom(l)
		rd.OnEnd(func() {
			s.connectResponseChannel <- softwareSDK.ConnectResponse{
				Type: softwareSDK.ConnectResponse_DRIVER,
				DriverData: &softwareSDK.ConnectResponse_DriverData{
					Action: softwareSDK.ConnectResponse_DriverData_END,
					UUID:   uuid,
				},
			}
		})
		s.randomDrivers[uuid] = rd
	case softwareSDK.CreateRequest_DriverData_CARACTER:
		cd := drivers.NewCaracter(l, *driverData.Font)
		cd.OnEnd(func() {
			s.connectResponseChannel <- softwareSDK.ConnectResponse{
				Type: softwareSDK.ConnectResponse_DRIVER,
				DriverData: &softwareSDK.ConnectResponse_DriverData{
					Action: softwareSDK.ConnectResponse_DriverData_END,
					UUID:   uuid,
				},
			}
		})
		s.caracterDrivers[uuid] = cd
	case softwareSDK.CreateRequest_DriverData_TEXT:
		td := drivers.NewText(l, *driverData.Font)
		td.OnEnd(func() {
			s.connectResponseChannel <- softwareSDK.ConnectResponse{
				Type: softwareSDK.ConnectResponse_DRIVER,
				DriverData: &softwareSDK.ConnectResponse_DriverData{
					Action: softwareSDK.ConnectResponse_DriverData_END,
					UUID:   uuid,
				},
			}
		})
		td.OnStep(func(total, current uint64) {
			s.connectResponseChannel <- softwareSDK.ConnectResponse{
				Type: softwareSDK.ConnectResponse_DRIVER,
				DriverData: &softwareSDK.ConnectResponse_DriverData{
					Action:  softwareSDK.ConnectResponse_DriverData_STEP,
					UUID:    uuid,
					Total:   total,
					Current: current,
				},
			}
		})
		s.textDrivers[uuid] = td
	case softwareSDK.CreateRequest_DriverData_IMAGE:
		id := drivers.NewImage(l)
		id.OnEnd(func() {
			s.connectResponseChannel <- softwareSDK.ConnectResponse{
				Type: softwareSDK.ConnectResponse_DRIVER,
				DriverData: &softwareSDK.ConnectResponse_DriverData{
					Action: softwareSDK.ConnectResponse_DriverData_END,
					UUID:   uuid,
				},
			}
		})
		s.imageDrivers[uuid] = id
	}

	return uuid
}

func (s *software) GetLayerByUUID(uuid string) *render.Frame {
	if l, ok := s.layers[uuid]; ok {
		return l
	}

	return nil
}

// SoftwareMeta contains metadata for a software.
type SoftwareMeta struct {
	UUID                           string
	Logo                           softwareSDK.Image
	MinPlayerCount, MaxPlayerCount uint64
}
