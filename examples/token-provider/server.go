package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pusher/chatkit-server-go"
)

const grantType = "client_credentials"

// This is a very simple auth server which shows basic usage of the `Authenticate` method
func main() {
	http.HandleFunc(
		"/auth",
		func(w http.ResponseWriter, r *http.Request) {
			client, err := chatkit.NewClient(
				"your:instance:locator",
				"your:key",
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userID := "your-user-id"
			authRes, err := client.Authenticator().Authenticate(chatkit.AuthenticationPayload{
				GrantType: grantType,
			}, chatkit.AuthenticationOptions{
				UserID: &userID,
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			bytes, err := json.Marshal(authRes)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(bytes)
		})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
