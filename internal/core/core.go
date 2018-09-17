package core

import (
	"github.com/pusher/pusher-platform-go/instance"
)

// Exposes methods to interact with the core chatkit service.
// This allows interacting with the messages, rooms and users API.
type Service interface {
	// Users
	GetUser(userID string) (User, error)
	GetUsers() ([]User, error)
	GetUsersByIDs(userIDs []string) ([]User, error)
	CreateUser(options CreateUserOptions) error
	CreateUsers(users []CreateUserOptions) error
	UpdateUser(options UpdateUserOptions) error
	DeleteUser(userID string) error

	// Rooms
	GetRoom(roomID uint) (Room, error)
	GetRooms() ([]Room, error)
	GetUserRooms(userID string) ([]Room, error)
	GetUserJoinableRooms(userID string) ([]Room, error)
	CreateRoom(options CreateRoomOptions) error
	UpdateRoom(options UpdateRoomOptions) error
	DeleteRoom(roomID uint) error
	AddUsersToRoom(roomID uint, userIDs []string) error
	RemoveUsersFromRoom(roomID uint, userIds []string) error

	// Messages
	SendMessage(options CreateMessageOptions) (uint, error)
	GetRoomMessages(options GetRoomMessagesOptions) ([]Message, error)
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
