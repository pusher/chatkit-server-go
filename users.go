package chatkit

import "net/http"

// User is a type used by the CreateUser method to define a new user
type User struct {
	ID         string      `json:"id"`                    // (string| required): REQUIRED User id assigned to the user in your app.
	Name       string      `json:"name"`                  // (string| required): Name of the new user.
	AvatarURL  string      `json:"avatar_url,omitempty"`  // (string| optional): A link to the user’s photo/image.
	CustomData interface{} `json:"custom_data,omitempty"` // (object| optional): Custom data that may be associated with a user.
}

func (csc *client) CreateUser(user User) error {
	path := "/users"
	req, err := csc.newRequest(http.MethodPost, chatkitService, path, user)
	if err != nil {
		return err
	}

	return csc.do(req, nil)
}

func (csc *client) DeleteUser(userID string) error {
	path := "/users/" + userID
	req, err := csc.newRequest(http.MethodDelete, chatkitService, path, nil)
	if err != nil {
		return err
	}

	return csc.do(req, nil)
}
