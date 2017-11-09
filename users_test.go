package chatkit

import (
	"fmt"
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

func TestGetUsersSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
			{
				"id": "123abc",
				"name": "gary",
				"avatar_url": "https://gravatar.com/img/8819",
				"custom_data": {
				"email": "a@b.com"
				},
				"created_at":"2017-03-23T11:36:42Z",
				"updated_at":"2017-03-23T11:36:42Z"
			},
			{
				"id": "john",
				"name": "John Doe",
				"avatar_url": "https://gravatar.com/img/8819",
				"custom_data": {
				"email": "a@b.com"
				},
				"created_at":"2017-03-23T11:36:42Z",
				"updated_at":"2017-03-23T11:36:42Z"
			}
		]`)
	}))
	defer testServer.Close()

	users, err := testClient.GetUsers()
	assert.Nil(t, err, "expected no error")
	assert.Len(t, users, 2, "expected to contain 2 elements")
	assert.Equal(t, "gary", users[0].Name, "expected first user element to have Name of gary")
}

func TestGetUsersFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	users, err := testClient.GetUsers()
	assert.Error(t, err, "expected an error when the request returns 400")
	assert.Nil(t, users)
}
