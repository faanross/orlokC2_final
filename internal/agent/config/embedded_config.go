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
	AgentUUID:  "5499dab6-6042-4192-a158-be7e9d749f8b",
}
