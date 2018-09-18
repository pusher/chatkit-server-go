package cursors

import (
	"context"
	"errors"

	"github.com/pusher/pusher-platform-go/instance"
)

// Exposes methods to interact with the cursors service
type Service interface {
	GetUserReadCursors(ctx context.Context, userID string) ([]Cursor, error)
	SetReadCursor(ctx context.Context, options SetReadCursorOptions) error
	GetReadCursorsForRoom(ctx context.Context, roomID uint) ([]Cursor, error)
	GetReadCursor(ctx context.Context, options GetReadCursorOptions) (Cursor, error)
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
