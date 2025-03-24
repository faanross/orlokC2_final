package h2tls

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"math/rand"
	"net/http"
	"orlokC2_final/internal/listener"
	"orlokC2_final/internal/router"
	"orlokC2_final/internal/types"
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
	fmt.Printf("|CREATE|-> Listener %s configured on %s using %s\n", id, serverAddr, types.ProtocolStringFromType(types.H2TLS))

	return &listener.ConcreteListener{
		ID:       id,
		Addr:     serverAddr,
		Router:   r,
		Server:   serverInstance,
		Protocol: types.H1C,
	}
}
