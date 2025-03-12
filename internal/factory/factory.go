package factory

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"math/rand"
	"net/http"
	"orlokC2_final/internal/router"
)

type ListenerFactory struct{}

func NewListenerFactory() *ListenerFactory {
	return &ListenerFactory{}
}

type Listener struct {
	ID     string
	Addr   string
	Router *chi.Mux
}

func (lf *ListenerFactory) CreateListener(serverAddr string) *Listener {
	id := fmt.Sprintf("listener_%06d", rand.Intn(1000000))

	r := chi.NewRouter()

	router.SetupRoutes(r)

	fmt.Printf("|CREATE|-> Listener %s configured on %s\n", id, serverAddr)

	return &Listener{
		ID:     id,
		Addr:   serverAddr,
		Router: r,
	}

}

func (l *Listener) Start() error {
	fmt.Printf("|START|-> Listener %s serving on %s\n", l.ID, l.Addr)
	return http.ListenAndServe(l.Addr, l.Router)
}
