package chatkitServerGo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildServiceEndpoint(t *testing.T) {
	endpoint := buildServiceEndpoint("host", chatkitAuthService, "v1", "abc123")

	assert.Equal(t, "https://host/services/chatkit_authorizer/v1/abc123", endpoint, "Should return a correctly fomatted endpoint")
}

func TestNewRequest(t *testing.T) {
	type args struct {
		method  string
		service string
		path    string
		body    interface{}
	}
	tests := []struct {
		name          string
		client        *chatkitServerClient
		args          args
		expectRequest bool
		expectErr     bool
	}{
		{
			name: "valid with body AUTH",
			client: &chatkitServerClient{
				authEndpoint: "https://host/services/chatkit_authorizer/v1/abc123",
				tokenManager: newMockTokenManager(),
			},
			args: args{
				method:  "GET",
				service: chatkitAuthService,
				path:    "/roles",
				body:    "request_body!",
			},
			expectRequest: true,
			expectErr:     false,
		},
		{
			name: "valid without body AUTH",
			client: &chatkitServerClient{
				authEndpoint: "https://host/services/chatkit_authorizer/v1/abc123",
				tokenManager: newMockTokenManager(),
			},
			args: args{
				method:  "GET",
				service: chatkitAuthService,
				path:    "/roles",
				body:    nil,
			},
			expectRequest: true,
			expectErr:     false,
		},
		{
			name: "valid without body SERVER",
			client: &chatkitServerClient{
				serverEndpoint: "https://host/services/chatkit/v1/abc123",
				tokenManager:   newMockTokenManager(),
			},
			args: args{
				method:  "GET",
				service: chatkitService,
				path:    "/roles",
				body:    nil,
			},
			expectRequest: true,
			expectErr:     false,
		},
		{
			name: "invalid no service",
			client: &chatkitServerClient{
				serverEndpoint: "https://host/services/chatkit/v1/abc123",
				tokenManager:   newMockTokenManager(),
			},
			args: args{
				method:  "GET",
				service: "invalid service",
				path:    "/roles",
				body:    nil,
			},
			expectRequest: false,
			expectErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRequest, err := tt.client.newRequest(tt.args.method, tt.args.service, tt.args.path, tt.args.body)
			assert.Equal(t, tt.expectErr, (err != nil), "Unexpected error returned")
			assert.Equal(t, tt.expectRequest, (actualRequest != nil), "Unexpected request returned")
		})
	}
}

func newTestClientAndServer(handler http.Handler) (*chatkitServerClient, *httptest.Server) {
	testServer := httptest.NewServer(handler)

	testClient := &chatkitServerClient{
		Client:         http.Client{},
		authEndpoint:   testServer.URL,
		serverEndpoint: testServer.URL,
		tokenManager:   newMockTokenManager(),
	}

	return testClient, testServer
}

func TestGetinstanceLocatorComponents(t *testing.T) {
	tests := []struct {
		name               string
		instanceLocator    string
		expectedAPIVersion string
		expectedHost       string
		expectedAppID      string
		expectError        bool
	}{
		{
			name:            "empty string",
			instanceLocator: "",
			expectError:     true,
		},
		{
			name:            "no colons",
			instanceLocator: "notValid",
			expectError:     true,
		},
		{
			name:            "not enough colons",
			instanceLocator: "stillNot:Valid",
			expectError:     true,
		},
		{
			name:            "correct colons with no text",
			instanceLocator: "::",
			expectError:     true,
		},
		{
			name:            "empty components",
			instanceLocator: ":incorrect:",
			expectError:     true,
		},
		{
			name:               "valid",
			instanceLocator:    "v1:us:abc123",
			expectedAPIVersion: "v1",
			expectedHost:       "us",
			expectedAppID:      "abc123",
			expectError:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAPIVersion, gotHost, gotAppID, err := getinstanceLocatorComponents(tt.instanceLocator)
			assert.Equal(t, tt.expectError, (err != nil), "Unexpected error returned")
			assert.Equal(t, tt.expectedAPIVersion, gotAPIVersion, "Unexpected APIVersion returned")
			assert.Equal(t, tt.expectedAppID, gotAppID, "Unexpected AppID returned")
			assert.Equal(t, tt.expectedHost, gotHost, "Unexpected Host returned")
		})
	}
}

func TestGetKeyComponents(t *testing.T) {
	tests := []struct {
		name              string
		key               string
		expectedKeyID     string
		expectedKeySecret string
		expectError       bool
	}{
		{
			name:        "empty string",
			key:         "",
			expectError: true,
		},
		{
			name:        "no colons",
			key:         "notValid",
			expectError: true,
		},
		{
			name:        "correct colons with no text",
			key:         ":",
			expectError: true,
		},
		{
			name:        "empty components",
			key:         ":incorrect",
			expectError: true,
		},
		{
			name:              "valid",
			key:               "id:secret",
			expectedKeyID:     "id",
			expectedKeySecret: "secret",
			expectError:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeyID, gotKeySecret, err := getKeyComponents(tt.key)
			assert.Equal(t, tt.expectError, (err != nil), "Unexpected error returned")
			assert.Equal(t, tt.expectedKeyID, gotKeyID, "Unexpected keyID returned")
			assert.Equal(t, tt.expectedKeySecret, gotKeySecret, "Unexpected keySecret returned")
		})
	}
}
