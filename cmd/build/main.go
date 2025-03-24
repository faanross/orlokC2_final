package main

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// BuildConfig holds all the build configuration options
type BuildConfig struct {
	Binary struct {
		Name      string `yaml:"name"`
		OutputDir string `yaml:"output_dir"`
		OS        string `yaml:"os"`
		Arch      string `yaml:"arch"`
	} `yaml:"binary"`

	Agent struct {
		TargetHost string `yaml:"target_host"`
		TargetPort int    `yaml:"target_port"`
		Protocol   string `yaml:"protocol"`
		Sleep      int    `yaml:"sleep"`
		Jitter     int    `yaml:"jitter"`
	} `yaml:"agent"`
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	config, err := loadConfig("./cmd/build/config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Generate UUID for this build
	agentUUID := uuid.New().String()
	fmt.Printf("Generated UUID: %s\n", agentUUID)

	err = writeEmbeddedConfig(config, agentUUID)
	if err != nil {
		fmt.Printf("Error writing embedded config: %v\n", err)
		os.Exit(1)
	}

	// Determine binary name
	binaryName := config.Binary.Name
	if binaryName == "" {
		binaryName = fmt.Sprintf("agent_%06d", rand.Intn(1000000))
	}

	// Add extension for Windows
	if config.Binary.OS == "windows" {
		binaryName += ".exe"
	}

	// Create output directory if it doesn't exist
	err = os.MkdirAll(config.Binary.OutputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	outputPath := filepath.Join(config.Binary.OutputDir, binaryName)

	// Build the binary
	fmt.Printf("Building agent binary: %s\n", outputPath)
	cmd := exec.Command("go", "build", "-o", outputPath, "./cmd/agent")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", config.Binary.OS),
		fmt.Sprintf("GOARCH=%s", config.Binary.Arch),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Build failed: %v\n%s\n", err, output)
		os.Exit(1)
	}

	fmt.Printf("Successfully built agent: %s\n", outputPath)
}

// loadConfig loads the build configuration from a YAML file
func loadConfig(path string) (*BuildConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config BuildConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// writeEmbeddedConfig writes config values to embedded_config.go
func writeEmbeddedConfig(config *BuildConfig, agentUUID string) error {
	// Template for the embedded config file
	template := `// File: internal/agent/config/embedded_config.go
package config

// EmbeddedValues contains the configuration values embedded during build
// These values come directly from config.yaml
var EmbeddedValues = struct {
	TargetHost string
	TargetPort int
	Protocol   string
	Sleep      int
	Jitter     int
	AgentUUID  string
}{
	TargetHost: "%s",
	TargetPort: %d,
	Protocol:   "%s",
	Sleep:      %d,
	Jitter:     %d,
	AgentUUID:  "%s",
}
`
	// Format template with actual values
	content := fmt.Sprintf(
		template,
		config.Agent.TargetHost,
		config.Agent.TargetPort,
		config.Agent.Protocol,
		config.Agent.Sleep,
		config.Agent.Jitter,
		agentUUID,
	)

	// Write to file
	filePath := "./internal/agent/config/embedded_config.go"
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write embedded config: %v", err)
	}

	fmt.Println("Configuration values written to embedded_config.go")
	return nil
}
