package authorizer

import (
	"github.com/pusher/pusher-platform-go/instance"
)

// Exposes methods to interact with the roles and permissions API.
type Service interface {
	// Roles
	GetRoles() ([]Role, error)
	CreateGlobalRole(Role) error
	CreateRoomRole(Role) error
	DeleteGlobalRole(roleName string) error
	DeleteRoomRole(roleName string) error

	// Permissions
	GetPermissionsForGlobalRole(roleName string) (RolePermissions, error)
	GetPermissionsForRoomRole(roleName string) (RolePermissions, error)
	UpdatePermissionsForGlobalRole(roleName string, params UpdateRolePermissionsOptions) error
	UpdatePermissionsForRoomRole(roleName string, params UpdateRolePermissionsOptions) error

	// User roles
	GetUserRoles(userID string) ([]Role, error)
	AssignGlobalRoleToUser(userID string, role UserRole) error
	AssignRoomRoleToUser(userID string, role UserRole) error
	RemoveGlobalRoleForUser(userID string) error
	RemoveRoomRoleForUser(userID string, roomID string) error
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
