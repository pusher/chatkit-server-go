package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	instanceLocator := os.Getenv("CHATKIT_INSTANCE_ID")
	key := os.Getenv("CHATKIT_KEY")
	if instanceLocator == "" || key == "" {
		log.Fatalln("Please set the CHATKIT_INSTANCE_ID and CHATKIT_KEY environment variables to run the example")
	}

	serverClient, err := chatkitServerGo.NewClient(instanceLocator, key)
	handleErr("Instatiating Client", err)

	log.Println("Creating User")
	newUser := chatkitServerGo.User{
		Name: "testUser",
		ID:   "testUser",
	}
	err = serverClient.CreateUser(newUser)
	handleErr("Creating User", err)

	log.Println("Getting User Roles")
	userRoles, err := serverClient.GetUserRoles("testUser")
	handleErr("Getting User Roles", err)
	for _, role := range userRoles {
		fmt.Println(role)
	}

	log.Println("Creating New Role")
	newRole := chatkitServerGo.Role{
		Name:        "testRole",
		Permissions: []string{"room:join", "room:leave", "message:create", "room:delete"},
		Scope:       "global",
	}
	err = serverClient.CreateRole(newRole)
	handleErr("Creating New Role", err)

	log.Println("Getting All User Roles")
	roles, err := serverClient.GetRoles()
	handleErr("Getting All User Roles", err)
	for _, role := range roles {
		fmt.Println(role)
	}

	log.Println("Assigning New Role To User")
	err = serverClient.UpdateUserRole("testUser", chatkitServerGo.UserRole{
		Name: "testRole",
	})
	handleErr("Assigning New Role To User", err)

	log.Println("Getting User Roles")
	userRoles, err = serverClient.GetUserRoles("testUser")
	handleErr("Getting User Roles", err)
	for _, role := range userRoles {
		fmt.Println(role)
	}

	log.Println("Delete User")
	err = serverClient.DeleteUser("testUser")
	handleErr("Delete User", err)

	log.Println("Delete Role")
	err = serverClient.DeleteRole("testRole", "global")
	handleErr("Delete Role", err)
}

func handleErr(descrip string, err error) {
	if err != nil {
		log.Fatalln("Error ", descrip, err)
	}
}
