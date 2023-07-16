package app

var (
	_ EventSource    = NoopEventSource{}
	_ EventForwarder = NoopEventSource{}
)

type NoopEventSource struct{}

// Close implements EventSource.
func (NoopEventSource) Close() error {
	return nil
}

// Forward implements EventSource.
func (NoopEventSource) Forward(channel string, payload []byte) error {
	return nil
}

// Notify implements EventSource.
func (NoopEventSource) Notify(chan *Event, chan error) {}
