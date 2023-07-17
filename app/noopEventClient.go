package app

var (
	_ EventClient    = NoopEventClient{}
	_ EventForwarder = NoopEventClient{}
)

type NoopEventClient struct{}

// Close implements EventSource.
func (NoopEventClient) Close() error {
	return nil
}

// Forward implements EventSource.
func (NoopEventClient) Forward(channel string, payload []byte) error {
	return nil
}

// Start implements EventSource.
func (NoopEventClient) Start(*EventPipe) {}

// Stop implements EventSource.
func (NoopEventClient) Stop() {}
