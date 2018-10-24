// Package cursors exposes an interface that allows making requests to the Chatkit cursors service.
package cursors

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/pusher/chatkit-server-go/internal/common"

	"github.com/pusher/pusher-platform-go/client"
	"github.com/pusher/pusher-platform-go/instance"
)

const readCursorType = 0

// Exposes methods to interact with the cursors API.
type Service interface {
	GetUserReadCursors(ctx context.Context, userID string) ([]Cursor, error)
	SetReadCursor(ctx context.Context, userID string, roomID string, position uint) error
	GetReadCursorsForRoom(ctx context.Context, roomID string) ([]Cursor, error)
	GetReadCursor(ctx context.Context, userID string, roomID string) (Cursor, error)

	// Generic requests
	Request(ctx context.Context, options client.RequestOptions) (*http.Response, error)
}

type cursorsService struct {
	underlyingInstance instance.Instance
}

// Returns a new cursorsService instance conforming to
// the Service interface
func NewService(platformInstance instance.Instance) Service {
	return &cursorsService{
		underlyingInstance: platformInstance,
	}
}

// GetUserReadCursors retrieves cursors for a user.
func (cs *cursorsService) GetUserReadCursors(ctx context.Context, userID string) ([]Cursor, error) {
	if userID == "" {
		return nil, errors.New("You must provide the ID of the user whos read cursors you want to fetch")
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/cursors/%d/users/%s", readCursorType, userID),
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var cursors []Cursor
	err = common.DecodeResponseBody(response.Body, &cursors)
	if err != nil {
		return nil, err
	}

	return cursors, nil
}

// SetReadCursor sets a read cursor for a given room and user.
func (cs *cursorsService) SetReadCursor(
	ctx context.Context,
	userID string,
	roomID string,
	position uint,
) error {
	if userID == "" {
		return errors.New("You must provide the ID of the user whose read cursor you want to set")
	}

	requestBody, err := common.CreateRequestBody(map[string]uint{"position": position})
	if err != nil {
		return err
	}

	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodPut,
		Path: fmt.Sprintf(
			"/cursors/%d/rooms/%s/users/%s",
			readCursorType,
			roomID,
			userID,
		),
		Body: requestBody,
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

// GetReadCursorsForRoom retrieves read cursors for a given room.
func (cs *cursorsService) GetReadCursorsForRoom(ctx context.Context, roomID string) ([]Cursor, error) {
	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/cursors/%d/rooms/%s", readCursorType, roomID),
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var cursors []Cursor
	err = common.DecodeResponseBody(response.Body, &cursors)
	if err != nil {
		return nil, err
	}

	return cursors, nil
}

// GetReadCursor fetches a single cursor for a given user and room.
func (cs *cursorsService) GetReadCursor(
	ctx context.Context,
	userID string,
	roomID string,
) (Cursor, error) {
	response, err := common.RequestWithSuToken(cs.underlyingInstance, ctx, client.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/cursors/%d/rooms/%s/users/%s", readCursorType, roomID, userID),
	})
	if err != nil {
		return Cursor{}, nil
	}
	defer response.Body.Close()

	var cursor Cursor
	err = common.DecodeResponseBody(response.Body, &cursor)
	if err != nil {
		return Cursor{}, nil
	}

	return cursor, nil
}

// Request allows performing requests to the cursors service and returns the raw http response.
func (cs *cursorsService) Request(
	ctx context.Context,
	options client.RequestOptions,
) (*http.Response, error) {
	return cs.underlyingInstance.Request(ctx, options)
}
