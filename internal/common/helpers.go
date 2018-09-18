package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pusher/pusher-platform-go/auth"
	"github.com/pusher/pusher-platform-go/client"
	"github.com/pusher/pusher-platform-go/instance"
)

// DecodeResponseBody takes an io.Reader and decodes the body into a destination struct
func DecodeResponseBody(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(dest)
	if err != nil {
		return fmt.Errorf("Failed to decode response body: %s", err.Error())
	}

	return nil
}

// CreateRequestBody takes a struct/ map and converts it into an io.Reader
func CreateRequestBody(target interface{}) (io.Reader, error) {
	bodyBytes, err := json.Marshal(target)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal into json: %v", err)
	}

	return bytes.NewReader(bodyBytes), nil
}

// generateTokenFromInstance generates a token with the given options.
func generateTokenFromInstance(inst instance.Instance, options auth.Options) (string, error) {
	tokenWithExpiry, err := inst.GenerateAccessToken(options)
	if err != nil {
		return "", fmt.Errorf("Failed to generate token: %s", err.Error())
	}

	return tokenWithExpiry.Token, nil
}

// RequestWithToken makes a request and includes token generation as a part of it.
// It generates a token with the `su` claim.
func RequestWithSuToken(
	inst instance.Instance,
	ctx context.Context,
	options client.RequestOptions,
) (*http.Response, error) {
	token, err := generateTokenFromInstance(inst, auth.Options{Su: true})
	if err != nil {
		return nil, err
	}

	return inst.Request(ctx, client.RequestOptions{
		Method:      options.Method,
		Path:        options.Path,
		Body:        options.Body,
		Headers:     options.Headers,
		QueryParams: options.QueryParams,
		Jwt:         &token,
	})
}

// RequestWithUserToken makes a request and includes the user id as part of the `sub` claim.
func RequestWithUserToken(
	inst instance.Instance,
	ctx context.Context,
	userID string,
	options client.RequestOptions,
) (*http.Response, error) {
	token, err := generateTokenFromInstance(inst, auth.Options{UserID: &userID})
	if err != nil {
		return nil, err
	}

	return inst.Request(ctx, client.RequestOptions{
		Method:      options.Method,
		Path:        options.Path,
		Body:        options.Body,
		Headers:     options.Headers,
		QueryParams: options.QueryParams,
		Jwt:         &token,
	})
}
