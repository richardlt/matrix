package emulator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/internal/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/display"
	"github.com/richardlt/matrix/websocket"
)

type frame struct {
	Number int     `json:"number"`
	Pixels []pixel `json:"pixels"`
}

type pixel struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

// Start the emulator server.
func Start(port int, uri string) error {
	frameChannel := make(chan frame)
	defer close(frameChannel)

	s := websocket.NewServer()

	go func() {
		for f := range frameChannel {
			if err := s.Broadcast("frame", f); err != nil {
				logrus.Errorf("%+v", err)
			}
		}
	}()

	go func() {
		if err := display.Connect(uri, emulator{frameChannel}, true); err != nil {
			logrus.Errorf("%+v", errors.Errorf("connecting emulator display to %s: %w", uri, err))
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/websocket", s)
	mux.Handle("/", http.FileServer(http.Dir("./emulator/public")))

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	logrus.Infof("Start emulator on port %d\n", port)
	if err := srv.ListenAndServe(); err != nil {
		return errors.Errorf("serving emulator on port %d: %w", port, err)
	}
	return nil
}

type emulator struct{ frameChannel chan frame }

func (e emulator) FramesReceived(fs []*common.Frame) {
	for i, f := range fs {
		frame := frame{
			Number: i,
			Pixels: make([]pixel, len(f.Pixels)),
		}
		for i, c := range f.Pixels {
			frame.Pixels[i].R = int(c.R)
			frame.Pixels[i].G = int(c.G)
			frame.Pixels[i].B = int(c.B)
		}
		e.frameChannel <- frame
	}
}
