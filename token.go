package chatkitServerGo

import (
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// NewChatkitSUToken produces a new jwt su token that will be compatible with Chatkit
func NewChatkitSUToken(appID string, keyID string, keySecret string, expiryDuration time.Duration) (tokenString string, expiry time.Time, err error) {
	jwtClaims, tokenExpiry := getGenericTokenClaims(appID, keyID, expiryDuration)

	jwtClaims["su"] = true

	tokenString, err = signToken(keySecret, jwtClaims)
	return tokenString, tokenExpiry, err
}

// NewChatkitUserToken produces a new jwt user token that will be compatible with Chatkit
func NewChatkitUserToken(appID string, keyID string, keySecret string, userID string, expiryDuration time.Duration) (tokenString string, expiry time.Time, err error) {
	jwtClaims, tokenExpiry := getGenericTokenClaims(appID, keyID, expiryDuration)

	jwtClaims["sub"] = userID

	tokenString, err = signToken(keySecret, jwtClaims)
	return tokenString, tokenExpiry, err
}

func getGenericTokenClaims(appID string, keyID string, expiryDuration time.Duration) (jwtClaims jwt.MapClaims, tokenExpiry time.Time) {
	timeNow := time.Now()
	tokenExpiry = timeNow.Add(expiryDuration)

	jwtClaims = jwt.MapClaims{
		"app": appID,
		"iss": "api_keys/" + keyID,
		"iat": timeNow.Unix(),
		"exp": tokenExpiry.Unix(),
	}

	return jwtClaims, tokenExpiry
}

func signToken(keySecret string, jwtClaims jwt.MapClaims) (tokenString string, err error) {
	// Create a new access token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	// Sign using the keySecret and get the complete encoded token as a string
	tokenString, err = token.SignedString([]byte(keySecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type tokenManager interface {
	getToken() (string, error)
}

type suTokenManager struct {
	tokenExpiry time.Time
	token       string
	mutex       sync.Mutex

	appID     string
	keyID     string
	keySecret string
}

func newTokenManager(appID string, keyID string, keySecret string) tokenManager {
	return &suTokenManager{
		tokenExpiry: time.Now().Add(-time.Minute),
		mutex:       sync.Mutex{},

		appID:     appID,
		keyID:     keyID,
		keySecret: keySecret,
	}
}

func (stm *suTokenManager) getToken() (string, error) {
	stm.mutex.Lock()
	defer stm.mutex.Unlock()

	if time.Now().After(stm.tokenExpiry) {
		tokenString, tokenExpiry, err := NewChatkitSUToken(stm.appID, stm.keyID, stm.keySecret, time.Hour*24)
		if err != nil {
			return "", err
		}
		stm.tokenExpiry = tokenExpiry
		stm.token = tokenString
	}
	return stm.token, nil
}
