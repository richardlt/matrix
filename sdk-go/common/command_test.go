package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandToButtonState(t *testing.T) {
	assert := assert.New(t)

	button, pressed := CommandToButtonState(Command_A_UP)
	assert.Equal(Button_A, button)
	assert.False(pressed)

	button, pressed = CommandToButtonState(Command_B_DOWN)
	assert.Equal(Button_B, button)
	assert.True(pressed)
}
