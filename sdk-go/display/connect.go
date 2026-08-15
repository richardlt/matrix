package display

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/richardlt/matrix/internal/errors"
	common "github.com/richardlt/matrix/sdk-go/common"
)

type Display interface {
	FramesReceived([]*common.Frame)
}

// Connect initializes a new connection.
func Connect(uri string, d Display, reconnect bool) error {
	err := connect(uri, d)

	if !reconnect {
		return err
	}

	logrus.Debug("Display will reconnect in 1 sec")
	time.Sleep(time.Second)
	return Connect(uri, d, true)
}

func connect(uri string, d Display) error {
	conn, err := grpc.Dial(uri, grpc.WithInsecure())
	if err != nil {
		return errors.Errorf("dialing core at %s: %w", uri, err)
	}
	defer func() { _ = conn.Close() }()

	c := NewDisplayClient(conn)

	st, err := c.Connect(context.Background())
	if err != nil {
		return errors.Errorf("opening display stream to %s: %w", uri, err)
	}

	requestChannel := make(chan *Request)
	defer close(requestChannel)

	// send event to the matrix core from channel
	go func() {
		for cr := range requestChannel {
			if err := st.Send(cr); err != nil {
				logrus.Errorf("%+v", errors.Errorf("sending display request: %w", err))
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
				requestChannel <- &Request{Type: Request_PING}
			}
		}
	}()

	for {
		res, err := st.Recv()
		if err != nil {
			return errors.Errorf("receiving display response: %w", err)
		}

		processResponse(d, res)
	}
}

func processResponse(d Display, res *Response) {
	switch res.Type {
	case Response_DISPLAY:
		switch res.DisplayData.Action {
		case Response_DisplayData_FRAMES:
			d.FramesReceived(res.DisplayData.Frames)
		}
	}
}
