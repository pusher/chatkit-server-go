// Package chatkit is the Golang server SDK for Pusher Chatkit.
// This package provides functionality to interact with various Chatkit services.
//
// More information can be found in the Chatkit docs: https://docs.pusher.com/chatkit/overview/.
//
// Please report any bugs or feature requests at: https://github.com/pusher/chatkit-server-go.
package chatkit

import (
	"github.com/pusher/chatkit-server-go/internal/authenticator"
	"github.com/pusher/chatkit-server-go/internal/authorizer"
	"github.com/pusher/chatkit-server-go/internal/core"
	"github.com/pusher/chatkit-server-go/internal/cursors"

	platformclient "github.com/pusher/pusher-platform-go/client"
	instance "github.com/pusher/pusher-platform-go/instance"
)

const (
	chatkitAuthorizerServiceName    = "chatkit_authorizer"
	chatkitAuthorizerServiceVersion = "v1"
	chatkitServiceName              = "chatkit"
	chatkitServiceVersion           = "v2"
	chatkitCursorsServiceName       = "chatkit_cursors"
	chatkitCursorsServiceVersion    = "v1"
)

// Public interface for the library.
// It allows interacting with different Chatkit services.
type Client struct {
	coreService          core.Service
	authorizerService    authorizer.Service
	cursorsService       cursors.Service
	authenticatorService authenticator.Service
}

// NewClient returns an instantiated instance that fulfils the Client interface.
func NewClient(instanceLocator string, key string) (*Client, error) {
	locatorComponents, err := instance.ParseInstanceLocator(instanceLocator)
	if err != nil {
		return nil, err
	}

	keyComponents, err := instance.ParseKey(key)
	if err != nil {
		return nil, err
	}

	baseClient := platformclient.New(platformclient.Options{
		Host: locatorComponents.Host(),
	})

	coreInstance, err := instance.New(instance.Options{
		Locator:        instanceLocator,
		Key:            key,
		ServiceName:    chatkitServiceName,
		ServiceVersion: chatkitServiceVersion,
		Client:         baseClient,
	})
	if err != nil {
		return nil, err
	}

	authorizerInstance, err := instance.New(instance.Options{
		Locator:        instanceLocator,
		Key:            key,
		ServiceName:    chatkitAuthorizerServiceName,
		ServiceVersion: chatkitAuthorizerServiceVersion,
		Client:         baseClient,
	})
	if err != nil {
		return nil, err
	}

	cursorsInstance, err := instance.New(instance.Options{
		Locator:        instanceLocator,
		Key:            key,
		ServiceName:    chatkitCursorsServiceName,
		ServiceVersion: chatkitCursorsServiceVersion,
		Client:         baseClient,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		coreService:       core.NewService(coreInstance),
		authorizerService: authorizer.NewService(authorizerInstance),
		cursorsService:    cursors.NewService(cursorsInstance),
		authenticatorService: authenticator.NewService(
			locatorComponents.InstanceID,
			keyComponents.Key,
			keyComponents.Secret,
		),
	}, nil
}
