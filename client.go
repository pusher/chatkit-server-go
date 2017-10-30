/*
Package chatkit is the Golang server SDK for Pusher Chatkit.

This package provides the Client type for managing Chatkit users and
interacting with roles and permissions of those users. It also contains some helper
functions for creating your own JWT tokens for authentication with the Chatkit
service.

More information can be found in the Chatkit docs: https://docs.pusher.com/chatkit/overview/

Please report any bugs or feature requests at: https://github.com/pusher/chatkit-server-go
*/
package chatkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	chatkitAuthService = "chatkit_authorizer"
	chatkitService     = "chatkit"
)

// Client is the public interface of the Chatkit Server Client.
// It contains methods for creating and deleting users and managing those user's
// roles and permissions.
type Client interface {
	// Chatkit Roles and Permissions methods
	GetRoles() ([]Role, error)
	CreateRole(Role) error
	DeleteRole(roleName string, scopeType string) error

	GetUserRoles(userID string) ([]Role, error)
	CreateUserRole(userID string, userRole UserRole) error
	UpdateUserRole(userID string, userRole UserRole) error
	DeleteUserRole(userID string, roomID *string) error

	CreateRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error
	GetRolePermissions(roleName string, scopeName string) (*RolePermissions, error)
	EditRolePermissions(roleName string, scopeName string, rolePerms RolePermissions) error

	// Chatkit User methods
	CreateUser(user User) error
	DeleteUser(userID string) error
}

// NewClient returns an instantiated instance that fulfils the Client interface
func NewClient(instanceLocator string, key string) (Client, error) {
	apiVersion, host, appID, err := getinstanceLocatorComponents(instanceLocator)
	if err != nil {
		return nil, err
	}

	keyID, keySecret, err := getKeyComponents(key)
	if err != nil {
		return nil, err
	}

	tokenManager := newTokenManager(appID, keyID, keySecret)

	return newClient(host, apiVersion, appID, tokenManager), nil
}

type client struct {
	Client http.Client

	tokenManager tokenManager

	authEndpoint   string
	serverEndpoint string
}

func newClient(host string, apiVersion string, appID string, tokenManager tokenManager) *client {
	return &client{
		Client:         http.Client{},
		authEndpoint:   buildServiceEndpoint(host, chatkitAuthService, apiVersion, appID),
		serverEndpoint: buildServiceEndpoint(host, chatkitService, apiVersion, appID),
		tokenManager:   tokenManager,
	}
}

func (csc *client) newRequest(method, serviceName, path string, body interface{}) (*http.Request, error) {
	var url string
	switch serviceName {
	case chatkitAuthService:
		url = csc.authEndpoint + path
	case chatkitService:
		url = csc.serverEndpoint + path
	default:
		return nil, errors.New("no service was provided to newRequest")
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	token, err := csc.tokenManager.getToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

func (csc *client) do(req *http.Request, responseBody interface{}) error {
	resp, err := csc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New(resp.Status)
		}
		return errors.New(resp.Status + ": " + string(body))
	}

	if responseBody != nil {
		return json.NewDecoder(resp.Body).Decode(responseBody)
	}
	return nil
}

func buildServiceEndpoint(host string, service string, apiVersion string, appID string) string {
	return fmt.Sprint("https://", host, ".pusherplatform.io/services/", service, "/", apiVersion, "/", appID)
}

func getinstanceLocatorComponents(instanceLocator string) (apiVersion string, host string, appID string, err error) {
	components, err := getColonSeperatedComponents(instanceLocator, 3)
	if err != nil {
		return "", "", "", errors.New("Incorrect instanceLocator format given, please get your app instanceLocator from your user dashboard")
	}
	return components[0], components[1], components[2], nil
}

func getKeyComponents(key string) (keyID string, keySecret string, err error) {
	components, err := getColonSeperatedComponents(key, 2)
	if err != nil {
		return "", "", errors.New("Incorrect key format given, please get your app key from your user dashboard")
	}
	return components[0], components[1], nil
}

func getColonSeperatedComponents(s string, expectedComponents int) ([]string, error) {
	if s == "" {
		return nil, errors.New("empty string")
	}

	components := strings.Split(s, ":")
	if len(components) != expectedComponents {
		return nil, errors.New("incorrect format")
	}

	for _, component := range components {
		if component == "" {
			return nil, errors.New("incorrect format")
		}
	}

	return components, nil
}
