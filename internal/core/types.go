package core

import (
	"io"
	"time"
)

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
	RoomWithoutMembers
	MemberUserIDs []string `json:"member_user_ids,omitempty"` // List of user id's in the room
}

// RoomWithoutMembers represents a chatkit room without listing its members.
type RoomWithoutMembers struct {
	ID                            string      `json:"id"`                                         // ID assigned to a room
	CreatedByID                   string      `json:"created_by_id"`                              // User ID that created the room
	Name                          string      `json:"name"`                                       // Name assigned to the room
	PushNotificationTitleOverride *string     `json:"push_notification_title_override,omitempty"` // Optionally override Push Notification title
	Private                       bool        `json:"private"`                                    // Indicates if room is private or not
	CustomData                    interface{} `json:"custom_data,omitempty"`                      // Custom data that can be added to rooms
	CreatedAt                     time.Time   `json:"created_at"`                                 // Creation timestamp
	UpdatedAt                     time.Time   `json:"updated_at"`                                 // Updation timestamp
}

type messageIsh interface {
	isMessageIsh()
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

func (Message) isMessageIsh() {}

// MultipartMessage represents a message sent to a chatkit room.
type MultipartMessage struct {
	ID        uint      `json:"id"`         // Message ID
	UserID    string    `json:"user_id"`    // User that sent the message
	RoomID    string    `json:"room_id"`    // Room the message was sent to
	Parts     []Part    `json:"parts"`      // Parts composing the message
	CreatedAt time.Time `json:"created_at"` // Creation timestamp
	UpdatedAt time.Time `json:"updated_at"` // Updation timestamp
}

func (MultipartMessage) isMessageIsh() {}

type Part struct {
	Type       string      `json:"type"`
	Content    *string     `json:"content,omitempty"`
	URL        *string     `json:"url,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

type Attachment struct {
	ID          string      `json:"id"`
	DownloadURL string      `json:"download_url"`
	RefreshURL  string      `json:"refresh_url"`
	Expiration  time.Time   `json:"expiration"`
	Name        string      `json:"name"`
	CustomData  interface{} `json:"custom_data,omitempty"`
	Size        uint        `json:"size"`
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
	ID                            *string     `json:"id,omitempty"`
	Name                          string      `json:"name"`
	PushNotificationTitleOverride *string     `json:"push_notification_title_override,omitempty"`
	Private                       bool        `json:"private"`
	UserIDs                       []string    `json:"user_ids,omitempty"` // User ID's to be added to the room during creation
	CustomData                    interface{} `json:"custom_data,omitempty"`
	CreatorID                     string
}

// GetRoomsOptions contains parameters to pass to fetch rooms.
type GetRoomsOptions struct {
	FromID         *string `json:"from_id,omitempty"`
	IncludePrivate bool    `json:"include_private"`
}

// UpdateRoomOptions contains parameters to pass when updating a room.
type UpdateRoomOptions struct {
	Name                          *string     `json:"name,omitempty"`
	PushNotificationTitleOverride *string     `json:"push_notification_title_override,omitempty"`
	Private                       *bool       `json:"private,omitempty"`
	CustomData                    interface{} `json:"custom_data,omitempty"`
}

// ExplicitlyResetPushNotificationTitleOverride when used in the UpdateRoomOptions
// signifies that the override is to be removed entirely
var ExplicitlyResetPushNotificationTitleOverride = "null"

// SendMessageOptions contains parameters to pass when sending a new message.
type SendMessageOptions = SendSimpleMessageOptions

// SendMultipartMessageOptions contains parameters to pass when sending a new message.
type SendMultipartMessageOptions struct {
	RoomID   string
	SenderID string
	Parts    []NewPart
}

// SendSimpleMessageOptions contains parameters to pass when sending a new message.
type SendSimpleMessageOptions struct {
	RoomID   string
	Text     string
	SenderID string
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

// NewAttachmentPart has no JSON annotations because it cannot be sent directly
// to the backend. The attachment must first be uploaded and a
// newAttachmentPartUploaded sent instead.
type NewAttachmentPart struct {
	Type       string
	Name       *string
	CustomData interface{}
	File       io.Reader
}

func (p NewAttachmentPart) isNewPart() {}

type newAttachmentPartUploaded struct {
	Type       string             `json:"type"`
	Attachment uploadedAttachment `json:"attachment"`
}

type uploadedAttachment struct {
	ID string `json:"id"`
}

type fetchMessagesOptions struct {
	InitialID *uint   // Starting ID of messages to retrieve
	Direction *string // One of older or newer
	Limit     *uint   // Number of messages to retrieve
}

// FetchMultipartMessagesOptions contains parameters to pass when fetching messages from a room.
type FetchMultipartMessagesOptions = fetchMessagesOptions

// GetRoomMessagesOptions contains parameters to pass when fetching messages from a room.
type GetRoomMessagesOptions = fetchMessagesOptions

type DeleteMessageOptions struct {
	RoomID    string
	MessageID uint
}
