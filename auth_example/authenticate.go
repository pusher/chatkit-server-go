package main

import (
	"log"
	"os"

	"github.com/pusher/chatkit-server-go"
)

func main() {
	instanceLocator := os.Getenv("CHATKIT_INSTANCE_LOCATOR")
	key := os.Getenv("CHATKIT_KEY")
	if instanceLocator == "" || key == "" {
		log.Fatalln("Please set the CHATKIT_INSTANCE_LOCATOR and CHATKIT_KEY environment variables to run the example")
	}

	serverClient, err := chatkit.NewClient(instanceLocator, key)
	if err != nil {
		log.Fatalf("Error instantiating client: %s", err.Error())
	}

	log.Println("Authenticating")
	authRes := serverClient.Authenticate("ham")
	log.Println(authRes.Status, authRes.Headers, authRes.Body)
}
