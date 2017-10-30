# chatkit-server-go [![GoDoc](https://godoc.org/github.com/pusher/chatkit-server-go?status.svg)](http://godoc.org/github.com/pusher/chatkit-server-go)

Package chatkitServerGo is the Golang server SDK for [Pusher Chatkit](https://pusher.com/chatkit).

This package provides the ChatkitServerClient type for managing Chatkit users and
interacting with roles and permissions of those users. It also contains some helper
functions for creating your own JWT tokens for authentication with the Chatkit
service.

Please report any bugs or feature requests via a Github issue on this repo.

## Interface

```go
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

    // Chatkit User methods
    CreateUser(user User) error
    DeleteUser(userID string) error
}

// NewChatkitServerClient instantiates a new ChatkitServerClient instance 
func NewChatkitServerClient(instanceLocator string, key string) (ChatkitServerClient, error)

// NewChatkitUserToken and NewChatkitSUToken are Chatkit JWT token generation helper functions
func NewChatkitUserToken(appID string, keyID string, keySecret string, userID string, expiryDuration time.Duration) (tokenString string, expiry time.Time, err error)
func NewChatkitSUToken(appID string, keyID string, keySecret string, expiryDuration time.Duration) (tokenString string, expiry time.Time, err error)
```

## Installation

    $ go get pusher/chatkit-server-go

## Getting Started

Please refer to the /example directory.

## Tests

    $ go test -v -cover

## Documentation

Available in the [Pusher Docs](https://docs.pusher.com/chatkit/overview/).

## License

This code is free to use under the terms of the MIT license. Please refer to LICENSE.md for more information.