package chatkit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
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

func createUser(
	client *Client,
) (string, error) {
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
	userID, err := createUser(client)
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
			userID, err := createUser(client)
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

				Convey("and to get that role", func() {
					roles, err := client.GetUserRoles(context.Background(), userID)
					So(err, ShouldBeNil)
					So(roles, ShouldContain, Role{
						Name:        globalRoleName,
						Permissions: globalPermissions,
						Scope:       "global",
					})
				})

				Convey("and remove it again", func() {
					err := client.RemoveGlobalRoleForUser(context.Background(), userID)
					So(err, ShouldBeNil)
				})
			})

			Convey("it should be possible to assign a room scoped role to a user", func() {
				err := client.AssignRoomRoleToUser(context.Background(), userID, room.ID, roomRoleName)
				So(err, ShouldBeNil)

				Convey("and to get that role", func() {
					roles, err := client.GetUserRoles(context.Background(), userID)
					So(err, ShouldBeNil)
					So(roles, ShouldContain, Role{
						Name:        roomRoleName,
						Permissions: roomPermissions,
						Scope:       "room",
					})
				})

				Convey("and remove it again", func() {
					err := client.RemoveRoomRoleForUser(context.Background(), userID, room.ID)
					So(err, ShouldBeNil)
				})
			})

		})

		Reset(func() {
			deleteAllResources(client)
		})
	})
}

func TestUsers(t *testing.T) {
	ctx := context.Background()

	config, err := getConfig()
	if err != nil {
		t.Fatalf("Failed to get test config: %s", err.Error())
	}

	client, err := NewClient(config.instanceLocator, config.key)
	if err != nil {
		t.Fatalf("Failed to create client: %s", err.Error())
	}

	Convey("We can create a user", t, func() {
		userID := randomString()
		avatarURL := "https://" + randomString()
		err := client.CreateUser(ctx, CreateUserOptions{
			ID:         userID,
			Name:       "integration-test-user",
			AvatarURL:  &avatarURL,
			CustomData: json.RawMessage([]byte(`{"foo":"hello","bar":42}`)),
		})
		So(err, ShouldBeNil)

		Convey("and we can can get them", func() {
			user, err := client.GetUser(ctx, userID)
			So(err, ShouldBeNil)
			So(user.ID, ShouldEqual, userID)
			So(user.Name, ShouldEqual, "integration-test-user")
			So(user.AvatarURL, ShouldEqual, avatarURL)
			So(len(user.CustomData), ShouldEqual, 2)
			So(user.CustomData["foo"], ShouldEqual, "hello")
			So(user.CustomData["bar"], ShouldEqual, 42)
		})

		Convey("and we can update them", func() {
			newName := randomString()
			newAvatarURL := "https://" + randomString()

			err := client.UpdateUser(ctx, userID, UpdateUserOptions{
				Name:       &newName,
				AvatarUrl:  &newAvatarURL,
				CustomData: json.RawMessage([]byte(`{"foo":"goodbye"}`)),
			})
			So(err, ShouldBeNil)

			Convey("and get them again", func() {
				user, err := client.GetUser(ctx, userID)
				So(err, ShouldBeNil)
				So(user.ID, ShouldEqual, userID)
				So(user.Name, ShouldEqual, newName)
				So(user.AvatarURL, ShouldEqual, newAvatarURL)
				So(len(user.CustomData), ShouldEqual, 1)
				So(user.CustomData["foo"], ShouldEqual, "goodbye")
			})
		})

		Convey("and we can delete them", func() {
			err := client.DeleteUser(ctx, userID)
			So(err, ShouldBeNil)

			Convey("and can't get them any more", func() {
				_, err := client.GetUser(ctx, userID)
				So(err.(*ErrorResponse).Status, ShouldEqual, 404)
				So(
					err.(*ErrorResponse).Info.(map[string]interface{})["error"],
					ShouldEqual,
					"services/chatkit/not_found/user_not_found",
				)
			})
		})
	})

	Convey("We can create a batch of users", t, func() {
		ids := []string{randomString(), randomString(), randomString(), randomString()}
		sort.Strings(ids)

		avatarURLs := make([]string, 4)
		for i, id := range ids {
			avatarURLs[i] = "https://avatars/" + id
		}

		err := client.CreateUsers(ctx, []CreateUserOptions{
			CreateUserOptions{
				ID:         ids[0],
				Name:       "Alice",
				AvatarURL:  &avatarURLs[0],
				CustomData: json.RawMessage(`{"a":"aaa"}`),
			},
			CreateUserOptions{
				ID:         ids[1],
				Name:       "Bob",
				AvatarURL:  &avatarURLs[1],
				CustomData: json.RawMessage(`{"b":"bbb"}`),
			},
			CreateUserOptions{
				ID:         ids[2],
				Name:       "Carol",
				AvatarURL:  &avatarURLs[2],
				CustomData: json.RawMessage(`{"c":"ccc"}`),
			},
			CreateUserOptions{
				ID:         ids[3],
				Name:       "Dave",
				AvatarURL:  &avatarURLs[3],
				CustomData: json.RawMessage(`{"d":"ddd"}`),
			},
		})
		So(err, ShouldBeNil)

		Convey("and get them all by ID", func() {
			users, err := client.GetUsersByID(ctx, ids)
			So(err, ShouldBeNil)

			sort.Slice(users, func(i, j int) bool {
				return users[i].ID < users[j].ID
			})

			So(len(users), ShouldResemble, 4)

			So(users[0].ID, ShouldEqual, ids[0])
			So(users[0].Name, ShouldEqual, "Alice")
			So(users[0].AvatarURL, ShouldEqual, avatarURLs[0])
			So(len(users[0].CustomData), ShouldEqual, 1)
			So(users[0].CustomData["a"], ShouldEqual, "aaa")

			So(users[1].ID, ShouldEqual, ids[1])
			So(users[1].Name, ShouldEqual, "Bob")
			So(users[1].AvatarURL, ShouldEqual, avatarURLs[1])
			So(len(users[1].CustomData), ShouldEqual, 1)
			So(users[1].CustomData["b"], ShouldEqual, "bbb")

			So(users[2].ID, ShouldEqual, ids[2])
			So(users[2].Name, ShouldEqual, "Carol")
			So(users[2].AvatarURL, ShouldEqual, avatarURLs[2])
			So(len(users[2].CustomData), ShouldEqual, 1)
			So(users[2].CustomData["c"], ShouldEqual, "ccc")

			So(users[3].ID, ShouldEqual, ids[3])
			So(users[3].Name, ShouldEqual, "Dave")
			So(users[3].AvatarURL, ShouldEqual, avatarURLs[3])
			So(len(users[3].CustomData), ShouldEqual, 1)
			So(users[3].CustomData["d"], ShouldEqual, "ddd")
		})

		Convey("and get them all (paginated)", func() {
			users, err := client.GetUsers(ctx, nil)
			So(err, ShouldBeNil)

			sort.Slice(users, func(i, j int) bool {
				return users[i].ID < users[j].ID
			})

			So(users[0].ID, ShouldEqual, ids[0])
			So(users[0].Name, ShouldEqual, "Alice")
			So(users[0].AvatarURL, ShouldEqual, avatarURLs[0])
			So(len(users[0].CustomData), ShouldEqual, 1)
			So(users[0].CustomData["a"], ShouldEqual, "aaa")

			So(users[1].ID, ShouldEqual, ids[1])
			So(users[1].Name, ShouldEqual, "Bob")
			So(users[1].AvatarURL, ShouldEqual, avatarURLs[1])
			So(len(users[1].CustomData), ShouldEqual, 1)
			So(users[1].CustomData["b"], ShouldEqual, "bbb")

			So(users[2].ID, ShouldEqual, ids[2])
			So(users[2].Name, ShouldEqual, "Carol")
			So(users[2].AvatarURL, ShouldEqual, avatarURLs[2])
			So(len(users[2].CustomData), ShouldEqual, 1)
			So(users[2].CustomData["c"], ShouldEqual, "ccc")

			So(users[3].ID, ShouldEqual, ids[3])
			So(users[3].Name, ShouldEqual, "Dave")
			So(users[3].AvatarURL, ShouldEqual, avatarURLs[3])
			So(len(users[3].CustomData), ShouldEqual, 1)
			So(users[3].CustomData["d"], ShouldEqual, "ddd")
		})

		Reset(func() {
			deleteAllResources(client)
		})
	})
}
