package uuid

import (
	"github.com/google/uuid"
)

// GenerateUUID creates a new random UUID
func GenerateUUID() string {
	return uuid.New().String()
}
