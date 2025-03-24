package types

type ProtocolType int

const (
	H1C ProtocolType = iota + 1
	H1TLS
	H2C
	H2TLS
	H3
)

func ProtocolStringFromType(t ProtocolType) string {
	switch t {
	case H1C:
		return "HTTP/1.1 Clear"
	case H1TLS:
		return "HTTP/1.1 TLS"
	case H2C:
		return "HTTP/2.0 Clear"
	case H2TLS:
		return "HTTP/2.0 TLS"
	default:
		return "Unknown Protocol Type"
	}
}
