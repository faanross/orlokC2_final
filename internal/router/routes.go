package router

import (
	"github.com/go-chi/chi/v5"
	"orlokC2_final/internal/middleware"
)

func SetupRoutes(r chi.Router) {
	// Apply the middleware to all routes
	r.Use(middleware.AgentIdentifierMiddleware)

	r.Get("/command", CommandHandler)
	r.Post("/result", ResultHandler)
}
