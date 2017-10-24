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
func NewChatkitServerClient(host string, apiVersion string, appID string) (ChatkitServerClient, error) {
	return newChatkitServerClient(host, apiVersion, appID)
}

// // NewChatkitServerClient returns an instantiated instance that fulfils the ChatkitServerClient interface
// func NewChatkitServerClient(instanceID string, key string) (ChatkitServerClient, error) {
// 	apiVersion, host, appID, err := getIntanceIDComponents(instanceID)
// 	if err != nil {
// 		return err
// 	}

// 	return newChatkitServerClient(host, apiVersion, appID)
// }

type chatkitServerClient struct {
	Client http.Client

	endpoint string
}

func newChatkitServerClient(host string, apiVersion string, appID string) (*chatkitServerClient, error) {
	endpoint, err := buildEndpoint(host, apiVersion, appID)
	if err != nil {
		return nil, err
	}

	return &chatkitServerClient{
		Client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		endpoint: endpoint,
	}, nil
}

func (csc *chatkitServerClient) newRequest(method, path string, body interface{}) (*http.Request, error) {
	url := csc.endpoint + path

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

func buildEndpoint(host string, apiVersion string, appID string) (string, error) {
	if host == "" || apiVersion == "" || appID == "" {
		return "", errors.New("APIVersion, Cluster and AppID must be non zero length strings")
	}

	return fmt.Sprint("https://", host, "/services/chatkit_authorizer/", apiVersion, "/", appID), nil
}

func buildServiceEndpoint(host string, service string, apiVersion string, appID string) (string, error) {
	if host == "" || apiVersion == "" || appID == "" {
		return "", errors.New("APIVersion, Cluster and AppID must be non zero length strings")
	}

	return fmt.Sprint("https://", host, "/services/", service, "/", apiVersion, "/", appID), nil
}
