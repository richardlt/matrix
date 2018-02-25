package software

import (
	"context"
	"time"

	"github.com/pkg/errors"
	common "github.com/richardlt/matrix/sdk-go/common"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Software interface {
	Init(API) error
	Start(uint64)
	Close()
	ActionReceived(uint64, common.Command)
}

func init() {
	if err := loadImages(); err != nil {
		logrus.Errorf("%+v", err)
	}
	if err := loadThemes(); err != nil {
		logrus.Errorf("%+v", err)
	}
	if err := loadFonts(); err != nil {
		logrus.Errorf("%+v", err)
	}
}

// Connect initializes a new connection.
func Connect(uri string, s Software, reconnect bool) error {
	err := connect(uri, s)

	if !reconnect {
		return err
	}

	logrus.Debug("Software will reconnect in 1 sec")
	time.Sleep(time.Second)
	return Connect(uri, s, true)
}

func connect(uri string, s Software) error {
	conn, err := grpc.Dial(uri, grpc.WithInsecure())
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	c := NewSoftwareClient(conn)

	st, err := c.Connect(context.Background())
	if err != nil {
		return errors.WithStack(err)
	}

	connectRequestChannel := make(chan ConnectRequest)
	defer close(connectRequestChannel)

	// send event to the matrix core from channel
	go func() {
		for cr := range connectRequestChannel {
			if err := st.Send(&cr); err != nil {
				logrus.Errorf("%+v", errors.WithStack(err))
			}
		}
	}()

	// register the new software
	if err := st.Send(&ConnectRequest{
		Type: ConnectRequest_SOFTWARE,
		SoftwareData: &ConnectRequest_SoftwareData{
			Action: ConnectRequest_SoftwareData_REGISTER,
		},
	}); err != nil {
		return errors.WithStack(err)
	}

	// wait for the first response to obtain software uuid
	res, err := st.Recv()
	if err != nil {
		return errors.WithStack(err)
	}
	if res.Type != ConnectResponse_SOFTWARE ||
		res.SoftwareData.Action != ConnectResponse_SoftwareData_INIT {
		return errors.New("Error init software")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// every 3 seconds ping the matrix to test the conn
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				connectRequestChannel <- ConnectRequest{Type: ConnectRequest_PING}
			}
		}
	}()

	ct := newContext(connectRequestChannel, c)
	api := &api{ctx: ct, softwareUUID: res.SoftwareData.UUID}
	defer ct.Close()

	if err := s.Init(api); err != nil {
		return err
	}

	for {
		res, err := st.Recv()
		if err != nil {
			return errors.WithStack(err)
		}

		if res.Type == ConnectResponse_SOFTWARE {
			processResponse(s, res)
		} else {
			ct.ReceiveConnectResponse(*res)
		}
	}
}

func processResponse(s Software, res *ConnectResponse) {
	switch res.Type {
	case ConnectResponse_SOFTWARE:
		switch res.SoftwareData.Action {
		case ConnectResponse_SoftwareData_START:
			go s.Start(res.SoftwareData.PlayerCount)
		case ConnectResponse_SoftwareData_CLOSE:
			go s.Close()
		case ConnectResponse_SoftwareData_PLAYER_COMMAND:
			go s.ActionReceived(res.SoftwareData.Slot, res.SoftwareData.Command)
		}
	}
}
