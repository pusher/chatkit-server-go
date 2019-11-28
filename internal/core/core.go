// Package core exposes an interface that allows making requests to the core
// Chatkit API to allow operations to be performed against Users, Rooms and Messages.
package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/pusher/pusher-platform-go/client"
	"github.com/pusher/pusher-platform-go/instance"

	"github.com/pusher/chatkit-server-go/internal/common"
)

// Service exposes methods to interact with the core chatkit service.
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
	GetRoom(ctx context.Context, roomID string) (Room, error)
	GetRooms(ctx context.Context, options GetRoomsOptions) ([]RoomWithoutMembers, error)
	GetUserRooms(ctx context.Context, userID string) ([]Room, error)
	GetUserJoinableRooms(ctx context.Context, userID string) ([]Room, error)
	CreateRoom(ctx context.Context, options CreateRoomOptions) (Room, error)
	UpdateRoom(ctx context.Context, roomID string, options UpdateRoomOptions) error
	DeleteRoom(ctx context.Context, roomID string) error
	AddUsersToRoom(ctx context.Context, roomID string, userIDs []string) error
	RemoveUsersFromRoom(ctx context.Context, roomID string, userIds []string) error

	// Messages
	SendMessage(ctx context.Context, options SendMessageOptions) (uint, error)
	SendMultipartMessage(ctx context.Context, options SendMultipartMessageOptions) (uint, error)
	SendSimpleMessage(ctx context.Context, options SendSimpleMessageOptions) (uint, error)
	GetRoomMessages(
		ctx context.Context,
		roomID string,
		options GetRoomMessagesOptions,
	) ([]Message, error)
	FetchMultipartMessage(
		ctx context.Context,
		options FetchMultipartMessageOptions,
	) (MultipartMessage, error)
	FetchMultipartMessages(
		ctx context.Context,
		roomID string,
		options FetchMultipartMessagesOptions,
	) ([]MultipartMessage, error)
	DeleteMessage(ctx context.Context, options DeleteMessageOptions) error

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
	if len(users) == 0 {
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
func (cs *coreService) GetRoom(ctx context.Context, roomID string) (Room, error) {
	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/rooms/%s", roomID),
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
func (cs *coreService) GetRooms(ctx context.Context, options GetRoomsOptions) ([]RoomWithoutMembers, error) {
	queryParams := url.Values{}
	if options.FromID != nil {
		queryParams.Add("from_id", *options.FromID)
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

	var rooms []RoomWithoutMembers
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
func (cs *coreService) UpdateRoom(ctx context.Context, roomID string, options UpdateRoomOptions) error {
	var formattedOptions interface{} = options
	if options.PushNotificationTitleOverride == &ExplicitlyResetPushNotificationTitleOverride {
		type updateRoomOptionsWithExplicitPNTitleOverride struct {
			UpdateRoomOptions
			// overriding internal `UpdateRoomOptions.PushNotificationTitleOverride` by removing the `omitempty` tag
			PushNotificationTitleOverride *string `json:"push_notification_title_override"`
		}
		formattedOptions = updateRoomOptionsWithExplicitPNTitleOverride{
			UpdateRoomOptions:             options,
			PushNotificationTitleOverride: nil,
		}
	}

	requestBody, err := common.CreateRequestBody(&formattedOptions)
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/rooms/%s", roomID),
		Body:   requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// DeleteRoom deletes an existing room.
func (cs *coreService) DeleteRoom(ctx context.Context, roomID string) error {
	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/rooms/%s", roomID),
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// AddUsersToRoom adds users to an existing room.
// The maximum number of users that can be added in a single request is 10.
func (cs *coreService) AddUsersToRoom(ctx context.Context, roomID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return errors.New("You must provide a list of IDs of the users you want to add to the room")
	}

	requestBody, err := common.CreateRequestBody(map[string][]string{"user_ids": userIDs})
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/rooms/%s/users/add", roomID),
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
func (cs *coreService) RemoveUsersFromRoom(ctx context.Context, roomID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return errors.New("You must provide a list of IDs of the users you want to remove from the room")
	}

	requestBody, err := common.CreateRequestBody(map[string][]string{"user_ids": userIDs})
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/rooms/%s/users/remove", roomID),
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
			Path:   fmt.Sprintf("/rooms/%s/messages", options.RoomID),
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

// SendMultipartMessage publishes a multipart message to a room.
func (cs *coreService) SendMultipartMessage(
	ctx context.Context,
	options SendMultipartMessageOptions,
) (uint, error) {
	if len(options.Parts) == 0 {
		return 0, errors.New("You must provide at least one message part")
	}

	if options.SenderID == "" {
		return 0, errors.New("You must provide the ID of the user sending the message")
	}

	requestParts := make([]interface{}, len(options.Parts))
	g := errgroup.Group{}

	for i, part := range options.Parts {
		switch p := part.(type) {
		case NewAttachmentPart:
			i := i
			g.Go(func() error {
				uploadedPart, err := cs.uploadAttachment(ctx, options.SenderID, options.RoomID, p)
				requestParts[i] = uploadedPart
				return err
			})
		default:
			requestParts[i] = part
		}
	}

	if err := g.Wait(); err != nil {
		return 0, fmt.Errorf("Failed to upload attachment: %v", err)
	}

	requestBody, err := common.CreateRequestBody(
		map[string]interface{}{"parts": requestParts},
	)
	if err != nil {
		return 0, err
	}

	response, err := common.RequestWithUserToken(
		cs.underlyingInstance,
		ctx,
		options.SenderID,
		client.RequestOptions{
			Method: http.MethodPost,
			Path:   fmt.Sprintf("/rooms/%s/messages", options.RoomID),
			Body:   requestBody,
		},
	)
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

func (cs *coreService) uploadAttachment(
	ctx context.Context,
	senderID string,
	roomID string,
	part NewAttachmentPart,
) (newAttachmentPartUploaded, error) {
	// Unfortunately since we need to provide the content length up front, we
	// have to read the whole file in to memory.
	b, err := ioutil.ReadAll(part.File)
	if err != nil {
		return newAttachmentPartUploaded{}, err
	}

	url, attachmentID, err := cs.requestPresignedURL(
		ctx,
		senderID,
		roomID,
		part.Type,
		len(b),
		part.Name,
		part.CustomData,
	)
	if err != nil {
		return newAttachmentPartUploaded{}, err
	}

	if err := cs.uploadToURL(ctx, url, part.Type, len(b), bytes.NewReader(b)); err != nil {
		return newAttachmentPartUploaded{}, err
	}

	return newAttachmentPartUploaded{
		Type:       part.Type,
		Attachment: uploadedAttachment{attachmentID},
	}, nil
}

func (cs *coreService) requestPresignedURL(
	ctx context.Context,
	senderID string,
	roomID string,
	contentType string,
	contentLength int,
	name *string,
	customData interface{},
) (string, string, error) {
	body, err := common.CreateRequestBody(map[string]interface{}{
		"content_type":   contentType,
		"content_length": contentLength,
		"name":           name,
		"custom_data":    customData,
	})
	if err != nil {
		return "", "", err
	}

	res, err := common.RequestWithUserToken(
		cs.underlyingInstance,
		ctx,
		senderID,
		client.RequestOptions{
			Method: http.MethodPost,
			Path:   fmt.Sprintf("/rooms/%s/attachments", roomID),
			Body:   body,
		},
	)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	var resBody map[string]string
	if err := common.DecodeResponseBody(res.Body, &resBody); err != nil {
		return "", "", err
	}

	return resBody["upload_url"], resBody["attachment_id"], nil
}

func (cs *coreService) uploadToURL(
	ctx context.Context,
	url string,
	contentType string,
	contentLength int,
	body io.Reader,
) error {
	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", contentType)
	req.Header.Add("content-length", strconv.Itoa(contentLength))

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %v", res.Status)
	}

	return nil
}

// SendSimpleMessage publishes a simple message to a room.
func (cs *coreService) SendSimpleMessage(
	ctx context.Context,
	options SendSimpleMessageOptions,
) (uint, error) {
	return cs.SendMultipartMessage(ctx, SendMultipartMessageOptions{
		RoomID:   options.RoomID,
		SenderID: options.SenderID,
		Parts:    []NewPart{NewInlinePart{Type: "text/plain", Content: options.Text}},
	})
}

// DeleteMessage deletes a previously sent message.
func (cs *coreService) DeleteMessage(ctx context.Context, options DeleteMessageOptions) error {
	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/rooms/%s/messages/%d", options.RoomID, options.MessageID),
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
	roomID string,
	options GetRoomMessagesOptions,
) ([]Message, error) {
	messages := []Message{}
	err := cs.fetchMessages(ctx, roomID, options, &messages)
	return messages, err
}

// FetchMultipartMessage fetches a single message sent to a room based on the passed in options.
func (cs *coreService) FetchMultipartMessage(
	ctx context.Context,
	options FetchMultipartMessageOptions,
) (MultipartMessage, error) {
	response, err := common.RequestWithSuToken(
		cs.underlyingInstance,
		ctx,
		client.RequestOptions{
			Method: http.MethodGet,
			Path:   fmt.Sprintf("/rooms/%s/messages/%d", options.RoomID, options.MessageID),
		},
	)
	if err != nil {
		return MultipartMessage{}, err
	}
	defer response.Body.Close()

	var message MultipartMessage
	err = common.DecodeResponseBody(response.Body, &message)
	if err != nil {
		return MultipartMessage{}, err
	}

	return message, nil
}

// FetchMultipartMessages fetches messages sent to a room based on the passed in options.
func (cs *coreService) FetchMultipartMessages(
	ctx context.Context,
	roomID string,
	options FetchMultipartMessagesOptions,
) ([]MultipartMessage, error) {
	messages := []MultipartMessage{}
	err := cs.fetchMessages(ctx, roomID, options, &messages)
	return messages, err
}

func (cs *coreService) fetchMessages(
	ctx context.Context,
	roomID string,
	options fetchMessagesOptions,
	target interface{}, // poor man's generics
) error {
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
		Path:        fmt.Sprintf("/rooms/%s/messages", roomID),
		QueryParams: &queryParams,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = common.DecodeResponseBody(response.Body, target)
	if err != nil {
		return err
	}

	return nil
}

// Request allows performing requests to the core chatkit service and returns the raw http response.
func (cs *coreService) Request(
	ctx context.Context,
	options client.RequestOptions,
) (*http.Response, error) {
	return cs.underlyingInstance.Request(ctx, options)
}
