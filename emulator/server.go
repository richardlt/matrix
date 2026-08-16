package emulator

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/matrix/internal/errors"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/display"
	"github.com/richardlt/matrix/websocket"
)

// public holds the built web app, so the binary serves it without needing the files
// next to it at runtime. `make build-all` runs the web build before the Go one, which
// is the order this requires. The checked-in public/.gitkeep keeps this pattern
// matching before any web build has run; the binary then starts and serves 404s.
//
//go:embed all:public
var public embed.FS

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

	// Core only sends frames when something changes, so a browser opened after startup
	// would otherwise sit on a blank screen until a player pressed a button. Keep the
	// most recent frame per screen and replay them to each new client.
	var (
		lastLock   sync.RWMutex
		lastFrames = map[int]frame{}
	)

	s.OnConnect(func(c *websocket.Client) {
		lastLock.RLock()
		defer lastLock.RUnlock()
		for i := 0; i < len(lastFrames); i++ {
			f, ok := lastFrames[i]
			if !ok {
				continue
			}
			if err := c.Send("frame", f); err != nil {
				logrus.Errorf("%+v", errors.Errorf("replaying frame %d to a new client: %w", i, err))
			}
		}
	})

	go func() {
		for f := range frameChannel {
			lastLock.Lock()
			lastFrames[f.Number] = f
			lastLock.Unlock()

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

	assets, err := fs.Sub(public, "public/app")
	if err != nil {
		return errors.Errorf("opening embedded assets: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/websocket", s)
	mux.Handle("/", http.FileServer(http.FS(assets)))

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
