package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"orlokC2_final/internal/router"
)

const serverAddr = ":7777"

func main() {

	r := chi.NewRouter()

	router.SetupRoutes(r)

	log.Printf("Starting HTTP server on %v", serverAddr)

	err := http.ListenAndServe(serverAddr, r)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
