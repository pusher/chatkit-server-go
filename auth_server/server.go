package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pusher/chatkit-server-go"
)

func main() {
	http.HandleFunc(
    "/auth",
    func(w http.ResponseWriter, r *http.Request)
  {
		serverClient, _ := chatkit.NewClient(
			"your:instance:locator",
			"your:key",
		)

		authRes := serverClient.Authenticate("your-user-id")
		for key, value := range authRes.Headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(authRes.Status)
		bytes, _ := json.Marshal(authRes.Body)
		w.Write(bytes)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
