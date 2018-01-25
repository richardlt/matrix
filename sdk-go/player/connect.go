package player

import (
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type Player interface {
	Init(*API) error
}

// Connect initializes a new connection.
func Connect(uri string, p Player, reconnect bool) error {
	err := connect(uri, p)

	if !reconnect {
		return err
	}

	logrus.Debug("Player will reconnect in 1 sec")
	time.Sleep(time.Second)
	return Connect(uri, p, true)
}

func connect(uri string, p Player) error {
	conn, err := grpc.Dial(uri, grpc.WithInsecure())
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	c := NewPlayerClient(conn)

	st, err := c.Connect(context.Background())
	if err != nil {
		return errors.WithStack(err)
	}

	requestChannel := make(chan Request)
	defer close(requestChannel)

	api := &API{requestChannel: requestChannel}
	defer api.Close()

	// send event to the matrix core from channel
	go func() {
		for cr := range requestChannel {
			if err := st.Send(&cr); err != nil {
				logrus.Errorf("%+v", errors.WithStack(err))
			}
		}
	}()

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
				requestChannel <- Request{Type: Request_PING}
			}
		}
	}()

	if err := p.Init(api); err != nil {
		return err
	}

	for {
		_, err := st.Recv()
		if err != nil {
			return errors.WithStack(err)
		}
	}
}
