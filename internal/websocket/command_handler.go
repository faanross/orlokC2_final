package websocket

import (
	"fmt"
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
