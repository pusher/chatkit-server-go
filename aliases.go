package chatkit

import (
	"github.com/pusher/chatkit-server-go/internal/authorizer"
	"github.com/pusher/chatkit-server-go/internal/core"
	"github.com/pusher/chatkit-server-go/internal/cursors"

	auth "github.com/pusher/pusher-platform-go/auth"
	platformclient "github.com/pusher/pusher-platform-go/client"
)

const GrantTypeClientCredentials = auth.GrantTypeClientCredentials

type (
	AuthenticatePayload = auth.Payload
	AuthenticateOptions = auth.Options

	ErrorResponse  = platformclient.ErrorResponse
	RequestOptions = platformclient.RequestOptions

	CreateRoleOptions            = authorizer.CreateRoleOptions
	UpdateRolePermissionsOptions = authorizer.UpdateRolePermissionsOptions
	Role                         = authorizer.Role

	Cursor = cursors.Cursor

	GetUsersOptions               = core.GetUsersOptions
	CreateUserOptions             = core.CreateUserOptions
	UpdateUserOptions             = core.UpdateUserOptions
	GetRoomsOptions               = core.GetRoomsOptions
	CreateRoomOptions             = core.CreateRoomOptions
	UpdateRoomOptions             = core.UpdateRoomOptions
	SendMessageOptions            = core.SendMessageOptions
	SendMultipartMessageOptions   = core.SendMultipartMessageOptions
	SendSimpleMessageOptions      = core.SendSimpleMessageOptions
	NewPart                       = core.NewPart
	NewInlinePart                 = core.NewInlinePart
	NewURLPart                    = core.NewURLPart
	NewAttachmentPart             = core.NewAttachmentPart
	GetRoomMessagesOptions        = core.GetRoomMessagesOptions
	FetchMultipartMessagesOptions = core.FetchMultipartMessagesOptions
	User                          = core.User
	Room                          = core.Room
	Message                       = core.Message
	MultipartMessage              = core.MultipartMessage
	Part                          = core.Part
	Attachment                    = core.Attachment
)
