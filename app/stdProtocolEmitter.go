package app

var (
	_ ProtocolEmitter = StdProtocolEmitter
)

func StdProtocolEmitter(format MessageFormat, protocol string, payload []byte) []byte {
	return payload
}
