package cursors

import (
	"github.com/pusher/pusher-platform-go/instance"
)

// Exposes methods to interact with the cursors service
type Service interface {
	GetUserReadCursors(userID string) ([]Cursor, error)
	SetReadCursor(options SetReadCursorOptions) error
	GetReadCursorsForRoom(roomID uint) ([]Cursor, error)
	GetReadCursor(options GetReadCursorOptions) (Cursor, error)
}

type cursorsService struct {
	underlyingInstance instance.Instance
}

// Returns a new cursorsService instance conforming to
// the Service interface
func NewService(platformInstance instance.Instance) Service {
	return &cursorsService{
		underlyingInstance: platformInstance,
	}
}
