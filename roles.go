package chatkitServerGo

import (
	"errors"
	"fmt"
	"net/http"
)

// Role is a type used by the role methods to define a role
type Role struct {
	Name        string   `json:"name"`        // (string| required): Name of the new role.
	Permissions []string `json:"permissions"` // (array| required): Permissions assigned to the role.
	Scope       string   `json:"scope"`       // (string| required): Scope of the new role; one of global or room.
}

func (csc *chatkitServerClient) CreateRole(role Role) error {
	req, err := csc.newRequest(http.MethodPost, chatkitAuthService, "/roles", role)
	if err != nil {
		return err
	}

	return csc.do(req, nil)
}

func (csc *chatkitServerClient) GetRoles() ([]Role, error) {
	req, err := csc.newRequest(http.MethodGet, chatkitAuthService, "/roles", nil)
	if err != nil {
		return nil, err
	}

	roles := &[]Role{}
	err = csc.do(req, roles)

	return *roles, err
}

func (csc *chatkitServerClient) DeleteRole(roleName string, scopeType string) error {
	if roleName == "" || scopeType == "" {
		return errors.New("Role name and Scope type cannot be empty strings")
	}

	path := fmt.Sprint("/roles/", roleName, "/scope/", scopeType)
	req, err := csc.newRequest(http.MethodDelete, chatkitAuthService, path, nil)
	if err != nil {
		return err
	}

	return csc.do(req, nil)
}
