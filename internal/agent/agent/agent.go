package agent

import (
	"fmt"
	"io"
	"time"

	"orlokC2_final/internal/agent/config"
	"orlokC2_final/internal/agent/protocols"
)

// Agent manages the C2 communication
type Agent struct {
	// Configuration
	Config *config.AgentConfig

	// Protocol implementation
	Protocol protocols.Protocol

	// Control
	stopChan chan struct{}
	running  bool
}

// New creates a new Agent with the given configuration and protocol
func New(config *config.AgentConfig, protocol protocols.Protocol) *Agent {
	return &Agent{
		Config:   config,
		Protocol: protocol,
		stopChan: make(chan struct{}),
		running:  false,
	}
}

// Start begins the agent's operation in a goroutine
func (a *Agent) Start() error {
	if a.running {
		return fmt.Errorf("agent already running")
	}

	// Initialize protocol
	err := a.Protocol.Initialize(
		a.Config.TargetHost,
		a.Config.TargetPort,
		int(a.Config.ConnectionTimeout.Seconds()),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize protocol: %v", err)
	}

	// Start main loop in a goroutine
	a.running = true
	go a.runLoop()

	fmt.Printf("Agent started with UUID: %s\n", a.Config.AgentUUID)
	return nil
}

// Stop gracefully stops the agent's operation
func (a *Agent) Stop() error {
	if !a.running {
		return fmt.Errorf("agent not running")
	}

	// Signal the loop to stop
	close(a.stopChan)
	a.running = false

	// Disconnect from server
	if err := a.Protocol.Disconnect(); err != nil {
		return fmt.Errorf("error disconnecting: %v", err)
	}

	fmt.Println("Agent stopped")
	return nil
}

// runLoop is the main agent operation loop
func (a *Agent) runLoop() {
	for {
		select {
		case <-a.stopChan:
			// Stop signal received
			return
		default:
			// Try to connect if not connected
			err := a.Protocol.Connect()
			if err != nil {
				fmt.Printf("Connection error: %v\n", err)
				time.Sleep(a.Config.ReconnectDelay)
				continue
			}

			// Send request to root endpoint
			resp, err := a.Protocol.SendRequest("/")
			if err != nil {
				fmt.Printf("Request error: %v\n", err)
				time.Sleep(a.Config.ReconnectDelay)
				continue
			}

			// Process response
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				fmt.Printf("Response: %s\n", string(body))
			}

			// Sleep before next check-in
			sleepTime := a.Config.CalculateSleepWithJitter()
			time.Sleep(sleepTime)
		}
	}
}
