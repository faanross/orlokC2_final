package router

import (
	"fmt"
	"net/http"
	"orlokC2_final/internal/middleware"
	"orlokC2_final/internal/types"
	"time"
)

var currentTime = time.Now().Format("2006-01-02 15:04:05.000")

func RootHandler(w http.ResponseWriter, r *http.Request) {

	// Get the agent UUID from the request context
	agentUUID, _ := r.Context().Value(middleware.AgentUUIDKey).(string)

	// Log a message with timestamp on the server side
	fmt.Printf("[%s] Endpoint %s hit by agent %s\n", currentTime, r.URL.Path, agentUUID)

	// Send a response to the client
	w.Write([]byte("I'm Mister Derp!"))
}

// ResultsHandler processes command results from agents
func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the agent UUID from the request context
	agentUUID, _ := r.Context().Value(middleware.AgentUUIDKey).(string)

	// Get the result from the request context (added by middleware)
	result, ok := r.Context().Value(middleware.ResultKey).(types.CommandResult)
	if !ok {
		http.Error(w, "Invalid result data", http.StatusBadRequest)
		return
	}

	// Log the result
	fmt.Printf("[%s] Result from agent %s:\n", currentTime, agentUUID)
	fmt.Printf("  Command: %s\n", result.Command)
	fmt.Printf("  Output: %s\n", result.Output)
	fmt.Printf("  Executed: %s\n", result.Timestamp.Format("2006-01-02 15:04:05.000"))

	// Respond with success
	w.WriteHeader(http.StatusOK)
}
