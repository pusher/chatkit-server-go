// Package chatkit is the Golang server SDK for Pusher Chatkit.
// This package provides functionality to interact with various Chatkit services.
//
// More information can be found in the Chatkit docs: https://docs.pusher.com/chatkit/overview/.
//
// Please report any bugs or feature requests at: https://github.com/pusher/chatkit-server-go.
package chatkit

import (
	"context"
	"net/http"

	"github.com/pusher/chatkit-server-go/internal/authenticator"
	"github.com/pusher/chatkit-server-go/internal/authorizer"
	"github.com/pusher/chatkit-server-go/internal/core"
	"github.com/pusher/chatkit-server-go/internal/cursors"

	"github.com/pusher/pusher-platform-go/auth"
	platformclient "github.com/pusher/pusher-platform-go/client"
	"github.com/pusher/pusher-platform-go/instance"
)

// Public interface for the library.
// It allows interacting with different Chatkit services.
type Client struct {
	coreServiceV2        core.Service
	coreServiceV6        core.Service
	authorizerService    authorizer.Service
	cursorsService       cursors.Service
	authenticatorService authenticator.Service
}

// NewClient returns an instantiated instance that fulfils the Client interface.
func NewClient(instanceLocator string, key string) (*Client, error) {
	locatorComponents, err := instance.ParseInstanceLocator(instanceLocator)
	if err != nil {
		return nil, err
	}

	keyComponents, err := instance.ParseKey(key)
	if err != nil {
		return nil, err
	}

	baseClient := platformclient.New(platformclient.Options{
		Host: locatorComponents.Host(),
	})

	coreInstanceV2, err := instance.New(instance.Options{
		Locator:        instanceLocator,
		Key:            key,
		ServiceName:    "chatkit",
		ServiceVersion: "v2",
		Client:         baseClient,
	})
	if err != nil {
		return nil, err
	}

	coreInstanceV6, err := instance.New(instance.Options{
		Locator:        instanceLocator,
		Key:            key,
		ServiceName:    "chatkit",
		ServiceVersion: "v6",
		Client:         baseClient,
	})
	if err != nil {
		return nil, err
	}

	authorizerInstance, err := instance.New(instance.Options{
		Locator:        instanceLocator,
		Key:            key,
		ServiceName:    "chatkit_authorizer",
		ServiceVersion: "v2",
		Client:         baseClient,
	})
	if err != nil {
		return nil, err
	}

	cursorsInstance, err := instance.New(instance.Options{
		Locator:        instanceLocator,
		Key:            key,
		ServiceName:    "chatkit_cursors",
		ServiceVersion: "v2",
		Client:         baseClient,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		coreServiceV2:     core.NewService(coreInstanceV2),
		coreServiceV6:     core.NewService(coreInstanceV6),
		authorizerService: authorizer.NewService(authorizerInstance),
		cursorsService:    cursors.NewService(cursorsInstance),
		authenticatorService: authenticator.NewService(
			locatorComponents.InstanceID,
			keyComponents.Key,
			keyComponents.Secret,
		),
	}, nil
}

// GetUserReadCursors returns a list of cursors that have been set across different rooms
// for the user.
func (c *Client) GetUserReadCursors(ctx context.Context, userID string) ([]Cursor, error) {
	return c.cursorsService.GetUserReadCursors(ctx, userID)
}

// SetReadCursor sets the cursor position for a room for a user.
// The position points to the message ID of a message that was sent to that room.
func (c *Client) SetReadCursor(ctx context.Context, userID string, roomID string, position uint) error {
	return c.cursorsService.SetReadCursor(ctx, userID, roomID, position)
}

// GetReadCursorsForRoom returns a list of cursors that have been set for a room.
// This returns cursors irrespective of the user that set them.
func (c *Client) GetReadCursorsForRoom(ctx context.Context, roomID string) ([]Cursor, error) {
	return c.cursorsService.GetReadCursorsForRoom(ctx, roomID)
}

// GetReadCursor returns a single cursor that was set by a user in a room.
func (c *Client) GetReadCursor(ctx context.Context, userID string, roomID string) (Cursor, error) {
	return c.cursorsService.GetReadCursor(ctx, userID, roomID)
}

// CursorsRequest allows performing a request to the cursors service that returns a raw HTTP
// response.
func (c *Client) CursorsRequest(
	ctx context.Context,
	options platformclient.RequestOptions,
) (*http.Response, error) {
	return c.cursorsService.Request(ctx, options)
}

// GetRoles retrieves all roles associated with an instance.
func (c *Client) GetRoles(ctx context.Context) ([]Role, error) {
	return c.authorizerService.GetRoles(ctx)
}

// CreateGlobalRole allows creating a globally scoped role.
func (c *Client) CreateGlobalRole(ctx context.Context, options CreateRoleOptions) error {
	return c.authorizerService.CreateGlobalRole(ctx, options)
}

// CreateRoomRole allows creating a room scoped role.
func (c *Client) CreateRoomRole(ctx context.Context, options CreateRoleOptions) error {
	return c.authorizerService.CreateRoomRole(ctx, options)
}

// DeleteGlobalRole deletes a previously created globally scoped role.
func (c *Client) DeleteGlobalRole(ctx context.Context, roleName string) error {
	return c.authorizerService.DeleteGlobalRole(ctx, roleName)
}

// DeleteRoomRole deletes a previously created room scoped role.
func (c *Client) DeleteRoomRole(ctx context.Context, roleName string) error {
	return c.authorizerService.DeleteRoomRole(ctx, roleName)
}

// GetPermissionsForGlobalRole returns permissions associated with a previously created global role.
func (c *Client) GetPermissionsForGlobalRole(
	ctx context.Context,
	roleName string,
) ([]string, error) {
	return c.authorizerService.GetPermissionsForGlobalRole(ctx, roleName)
}

// GetPermissionsForRoomRole returns permisisons associated with a previously created room role.
func (c *Client) GetPermissionsForRoomRole(
	ctx context.Context,
	roleName string,
) ([]string, error) {
	return c.authorizerService.GetPermissionsForRoomRole(ctx, roleName)
}

// UpdatePermissionsForGlobalRole allows adding or removing permissions from a previosuly created
// globally scoped role.
func (c *Client) UpdatePermissionsForGlobalRole(
	ctx context.Context,
	roleName string,
	options UpdateRolePermissionsOptions,
) error {
	return c.authorizerService.UpdatePermissionsForGlobalRole(ctx, roleName, options)
}

// UpdatePermissionsForRoomROle allows adding or removing permissions from a previously created
// room scoped role.
func (c *Client) UpdatePermissionsForRoomRole(
	ctx context.Context,
	roleName string,
	options UpdateRolePermissionsOptions,
) error {
	return c.authorizerService.UpdatePermissionsForRoomRole(ctx, roleName, options)
}

// GetUserRoles returns roles assosciated with a user.
func (c *Client) GetUserRoles(ctx context.Context, userID string) ([]Role, error) {
	return c.authorizerService.GetUserRoles(ctx, userID)
}

// AssignGlobalRoleToUser assigns a previously created globally scoped role to a user.
func (c *Client) AssignGlobalRoleToUser(ctx context.Context, userID string, roleName string) error {
	return c.authorizerService.AssignGlobalRoleToUser(ctx, userID, roleName)
}

// AssignRoomRoleToUser assigns a previously created room scoped role to a user.
func (c *Client) AssignRoomRoleToUser(
	ctx context.Context,
	userID string,
	roomID string,
	roleName string,
) error {
	return c.authorizerService.AssignRoomRoleToUser(ctx, userID, roomID, roleName)
}

// RemoveGlobalRoleForUser removes a previously assigned globally scoped role from a user.
// Users can only have one globall scoped role associated at any point.
func (c *Client) RemoveGlobalRoleForUser(ctx context.Context, userID string) error {
	return c.authorizerService.RemoveGlobalRoleForUser(ctx, userID)
}

// RemoveRoomRoleForUser removes a previously assigned room scoped role from a user.
// Users can have multiple room roles associated with them, but only one role per room.
func (c *Client) RemoveRoomRoleForUser(ctx context.Context, userID string, roomID string) error {
	return c.authorizerService.RemoveRoomRoleForUser(ctx, userID, roomID)
}

// AuthorizerRequest allows performing requests to the authorizer service
// and returns a raw HTTP response.
func (c *Client) AuthorizerRequest(
	ctx context.Context,
	options platformclient.RequestOptions,
) (*http.Response, error) {
	return c.authorizerService.Request(ctx, options)
}

// GetUser retrieves a previously created Chatkit user.
func (c *Client) GetUser(ctx context.Context, userID string) (User, error) {
	return c.coreServiceV6.GetUser(ctx, userID)
}

// GetUsers retrieves a list of users based on the options provided.
func (c *Client) GetUsers(ctx context.Context, options *GetUsersOptions) ([]User, error) {
	return c.coreServiceV6.GetUsers(ctx, options)
}

// GetUsersByID retrieves a list of users for the given id's.
func (c *Client) GetUsersByID(ctx context.Context, userIDs []string) ([]User, error) {
	return c.coreServiceV6.GetUsersByID(ctx, userIDs)
}

// CreateUser creates a new chatkit user.
func (c *Client) CreateUser(ctx context.Context, options CreateUserOptions) error {
	return c.coreServiceV6.CreateUser(ctx, options)
}

// CreateUsers creates a batch of users.
func (c *Client) CreateUsers(ctx context.Context, users []CreateUserOptions) error {
	return c.coreServiceV6.CreateUsers(ctx, users)
}

// UpdateUser allows updating a previously created user.
func (c *Client) UpdateUser(ctx context.Context, userID string, options UpdateUserOptions) error {
	return c.coreServiceV6.UpdateUser(ctx, userID, options)
}

// DeleteUser deletes a previously created user.
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	return c.coreServiceV6.DeleteUser(ctx, userID)
}

// GetRoom retrieves an existing room.
func (c *Client) GetRoom(ctx context.Context, roomID string) (Room, error) {
	return c.coreServiceV6.GetRoom(ctx, roomID)
}

// GetRooms retrieves a list of rooms based on the options provided.
func (c *Client) GetRooms(ctx context.Context, options GetRoomsOptions) ([]Room, error) {
	return c.coreServiceV6.GetRooms(ctx, options)
}

// GetUserRooms retrieves a list of rooms the user is an existing member of.
func (c *Client) GetUserRooms(ctx context.Context, userID string) ([]Room, error) {
	return c.coreServiceV6.GetUserRooms(ctx, userID)
}

// GetUserJoinableRooms retrieves a list of rooms the use can join (not an existing member of)
// Private rooms are not returned as part of the response.
func (c *Client) GetUserJoinableRooms(ctx context.Context, userID string) ([]Room, error) {
	return c.coreServiceV6.GetUserJoinableRooms(ctx, userID)
}

// CreateRoom creates a new room.
func (c *Client) CreateRoom(ctx context.Context, options CreateRoomOptions) (Room, error) {
	return c.coreServiceV6.CreateRoom(ctx, options)
}

// UpdateRoom allows updating an existing room.
func (c *Client) UpdateRoom(ctx context.Context, roomID string, options UpdateRoomOptions) error {
	return c.coreServiceV6.UpdateRoom(ctx, roomID, options)
}

// DeleteRoom deletes an existing room.
func (c *Client) DeleteRoom(ctx context.Context, roomID string) error {
	return c.coreServiceV6.DeleteRoom(ctx, roomID)
}

// AddUsersToRoom adds new users to an exising room.
func (c *Client) AddUsersToRoom(ctx context.Context, roomID string, userIDs []string) error {
	return c.coreServiceV6.AddUsersToRoom(ctx, roomID, userIDs)
}

// RemoveUsersFromRoom removes existing members from a room.
func (c *Client) RemoveUsersFromRoom(ctx context.Context, roomID string, userIDs []string) error {
	return c.coreServiceV6.RemoveUsersFromRoom(ctx, roomID, userIDs)
}

// SendMessage publishes a new message to a room.
func (c *Client) SendMessage(ctx context.Context, options SendMessageOptions) (uint, error) {
	return c.coreServiceV2.SendMessage(ctx, options)
}

// SendMultipartMessage publishes a new multipart message to a room.
func (c *Client) SendMultipartMessage(
	ctx context.Context,
	options SendMultipartMessageOptions,
) (uint, error) {
	return c.coreServiceV6.SendMultipartMessage(ctx, options)
}

// SendSimpleMessage publishes a new simple multipart message to a room.
func (c *Client) SendSimpleMessage(
	ctx context.Context,
	options SendSimpleMessageOptions,
) (uint, error) {
	return c.coreServiceV6.SendSimpleMessage(ctx, options)
}

// GetRoomMessages retrieves messages previously sent to a room based on the options provided.
func (c *Client) GetRoomMessages(
	ctx context.Context,
	roomID string,
	options GetRoomMessagesOptions,
) ([]Message, error) {
	return c.coreServiceV2.GetRoomMessages(ctx, roomID, options)
}

// FetchMultipartMessages retrieves messages previously sent to a room based on
// the options provided.
func (c *Client) FetchMultipartMessages(
	ctx context.Context,
	roomID string,
	options GetRoomMessagesOptions,
) ([]MultipartMessage, error) {
	return c.coreServiceV6.FetchMultipartMessages(ctx, roomID, options)
}

// DeleteMessage allows a previously sent message to be deleted.
func (c *Client) DeleteMessage(ctx context.Context, options DeleteMessageOptions) error {
	return c.coreServiceV6.DeleteMessage(ctx, options)
}

// CoreRequest allows making requests to the core chatkit service and returns a raw HTTP response.
func (c *Client) CoreRequest(
	ctx context.Context,
	options platformclient.RequestOptions,
) (*http.Response, error) {
	return c.coreServiceV6.Request(ctx, options)
}

// Authenticate returns a token response along with headers and status code to be used within
// the context of a token provider.
// Currently, the only supported GrantType is GrantTypeClientCredentials.
func (c *Client) Authenticate(payload auth.Payload, options auth.Options) (*auth.Response, error) {
	return c.authenticatorService.Authenticate(payload, options)
}

// GenerateAccessToken generates a JWT token based on the options provided.
func (c *Client) GenerateAccessToken(options auth.Options) (auth.TokenWithExpiry, error) {
	return c.authenticatorService.GenerateAccessToken(options)
}

// GenerateSuToken generates a JWT token with the `su` claim.
func (c *Client) GenerateSUToken(options auth.Options) (auth.TokenWithExpiry, error) {
	return c.authenticatorService.GenerateSUToken(options)
}
