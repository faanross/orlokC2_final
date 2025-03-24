package router

import (
	"fmt"
	"net/http"
	"time"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Log a message with timestamp on the server side
	currentTime := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Printf("[%s] You hit the endpoint: %s\n", currentTime, r.URL.Path)

	// Send a response to the client
	w.Write([]byte("I'm Mister Derp!"))
}
