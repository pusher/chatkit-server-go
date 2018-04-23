package chatkit

import (
	"fmt"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestNewChatkitTokenWithSuAndNoSub(t *testing.T) {
	tokenBody, errorBody := NewChatkitToken("instanceID", "keyID", "keySecret", nil, true, time.Hour)
	assert.Nil(t, errorBody, "expect no error")

	token := tokenBody.AccessToken
	expiry := time.Now().Add(time.Duration(int(tokenBody.ExpiresIn)) * time.Second)
	assert.False(t, time.Now().After(expiry), "expiry should be after now")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing alg is HMAC-SHA256:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return the key to be parsed with
		return []byte("keySecret"), nil
	})

	assert.NoError(t, err, "expect no error when parsing the token")
	assert.True(t, parsedToken.Valid, "token produced was invalid")

	claimMap, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fail()
	}

	_, present := claimMap["sub"]
	assert.False(t, present, "sub claim should not be present")
}

func TestNewChatkitTokenWithSubAndSu(t *testing.T) {
	sub := "jane"
	tokenBody, errorBody := NewChatkitToken("instanceID", "keyID", "keySecret", &sub, true, time.Hour)
	assert.Nil(t, errorBody, "expect no error")

	token := tokenBody.AccessToken
	expiry := time.Now().Add(time.Duration(int(tokenBody.ExpiresIn)) * time.Second)
	assert.False(t, time.Now().After(expiry), "expiry should be after now")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing alg is HMAC-SHA256:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return the key to be parsed with
		return []byte("keySecret"), nil
	})

	assert.NoError(t, err, "expect no error when parsing the token")
	assert.True(t, parsedToken.Valid, "token produced was invalid")

	claimMap, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fail()
	}

	assert.Equal(t, sub, claimMap["sub"], "token did not contain a sub claim with supplied user name")
}

func TestNewChatkitTokenWithSub(t *testing.T) {
	userID := "bob"
	tokenBody, errorBody := NewChatkitToken("instanceID", "keyID", "keySecret", &userID, false, time.Hour)
	assert.Nil(t, errorBody, "expect no error")

	token := tokenBody.AccessToken
	expiry := time.Now().Add(time.Duration(int(tokenBody.ExpiresIn)) * time.Second)
	assert.False(t, time.Now().After(expiry), "expiry should be after now")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing alg is HMAC-SHA256:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return the key to be parsed with
		return []byte("keySecret"), nil
	})

	assert.NoError(t, err, "expect no error when parsing the token")
	assert.True(t, parsedToken.Valid, "token produced was invalid")
}

func TestAuthenticatorGetTokenNew(t *testing.T) {
	authenticator := newAuthenticator("testApp", "keyID", "keySecret")
	token, err := authenticator.getSUToken()
	assert.NoError(t, err, "expect no error")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing alg is HMAC-SHA256:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return the key to be parsed with
		return []byte("keySecret"), nil
	})

	assert.NoError(t, err, "expect no error when parsing the token")
	assert.True(t, parsedToken.Valid, "token produced was invalid")
}

func TestAuthenticatorGetTokenNotExpired(t *testing.T) {
	authenticator := newAuthenticator("testApp", "keyID", "keySecret")
	firstToken, err := authenticator.getSUToken()
	assert.NoError(t, err, "expect no error")

	secondToken, err := authenticator.getSUToken()
	assert.NoError(t, err, "expect no error")

	assert.Equal(t, firstToken, secondToken, "don't expect tokens to be regenerated if not expired")
}

func newMockAuthenticator() Authenticator {
	return &mockAuthenticator{}
}

type mockAuthenticator struct{}

func (ma *mockAuthenticator) getSUToken() (string, error) {
	return "", nil
}

func (ma *mockAuthenticator) authenticate(string) AuthenticationResponse {
	return AuthenticationResponse{
		Status:  200,
		Headers: map[string]string{},
		Body: TokenBody{
			AccessToken: "",
			TokenType:   "",
			ExpiresIn:   0,
		},
	}
}
