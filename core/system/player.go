package system

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/sdk-go/common"
	playerSDK "github.com/richardlt/matrix/sdk-go/player"
)

// NewPlayerServer returns new player server.
func NewPlayerServer() *PlayerServer { return &PlayerServer{} }

// PlayerServer expose RPC server for players.
type PlayerServer struct {
	players        []*player
	playerLock     sync.RWMutex
	actionCallback func(Action)
}

// ResetCallback removes existing callbacks.
func (p *PlayerServer) ResetCallback() { p.actionCallback = nil }

// OnAction allows to set action callback.
func (p *PlayerServer) OnAction(f func(Action)) { p.actionCallback = f }

// AddPlayer allows to push a new player in server.
func (p *PlayerServer) AddPlayer(pl *player) {
	p.playerLock.Lock()
	p.players = append(p.players, pl)
	p.playerLock.Unlock()
}

// RemovePlayer removes a existing player in server.
func (p *PlayerServer) RemovePlayer(pl *player) {
	p.playerLock.Lock()
	var ps []*player
	for _, ep := range p.players {
		if ep.UUID != pl.UUID {
			ps = append(ps, ep)
		}
	}
	p.players = ps
	p.playerLock.Unlock()
}

// Connect player action.
func (p *PlayerServer) Connect(stream playerSDK.Player_ConnectServer) error {
	chRes := make(chan playerSDK.Response)
	defer close(chRes)

	pl := newPlayer(chRes)

	logrus.Debugf("Player %s connect", pl.UUID)
	defer logrus.Debugf("Player %s disconnect", pl.UUID)

	p.AddPlayer(pl)
	defer p.RemovePlayer(pl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// every 3 seconds ping the player to test the conn
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				chRes <- playerSDK.Response{Type: playerSDK.Response_PING}
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
		p.processRequest(pl, *req)
	}
}

func (p *PlayerServer) processRequest(pl *player, req playerSDK.Request) {
	if req.Type == playerSDK.Request_PLAYER {
		switch req.PlayerData.Action {
		case playerSDK.Request_PlayerData_COMMAND:
			logrus.Debugf("Player %s send %s command for %d slot", pl.UUID, req.PlayerData.Command, int(req.PlayerData.Slot))
			if p.actionCallback != nil {
				p.actionCallback(Action{req.PlayerData.Slot, req.PlayerData.Command})
			}
		}
	}
}

func newPlayer(chRes chan playerSDK.Response) *player {
	return &player{uuid.NewV4().String(), chRes}
}

type player struct {
	UUID            string
	responseChannel chan playerSDK.Response
}

// Action is generated by player's command.
type Action struct {
	Slot    uint64
	Command common.Command
}
