package listener

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"orlokC2_final/internal/factory"
)

type ConcreteListener struct {
	ID       string
	Addr     string
	Router   *chi.Mux
	Server   *http.Server
	Protocol factory.ProtocolType // ADD THIS
}
