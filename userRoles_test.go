package chatkitServerGo

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserRolesSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{
			  "role_name": "admin",
			  "permissions": ["create_user", "update_user"],
			  "scope": "global"
			}
		  ]`)
	}))
	defer testServer.Close()

	roles, err := testClient.GetUserRoles("userID")
	assert.NotEmpty(t, roles, "Should be empty")
	assert.Nil(t, err, "expected no error")
}

func TestGetUserRolesFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	roles, err := testClient.GetUserRoles("userID")
	assert.Empty(t, roles, "Should be empty")
	assert.NotNil(t, err, "expected no error")
}

func TestCreateUserRolesSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	err := testClient.CreateUserRole("testUserID", UserRole{
		Name:   "testName",
		RoomID: 123,
	})
	assert.Nil(t, err, "expected no error")
}

func TestCreateUserRolesFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	err := testClient.CreateUserRole("testUserID", UserRole{
		Name:   "testName",
		RoomID: 123,
	})
	assert.NotNil(t, err, "expected no error")
}

func TestUpdateUserRolesSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	err := testClient.UpdateUserRole("testUserID", UserRole{
		Name:   "testName",
		RoomID: 123,
	})
	assert.Nil(t, err, "expected no error")
}

func TestUpdateUserRolesFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	err := testClient.UpdateUserRole("testUserID", UserRole{
		Name:   "testName",
		RoomID: 123,
	})
	assert.NotNil(t, err, "expected no error")
}

func TestDeleteUserRolesSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	err := testClient.DeleteUserRole("testUserID", nil)
	assert.Nil(t, err, "expected no error")
}

func TestDeleteUserRolesWithRoomIDSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	roomID := "testRoomID"
	err := testClient.DeleteUserRole("testUserID", &roomID)
	assert.Nil(t, err, "expected no error")
}

func TestDeleteUserRolesFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	err := testClient.DeleteUserRole("testUserID", nil)
	assert.NotNil(t, err, "expected no error")
}
