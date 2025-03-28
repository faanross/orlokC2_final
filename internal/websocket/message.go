package websocket

// MessageType defines the types of messages exchanged
type MessageType string

const (
	// Command is sent from UI to server to execute on agent
	CommandMessage MessageType = "command"
	// Response is sent from agent to UI with command results
	ResponseMessage MessageType = "response"
)

// Message represents the general message structure
type Message struct {
	Type    MessageType `json:"type"`
	Command string      `json:"command,omitempty"`
	Output  string      `json:"output,omitempty"`
	Status  string      `json:"status,omitempty"`
}
