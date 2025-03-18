package h1c

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"math/rand"
	"net/http"
	"orlokC2_final/internal/listener"
	"orlokC2_final/internal/router"
)

type ListenerFactory struct{}

func (lf *ListenerFactory) CreateListener(serverAddr string) *listener.ConcreteListener {
	id := fmt.Sprintf("listener_%06d", rand.Intn(1000000))

	r := chi.NewRouter()

	router.SetupRoutes(r)

	serverInstance := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	fmt.Printf("|CREATE|-> Listener %s configured on %s\n", id, serverAddr)

	return &listener.ConcreteListener{
		ID:     id,
		Addr:   serverAddr,
		Router: r,
		Server: serverInstance,
	}
}
