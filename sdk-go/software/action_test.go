package software

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/richardlt/matrix/sdk-go/common"
)

func TestPassThrough(t *testing.T) {
	assert := assert.New(t)

	var s uint64
	pt := passThrough{}
	pt.OnAction(func(slot uint64) { s = slot })

	pt.SendAction(3, common.Command_A_DOWN)
	assert.Equal(uint64(3), s)
}

func TestMultiPress(t *testing.T) {
	assert := assert.New(t)

	multi := make([]bool, 2)
	mp := NewMultiPress(common.Button_L, common.Button_R)
	mp.OnAction(func(slot uint64) { multi[slot] = true })

	mp.SendAction(0, common.Command_L_DOWN)
	mp.SendAction(1, common.Command_R_DOWN)
	assert.False(multi[0])
	assert.False(multi[1])

	mp.SendAction(0, common.Command_R_UP)
	assert.False(multi[0])
	assert.False(multi[1])

	mp.SendAction(0, common.Command_R_DOWN)
	assert.True(multi[0])
	assert.False(multi[1])

	mp.SendAction(1, common.Command_L_DOWN)
	assert.True(multi[0])
	assert.True(multi[1])
}

func TestLongPress(t *testing.T) {
	assert := assert.New(t)

	count := make([]int, 3)
	var mutex sync.Mutex
	lp := NewLongPress(common.Button_A, 500*time.Millisecond, 250*time.Millisecond)
	lp.OnAction(func(slot uint64) {
		mutex.Lock()
		count[slot]++
		mutex.Unlock()
	})

	lp.SendAction(0, common.Command_A_DOWN)
	time.Sleep(200 * time.Millisecond)

	lp.SendAction(1, common.Command_A_DOWN)
	time.Sleep(200 * time.Millisecond)

	lp.SendAction(1, common.Command_A_UP)
	lp.SendAction(2, common.Command_A_DOWN)

	time.Sleep(600 * time.Millisecond)
	lp.SendAction(0, common.Command_A_UP)
	lp.SendAction(2, common.Command_A_UP)

	mutex.Lock()
	assert.True(count[0] > 0 && count[1] == 0 && count[2] > 0)
	assert.True(count[0] > count[2])
	mutex.Unlock()
}

// A generator only learns a button came up from the command saying so, and a paused
// software receives none, so Reset has to end the repeat on its own.
func TestLongPressResetStopsRepeat(t *testing.T) {
	assert := assert.New(t)

	var mutex sync.Mutex
	var count int
	lp := NewLongPress(common.Button_A, 100*time.Millisecond, 100*time.Millisecond)
	lp.OnAction(func(uint64) {
		mutex.Lock()
		count++
		mutex.Unlock()
	})

	// The reset falls between two fires -- at 100, 200 and 300ms -- so that it cannot be
	// mistaken for one of them.
	lp.SendAction(0, common.Command_A_DOWN)
	time.Sleep(320 * time.Millisecond)
	lp.Reset()

	// A fire is handed to the callback on its own goroutine, so the last one before the
	// reset can still be on its way: give it time to arrive, but less than the fire delay
	// so a repeat that kept running could not have fired again.
	time.Sleep(50 * time.Millisecond)

	mutex.Lock()
	atReset := count
	mutex.Unlock()
	assert.True(atReset > 0, "the repeat should have started")

	// Long enough for three more fires.
	time.Sleep(300 * time.Millisecond)

	mutex.Lock()
	assert.Equal(atReset, count, "the repeat should have stopped at reset")
	mutex.Unlock()
}

// The same for a combination: a button left marked as down would fire the combination on
// the next press after play resumes.
func TestMultiPressReset(t *testing.T) {
	assert := assert.New(t)

	var fired bool
	mp := NewMultiPress(common.Button_L, common.Button_R)
	mp.OnAction(func(uint64) { fired = true })

	mp.SendAction(0, common.Command_L_DOWN)
	mp.Reset()
	mp.SendAction(0, common.Command_R_DOWN)

	assert.False(fired, "a reset combination must start over")
}

// Every press starts a goroutine, and on a console they come in by the thousand, so none
// of them may outlive the button.
func TestLongPressLeavesNoGoroutine(t *testing.T) {
	lp := NewLongPress(common.Button_A, 50*time.Millisecond, 50*time.Millisecond)
	lp.OnAction(func(uint64) {})

	before := runtime.NumGoroutine()

	for i := 0; i < 20; i++ {
		lp.SendAction(uint64(i), common.Command_A_DOWN)
		lp.SendAction(uint64(i), common.Command_A_UP) // released before it ever triggers
	}

	lp.SendAction(0, common.Command_A_DOWN)
	time.Sleep(150 * time.Millisecond) // long enough to be repeating
	lp.SendAction(0, common.Command_A_UP)

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= before+2 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("goroutines stranded: %d before, %d after", before, runtime.NumGoroutine())
}
