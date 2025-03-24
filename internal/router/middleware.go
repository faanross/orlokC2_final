package router

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"net/http"
	"orlokC2_final/internal/types"
)

// Key for storing agent UUID in request context
type contextKey string

const AgentUUIDKey contextKey = "agentUUID"

const ResultKey contextKey = "commandResult"

// AgentIdentifierMiddleware extracts the agent UUID from headers
func AgentIdentifierMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the UUID from the header
		agentUUID := r.Header.Get("X-Agent-ID")

		// If no UUID is provided, use "unknown" as a default
		if agentUUID == "" {
			agentUUID = "unknown"
		}

		// Add the UUID to the request context
		ctx := context.WithValue(r.Context(), AgentUUIDKey, agentUUID)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// init registers our CommandResult type with the gob package
func init() {
	// Register the CommandResult type with gob
	gob.Register(types.CommandResult{})
}

// GobDecodingMiddleware decodes GOB-encoded request bodies into CommandResult structures
func GobDecodingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only process the body if it's a POST request
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Decode the GOB data
		var result types.CommandResult
		decoder := gob.NewDecoder(bytes.NewReader(body))
		if err := decoder.Decode(&result); err != nil {
			http.Error(w, "Failed to decode GOB data", http.StatusBadRequest)
			return
		}

		// Add the decoded result to the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, ResultKey, result)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
