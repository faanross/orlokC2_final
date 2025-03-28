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
	AgentUUID:  "6fa76a3a-c82e-4bd2-8405-ec9b26205a9d",
}
