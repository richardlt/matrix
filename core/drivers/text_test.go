package drivers

import (
	"runtime"
	"testing"
	"time"

	"github.com/richardlt/matrix/core/render"
	"github.com/richardlt/matrix/sdk-go/common"
	"github.com/richardlt/matrix/sdk-go/software"
)

// A stopped scroll must not leave its goroutine behind. The paused screen starts and
// stops one every time play is suspended, so anything left running would pile up for as
// long as the console is on.
func TestStopEndsScrollGoroutine(t *testing.T) {
	f := render.NewFrame(16, 9)
	before := runtime.NumGoroutine()

	for i := 0; i < 20; i++ {
		td := NewText(&f, &software.Font{Height: 5})
		td.Render("PAUSE", &common.Coord{X: 4, Y: 4}, &common.Color{}, &common.Color{}, true)
		td.Stop()
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= before+2 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("goroutines stranded: %d before, %d after", before, runtime.NumGoroutine())
}

// Stopping a scroll that never started, or stopping one twice, happens whenever a state
// is left without a pause in progress.
func TestStopIsIdempotent(t *testing.T) {
	f := render.NewFrame(16, 9)
	td := NewText(&f, &software.Font{Height: 5})

	td.Stop()
	td.Render("PAUSE", &common.Coord{X: 4, Y: 4}, &common.Color{}, &common.Color{}, true)
	td.Stop()
	td.Stop()
}
