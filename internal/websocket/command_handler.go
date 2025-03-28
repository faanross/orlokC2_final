package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// In command_handler_euzqrkrx.go
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

// In command_handler_euzqrkrx.go
func ResultEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] Result endpoint hit\n",
		time.Now().Format("2006-01-02 15:04:05.000"))

	// Parse the incoming result
	var result struct {
		Command string `json:"command"`
		Output  string `json:"output"`
		Status  string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("[%s] Received result:\n  Command: %s\n  Status: %s\n  Output: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		result.Command,
		result.Status,
		result.Output)

	// Forward the result to the WebSocket clients
	if GlobalWSServer != nil {
		GlobalWSServer.SendCommandResult(result.Command, result.Output)
	}

	// Send acknowledgment
	w.WriteHeader(http.StatusOK)
}
