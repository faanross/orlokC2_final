package modules

import (
	"os"
	"os/exec"
	"strings"
)

// GetWhoami executes the whoami command
func GetWhoami() (string, error) {
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetHostname returns the hostname of the system
func GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

// GetPwd returns the current working directory
func GetPwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}
