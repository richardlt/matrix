package device

import "github.com/richardlt/matrix/sdk-go/common"

func commandFromString(cmd string) common.Command {
	switch cmd {
	case "left:release":
		return common.Command_LEFT_UP
	case "left:press":
		return common.Command_LEFT_DOWN
	case "up:release":
		return common.Command_UP_UP
	case "up:press":
		return common.Command_UP_DOWN
	case "right:release":
		return common.Command_RIGHT_UP
	case "right:press":
		return common.Command_RIGHT_DOWN
	case "down:release":
		return common.Command_DOWN_UP
	case "down:press":
		return common.Command_DOWN_DOWN
	case "x:release":
		return common.Command_X_UP
	case "x:press":
		return common.Command_X_DOWN
	case "y:release":
		return common.Command_Y_UP
	case "y:press":
		return common.Command_Y_DOWN
	case "a:release":
		return common.Command_A_UP
	case "a:press":
		return common.Command_A_DOWN
	case "b:release":
		return common.Command_B_UP
	case "b:press":
		return common.Command_B_DOWN
	case "l:release":
		return common.Command_L_UP
	case "l:press":
		return common.Command_L_DOWN
	case "r:release":
		return common.Command_R_UP
	case "r:press":
		return common.Command_R_DOWN
	case "select:release":
		return common.Command_SELECT_UP
	case "select:press":
		return common.Command_SELECT_DOWN
	case "start:release":
		return common.Command_START_UP
	case "start:press":
		return common.Command_START_DOWN
	}
	return 0
}
