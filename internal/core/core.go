// Package core exposes an interface that allows making requests to the core
// Chatkit API to allow operations to be performed against Users, Rooms and Messages.
package core

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

// Exposes methods to interact with the core chatkit service.
// This allows interacting with the messages, rooms and users API.
type Service interface {
	// Users
	GetUser(ctx context.Context, userID string) (User, error)
	GetUsers(ctx context.Context, options *GetUsersOptions) ([]User, error)
	GetUsersByID(ctx context.Context, userIDs []string) ([]User, error)
	CreateUser(ctx context.Context, options CreateUserOptions) error
	CreateUsers(ctx context.Context, users []CreateUserOptions) error
	UpdateUser(ctx context.Context, userID string, options UpdateUserOptions) error
	DeleteUser(ctx context.Context, userID string) error

	// Rooms
	GetRoom(ctx context.Context, roomID uint) (Room, error)
	GetRooms(ctx context.Context, options GetRoomsOptions) ([]Room, error)
	GetUserRooms(ctx context.Context, userID string) ([]Room, error)
	GetUserJoinableRooms(ctx context.Context, userID string) ([]Room, error)
	CreateRoom(ctx context.Context, options CreateRoomOptions) (Room, error)
	UpdateRoom(ctx context.Context, roomID uint, options UpdateRoomOptions) error
	DeleteRoom(ctx context.Context, roomID uint) error
	AddUsersToRoom(ctx context.Context, roomID uint, userIDs []string) error
	RemoveUsersFromRoom(ctx context.Context, roomID uint, userIds []string) error

	// Messages
	SendMessage(ctx context.Context, options SendMessageOptions) (uint, error)
	GetRoomMessages(ctx context.Context, roomID uint, options GetRoomMessagesOptions) ([]Message, error)
	DeleteMessage(ctx context.Context, messageID uint) error

	// Generic requests
	Request(ctx context.Context, options client.RequestOptions) (*http.Response, error)
}

type coreService struct {
	underlyingInstance instance.Instance
}

// Returns a new coreService instance that conforms to the Service interface.
func NewService(platformInstance instance.Instance) Service {
	return &coreService{
		underlyingInstance: platformInstance,
	}
}

// GetUser retrieves a user for the given user id.
func (cs *coreService) GetUser(ctx context.Context, userID string) (User, error) {
	if userID == "" {
		return User{}, errors.New("You must provide the ID of the user you want to fetch")
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/users/%s", userID),
	})
	if err != nil {
		return User{}, err
	}
	defer response.Body.Close()

	var user User
	err = common.DecodeResponseBody(response.Body, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetUsers retrieves a batch of users depending on the optionally passed in parameters.
// If not options are passed in, the server will return the default limit of 20 users.
func (cs *coreService) GetUsers(ctx context.Context, options *GetUsersOptions) ([]User, error) {
	queryParams := url.Values{}
	if options != nil {
		queryParams.Add("from_ts", options.FromTimestamp)
		queryParams.Add("limit", strconv.Itoa(int(options.Limit)))
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method:      http.MethodGet,
		Path:        "/users",
		QueryParams: &queryParams,
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var users []User
	err = common.DecodeResponseBody(response.Body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersByID returns a list of users whose ID's are supplied.
func (cs *coreService) GetUsersByID(ctx context.Context, userIDs []string) ([]User, error) {
	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodGet,
		Path:   "/users_by_ids",
		QueryParams: &url.Values{
			"id": userIDs,
		},
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var users []User
	err = common.DecodeResponseBody(response.Body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// CreateUser creates a new chatkit user based on the provided options.
func (cs *coreService) CreateUser(ctx context.Context, options CreateUserOptions) error {
	if options.ID == "" {
		return errors.New("You must provide the ID of the user to create")
	}

	if options.Name == "" {
		return errors.New("You must provide the name of the user to create")
	}

	requestBody, err := common.CreateRequestBody(&options)
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPost,
		Path:   "/users",
		Body:   requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// CreateUsers creates users in a batch.
// A maximum of 10 users can be created per batch.
func (cs *coreService) CreateUsers(ctx context.Context, users []CreateUserOptions) error {
	if users == nil || len(users) == 0 {
		return errors.New("You must provide a list of users to create")
	}

	requestBody, err := common.CreateRequestBody(map[string]interface{}{"users": users})
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPost,
		Path:   "/batch_users",
		Body:   requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// UpdateUser updates an existing user depending on the options provided.
func (cs *coreService) UpdateUser(
	ctx context.Context,
	userID string,
	options UpdateUserOptions,
) error {
	if userID == "" {
		return errors.New("You must provide the ID of the user to update")
	}

	requestBody, err := common.CreateRequestBody(&options)
	if err != nil {
		return err
	}

	response, err := common.RequestWithUserToken(cs.underlyingInstance, ctx, userID, client.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/users/%s", userID),
		Body:   requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// DeleteUser deles an existing user.
// Users can only be deleted with a sudo token.
func (cs *coreService) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("You must provide the ID of the user to delete")
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/users/%s", userID),
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// GetRoom retrieves a room with the given id.
func (cs *coreService) GetRoom(ctx context.Context, roomID uint) (Room, error) {
	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/rooms/%d", roomID),
	})
	if err != nil {
		return Room{}, err
	}
	defer response.Body.Close()

	var room Room
	err = common.DecodeResponseBody(response.Body, &room)
	if err != nil {
		return Room{}, err
	}

	return room, nil
}

// GetRooms retrieves a list of rooms with the given parameters.
func (cs *coreService) GetRooms(ctx context.Context, options GetRoomsOptions) ([]Room, error) {
	queryParams := url.Values{}
	if options.FromID != nil {
		queryParams.Add("from_id", strconv.Itoa(int(*options.FromID)))
	}

	strIncludePrivate := "false"
	if options.IncludePrivate {
		strIncludePrivate = "true"
	}
	queryParams.Add("include_private", strIncludePrivate)

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method:      http.MethodGet,
		Path:        "/rooms",
		QueryParams: &queryParams,
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var rooms []Room
	err = common.DecodeResponseBody(response.Body, &rooms)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

// GetUserRooms retrieves a list of rooms that the user is currently a member of.
func (cs *coreService) GetUserRooms(ctx context.Context, userID string) ([]Room, error) {
	return cs.getRoomsForUser(ctx, userID, false)
}

// GetUserJoinable rooms returns a list of rooms that the user can join (has not joined previously).
func (cs *coreService) GetUserJoinableRooms(ctx context.Context, userID string) ([]Room, error) {
	return cs.getRoomsForUser(ctx, userID, true)
}

// getRoomsForUser is used by GetUserRooms and GetUserJoinableRooms
func (cs *coreService) getRoomsForUser(
	ctx context.Context,
	userID string,
	joinable bool,
) ([]Room, error) {
	if userID == "" {
		return nil, errors.New("You must privde the ID of the user to retrieve rooms for")
	}

	strJoinable := "false"
	if joinable {
		strJoinable = "true"
	}
	queryParams := url.Values{"joinable": []string{strJoinable}}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method:      http.MethodGet,
		Path:        fmt.Sprintf("/users/%s/rooms", userID),
		QueryParams: &queryParams,
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var rooms []Room
	err = common.DecodeResponseBody(response.Body, &rooms)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

// CreateRoom creates a room.
func (cs *coreService) CreateRoom(ctx context.Context, options CreateRoomOptions) (Room, error) {
	if options.CreatorID == "" {
		return Room{}, errors.New("Yout must provide the ID of the user creating the room")
	}

	if options.Name == "" {
		return Room{}, errors.New("You must provide a name for the room")
	}

	requestBody, err := common.CreateRequestBody(&options)
	if err != nil {
		return Room{}, err
	}

	response, err := common.RequestWithUserToken(
		cs.underlyingInstance,
		ctx,
		options.CreatorID, client.RequestOptions{
			Method: http.MethodPost,
			Path:   "/rooms",
			Body:   requestBody,
		},
	)
	if err != nil {
		return Room{}, err
	}
	defer response.Body.Close()

	var room Room
	err = common.DecodeResponseBody(response.Body, &room)
	if err != nil {
		return Room{}, err
	}

	return room, nil
}

// UpdateRoom updates an existing room based on the options provided.
func (cs *coreService) UpdateRoom(ctx context.Context, roomID uint, options UpdateRoomOptions) error {
	requestBody, err := common.CreateRequestBody(&options)
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/rooms/%d", roomID),
		Body:   requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// DeleteRoom deletes an existing room.
func (cs *coreService) DeleteRoom(ctx context.Context, roomID uint) error {
	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/rooms/%d", roomID),
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// AddUsersToRoom adds users to an existing room.
// The maximum number of users that can be added in a single request is 10.
func (cs *coreService) AddUsersToRoom(ctx context.Context, roomID uint, userIDs []string) error {
	if userIDs == nil || len(userIDs) == 0 {
		return errors.New("You must provide a list of IDs of the users you want to add to the room")
	}

	requestBody, err := common.CreateRequestBody(map[string][]string{"user_ids": userIDs})
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/rooms/%d/users/add", roomID),
		Body:   requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// RemoveUsersFromRoom removes a list of users from the room.
// The maximum number of users that can be removed in a single request is 10.
func (cs *coreService) RemoveUsersFromRoom(ctx context.Context, roomID uint, userIDs []string) error {
	if userIDs == nil || len(userIDs) == 0 {
		return errors.New("You must provide a list of IDs of the users you want to remove from the room")
	}

	requestBody, err := common.CreateRequestBody(map[string][]string{"user_ids": userIDs})
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/rooms/%d/users/remove", roomID),
		Body:   requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// SendMessage publishes a message to a room.
func (cs *coreService) SendMessage(ctx context.Context, options SendMessageOptions) (uint, error) {
	if options.Text == "" {
		return 0, errors.New("You must provide some text for the message")
	}

	if options.SenderID == "" {
		return 0, errors.New("You must provide the ID of the user sending the message")
	}

	requestBody, err := common.CreateRequestBody(map[string]string{"text": options.Text})
	if err != nil {
		return 0, err
	}

	response, err := common.RequestWithUserToken(
		cs.underlyingInstance,
		ctx,
		options.SenderID,
		client.RequestOptions{
			Method: http.MethodPost,
			Path:   fmt.Sprintf("/rooms/%d/messages", options.RoomID),
			Body:   requestBody,
		})
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	var messageResponse map[string]uint
	err = common.DecodeResponseBody(response.Body, &messageResponse)
	if err != nil {
		return 0, err
	}

	return messageResponse["message_id"], nil
}

// DeleteMessage deletes a previously sent message.
func (cs *coreService) DeleteMessage(ctx context.Context, messageID uint) error {
	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/messages/%d", messageID),
	})
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	return nil
}

// GetRoomMessages fetches messages sent to a room based on the passed in options.
func (cs *coreService) GetRoomMessages(
	ctx context.Context,
	roomID uint,
	options GetRoomMessagesOptions,
) ([]Message, error) {
	queryParams := url.Values{}
	if options.Direction != nil {
		queryParams.Add("direction", *options.Direction)
	}

	if options.InitialID != nil {
		queryParams.Add("initial_id", strconv.Itoa(int(*options.InitialID)))
	}

	if options.Limit != nil {
		queryParams.Add("limit", strconv.Itoa(int(*options.Limit)))
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method:      http.MethodGet,
		Path:        fmt.Sprintf("/rooms/%d/messages", roomID),
		QueryParams: &queryParams,
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var messages []Message
	err = common.DecodeResponseBody(response.Body, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// Request allows performing requests to the core chatkit service and returns the raw http response.
func (cs *coreService) Request(
	ctx context.Context,
	options client.RequestOptions,
) (*http.Response, error) {
	return cs.underlyingInstance.Request(ctx, options)
}
