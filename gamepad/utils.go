package gamepad

import "github.com/richardlt/matrix/sdk-go/common"

func commandFromString(cmd string) common.Command {
	switch cmd {
	case "left_up":
		return common.Command_LEFT_UP
	case "left_down":
		return common.Command_LEFT_DOWN
	case "up_up":
		return common.Command_UP_UP
	case "up_down":
		return common.Command_UP_DOWN
	case "right_up":
		return common.Command_RIGHT_UP
	case "right_down":
		return common.Command_RIGHT_DOWN
	case "down_up":
		return common.Command_DOWN_UP
	case "down_down":
		return common.Command_DOWN_DOWN
	case "x_up":
		return common.Command_X_UP
	case "x_down":
		return common.Command_X_DOWN
	case "y_up":
		return common.Command_Y_UP
	case "y_down":
		return common.Command_Y_DOWN
	case "a_up":
		return common.Command_A_UP
	case "a_down":
		return common.Command_A_DOWN
	case "b_up":
		return common.Command_B_UP
	case "b_down":
		return common.Command_B_DOWN
	case "l_up":
		return common.Command_L_UP
	case "l_down":
		return common.Command_L_DOWN
	case "r_up":
		return common.Command_R_UP
	case "r_down":
		return common.Command_R_DOWN
	case "select_up":
		return common.Command_SELECT_UP
	case "select_down":
		return common.Command_SELECT_DOWN
	case "start_up":
		return common.Command_START_UP
	case "start_down":
		return common.Command_START_DOWN
	}
	return 0
}
