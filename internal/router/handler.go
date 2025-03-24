package router

import (
	"fmt"
	"net/http"
	"time"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {

	// Get the agent UUID from the request context
	agentUUID, _ := r.Context().Value(AgentUUIDKey).(string)

	// Log a message with timestamp on the server side
	currentTime := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Printf("[%s] Endpoint %s hit by agent %s\n", currentTime, r.URL.Path, agentUUID)

	// Send a response to the client
	w.Write([]byte("I'm Mister Derp!"))
}
