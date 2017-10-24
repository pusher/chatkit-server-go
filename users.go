package chatkitServerGo

// User is a type used by the CreateUser method to create a new user
type User struct {
	ID         string      `json:"id"`                     // (string| required): REQUIRED User id assigned to the user in your app.
	Name       string      `json:"name"`                   // (string| required): Name of the new user.
	AvatarURL  string      `json:"avatar_url,omitempty"`   // (string| optional): A link to the user’s photo/image.
	CustomData interface{} `json:"custom_datas,omitempty"` // (object| optional): Custom data that may be associated with a user.
}

func (csc *chatkitServerClient) CreateUser(user User) error {
	return nil
}

func (csc *chatkitServerClient) DeleteUser(userID string) error {
	return nil
}
