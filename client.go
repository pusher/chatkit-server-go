// Package chatkit is the Golang server SDK for Pusher Chatkit.
//
// This package provides the Client type for managing Chatkit users and
// interacting with roles and permissions of those users. It also contains some helper
// functions for creating your own JWT tokens for authentication with the Chatkit
// service.
//
// More information can be found in the Chatkit docs: https://docs.pusher.com/chatkit/overview/
//
// Please report any bugs or feature requests at: https://github.com/pusher/chatkit-server-go
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
	chatkitServiceVersion           = "v1"
	chatkitCursorsServiceName       = "chatkit_cursors"
	chatkitCursorsServiceVersion    = "v1"
)

// Public interface for the library.
// It allows interacting with different Chatkit services.
type Client interface {
	authorizer.Service    // Provides access to the Roles and Permissions API
	core.Service          // Provides access to the core (rooms, messages, users) chatkit API
	cursors.Service       // Provides access to the Cursors API
	authenticator.Service // Token generation and Authentication
}

type client struct {
	coreService          core.Service
	authorizerService    authorizer.Service
	cursorsService       cursors.Service
	authenticatorService authenticator.Service
}

// NewClient returns an instantiated instance that fulfils the Client interface.
func NewClient(instanceLocator string, key string) (Client, error) {
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

	coreService, err := core.NewService(instance.Options{
		Locator:        instanceLocator,
		Key:            key,
		ServiceName:    chatkitServiceName,
		ServiceVersion: chatkitServiceVersion,
		Client:         baseClient,
	})
	if err != nil {
		return nil, err
	}

	authorizerService, err := authorizer.NewService(instance.Options{
		Locator:        instanceLocator,
		Key:            key,
		ServiceName:    chatkitAuthorizerServiceName,
		ServiceVersion: chatkitAuthorizerServiceVersion,
		Client:         baseClient,
	})
	if err != nil {
		return nil, err
	}

	cursorsService, err := cursors.NewService(instance.Options{
		Locator:        instanceLocator,
		Key:            key,
		ServiceName:    chatkitCursorsServiceName,
		ServiceVersion: chatkitCursorsServiceVersion,
		Client:         baseClient,
	})
	if err != nil {
		return nil, err
	}

	return &client{
		coreService:       coreService,
		authorizerService: authorizerService,
		cursorsService:    cursorsService,
		authenticatorService: authenticator.NewService(
			locatorComponents.InstanceID,
			keyComponents.Key,
			keyComponents.Secret,
		),
	}, nil
}
