package factory

import (
	"orlokC2_final/internal/interfaces"
	"orlokC2_final/internal/protocols/h1c"
	"orlokC2_final/internal/types"
)

// AbstractFactory decides which protocol-specific factory to use
type AbstractFactory struct {
	factories map[types.ProtocolType]interfaces.ListenerFactory
}

// NewAbstractFactory creates a new AbstractFactory with all registered protocol factories
func NewAbstractFactory() *AbstractFactory {
	return &AbstractFactory{
		factories: map[types.ProtocolType]interfaces.ListenerFactory{
			types.H1C: &h1c.ListenerFactory{},
			// Other protocols will be added here as they are implemented
			// H1TLS: &h1tls.Factory{},
			// H2C: &h2c.Factory{},
			// H2TLS: &h2tls.Factory{},
			// H3: &h3.Factory{},
		},
	}
}
