package chatkitServerGo

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRolePermissionsSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	err := testClient.CreateRolePermissions("testRole", "testScope", RolePermissions{
		[]string{"testPermission"},
	})
	assert.Nil(t, err, "expected no error")
}

func TestCreateRolePermissionsFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	err := testClient.CreateRolePermissions("testRole", "testScope", RolePermissions{
		[]string{"testPermission"},
	})
	assert.Error(t, err, "expected an error")
}

func TestGetRolePermissionsFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	rolePerms, err := testClient.GetRolePermissions("testRole", "testScope")
	assert.Error(t, err, "expected an error")
	assert.Empty(t, rolePerms.Permissions, "Should be empty")
}

func TestGetRolePermissionsSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `["update_user", "create_user"]`)
	}))
	defer testServer.Close()

	rolePerms, err := testClient.GetRolePermissions("testRole", "testScope")
	assert.NoError(t, err, "expected an error")
	assert.NotEmpty(t, rolePerms.Permissions, "Should be empty")
}

func TestEditRolePermissionsFail(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	err := testClient.EditRolePermissions("testRole", "testScope", RolePermissions{
		[]string{"testPermission"},
	})
	assert.Error(t, err, "expected an error")
}

func TestEditRolePermissionsSuccess(t *testing.T) {
	testClient, testServer := newTestClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	err := testClient.EditRolePermissions("testRole", "testScope", RolePermissions{
		[]string{"testPermission"},
	})
	assert.NoError(t, err, "expected an error")
}
