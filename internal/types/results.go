package types

import (
	"time"
)

// CommandResult represents the result of a command execution
type CommandResult struct {
	Command   string    // The command that was executed
	Output    string    // The output of the command
	Timestamp time.Time // When the command was executed
	AgentUUID string    // UUID of the agent that executed the command
}
