package h1c

import (
	"fmt"
	"net/http"
	"time"

	"orlokC2_final/internal/types"
)

// H1CProtocol implements the HTTP/1.1 Clear protocol
type H1CProtocol struct {
	// Configuration
	targetHost string
	targetPort int
	timeout    time.Duration
	reqTimeout time.Duration

	// HTTP client for making requests
	client *http.Client

	// State tracking
	isConnected bool
}

// Initialize sets up the protocol with the given configuration
func (p *H1CProtocol) Initialize(targetHost string, targetPort int, timeout int) error {
	p.targetHost = targetHost
	p.targetPort = targetPort
	p.timeout = time.Duration(timeout) * time.Second
	p.reqTimeout = time.Duration(timeout) * time.Second

	// Create HTTP client with appropriate timeouts
	p.client = &http.Client{
		Timeout: p.reqTimeout,
	}

	p.isConnected = false

	return nil
}

// Connect establishes a connection to the C2 server
func (p *H1CProtocol) Connect() error {
	// Connect checks if the C2 server is reachable
	// HTTP/1.1 supports persistent connections by default,
	// and Go's http.Client will automatically handle connection pooling
	url := fmt.Sprintf("http://%s:%d/", p.targetHost, p.targetPort)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	// If we get here, we successfully connected
	p.isConnected = true
	return nil
}

// Disconnect closes the connection to the C2 server
func (p *H1CProtocol) Disconnect() error {
	// For HTTP/1.1, there's no persistent connection to close,
	// but we update our state
	p.isConnected = false
	return nil
}

// SendRequest sends a request to the C2 server and returns the response
func (p *H1CProtocol) SendRequest(endpoint string) (*http.Response, error) {
	// Check if we're connected
	if !p.isConnected {
		return nil, fmt.Errorf("not connected to server")
	}

	// Create the full URL
	url := fmt.Sprintf("http://%s:%d%s", p.targetHost, p.targetPort, endpoint)

	// Create and send the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add basic headers
	req.Header.Set("User-Agent", "Mozilla/5.0")

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		p.isConnected = false
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}

// GetType returns the type of protocol being used
func (p *H1CProtocol) GetType() types.ProtocolType {
	return types.H1C
}
