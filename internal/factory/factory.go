package factory

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"math/rand"
	"net/http"
	"orlokC2_final/internal/router"
	"time"
)

type ListenerFactory struct{}

func NewListenerFactory() *ListenerFactory {
	return &ListenerFactory{}
}

func (lf *ListenerFactory) CreateListener(serverAddr string) *Listener {
	id := fmt.Sprintf("listener_%06d", rand.Intn(1000000))

	r := chi.NewRouter()

	router.SetupRoutes(r)

	serverInstance := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	fmt.Printf("|CREATE|-> Listener %s configured on %s\n", id, serverAddr)

	return &Listener{
		ID:     id,
		Addr:   serverAddr,
		Router: r,
		Server: serverInstance,
	}
}

func (l *Listener) Start() error {
	fmt.Printf("|START|-> Listener %s serving on %s\n", l.ID, l.Addr)
	return l.Server.ListenAndServe()
}

// Stop gracefully shuts down the HTTP server
func (l *Listener) Stop() error {
	if l.Server == nil {
		return fmt.Errorf("server not started")
	}

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	err := l.Server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down listener %s: %v", l.ID, err)
	}

	fmt.Printf("|STOP| Listener %s on %s shut down successfully\n", l.ID, l.Addr)
	return nil
}
