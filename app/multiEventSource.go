package app

var (
	_ EventSource    = new(MultiEventSource)
	_ EventForwarder = new(MultiEventSource)
)

type MultiEventSource map[string]EventSource

// Close implements EventBroker.
func (hub MultiEventSource) Close() error {
	for _, s := range hub {
		_ = s.Close()
	}
	return nil
}

// Stop implements EventBroker.
func (hub MultiEventSource) Stop() {
	for _, s := range hub {
		s.Stop()
	}
}

// Forward implements EventBroker.
func (hub MultiEventSource) Forward(channel string, payload []byte) error {
	if hub != nil {
		s, ok := hub[channel]
		if ok {
			return s.Forward(channel, payload)
		}

		if s, ok := hub[InvalidChannel]; ok {
			return s.Forward(channel, payload)
		}
	}
	return Nop
}

// Start implements EventBroker.
func (hub MultiEventSource) Start(observer chan *Event, err chan error) {
	for _, s := range hub {
		s.Start(observer, err)
	}
}
