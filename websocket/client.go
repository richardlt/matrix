package websocket

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		ID:   uuid.NewV4().String(),
		conn: conn,
	}
}

type Client struct {
	ID                   string
	conn                 *websocket.Conn
	onEventCallback      func(eventType string, data interface{})
	onDisconnectCallback func()
}

func (c *Client) Listen(ctx context.Context) error {
	logrus.Debugf("websocket client %q connected", c.ID)
	c.conn.EnableWriteCompression(true)
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var m message
		if err := c.conn.ReadJSON(&m); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return errors.WithStack(err)
			}
			logrus.Debugf("%+v", err)
			break
		}
		if c.onEventCallback != nil {
			c.onEventCallback(m.Type, m.Data)
		}
	}
	logrus.Debugf("websocket client %q disconnected", c.ID)
	return nil
}

func (c *Client) Send(eventType string, data interface{}) error {
	if err := c.conn.WriteJSON(message{
		Type: eventType,
		Data: data,
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Client) OnEvent(f func(eventType string, data interface{})) {
	c.onEventCallback = f
}

func (c *Client) OnDisconnect(f func()) {
	c.onDisconnectCallback = f
}

type message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}
