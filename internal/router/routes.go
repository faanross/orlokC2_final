package router

import (
	"github.com/go-chi/chi/v5"
	"orlokC2_final/internal/middleware"
	"orlokC2_final/internal/websocket"
)

func SetupRoutes(r chi.Router) {
	// Apply the middleware to all routes
	r.Use(middleware.AgentIdentifierMiddleware)

	r.Get("/command", websocket.CommandEndpoint)
	r.Post("/result", websocket.ResultEndpoint)
}
