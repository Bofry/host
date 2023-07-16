package app

var (
	_ EventSource    = NoopEventSource{}
	_ EventForwarder = NoopEventSource{}
)

type NoopEventSource struct{}

// Stop implements EventSource.
func (NoopEventSource) Stop() error {
	return nil
}

// Forward implements EventSource.
func (NoopEventSource) Forward(channel string, payload []byte) error {
	return nil
}

// Start implements EventSource.
func (NoopEventSource) Start(chan *Event, chan error) {}
