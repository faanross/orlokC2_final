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
	"strings"
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

	// Embed UUID in source code
	err = embedUUID(agentUUID)
	if err != nil {
		fmt.Printf("Error embedding UUID: %v\n", err)
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

// embedUUID embeds the UUID in the agent source code
func embedUUID(uuid string) error {
	// Path to the UUID file
	uuidFile := "./internal/agent/uuid/uuid.go"

	// Read the file
	content, err := ioutil.ReadFile(uuidFile)
	if err != nil {
		return fmt.Errorf("failed to read UUID file: %v", err)
	}

	// Replace the placeholder with the actual UUID
	newContent := strings.Replace(
		string(content),
		`var AgentUUID = "00000000-0000-0000-0000-000000000000"`,
		fmt.Sprintf(`var AgentUUID = "%s"`, uuid),
		1,
	)

	// Write back to the file
	err = ioutil.WriteFile(uuidFile, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write UUID file: %v", err)
	}

	return nil
}
