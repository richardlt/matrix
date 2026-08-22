package websocket

import (
	"context"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/internal/errors"
)

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		ID:   uuid.NewString(),
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
				return errors.Errorf("reading websocket message from client %s: %w", c.ID, err)
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
		return errors.Errorf("sending %q event to client %s: %w", eventType, c.ID, err)
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
