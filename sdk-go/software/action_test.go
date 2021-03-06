package software

import (
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
