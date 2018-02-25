package software

import (
	"testing"

	"github.com/richardlt/matrix/sdk-go/common"

	"github.com/stretchr/testify/assert"
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
