package uuid

import (
	"github.com/google/uuid"
)

// GenerateUUID creates a new random UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// This will be replaced during build
var AgentUUID = "57d4fb03-2828-465f-904e-3336bb532f65"
