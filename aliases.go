package chatkit

import (
	auth "github.com/pusher/pusher-platform-go/auth"
	platformclient "github.com/pusher/pusher-platform-go/client"
)

// Aliases some platform types
type (
	AuthenticatePayload    = auth.Payload
	AuthenticateOptions    = auth.Options
	AuthenticationResponse = auth.Response

	ErrorResponse  = platformclient.ErrorResponse
	RequestOptions = platformclient.RequestOptions
)
