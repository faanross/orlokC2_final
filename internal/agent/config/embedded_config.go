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
	Jitter:     50,
	AgentUUID:  "759f2e2a-276b-4a5a-8bac-70ae4e54c896",
}
