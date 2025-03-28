// File: internal/agent/config/embedded_config.go
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
	TargetHost: "localhost",
	TargetPort: 7777,
	Protocol:   "h1c",
	Sleep:      5,
	Jitter:     20,
	AgentUUID:  "4a3f2943-0dff-4c4e-9e9d-c197d5d6d658",
}
