package chatkitServerGo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEndpointValid(t *testing.T) {

	endpoint, err := buildEndpoint("cluster.pusherplatform.io", "api_version", "app_id")

	assert.NoError(t, err, "Should not return error with valid client endpoint variables")
	assert.Equal(t, "https://cluster.pusherplatform.io/services/chatkit_authorizer/api_version/app_id", endpoint, "Should return a correctly fomatted endpoint")
}

func TestBuildEndpointInvalid(t *testing.T) {
	endpoint, err := buildEndpoint("", "", "")

	assert.Error(t, err, "Should return error when testClient endpoint variables are empty")
	assert.Equal(t, "", endpoint, "Should return a an empty string")
}

func TestNewRequest(t *testing.T) {
	type args struct {
		method string
		path   string
		body   interface{}
	}
	tests := []struct {
		name          string
		client        *chatkitServerClient
		args          args
		expectRequest bool
		expectErr     bool
	}{
		{
			name: "valid with body",
			client: &chatkitServerClient{
				endpoint: "pusher.com",
			},
			args: args{
				method: "GET",
				path:   "/roles",
				body:   "request_body!",
			},
			expectRequest: true,
			expectErr:     false,
		},
		{
			name: "valid without body",
			client: &chatkitServerClient{
				endpoint: "pusher.com",
			},
			args: args{
				method: "GET",
				path:   "/roles",
				body:   nil,
			},
			expectRequest: true,
			expectErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRequest, err := tt.client.newRequest(tt.args.method, tt.args.path, tt.args.body)
			assert.Equal(t, tt.expectErr, (err != nil), "Unexpected error returned")
			assert.Equal(t, tt.expectRequest, (actualRequest != nil), "Unexpected request returned")
		})
	}
}

func newTestClientAndServer(handler http.Handler) (*chatkitServerClient, *httptest.Server) {
	testServer := httptest.NewServer(handler)

	testClient := &chatkitServerClient{
		Client:   http.Client{},
		endpoint: testServer.URL,
	}

	return testClient, testServer
}

func TestGetInstanceIDComponents(t *testing.T) {
	tests := []struct {
		name               string
		instanceID         string
		expectedAPIVersion string
		expectedHost       string
		expectedAppID      string
		expectError        bool
	}{
		{
			name:        "empty string",
			instanceID:  "",
			expectError: true,
		},
		{
			name:        "no colons",
			instanceID:  "notValid",
			expectError: true,
		},
		{
			name:        "not enough colons",
			instanceID:  "stillNot:Valid",
			expectError: true,
		},
		{
			name:        "correct colons with no text",
			instanceID:  "::",
			expectError: true,
		},
		{
			name:        "empty components",
			instanceID:  ":incorrect:",
			expectError: true,
		},
		{
			name:               "valid",
			instanceID:         "v1:us:abc123",
			expectedAPIVersion: "v1",
			expectedHost:       "us",
			expectedAppID:      "abc123",
			expectError:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAPIVersion, gotHost, gotAppID, err := getInstanceIDComponents(tt.instanceID)
			assert.Equal(t, tt.expectError, (err != nil), "Unexpected error returned")
			assert.Equal(t, tt.expectedAPIVersion, gotAPIVersion, "Unexpected APIVersion returned")
			assert.Equal(t, tt.expectedAppID, gotAppID, "Unexpected AppID returned")
			assert.Equal(t, tt.expectedHost, gotHost, "Unexpected Host returned")
		})
	}
}
