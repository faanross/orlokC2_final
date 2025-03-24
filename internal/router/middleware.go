package router

import (
	"context"
	"net/http"
)

// Key for storing agent UUID in request context
type contextKey string

const AgentUUIDKey contextKey = "agentUUID"

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
