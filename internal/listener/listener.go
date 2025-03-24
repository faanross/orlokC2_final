package listener

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"orlokC2_final/internal/types"
	"time"
)

type ConcreteListener struct {
	ID       string
	Addr     string
	Router   *chi.Mux
	Server   *http.Server
	Protocol types.ProtocolType
	CertFile string
	Keyfile  string
}

func (l *ConcreteListener) Start() error {
	fmt.Printf("|START|-> Listener %s serving on %s using %s\n", l.ID, l.Addr, l.GetProtocol())

	if l.Protocol == types.H2TLS || l.Protocol == types.H1TLS {
		return l.Server.ListenAndServeTLS(l.CertFile, l.Keyfile)
	}

	return l.Server.ListenAndServe()
}

func (l *ConcreteListener) Stop() error {
	if l.Server == nil {
		return fmt.Errorf("server not started")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := l.Server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down listener %s: %v", l.ID, err)
	}

	fmt.Printf("|STOP|-> Listener %s on %s using %s shut down successfully\n", l.ID, l.Addr, l.GetProtocol())
	return nil
}

func (l *ConcreteListener) GetProtocol() string {
	switch l.Protocol {
	case types.H1C:
		return "HTTP/1.1"
	case types.H1TLS:
		return "HTTP/1.1 (TLS)"
	case types.H2C:
		return "HTTP/2"
	case types.H2TLS:
		return "HTTP/2 (TLS)"
	case types.H3:
		return "HTTP/3"
	default:
		return "Unknown"
	}
}

func (l *ConcreteListener) GetAddr() string {
	return l.Addr
}

func (l *ConcreteListener) GetID() string {
	return l.ID
}
