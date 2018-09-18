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
