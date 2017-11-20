package chatkit

import (
	"errors"
	"fmt"
	"net/http"
)

// UserRole is a type used by the CreateUserRole method to define a user role
type UserRole struct {
	Name   string `json:"name"`              // (string| required): Name of the role. If no room_id value is provided then the scope will be global.
	RoomID int    `json:"room_id,omitempty"` // (integer| optional): The ID of the room you want to create the user role for.
}

func (csc *client) GetUserRoles(userID string) ([]Role, error) {
	if userID == "" {
		return nil, errors.New("User ID cannot be an empty string")
	}

	path := fmt.Sprint("/users/", userID, "/roles")
	req, err := csc.newRequest(http.MethodGet, chatkitAuthService, path, nil)
	if err != nil {
		return nil, err
	}

	roles := &[]Role{}
	err = csc.do(req, roles)

	return *roles, err
}

func (csc *client) SetUserRole(userID string, userRole UserRole) error {
	path := fmt.Sprint("/users/", userID, "/roles")
	req, err := csc.newRequest(
		http.MethodPut,
		chatkitAuthService,
		path,
		userRole,
	)
	if err != nil {
		return err
	}

	return csc.do(req, nil)
}

func (csc *client) DeleteUserRole(userID string, roomID *string) error {
	path := fmt.Sprint("/users/", userID, "/roles")
	if roomID != nil {
		path = path + "?room_id=" + *roomID
	}
	req, err := csc.newRequest(http.MethodDelete, chatkitAuthService, path, nil)
	if err != nil {
		return err
	}

	return csc.do(req, nil)
}
