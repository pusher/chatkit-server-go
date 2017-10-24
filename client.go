package chatkitServerGo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// ChatkitServerClient is the public interface of the Chatkit Server Client
type ChatkitServerClient interface {
	// Chatkit Roles and Permissions services
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

	// Chatkit Server services
	CreateUser(user User) error
	DeleteUser(userID string) error
}

// NewChatkitServerClient returns an instantiated instance that fulfils the ChatkitServerClient interface
func NewChatkitServerClient(instanceID string, key string) (ChatkitServerClient, error) {
	apiVersion, host, appID, err := getInstanceIDComponents(instanceID)
	if err != nil {
		return nil, err
	}

	keyID, keySecret, err := getKeyComponents(key)
	if err != nil {
		return nil, err
	}

	tokenManager := newTokenManager(appID, keyID, keySecret)

	return newChatkitServerClient(host, apiVersion, appID, tokenManager), nil
}

type chatkitServerClient struct {
	Client http.Client

	tokenManager *tokenManager

	authEndpoint   string
	serverEndpoint string
}

func newChatkitServerClient(host string, apiVersion string, appID string, tokenManager *tokenManager) *chatkitServerClient {
	return &chatkitServerClient{
		Client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		authEndpoint:   buildServiceEndpoint(host, CHATKIT_AUTH, apiVersion, appID),
		serverEndpoint: buildServiceEndpoint(host, CHATKIT_SERVER, apiVersion, appID),
		tokenManager:   tokenManager,
	}
}

func (csc *chatkitServerClient) newRequest(method, service, path string, body interface{}) (*http.Request, error) {
	var url string
	switch service {
	case CHATKIT_AUTH:
		url = csc.authEndpoint + path
	case CHATKIT_SERVER:
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

	return req, nil
}

func (csc *chatkitServerClient) do(req *http.Request, v interface{}) error {
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

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

const (
	CHATKIT_AUTH   = "chatkit_authorizer"
	CHATKIT_SERVER = "chatkit"
)

func buildServiceEndpoint(host string, service string, apiVersion string, appID string) string {
	return fmt.Sprint("https://", host, "/services/", service, "/", apiVersion, "/", appID)
}

func getInstanceIDComponents(instanceID string) (apiVersion string, host string, appID string, err error) {
	components, err := getColonSeperatedComponents(instanceID, 3)
	if err != nil {
		return "", "", "", err
	}
	return components[0], components[1], components[2], nil
}

func getKeyComponents(key string) (keyID string, keySecret string, err error) {
	components, err := getColonSeperatedComponents(key, 2)
	if err != nil {
		return "", "", err
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
