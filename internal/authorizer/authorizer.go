// Package core expoeses the Authorizer API that allows making requests to the
// Chatkit Authorizer service. This allows manipulation of Roles and Permissions.
//
// All methods that are part of the interface require a JWT token with a
// `su` to be able to access them.
package authorizer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pusher/chatkit-server-go/internal/common"

	"github.com/pusher/pusher-platform-go/client"
	"github.com/pusher/pusher-platform-go/instance"
)

const (
	scopeGlobal = "global"
	scopeRoom   = "room"
)

// Exposes methods to interact with the roles and permissions API.
type Service interface {
	// Roles
	GetRoles(ctx context.Context) ([]Role, error)
	CreateGlobalRole(ctx context.Context, options CreateRoleOptions) error
	CreateRoomRole(ctx context.Context, options CreateRoleOptions) error
	DeleteGlobalRole(ctx context.Context, roleName string) error
	DeleteRoomRole(ctx context.Context, roleName string) error

	// Permissions
	GetPermissionsForGlobalRole(ctx context.Context, roleName string) ([]string, error)
	GetPermissionsForRoomRole(ctx context.Context, roleName string) ([]string, error)
	UpdatePermissionsForGlobalRole(
		ctx context.Context,
		roleName string,
		options UpdateRolePermissionsOptions,
	) error
	UpdatePermissionsForRoomRole(
		ctx context.Context,
		roleName string,
		options UpdateRolePermissionsOptions,
	) error

	// User roles
	GetUserRoles(ctx context.Context, userID string) ([]Role, error)
	AssignGlobalRoleToUser(ctx context.Context, userID string, roleName string) error
	AssignRoomRoleToUser(ctx context.Context, userID string, roomID uint, roleName string) error
	RemoveGlobalRoleForUser(ctx context.Context, userID string) error
	RemoveRoomRoleForUser(ctx context.Context, userID string, roomID uint) error

	// Generic requests
	Request(ctx context.Context, options client.RequestOptions) (*http.Response, error)
}

type authorizerService struct {
	underlyingInstance instance.Instance
}

// Returns an new authorizerService instance conforming to the Service interface.
func NewService(platformInstance instance.Instance) Service {
	return &authorizerService{
		underlyingInstance: platformInstance,
	}
}

// GetRoles fetches all roles for an instance.
func (as *authorizerService) GetRoles(ctx context.Context) ([]Role, error) {
	response, err := common.RequestWithSuToken(as.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodGet,
		Path:   "/roles",
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var roles []Role
	err = common.DecodeResponseBody(response.Body, &roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// CreateGlobalRole creates a global role for an instance.
func (as *authorizerService) CreateGlobalRole(ctx context.Context, options CreateRoleOptions) error {
	return as.createRole(ctx, Role{
		Name:        options.Name,
		Permissions: options.Permissions,
		Scope:       scopeGlobal,
	})
}

// CreateRoomRole creates a room role for an instance.
func (as *authorizerService) CreateRoomRole(ctx context.Context, options CreateRoleOptions) error {
	return as.createRole(ctx, Role{
		Name:        options.Name,
		Permissions: options.Permissions,
		Scope:       scopeRoom,
	})
}

// createRole is used by CreateGlobalRole and CreateRoomRole.
func (as *authorizerService) createRole(ctx context.Context, role Role) error {
	if role.Name == "" {
		return errors.New("You must provide a name for the role")
	}

	if role.Permissions == nil {
		return errors.New("You must provide permissions of the role")
	}

	requestBody, err := common.CreateRequestBody(&role)
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(as.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPost,
		Path:   "/roles",
		Body:   requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// DeleteGlobalRole deletes a role with the given name at the global scope.
func (as *authorizerService) DeleteGlobalRole(ctx context.Context, roleName string) error {
	return as.deleteRole(ctx, roleName, scopeGlobal)
}

// DeleteRoomRole deletes a role with the given name at the room scope.
func (as *authorizerService) DeleteRoomRole(ctx context.Context, roleName string) error {
	return as.deleteRole(ctx, roleName, scopeRoom)
}

// deleteRole is used by DeleteGlobalRole and DeleteRoomRole.
func (as *authorizerService) deleteRole(ctx context.Context, roleName string, scope string) error {
	response, err := common.RequestWithSuToken(as.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/roles/%s/scope/%s", roleName, scope),
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// GetPermissionsForGlobalRole retrieves a list of permissions associated with the role.
func (as *authorizerService) GetPermissionsForGlobalRole(
	ctx context.Context,
	roleName string,
) ([]string, error) {
	return as.getPermissions(ctx, roleName, scopeGlobal)
}

// GetPermissionsForRoomRole retrieves a list of permissions associated with the role.
func (as *authorizerService) GetPermissionsForRoomRole(
	ctx context.Context,
	roleName string,
) ([]string, error) {
	return as.getPermissions(ctx, roleName, scopeRoom)
}

// getPermissions is used by GetPermissionsForGlobalRole and GetPermissionsForRoomRole.
func (as *authorizerService) getPermissions(
	ctx context.Context,
	roleName string,
	scope string,
) ([]string, error) {
	response, err := common.RequestWithSuToken(as.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/roles/%s/scope/%s/permissions", roleName, scope),
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var rolePermissions []string
	err = common.DecodeResponseBody(response.Body, &rolePermissions)
	if err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

// UpdatePermissionsForGlobalRole allows updating permissions associated with a global role.
func (as *authorizerService) UpdatePermissionsForGlobalRole(
	ctx context.Context,
	roleName string,
	options UpdateRolePermissionsOptions,
) error {
	return as.updatePermissions(ctx, roleName, options, scopeGlobal)
}

// UpdatePermissionsForRoomRole allows updating permissions associated with a room role.
func (as *authorizerService) UpdatePermissionsForRoomRole(
	ctx context.Context,
	roleName string,
	options UpdateRolePermissionsOptions,
) error {
	return as.updatePermissions(ctx, roleName, options, scopeRoom)
}

// updatePermissions is used by UpdatePermissionsForGlobalRole and UpdatePermissionsForRoomRole.
func (as *authorizerService) updatePermissions(
	ctx context.Context,
	roleName string,
	options UpdateRolePermissionsOptions,
	scope string,
) error {
	if (options.PermissionsToAdd == nil || len(options.PermissionsToAdd) == 0) &&
		(options.PermissionsToRemove == nil || len(options.PermissionsToRemove) == 0) {
		return errors.New("PermissionsToAdd and PermissionsToRemove cannot both be empty")
	}

	requestBody, err := common.CreateRequestBody(&options)
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(as.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/roles/%s/scope/%s/permissions", roleName, scope),
		Body:   requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// GetUserRoles fetches roles associated with a user.
func (as *authorizerService) GetUserRoles(ctx context.Context, userID string) ([]Role, error) {
	if userID == "" {
		return nil, errors.New("You must provide the ID of the user whose roles you want to fetch")
	}

	response, err := common.RequestWithSuToken(as.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/users/%s/roles", userID),
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var roles []Role
	err = common.DecodeResponseBody(response.Body, &roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// AssignGlobalRoleToUser assigns a previously created global role to a user.
func (as *authorizerService) AssignGlobalRoleToUser(
	ctx context.Context,
	userID string,
	roleName string,
) error {
	return as.assignRoleToUser(ctx, userID, roleName, nil)
}

// AssignRoomRoleToUser assigns a previously created room role to a user.
func (as *authorizerService) AssignRoomRoleToUser(
	ctx context.Context,
	userID string,
	roomID uint,
	roleName string,
) error {
	return as.assignRoleToUser(ctx, userID, roleName, &roomID)
}

// assignRoleToUser is used by AssignGlobalRoleToUser and AssignRoomRoleToUser.
func (as *authorizerService) assignRoleToUser(
	ctx context.Context,
	userID string,
	roleName string,
	roomID *uint,
) error {
	if userID == "" {
		return errors.New("You must provide the ID of the user you want to assign a role to")
	}

	if roleName == "" {
		return errors.New("You must provide the role name of the role you want to assign")
	}

	requestBody, err := common.CreateRequestBody(&UserRole{Name: roleName, RoomID: roomID})
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(as.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/users/%s/roles", userID),
		Body:   requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// RemoveGlobalRoleForUser removes a role that was previously assigned to a user.
// A user can only have one global role assigned at a time.
func (as *authorizerService) RemoveGlobalRoleForUser(ctx context.Context, userID string) error {
	return as.removeRoleForUser(ctx, userID, nil)
}

// RemoveRoomRoleForUser removes a role that was previously assigned to a user.
// One user can have several room roles assigned to them (1 per room)
func (as *authorizerService) RemoveRoomRoleForUser(
	ctx context.Context,
	userID string,
	roomID uint,
) error {
	return as.removeRoleForUser(ctx, userID, &roomID)
}

// removeRole is used by RemoveGlobalRoleForUser and RemoveRoomRoleForUser
func (as *authorizerService) removeRoleForUser(
	ctx context.Context,
	userID string,
	roomID *uint,
) error {
	if userID == "" {
		return errors.New("You must provide the ID of the user you want to remove a role for")
	}

	queryParams := url.Values{}
	if roomID != nil {
		queryParams.Add("room_id", strconv.Itoa(int(*roomID)))
	}

	response, err := common.RequestWithSuToken(as.underlyingInstance, ctx, client.RequestOptions{
		Method:      http.MethodDelete,
		Path:        fmt.Sprintf("/users/%s/roles", userID),
		QueryParams: &queryParams,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// Request allows performing requests to the authorizer service that return a raw http response.
func (as *authorizerService) Request(
	ctx context.Context,
	options client.RequestOptions,
) (*http.Response, error) {
	return as.underlyingInstance.Request(ctx, options)
}
