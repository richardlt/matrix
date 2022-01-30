package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]*Client),
	}
}

type Server struct {
	mutex             sync.RWMutex
	clients           map[string]*Client
	onConnectCallback func(*Client)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := s.Serve(w, r); err != nil {
		logrus.Errorf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) error {
	c, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer c.Close()

	client := NewClient(c)

	s.AddClient(client)
	defer s.RemoveClient(client)

	return client.Listen(r.Context())
}

func (s *Server) AddClient(c *Client) {
	if s.onConnectCallback != nil {
		s.onConnectCallback(c)
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.clients[c.ID] = c
}

func (s *Server) RemoveClient(c *Client) {
	if c.onDisconnectCallback != nil {
		c.onDisconnectCallback()
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.clients, c.ID)
}

func (s *Server) Broadcast(eventType string, data interface{}) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for i := range s.clients {
		if err := s.clients[i].Send(eventType, data); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) OnConnect(f func(*Client)) {
	s.onConnectCallback = f
}
