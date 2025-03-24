package factory

import (
	"orlokC2_final/internal/interfaces"
	"orlokC2_final/internal/listener"
	"orlokC2_final/internal/protocols/h1c"
	"orlokC2_final/internal/protocols/h2tls"
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
			types.H1C:   &h1c.ListenerFactory{},
			types.H2TLS: &h2tls.ListenerFactory{},
			// Other protocols will be added here as they are implemented
		},
	}
}

// CreateListener creates a listener with the specified protocol type
func (af *AbstractFactory) CreateListener(protocol types.ProtocolType, addr string) *listener.ConcreteListener {
	factory := af.factories[protocol]
	return factory.CreateListener(addr)
}
