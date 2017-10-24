package chatkitServerGo

import (
	"fmt"
	"net/http"
)

// RolePermissions is a type used to create, get and edit role permissions
type RolePermissions struct {
	Permissions []string `json:"permissions"` // (array| required): The permissions that you want to attach to the specified role (at the specified scope).
}

func (csc *chatkitServerClient) CreateRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error {
	path := fmt.Sprint("/roles/", roleName, "/scope/", scopeName, "/permissions")
	req, err := csc.newRequest(http.MethodPost, CHATKIT_AUTH, path, rolePerms)
	if err != nil {
		return err
	}

	return csc.do(req, nil)
}

func (csc *chatkitServerClient) GetRolePermissions(roleName string, scopeName string) (*RolePermissions, error) {
	path := fmt.Sprint("/roles/", roleName, "/scope/", scopeName, "/permissions")
	req, err := csc.newRequest(http.MethodGet, CHATKIT_AUTH, path, nil)
	if err != nil {
		return nil, err
	}

	stringSlice := &[]string{}
	err = csc.do(req, stringSlice)

	return &RolePermissions{*stringSlice}, err
}

func (csc *chatkitServerClient) EditRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error {
	path := fmt.Sprint("/roles/", roleName, "/scope/", scopeName, "/permissions")
	req, err := csc.newRequest(http.MethodPut, CHATKIT_AUTH, path, rolePerms)
	if err != nil {
		return err
	}

	return csc.do(req, nil)
}
