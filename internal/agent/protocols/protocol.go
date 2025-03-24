package protocols

import (
	"net/http"
	"orlokC2_final/internal/types"
)

// Protocol defines the interface that all protocol implementations must satisfy
type Protocol interface {
	// Initialize sets up the protocol with the given configuration
	Initialize(targetHost string, targetPort int, timeout int) error

	// Connect establishes a connection to the C2 server
	Connect() error

	// Disconnect closes the connection to the C2 server
	Disconnect() error

	// SendRequest sends a request to the C2 server and returns the response
	SendRequest(endpoint string) (*http.Response, error)

	// GetType returns the type of protocol being used
	GetType() types.ProtocolType
}
