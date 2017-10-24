package chatkitServerGo

// import (
// 	"context"

// 	_ "github.com/pusher/elements-client-go"
// )

// type ChatkitServerClientFull interface {
// 	// Chatkit Authorizer Methods
// 	GetRoles() ([]Role, error)
// 	CreateRole(Role) error
// 	DeleteRole(roleName string, scopeType string) error

// 	GetUserRoles(userID string) ([]Role, error)
// 	CreateUserRole(userID string, userRole UserRole) error
// 	UpdateUserRole(userID string, userRole UserRole) error
// 	DeleteUserRole(userID string, roomID *string) error

// 	CreateRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error
// 	GetRolePermissions(roleName string, scopeName string) (*RolePermissions, error)
// 	EditRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error

// 	// Chatkit Server Methods
// 	GetUser(userID string) (User, error)
// 	CreateUser(user User) error
// 	UpdateUser(user User) error
// 	DeleteUser(userID string) error
// 	UserGetRooms(userID string, joinable bool) ([]Room, error)
// 	UserJoinRoom(userID string, roomID string) error
// 	UserLeaveRoom(userID string, roomID string) error

// 	CreateRoom(newRoom Room) error
// 	GetRoom(roomID string, withMessages int) (Room, error)
// 	DeleteRoom(roomID string) error
// 	GetAppRooms(token string) ([]Room, error)
// 	RoomRename(roomID string, newName string) error
// 	RoomAddUsers(roomID string, userIDs []string) error
// 	RoomRemoveUsers(roomID string, userIDs []string) error
// 	SendMessage(roomID string, text string) error
// 	RoomGetMessages(roomID string, options *RoomGetMessageOptions) ([]Message, error)

// 	SendTypingIndicator(roomID string, startedTyping bool) error

// 	SubscribeUserEvents(ctx context.Context) (elementsClient.Subscription, error)
// 	SubscribeUserPresence(userID string, ctx context.Context) (elementsClient.Subscription, error)
// 	SubscribeRoomMessages(roomID string, messageLimit int, ctx context.Context) (elementsClient.Subscription, error)
// }

// type Role struct{}
// type RolePermissions struct{}

// type User struct{}
// type Room struct{}
// type Message struct{}
// type RoomGetMessageOptions struct{}

// // Creating a Role with the Room scope
// // Creating a Role with the Global scope
// // Deleting a Room scoped Role
// // Deleting a Global scoped Role
// // Assigning a Room scoped role to a User
// // Assigning a Global scoped role a User
// // Removing a Room scoped role for a User
// // Removing a Global scoped role for a User
// // List all Roles for an Instance
// // List Roles for a User
// // List Permissions associated with a Room scoped Role
// // List permissions associated with a Global scoped Role

// type ChatkitServerClient interface {
// 	// Chatkit Authorizer Methods
// 	GetRoles() ([]Role, error)
// 	CreateRole(Role) error
// 	DeleteRole(roleName string, scopeType string) error

// 	GetUserRoles(userID string) ([]Role, error)
// 	CreateUserRole(userID string, userRole UserRole) error
// 	UpdateUserRole(userID string, userRole UserRole) error
// 	DeleteUserRole(userID string, roomID *string) error

// 	CreateRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error
// 	GetRolePermissions(roleName string, scopeName string) (*RolePermissions, error)
// 	EditRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error

// 	// Chatkit Server Methods
// 	// Creating a User
// 	CreateUser(user User) error
// 	// Deleting a User
// 	DeleteUser(userID string) error
// }
