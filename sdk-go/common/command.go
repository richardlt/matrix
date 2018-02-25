package common

import (
	"strings"
)

// CommandToButtonState returns button from given command with press state.
func CommandToButtonState(cmd Command) (button Button, pressed bool) {
	t := strings.Split(cmd.String(), "_")
	if v, ok := Button_value[t[0]]; ok {
		button = Button(v)
		pressed = t[1] == "DOWN"
	}
	return
}
