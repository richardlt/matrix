package software

import (
	"sync"

	"github.com/pkg/errors"
	context "golang.org/x/net/context"
)

func newContext(connectRequestChannel chan ConnectRequest,
	client SoftwareClient) *ctx {
	return &ctx{
		connectRequestChannel: connectRequestChannel,
		client:                client,
		layers:                make(map[string]*layer),
		drivers:               make(map[string]driver),
	}
}

type ctx struct {
	client                SoftwareClient
	connectRequestChannel chan ConnectRequest
	layers                map[string]*layer
	layerLock             sync.RWMutex
	drivers               map[string]driver
	driverLock            sync.RWMutex
}

func (c *ctx) AddDriver(uuid string, d driver) {
	c.driverLock.Lock()
	c.drivers[uuid] = d
	c.driverLock.Unlock()
}

func (c *ctx) AddLayer(uuid string, l *layer) {
	c.layerLock.Lock()
	c.layers[uuid] = l
	c.layerLock.Unlock()
}

func (c *ctx) SendConnectRequest(req ConnectRequest) error {
	if c.connectRequestChannel == nil {
		return errors.New("API is closed")
	}

	c.connectRequestChannel <- req
	return nil
}

func (c *ctx) SendCreateRequest(req CreateRequest) (*CreateResponse, error) {
	if c.client == nil {
		return nil, errors.New("API is closed")
	}

	res, err := c.client.Create(context.Background(), &req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (c *ctx) SendLoadRequest(req LoadRequest) (*LoadResponse, error) {
	if c.client == nil {
		return nil, errors.New("API is closed")
	}

	res, err := c.client.Load(context.Background(), &req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (c *ctx) Close() {
	c.connectRequestChannel = nil
	c.client = nil
}

func (c *ctx) ReceiveConnectResponse(res ConnectResponse) {
	switch res.Type {
	case ConnectResponse_DRIVER:
		c.driverLock.RLock()
		if d, ok := c.drivers[res.DriverData.UUID]; ok {
			switch res.DriverData.Action {
			case ConnectResponse_DriverData_END:
				go d.End()
			case ConnectResponse_DriverData_STEP:
				go d.Step(res.DriverData.Total, res.DriverData.Current)
			}
		}
		c.driverLock.RUnlock()
	}
}
