package authenticator

import (
	auth "github.com/pusher/pusher-platform-go/auth"
)

// Exposes helper methods for authentication.
// Extends the `auth.Authenticator` interface.
type Service interface {
	Authenticate(payload auth.Payload, options auth.Options) (*auth.Response, error)
	GenerateAccessToken(options auth.Options) (auth.TokenWithExpiry, error)
	GenerateSUToken(options auth.Options) (auth.TokenWithExpiry, error)
}

type authenticator struct {
	platformAuthenticator auth.Authenticator
}

// NewService returns a new instance of an authenticator that conforms to the `Service` interface.
func NewService(
	instanceID string,
	keyID string,
	keySecret string,
) Service {
	return &authenticator{
		platformAuthenticator: auth.New(instanceID, keyID, keySecret),
	}
}

// Authenticate should be used within a token providing endpoint to
// generate access tokens for a user.
func (a *authenticator) Authenticate(
	payload auth.Payload,
	options auth.Options,
) (*auth.Response, error) {
	return a.platformAuthenticator.Do(payload, options)
}

// GenerateAccessToken returns a TokenWithExpiry based on the options provided.
func (a *authenticator) GenerateAccessToken(
	options auth.Options,
) (auth.TokenWithExpiry, error) {
	return a.platformAuthenticator.GenerateAccessToken(options)
}

// GenerateSUToken returns a TokenWithExpiry with the `su` claim set to true.
func (a *authenticator) GenerateSUToken(options auth.Options) (auth.TokenWithExpiry, error) {
	return a.GenerateAccessToken(auth.Options{
		UserID:        options.UserID,
		ServiceClaims: options.ServiceClaims,
		Su:            true,
		TokenExpiry:   options.TokenExpiry,
	})
}
