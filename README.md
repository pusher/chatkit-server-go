# chatkit-server-go [![GoDoc](https://godoc.org/github.com/pusher/chatkit-server-go?status.svg)](http://godoc.org/github.com/pusher/chatkit-server-go)

package chatkit is the Golang server SDK for [Pusher Chatkit](https://pusher.com/chatkit).

This package provides the Client type for managing Chatkit users and
interacting with roles and permissions of those users. It also contains some helper
functions for creating your own JWT tokens for authentication with the Chatkit
service.

Please report any bugs or feature requests via a GitHub issue on this repo.

## Interface

```go
// Client is the public interface of the Chatkit Server Client
type Client interface {
    // Authentication method
    Authenticate(userID string) AuthenticationResponse

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

    // Chatkit User methods
    CreateUser(user User) error
    DeleteUser(userID string) error
    GetUsers() ([]User, error)
}

// NewClient instantiates a new Client instance
func NewClient(instanceLocator string, key string) (Client, error)

// NewChatkitToken is a Chatkit JWT token generation helper function
func NewChatkitToken(instanceID string, keyID string, keySecret string, userID *string, su bool, expiryDuration time.Duration) (tokenBody *TokenBody, errorBody *ErrorBody)
```

## Installation

    $ go get github.com/pusher/chatkit-server-go

## Getting Started

Please refer to the [`/example`](https://github.com/pusher/chatkit-server-go/tree/master/example) directory for examples that relate to managing users, roles, and permisisons.

Check the [`/auth_example`](https://github.com/pusher/chatkit-server-go/tree/master/auth_example) directory for an example that shows authentication.

## Tests

    $ go test -v -cover

## Documentation

Available in the [Pusher Docs](https://docs.pusher.com/chatkit).

## License

This code is free to use under the terms of the MIT license. Please refer to LICENSE.md for more information.
