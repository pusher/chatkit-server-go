package chatkit

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRoleSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	testRole := Role{
		Name:        "test_role",
		Permissions: []string{"join_room", "leave_room"},
		Scope:       "global",
	}

	err := testClient.CreateRole(testRole)
	assert.Nil(t, err, "expected no error")
}

func TestCreateRoleFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	testRole := Role{
		Name:        "test_role",
		Permissions: []string{"join_room", "leave_room"},
		Scope:       "global",
	}

	err := testClient.CreateRole(testRole)
	assert.Error(t, err, "expected an error on non 2xx return code")
}

func TestGetRolesSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w,
			`[
				{
				"role_name": "admin",
				"permissions": ["create_user", "update_user"],
				"scope": "global"
				}
			]`,
		)
	}))
	defer testServer.Close()

	roles, err := testClient.GetRoles()
	assert.Nil(t, err, "expected no error")
	assert.NotEmpty(t, roles, "Should return a non empty type")
}

func TestGetRolesFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	roles, err := testClient.GetRoles()
	assert.Error(t, err, "expected an error")
	assert.Empty(t, roles, "Should return a empty []Roles slice")
}

func TestDeleteRolesInputFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	err := testClient.DeleteRole("", "")
	assert.Error(t, err, "expected an error")
}

func TestDeleteRolesResponseFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	err := testClient.DeleteRole("a", "a")
	assert.Error(t, err, "expected an error")
}

func TestDeleteRolesSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	err := testClient.DeleteRole("a", "a")
	assert.NoError(t, err, "expected an error")
}
