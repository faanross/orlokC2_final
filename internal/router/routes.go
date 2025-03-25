package router

import "github.com/go-chi/chi/v5"

func SetupRoutes(r chi.Router) {
	// Apply the middleware to all routes
	r.Use(AgentIdentifierMiddleware)

	r.Get("/", RootHandler)

	// Results endpoint - with GOB decoding middleware
	r.With(GobDecodingMiddleware).Post("/results", ResultsHandler)
}
