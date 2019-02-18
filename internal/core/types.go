package core

import "time"

// User represents a chatkit user.
type User struct {
	ID         string                 `json:"id"`                    // ID of the user
	Name       string                 `json:"name"`                  // Name associated with the user
	AvatarURL  string                 `json:"avatar_url,omitempty"`  // Link to a photo/ image of the user
	CustomData map[string]interface{} `json:"custom_data,omitempty"` // A custom data object associated with the user
	CreatedAt  time.Time              `json:"created_at"`            // Creation timestamp
	UpdatedAt  time.Time              `json:"updated_at"`            // Updating timestamp
}

// Room represents a chatkit room.
type Room struct {
	ID            string      `json:"id"`                        // ID assigned to a room
	CreatedByID   string      `json:"created_by_id"`             // User ID that created the room
	Name          string      `json:"name"`                      // Name assigned to the room
	Private       bool        `json:"private"`                   // Indicates if room is private or not
	MemberUserIDs []string    `json:"member_user_ids,omitempty"` // List of user id's in the room
	CustomData    interface{} `json:"custom_data,omitempty"`     // Custom data that can be added to rooms
	CreatedAt     time.Time   `json:"created_at"`                // Creation timestamp
	UpdatedAt     time.Time   `json:"updated_at"`                // Updation timestamp
}

// Message represents a message sent to a chatkit room.
type Message struct {
	ID        uint      `json:"id"`         // Message ID
	UserID    string    `json:"user_id"`    // User that sent the message
	RoomID    string    `json:"room_id"`    // Room the message was sent to
	Text      string    `json:"text"`       // Content of the message
	CreatedAt time.Time `json:"created_at"` // Creation timestamp
	UpdatedAt time.Time `json:"updated_at"` // Updation timestamp
}

// GetUsersOptions contains parameters to pass when fetching users.
type GetUsersOptions struct {
	FromTimestamp string
	Limit         uint
}

// CreateUserOptions contains parameters to pass when creating a new user.
type CreateUserOptions struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	AvatarURL  *string     `json:"avatar_url,omitempty"`
	CustomData interface{} `json:"custom_data,omitempty"`
}

// UpdateUserOptions contains parameters to pass when updating a user.
type UpdateUserOptions struct {
	Name       *string     `json:"name,omitempty"`
	AvatarUrl  *string     `json:"avatar_url,omitempty"`
	CustomData interface{} `json:"custom_data,omitempty"`
}

// CreateRoomOptions contains parameters to pass when creating a new room.
type CreateRoomOptions struct {
	Name       string      `json:"name"`
	Private    bool        `json:"private"`
	UserIDs    []string    `json:"user_ids,omitempty"` // User ID's to be added to the room during creation
	CustomData interface{} `json:"custom_data,omitempty"`
	CreatorID  string
}

// GetRoomsOptions contains parameters to pass to fetch rooms.
type GetRoomsOptions struct {
	FromID         *string `json:"from_id,omitempty"`
	IncludePrivate bool    `json:"include_private"`
}

// UpdateRoomOptions contains parameters to pass when updating a room.
type UpdateRoomOptions struct {
	Name       *string     `json:"name,omitempty"`
	Private    *bool       `json:"private,omitempty"`
	CustomData interface{} `json:"custom_data,omitempty"`
}

// SendMessageOptions contains parameters to pass when sending a new message.
type SendMessageOptions struct {
	RoomID   string
	Text     string
	SenderID string
}

// SendMultipartMessageOptions contains parameters to pass when sending a new message.
type SendMultipartMessageOptions struct {
	RoomID   string
	SenderID string
	Parts    []NewPart
}

type NewPart interface {
	isNewPart()
}

type NewInlinePart struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (p NewInlinePart) isNewPart() {}

type NewURLPart struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

func (p NewURLPart) isNewPart() {}

// GetRoomMessagesOptions contains parameters to pass when fetching messages from a room.
type GetRoomMessagesOptions struct {
	InitialID *uint   // Starting ID of messages to retrieve
	Direction *string // One of older or newer
	Limit     *uint   // Number of messages to retrieve
}
