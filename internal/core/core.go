package core

import (
	"context"

	"github.com/pusher/pusher-platform-go/instance"
)

// Exposes methods to interact with the core chatkit service.
// This allows interacting with the messages, rooms and users API.
type Service interface {
	// Users
	GetUser(ctx context.Context, userID string) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetUsersByIDs(ctx context.Context, userIDs []string) ([]User, error)
	CreateUser(ctx context.Context, options CreateUserOptions) error
	CreateUsers(ctx context.Context, users []CreateUserOptions) error
	UpdateUser(ctx context.Context, options UpdateUserOptions) error
	DeleteUser(ctx context.Context, userID string) error

	// Rooms
	GetRoom(ctx context.Context, roomID uint) (Room, error)
	GetRooms(ctx context.Context) ([]Room, error)
	GetUserRooms(ctx context.Context, userID string) ([]Room, error)
	GetUserJoinableRooms(ctx context.Context, userID string) ([]Room, error)
	CreateRoom(ctx context.Context, options CreateRoomOptions) error
	UpdateRoom(ctx context.Context, options UpdateRoomOptions) error
	DeleteRoom(ctx context.Context, roomID uint) error
	AddUsersToRoom(ctx context.Context, roomID uint, userIDs []string) error
	RemoveUsersFromRoom(ctx context.Context, roomID uint, userIds []string) error

	// Messages
	SendMessage(ctx context.Context, options CreateMessageOptions) (uint, error)
	GetRoomMessages(ctx context.Context, options GetRoomMessagesOptions) ([]Message, error)
}

type coreService struct {
	underlyingInstance instance.Instance
}

// Returns a new coreService instance that conforms to the Service interface.
func NewService(platformInstance instance.Instance) Service {
	return &coreService{
		underlyingInstance: platformInstance,
	}
}
