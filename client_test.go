package chatkit

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pusher/pusher-platform-go/auth"
	platformclient "github.com/pusher/pusher-platform-go/client"
	. "github.com/smartystreets/goconvey/convey"
)

// Helpers

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// randomString generates a random string of length 10.
func randomString() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// config represents test config.
type config struct {
	instanceLocator string
	key             string
}

// GetConfig retrieves instance specific config from the ENV.
func getConfig() (*config, error) {
	instanceLocator := os.Getenv("TEST_INSTANCE_LOCATOR")
	if instanceLocator == "" {
		return nil, errors.New("TEST_INSTANCE_LOCATOR not set")
	}

	key := os.Getenv("TEST_KEY")
	if key == "" {
		return nil, errors.New("TEST_KEY not set")
	}

	return &config{instanceLocator, key}, nil
}

func createUser() (string, error) {
	userID := randomString()
	err := client.CreateUser(context.Background(), CreateUserOptions{
		ID:   userID,
		Name: "integration-test-user",
	})
	if err != nil {
		return "", err
	}

	return userID, nil
}

// createUserRoleWithGlobalPermissions generates a random user id, creates a user with that id
// and assigns globally scoped permissions to them.
// It returns the generated user id or an error.
func createUserWithGlobalPermissions(
	client *Client,
	permissions []string,
) (string, error) {
	userID, err := createUser()
	if err != nil {
		return userID, err
	}
	assignGlobalPermissionsToUser(client, userID, permissions)

	return userID, nil

}

// assignGlobalPermissionsToUser assigns permissions to a user
// and returns the role name
func assignGlobalPermissionsToUser(
	client *Client,
	userID string,
	permissions []string,
) (string, error) {
	roleName := randomString()

	err := client.CreateGlobalRole(context.Background(), CreateRoleOptions{
		Name:        roleName,
		Permissions: permissions,
	})
	if err != nil {
		return "", err
	}

	err = client.AssignGlobalRoleToUser(context.Background(), userID, roleName)
	if err != nil {
		return "", err
	}

	return roleName, nil
}

// DeleteResources deletes all resources associated with an instance.
// This allows tearing down resources after a test.
func deleteAllResources(client *Client) error {
	tokenWithExpiry, err := client.GenerateSUToken(auth.Options{})
	if err != nil {
		return err
	}

	_, err = client.CoreRequest(context.Background(), platformclient.RequestOptions{
		Method: http.MethodDelete,
		Path:   "/resources",
		Jwt:    &tokenWithExpiry.Token,
	})
	if err != nil {
		return err
	}

	return nil
}

func TestCursors(t *testing.T) {
	config, err := getConfig()
	if err != nil {
		t.Fatalf("Failed to get test config: %s", err.Error())
	}

	Convey("A chatkit user, having sent a message to a room", t, func() {
		client, err := NewClient(config.instanceLocator, config.key)
		So(err, ShouldBeNil)

		userID, err := createUserWithGlobalPermissions(client, []string{
			"cursors:read:set",
			"cursors:read:get",
			"room:create",
			"message:create",
		})
		So(err, ShouldBeNil)

		room, err := client.CreateRoom(context.Background(), CreateRoomOptions{
			Name:      randomString(),
			CreatorID: userID,
		})
		So(err, ShouldBeNil)

		messageID, err := client.SendMessage(context.Background(), SendMessageOptions{
			RoomID:   room.ID,
			Text:     "Hello!",
			SenderID: userID,
		})
		So(err, ShouldBeNil)

		Convey("and has set the read cursor", func() {
			err = client.SetReadCursor(context.Background(), userID, room.ID, messageID)
			So(err, ShouldBeNil)

			Convey("it should be possible to get back cursors for the user", func() {
				userCursors, err := client.GetUserReadCursors(context.Background(), userID)
				So(err, ShouldBeNil)

				So(userCursors[0].CursorType, ShouldEqual, 0)
				So(userCursors[0].RoomID, ShouldEqual, room.ID)
				So(userCursors[0].UserID, ShouldEqual, userID)
				So(userCursors[0].Position, ShouldEqual, messageID)
			})

			Convey("it should be possible to get back cursors for a room", func() {
				roomCursors, err := client.GetReadCursorsForRoom(context.Background(), room.ID)
				So(err, ShouldBeNil)

				So(roomCursors[0].CursorType, ShouldEqual, 0)
				So(roomCursors[0].RoomID, ShouldEqual, room.ID)
				So(roomCursors[0].UserID, ShouldEqual, userID)
				So(roomCursors[0].Position, ShouldEqual, messageID)
			})
		})

		Convey("On sending a new message and setting the read cursor", func() {
			latestMessageID, err := client.SendMessage(context.Background(), SendMessageOptions{
				RoomID:   room.ID,
				Text:     "Hello!",
				SenderID: userID,
			})
			So(err, ShouldBeNil)

			err = client.SetReadCursor(context.Background(), userID, room.ID, latestMessageID)
			So(err, ShouldBeNil)

			Convey("it should be possible to get the latest read cursor for a user and room", func() {
				cursor, err := client.GetReadCursor(context.Background(), userID, room.ID)
				So(err, ShouldBeNil)

				So(cursor.CursorType, ShouldEqual, 0)
				So(cursor.RoomID, ShouldEqual, room.ID)
				So(cursor.UserID, ShouldEqual, userID)
				So(cursor.Position, ShouldEqual, latestMessageID)
			})

			Convey("it should be possible to make a raw request to cursors", func() {
				tokenWithExpiry, err := client.GenerateSUToken(AuthenticateOptions{})
				So(err, ShouldBeNil)

				// Get back last set cursor
				response, err := client.CursorsRequest(context.Background(), RequestOptions{
					Method: http.MethodGet,
					Path:   fmt.Sprintf("/cursors/0/rooms/%d/users/%s", room.ID, userID),
					Jwt:    &tokenWithExpiry.Token,
				})
				So(err, ShouldBeNil)

				So(response.StatusCode, ShouldEqual, http.StatusOK)
			})
		})

		Reset(func() {
			deleteAllResources(client)
		})
	})
}

func TestAuthorizer(t *testing.T) {
	config, err := getConfig()
	if err != nil {
		t.Fatalf("Failed to get test config: %s", err.Error())
	}

	Convey("Given global and room roles have been created", t, func() {
		client, err := NewClient(config.instanceLocator, config.key)
		So(err, ShouldBeNil)

		globalRoleName := randomString()
		globalPermissions := []string{"message:create", "room:create"}
		err = client.CreateGlobalRole(context.Background(), CreateRoleOptions{
			Name:        globalRoleName,
			Permissions: globalPermissions,
		})
		So(err, ShouldBeNil)

		roomRoleName := randomString()
		roomPermissions := []string{"message:create"}
		err = client.CreateRoomRole(context.Background(), CreateRoleOptions{
			Name:        roomRoleName,
			Permissions: roomPermissions,
		})

		Convey("it should be possible to fetch them", func() {
			roles, err := client.GetRoles(context.Background())
			So(err, ShouldBeNil)

			So(roles, ShouldHaveLength, 2)
			So(roles, ShouldContain, Role{
				Name:        globalRoleName,
				Permissions: globalPermissions,
				Scope:       "global",
			})

			So(roles, ShouldContain, Role{
				Name:        roomRoleName,
				Permissions: roomPermissions,
				Scope:       "room",
			})
		})

		Convey("it should be possible to delete a global scoped role", func() {
			err := client.DeleteGlobalRole(context.Background(), globalRoleName)
			So(err, ShouldBeNil)
		})

		Convey("it should be possible to delete a room scoped role", func() {
			err := client.DeleteRoomRole(context.Background(), roomRoleName)
			So(err, ShouldBeNil)
		})

		Convey("it should be possible to retreive permissions for a global scoped role", func() {
			permissions, err := client.GetPermissionsForGlobalRole(context.Background(), globalRoleName)
			So(err, ShouldBeNil)
			So(permissions, ShouldResemble, globalPermissions)
		})

		Convey("it should be possible to retrieve permissions for a room scoped role", func() {
			permissions, err := client.GetPermissionsForRoomRole(context.Background(), roomRoleName)
			So(err, ShouldBeNil)
			So(permissions, ShouldResemble, roomPermissions)
		})

		Convey("it should be possible to update permissions for a globally scoped role", func() {
			err := client.UpdatePermissionsForGlobalRole(
				context.Background(),
				globalRoleName,
				UpdateRolePermissionsOptions{
					PermissionsToAdd:    []string{"cursors:read:set", "cursors:read:get"},
					PermissionsToRemove: []string{"room:create"},
				},
			)
			So(err, ShouldBeNil)

			permissions, err := client.GetPermissionsForGlobalRole(context.Background(), globalRoleName)
			So(err, ShouldBeNil)
			So(permissions, ShouldResemble, []string{"message:create", "cursors:read:set", "cursors:read:get"})
		})

		Convey("it should be possible to update permissions for a room scoped role", func() {
			err := client.UpdatePermissionsForRoomRole(
				context.Background(),
				roomRoleName,
				UpdateRolePermissionsOptions{
					PermissionsToAdd:    []string{"cursors:read:set", "cursors:read:get"},
					PermissionsToRemove: []string{"message:create"},
				},
			)
			So(err, ShouldBeNil)

			permissions, err := client.GetPermissionsForRoomRole(context.Background(), roomRoleName)
			So(err, ShouldBeNil)
			So(permissions, ShouldResemble, []string{"cursors:read:set", "cursors:read:get"})
		})

		Convey("on creating a user, a room and a global and room role", func() {
			userID, err := createUser()
			So(err, ShouldBeNil)

			globalRoleName := randomString()
			globalPermissions := []string{"message:create", "room:create"}
			err = client.CreateGlobalRole(context.Background(), CreateRoleOptions{
				Name:        globalRoleName,
				Permissions: globalPermissions,
			})
			So(err, ShouldBeNil)

			room, err := client.CreateRoom(context.Background(), CreateRoomOptions{
				Name:      randomString(),
				CreatorID: userID,
			})
			So(err, ShouldBeNil)

			roomRoleName := randomString()
			roomPermissions := []string{"message:create"}
			err = client.CreateRoomRole(context.Background(), CreateRoleOptions{
				Name:        roomRoleName,
				Permissions: roomPermissions,
			})
			So(err, ShouldBeNil)

			Convey("it should be possible to assign a global scoped role to a user", func() {
				err := client.AssignGlobalRoleToUser(context.Background(), userID, globalRoleName)
				So(err, ShouldBeNil)
			})

			Convey("it should be possible to assign a room scoped role to a user", func() {
				err := client.AssignRoomRoleToUser(context.Background(), userID, room.ID, roomRoleName)
				So(err, ShouldBeNil)
			})

			Convey("it should be possible to get roles for a user", func() {

			})
		})

		Reset(func() {
			deleteAllResources(client)
		})
	})
}
