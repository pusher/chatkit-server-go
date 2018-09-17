package cursors

import (
	"time"
)

// Cursor represents a read cursor.
type Cursor struct {
	CursorType uint      `json:"cursor_type"`
	RoomID     uint      `json:"room_id"`
	UserID     string    `json:"user_id"`
	Position   uint      `json:"position"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// SetReadCursorOptions contains information to pass when setting a read cursor
// in a room for a message
type SetReadCursorOptions struct {
	UserID   string
	RoomID   string
	Position uint
}

// GetReadCursorOptions contains information to pass when fetching cursors for a room and user
type GetReadCursorOptions struct {
	UserID string
	RoomID uint
}
