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
	TargetHost: "PLACEHOLDER_HOST",
	TargetPort: 0,
	Protocol:   "PLACEHOLDER_PROTOCOL",
	Sleep:      0,
	Jitter:     0,
	AgentUUID:  "PLACEHOLDER_UUID",
}
