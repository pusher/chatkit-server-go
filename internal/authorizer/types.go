package authorizer

// Role represents a chatkit authorizer role.
type Role struct {
	Name        string   `json:"name"`        // Name of new role
	Permissions []string `json:"permissions"` // List of permissions for role
	Scope       string   `json:"scope"`       // Scope of the new role (global or room)
}

// CreateRoleOptions contains information to pass to the CreateRole method.
type CreateRoleOptions struct {
	Name        string
	Permissions []string
}

// AssignRoleOptions contains information to pass to the AssignRoleToUser method.
type AssignRoleOptions struct {
	CreateRoleOptions
}

// RolePermissions represents the permissions associated with a role.
type RolePermissions struct {
	Permissions []string `json:"permissions` // List of permissions associated with a role
}

// UserRole represents the type of role associated with a user.
type UserRole struct {
	Name   string `json:"name"`              // Name of the role
	RoomID *uint  `json:"room_id,omitempty"` // Optional room id. If empty, the scope is global
}

// UpdateRolePermissionsOptions contains permissions to add/remove
// permissions to/ from a role.
type UpdateRolePermissionsOptions struct {
	PermissionsToAdd    []string `json:"add_permissions,omitempty"`    // Permissions to add
	PermissionsToRemove []string `json:"remove_permissions,omitempty"` // Permissions to remove
}
