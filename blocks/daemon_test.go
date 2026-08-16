package blocks

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/richardlt/matrix/sdk-go/common"
)

// action must not strand a goroutine once the game loop has gone, which happens both on
// close and when the game ends on its own.
func TestActionDoesNotBlockAfterGameLoopExits(t *testing.T) {
	b := &blocks{}
	ctx, cancel := context.WithCancel(context.Background())
	b.running, b.cancel = ctx, cancel
	b.commandChan = make(chan common.Command) // unbuffered, nobody reading

	cancel() // the game loop has returned

	before := runtime.NumGoroutine()
	for i := 0; i < 50; i++ {
		go b.action(common.Command_A_UP)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= before+2 {
			return // they all returned
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("goroutines stranded: %d before, %d after", before, runtime.NumGoroutine())
}
