package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Execute runs the specified command and returns the output
func Execute(cmd string) (string, error) {
	// Trim any whitespace
	cmd = strings.TrimSpace(cmd)

	// Check which command to run
	switch cmd {
	case "pwd":
		return Pwd()
	case "whoami":
		return WhoAmI()
	case "hostname":
		return Hostname()
	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}
}

// Pwd returns the current working directory
func Pwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// WhoAmI returns the current user
func WhoAmI() (string, error) {
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Hostname returns the machine hostname
func Hostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}
