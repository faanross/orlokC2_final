package types

type ProtocolType int

const (
	H1C ProtocolType = iota + 1
	H1TLS
	H2C
	H2TLS
	H3
)
