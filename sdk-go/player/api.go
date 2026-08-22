package player

import (
	"github.com/richardlt/matrix/internal/errors"
	common "github.com/richardlt/matrix/sdk-go/common"
)

// ErrAPIClosed is returned once the API can no longer reach the matrix core.
// It is a sentinel so callers can match it with errors.Is.
var ErrAPIClosed = errors.New("player api is closed")

// API allows the player to send events to the matrix core.
type API struct{ requestChannel chan *Request }

// Command send a command to the matrix core.
func (a *API) Command(slot uint64, command common.Command) error {
	if a.requestChannel == nil {
		return ErrAPIClosed
	}

	a.requestChannel <- &Request{
		Type: Request_PLAYER,
		PlayerData: &Request_PlayerData{
			Action:  Request_PlayerData_COMMAND,
			Slot:    slot,
			Command: command,
		},
	}

	return nil
}

// Close will make the API unusable for any action.
func (a *API) Close() { a.requestChannel = nil }
