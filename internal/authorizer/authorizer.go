package authorizer

import (
	"context"

	"github.com/pusher/pusher-platform-go/instance"
)

// Exposes methods to interact with the roles and permissions API.
type Service interface {
	// Roles
	GetRoles(ctx context.Context) ([]Role, error)
	CreateGlobalRole(ctx context.Context, role Role) error
	CreateRoomRole(ctx context.Context, role Role) error
	DeleteGlobalRole(ctx context.Context, roleName string) error
	DeleteRoomRole(ctx context.Context, roleName string) error

	// Permissions
	GetPermissionsForGlobalRole(ctx context.Context, roleName string) (RolePermissions, error)
	GetPermissionsForRoomRole(ctx context.Context, roleName string) (RolePermissions, error)
	UpdatePermissionsForGlobalRole(
		ctx context.Context,
		roleName string, params UpdateRolePermissionsOptions,
	) error
	UpdatePermissionsForRoomRole(
		ctx context.Context,
		roleName string,
		params UpdateRolePermissionsOptions,
	) error

	// User roles
	GetUserRoles(ctx context.Context, userID string) ([]Role, error)
	AssignGlobalRoleToUser(ctx context.Context, userID string, role UserRole) error
	AssignRoomRoleToUser(ctx context.Context, userID string, role UserRole) error
	RemoveGlobalRoleForUser(ctx context.Context, userID string) error
	RemoveRoomRoleForUser(ctx context.Context, userID string, roomID string) error
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
