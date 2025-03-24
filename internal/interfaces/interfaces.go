package interfaces

import "orlokC2_final/internal/listener"

// Listener interface defines methods common to all protocols
type Listener interface {
	Start() error
	Stop() error
	GetProtocol() string
	GetAddr() string
	GetID() string
}

// ListenerFactory interface defines methods unique to all protocols
type ListenerFactory interface {
	CreateListener(addr string) *listener.ConcreteListener
}
