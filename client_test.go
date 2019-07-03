package chatkit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
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
	instanceLocator := os.Getenv("CHATKIT_INSTANCE_LOCATOR")
	if instanceLocator == "" {
		return nil, errors.New("CHATKIT_INSTANCE_LOCATOR not set")
	}

	key := os.Getenv("CHATKIT_INSTANCE_KEY")
	if key == "" {
		return nil, errors.New("CHATKIT_INSTANCE_KEY not set")
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

func interfaceToSliceOfInterfaces(s interface{}) []interface{} {
	t := reflect.ValueOf(s)
	if t.Kind() != reflect.Slice {
		panic("can't coerce non slice to []interface{}")
	}
	u := make([]interface{}, t.Len())
	for i := 0; i < t.Len(); i++ {
		u[i] = t.Index(i).Interface()
	}
	return u
}

func containsResembling(s []interface{}, x interface{}) bool {
	for _, y := range s {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}

// assumes neither actual nor expected contain repetitions
func shouldResembleUpToReordering(
	actual interface{},
	expected ...interface{},
) string {
	s := interfaceToSliceOfInterfaces(actual)
	t := interfaceToSliceOfInterfaces(expected[0])
	if len(s) != len(t) {
		return fmt.Sprintf("%s and %s are not the same length!", s, t)
	}
	for _, x := range s {
		if !containsResembling(t, x) {
			return fmt.Sprintf("%s contains %s, but %s does not", s, x, t)
		}
	}
	return ""
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
					Path:   fmt.Sprintf("/cursors/0/rooms/%s/users/%s", room.ID, userID),
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
			So(permissions, shouldResembleUpToReordering, globalPermissions)
		})

		Convey("it should be possible to retrieve permissions for a room scoped role", func() {
			permissions, err := client.GetPermissionsForRoomRole(context.Background(), roomRoleName)
			So(err, ShouldBeNil)
			So(permissions, shouldResembleUpToReordering, roomPermissions)
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
			So(
				permissions,
				shouldResembleUpToReordering,
				[]string{"message:create", "cursors:read:set", "cursors:read:get"},
			)
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
			So(
				permissions,
				shouldResembleUpToReordering,
				[]string{"cursors:read:set", "cursors:read:get"},
			)
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

			So(len(users), ShouldEqual, 4)

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

func TestRooms(t *testing.T) {
	ctx := context.Background()

	config, err := getConfig()
	if err != nil {
		t.Fatalf("Failed to get test config: %s", err.Error())
	}

	client, err := NewClient(config.instanceLocator, config.key)
	if err != nil {
		t.Fatalf("Failed to create client: %s", err.Error())
	}

	Convey("Given some users", t, func() {
		aliceID, err := createUser(client)
		So(err, ShouldBeNil)

		bobID, err := createUser(client)
		So(err, ShouldBeNil)

		carolID, err := createUser(client)
		So(err, ShouldBeNil)

		Convey("we can create a room without providing an ID", func() {
			roomName := randomString()

			room, err := client.CreateRoom(ctx, CreateRoomOptions{
				Name:       roomName,
				Private:    true,
				UserIDs:    []string{aliceID, bobID},
				CreatorID:  aliceID,
				CustomData: map[string]interface{}{"foo": "bar"},
			})
			So(err, ShouldBeNil)
			So(room.Name, ShouldEqual, roomName)
			So(room.Private, ShouldEqual, true)
			So(room.MemberUserIDs, shouldResembleUpToReordering, []string{aliceID, bobID})

			Convey("and get it", func() {
				r, err := client.GetRoom(ctx, room.ID)
				So(err, ShouldBeNil)
				So(r.ID, ShouldEqual, room.ID)
				So(r.Name, ShouldEqual, roomName)
				So(r.Private, ShouldEqual, true)
				So(r.MemberUserIDs, shouldResembleUpToReordering, []string{aliceID, bobID})
				So(r.CustomData, ShouldResemble, map[string]interface{}{"foo": "bar"})
			})

			Convey("and update it", func() {
				newRoomName := randomString()
				err := client.UpdateRoom(ctx, room.ID, UpdateRoomOptions{
					Name:       &newRoomName,
					CustomData: map[string]interface{}{"foo": "baz"},
				})
				So(err, ShouldBeNil)

				Convey("and get it again", func() {
					r, err := client.GetRoom(ctx, room.ID)
					So(err, ShouldBeNil)
					So(r.ID, ShouldEqual, room.ID)
					So(r.Name, ShouldEqual, newRoomName)
					So(r.Private, ShouldEqual, true)
					So(r.MemberUserIDs, shouldResembleUpToReordering, []string{aliceID, bobID})
					So(r.CustomData, ShouldResemble, map[string]interface{}{"foo": "baz"})
				})
			})

			Convey("and delete it", func() {
				err := client.DeleteRoom(ctx, room.ID)
				So(err, ShouldBeNil)

				Convey("and can't get it any more (via GetRooms)", func() {
					rooms, err := client.GetRooms(ctx, GetRoomsOptions{})
					So(err, ShouldBeNil)
					So(len(rooms), ShouldEqual, 0)
				})
			})

			Convey("and add users to it", func() {
				err := client.AddUsersToRoom(ctx, room.ID, []string{carolID})
				So(err, ShouldBeNil)

				r, err := client.GetRoom(ctx, room.ID)
				So(err, ShouldBeNil)
				So(r.MemberUserIDs, shouldResembleUpToReordering, []string{aliceID, bobID, carolID})
			})

			Convey("and remove users from it", func() {
				err := client.RemoveUsersFromRoom(ctx, room.ID, []string{bobID})
				So(err, ShouldBeNil)

				r, err := client.GetRoom(ctx, room.ID)
				So(err, ShouldBeNil)
				So(r.MemberUserIDs, shouldResembleUpToReordering, []string{aliceID})
			})
		})

		Convey("we can create a room providing an ID", func() {
			roomID := randomString()
			roomName := randomString()

			room, err := client.CreateRoom(ctx, CreateRoomOptions{
				ID:         &roomID,
				Name:       roomName,
				Private:    true,
				UserIDs:    []string{aliceID, bobID},
				CreatorID:  aliceID,
				CustomData: map[string]interface{}{"foo": "bar"},
			})
			So(err, ShouldBeNil)
			So(room.ID, ShouldEqual, roomID)
			So(room.Name, ShouldEqual, roomName)
			So(room.Private, ShouldEqual, true)
			So(room.MemberUserIDs, shouldResembleUpToReordering, []string{aliceID, bobID})

			Convey("and get it", func() {
				r, err := client.GetRoom(ctx, room.ID)
				So(err, ShouldBeNil)
				So(r.ID, ShouldEqual, room.ID)
				So(r.Name, ShouldEqual, roomName)
				So(r.Private, ShouldEqual, true)
				So(r.MemberUserIDs, shouldResembleUpToReordering, []string{aliceID, bobID})
				So(r.CustomData, ShouldResemble, map[string]interface{}{"foo": "bar"})
			})
		})

		Convey("we can create a couple of rooms", func() {
			room1, err := client.CreateRoom(ctx, CreateRoomOptions{
				Name:      randomString(),
				UserIDs:   []string{aliceID, bobID},
				CreatorID: aliceID,
			})
			So(err, ShouldBeNil)

			room2, err := client.CreateRoom(ctx, CreateRoomOptions{
				Name:      randomString(),
				UserIDs:   []string{aliceID},
				CreatorID: aliceID,
			})
			So(err, ShouldBeNil)

			Convey("and get them", func() {
				rooms, err := client.GetRooms(ctx, GetRoomsOptions{})
				So(err, ShouldBeNil)
				So(len(rooms), ShouldEqual, 2)
				So(
					[]string{rooms[0].ID, rooms[1].ID},
					shouldResembleUpToReordering,
					[]string{room1.ID, room2.ID},
				)
			})

			Convey("and get a user's rooms", func() {
				rooms, err := client.GetUserRooms(ctx, bobID)
				So(err, ShouldBeNil)
				So(len(rooms), ShouldEqual, 1)
				So(rooms[0].ID, ShouldEqual, room1.ID)
			})

			Convey("and get a user's joinable rooms", func() {
				rooms, err := client.GetUserJoinableRooms(ctx, bobID)
				So(err, ShouldBeNil)
				So(len(rooms), ShouldEqual, 1)
				So(rooms[0].ID, ShouldEqual, room2.ID)
			})
		})

		Reset(func() {
			deleteAllResources(client)
		})
	})
}

func TestMessages(t *testing.T) {
	ctx := context.Background()

	config, err := getConfig()
	if err != nil {
		t.Fatalf("Failed to get test config: %s", err.Error())
	}

	client, err := NewClient(config.instanceLocator, config.key)
	if err != nil {
		t.Fatalf("Failed to create client: %s", err.Error())
	}

	Convey("Given a user and a room", t, func() {
		userID, err := createUser(client)
		So(err, ShouldBeNil)

		room, err := client.CreateRoom(ctx, CreateRoomOptions{
			Name:      randomString(),
			CreatorID: userID,
		})
		So(err, ShouldBeNil)

		Convey("we can publish messages from that user", func() {
			messageID1, err := client.SendMessage(ctx, SendMessageOptions{
				RoomID:   room.ID,
				Text:     "one",
				SenderID: userID,
			})
			So(err, ShouldBeNil)

			messageID2, err := client.SendMessage(ctx, SendMessageOptions{
				RoomID:   room.ID,
				Text:     "two",
				SenderID: userID,
			})
			So(err, ShouldBeNil)

			messageID3, err := client.SendMultipartMessage(ctx, SendMultipartMessageOptions{
				RoomID:   room.ID,
				SenderID: userID,
				Parts: []NewPart{
					NewInlinePart{Type: "text/plain", Content: "three"},
				}})
			So(err, ShouldBeNil)

			messageID4, err := client.SendSimpleMessage(ctx, SendSimpleMessageOptions{
				RoomID:   room.ID,
				Text:     "four",
				SenderID: userID,
			})
			So(err, ShouldBeNil)

			Convey("and fetch them", func() {
				limit := uint(2)
				messagesPage1, err := client.GetRoomMessages(ctx, room.ID, GetRoomMessagesOptions{
					Limit: &limit,
				})
				So(err, ShouldBeNil)
				So(len(messagesPage1), ShouldEqual, 2)
				So(messagesPage1[0].ID, ShouldEqual, messageID4)
				So(messagesPage1[0].Text, ShouldEqual, "four")
				So(messagesPage1[1].ID, ShouldEqual, messageID3)
				So(messagesPage1[1].Text, ShouldEqual, "three")

				messagesPage2, err := client.GetRoomMessages(ctx, room.ID, GetRoomMessagesOptions{
					InitialID: &messageID3,
				})
				So(err, ShouldBeNil)
				So(len(messagesPage2), ShouldEqual, 2)
				So(messagesPage2[0].ID, ShouldEqual, messageID2)
				So(messagesPage2[0].Text, ShouldEqual, "two")
				So(messagesPage2[1].ID, ShouldEqual, messageID1)
				So(messagesPage2[1].Text, ShouldEqual, "one")
			})

			Convey("and delete one of them", func() {
				err := client.DeleteMessage(ctx, DeleteMessageOptions{
					RoomID:    room.ID,
					MessageID: messageID3,
				})
				So(err, ShouldBeNil)

				limit := uint(4)
				messages, err := client.GetRoomMessages(ctx, room.ID, GetRoomMessagesOptions{
					Limit: &limit,
				})
				So(err, ShouldBeNil)
				So(len(messages), ShouldEqual, 3)
				So(messages[0].ID, ShouldEqual, messageID4)
				So(messages[0].Text, ShouldEqual, "four")
				So(messages[1].ID, ShouldEqual, messageID2)
				So(messages[1].Text, ShouldEqual, "two")
				So(messages[2].ID, ShouldEqual, messageID1)
				So(messages[2].Text, ShouldEqual, "one")
			})
		})

		Convey("we can publish a multipart messages", func() {
			fileName := "cat.jpg"
			file, err := os.Open(fileName)
			So(err, ShouldBeNil)
			defer file.Close()

			messageID, err := client.SendMultipartMessage(ctx, SendMultipartMessageOptions{
				RoomID:   room.ID,
				SenderID: userID,
				Parts: []NewPart{
					NewInlinePart{Type: "text/plain", Content: "see attached"},
					NewURLPart{Type: "audio/ogg", URL: "https://example.com/audio.ogg"},
					NewURLPart{Type: "audio/ogg", URL: "https://example.com/audio2.ogg"},
					NewAttachmentPart{
						Type:       "application/json",
						File:       strings.NewReader(`{"hello":"world"}`),
						CustomData: "anything",
					},
					NewAttachmentPart{
						Type: "image/png",
						File: file,
						Name: &fileName,
					},
				},
			})
			So(err, ShouldBeNil)

			Convey("and fetch it (v2)", func() {
				limit := uint(1)
				messages, err := client.GetRoomMessages(ctx, room.ID, GetRoomMessagesOptions{
					Limit: &limit,
				})
				So(err, ShouldBeNil)
				So(len(messages), ShouldEqual, 1)
				So(messages[0].ID, ShouldEqual, messageID)
				So(
					messages[0].Text,
					ShouldEqual,
					"You have received a message which can't be represented in this version of the app. You will need to upgrade to read it.",
				)
			})

			Convey("and fetch it (v6)", func() {
				limit := uint(1)
				messages, err := client.FetchMultipartMessages(
					ctx,
					room.ID,
					FetchMultipartMessagesOptions{Limit: &limit},
				)
				So(err, ShouldBeNil)
				So(len(messages), ShouldEqual, 1)
				So(messages[0].ID, ShouldEqual, messageID)
				So(len(messages[0].Parts), ShouldEqual, 5)

				So(messages[0].Parts[0].Type, ShouldEqual, "text/plain")
				So(messages[0].Parts[0].Content, ShouldNotBeNil)
				So(*messages[0].Parts[0].Content, ShouldEqual, "see attached")
				So(messages[0].Parts[0].URL, ShouldBeNil)
				So(messages[0].Parts[0].Attachment, ShouldBeNil)

				So(messages[0].Parts[1].Type, ShouldEqual, "audio/ogg")
				So(messages[0].Parts[1].Content, ShouldBeNil)
				So(messages[0].Parts[1].URL, ShouldNotBeNil)
				So(*messages[0].Parts[1].URL, ShouldEqual, "https://example.com/audio.ogg")
				So(messages[0].Parts[1].Attachment, ShouldBeNil)

				So(messages[0].Parts[2].Type, ShouldEqual, "audio/ogg")
				So(messages[0].Parts[2].Content, ShouldBeNil)
				So(messages[0].Parts[2].URL, ShouldNotBeNil)
				So(*messages[0].Parts[2].URL, ShouldEqual, "https://example.com/audio2.ogg")
				So(messages[0].Parts[2].Attachment, ShouldBeNil)

				So(messages[0].Parts[3].Type, ShouldEqual, "application/json")
				So(messages[0].Parts[3].Content, ShouldBeNil)
				So(messages[0].Parts[3].URL, ShouldBeNil)
				So(messages[0].Parts[3].Attachment, ShouldNotBeNil)
				So(messages[0].Parts[3].Attachment.RefreshURL, ShouldNotEqual, "")
				So(messages[0].Parts[3].Attachment.Expiration, ShouldNotEqual, time.Time{})
				So(messages[0].Parts[3].Attachment.Name, ShouldNotEqual, "")
				So(messages[0].Parts[3].Attachment.Size, ShouldEqual, 17)
				So(messages[0].Parts[3].Attachment.CustomData, ShouldEqual, "anything")
				res, err := http.Get(messages[0].Parts[3].Attachment.DownloadURL)
				So(err, ShouldBeNil)
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				So(err, ShouldBeNil)
				So(string(body), ShouldEqual, `{"hello":"world"}`)

				So(messages[0].Parts[4].Type, ShouldEqual, "image/png")
				So(messages[0].Parts[4].Content, ShouldBeNil)
				So(messages[0].Parts[4].URL, ShouldBeNil)
				So(messages[0].Parts[4].Attachment, ShouldNotBeNil)
				So(messages[0].Parts[4].Attachment.RefreshURL, ShouldNotEqual, "")
				So(messages[0].Parts[4].Attachment.Expiration, ShouldNotEqual, time.Time{})
				So(messages[0].Parts[4].Attachment.Name, ShouldEqual, fileName)
				So(messages[0].Parts[4].Attachment.Size, ShouldEqual, 44043)
				So(messages[0].Parts[4].Attachment.CustomData, ShouldBeNil)
				So(messages[0].Parts[4].Attachment.DownloadURL, ShouldNotEqual, "")
			})

			Convey("and delete it", func() {
				err := client.DeleteMessage(ctx, DeleteMessageOptions{
					RoomID:    room.ID,
					MessageID: messageID,
				})
				So(err, ShouldBeNil)

				limit := uint(1)
				messages, err := client.GetRoomMessages(ctx, room.ID, GetRoomMessagesOptions{
					Limit: &limit,
				})
				So(err, ShouldBeNil)
				So(len(messages), ShouldEqual, 0)
			})
		})

		Reset(func() {
			deleteAllResources(client)
		})
	})
}
