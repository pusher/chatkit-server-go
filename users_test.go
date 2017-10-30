package chatkit

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	}))
	defer testServer.Close()

	err := testClient.CreateUser(User{
		ID:   "123",
		Name: "dave",
	})
	assert.Nil(t, err, "expected no error")
}

func TestCreateUserFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	err := testClient.CreateUser(User{
		ID:   "123",
		Name: "dave",
	})
	assert.Error(t, err, "expected an error when the request returns 400")
}

func TestDeleteUserSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	}))
	defer testServer.Close()

	err := testClient.DeleteUser("123")
	assert.Nil(t, err, "expected no error")
}

func TestDeleteUserFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	err := testClient.DeleteUser("123")
	assert.Error(t, err, "expected an error when the request returns 400")
}
