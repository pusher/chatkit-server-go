# chatkit-server-go

Golang server SDK for Pusher Chatkit.

## Interface

    // ChatkitServerClient is the public interface of the Chatkit Server Client
    type ChatkitServerClient interface {
        // Chatkit Roles and Permissions methods
        GetRoles() ([]Role, error)
        CreateRole(Role) error
        DeleteRole(roleName string, scopeType string) error

        GetUserRoles(userID string) ([]Role, error)
        CreateUserRole(userID string, userRole UserRole) error
        UpdateUserRole(userID string, userRole UserRole) error
        DeleteUserRole(userID string, roomID *string) error

        CreateRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error
        GetRolePermissions(roleName string, scopeName string) (*RolePermissions, error)
        EditRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error

        // Chatkit Server methods
        CreateUser(user User) error
        DeleteUser(userID string) error
    }

    // NewChatkitServerClient instantiates a new ChatkitServerClient instance 
    func NewChatkitServerClient(instanceID string, key string) (ChatkitServerClient, error)

    // Also included are some useful Chatkit jwt token generation helper functions
    func NewChatkitUserToken(appID string, keyID string, keySecret string, userID string, expiryDuration time.Duration) (tokenString string, expiry time.Time, err error)
    func NewChatkitSUToken(appID string, keyID string, keySecret string, expiryDuration time.Duration) (tokenString string, expiry time.Time, err error)

## Installation

    $ go get chatkit-server-go

## Tests

    $ go test -v -cover

## Examples

Please refer to the examples directory.

## Documentation

Available in the [Pusher Docs](https://docs.pusher.com/chatkit/overview/).

# Future Work

the ideal and compleate interface would look like this:

    type ChatkitServerClientFull interface {
        // Chatkit Authorizer Methods
        GetRoles() ([]Role, error)
        CreateRole(Role) error
        DeleteRole(roleName string, scopeType string) error

        GetUserRoles(userID string) ([]Role, error)
        CreateUserRole(userID string, userRole UserRole) error
        UpdateUserRole(userID string, userRole UserRole) error
        DeleteUserRole(userID string, roomID *string) error

        CreateRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error
        GetRolePermissions(roleName string, scopeName string) (*RolePermissions, error)
        EditRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error

        // Chatkit Server Methods
        GetUser(userID string) (User, error)
        CreateUser(user User) error
        UpdateUser(user User) error
        DeleteUser(userID string) error

        UserGetRooms(userID string, joinable bool) ([]Room, error)
        UserJoinRoom(userID string, roomID string) error
        UserLeaveRoom(userID string, roomID string) error

        CreateRoom(newRoom Room) error
        GetRoom(roomID string, withMessages int) (Room, error)
        DeleteRoom(roomID string) error
        GetAppRooms(token string) ([]Room, error)
        RoomRename(roomID string, newName string) error
        RoomAddUsers(roomID string, userIDs []string) error
        RoomRemoveUsers(roomID string, userIDs []string) error
        SendMessage(roomID string, text string) error
        RoomGetMessages(roomID string, options *RoomGetMessageOptions) ([]Message, error)

        SendTypingIndicator(roomID string, startedTyping bool) error

        SubscribeUserEvents(ctx context.Context) (elementsClient.Subscription, error)
        SubscribeUserPresence(userID string, ctx context.Context) (elementsClient.Subscription, error)
        SubscribeRoomMessages(roomID string, messageLimit int, ctx context.Context) (elementsClient.Subscription, error)
    }