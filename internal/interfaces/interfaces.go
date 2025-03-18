package interfaces

import "orlokC2_final/internal/listener"

// Listener interface defines methods that all listener types must implement
type Listener interface {
	Start() error
	Stop() error
	GetProtocol() string
	GetAddr() string
	GetID() string
}

// ListenerFactory interface defines methods for creating listeners
type ListenerFactory interface {
	CreateListener(addr string) *listener.ConcreteListener
}
