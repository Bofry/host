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

// Start implements EventSource.
func (NoopEventSource) Start(chan *Event, chan error) {}

// Stop implements EventSource.
func (NoopEventSource) Stop() {}
