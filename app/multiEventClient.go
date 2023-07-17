package app

var (
	_ EventClient    = new(MultiEventClient)
	_ EventForwarder = new(MultiEventClient)
)

type MultiEventClient map[string]EventClient

// Close implements EventBroker.
func (hub MultiEventClient) Close() error {
	for _, s := range hub {
		_ = s.Close()
	}
	return nil
}

// Stop implements EventBroker.
func (hub MultiEventClient) Stop() {
	for _, s := range hub {
		s.Stop()
	}
}

// Forward implements EventBroker.
func (hub MultiEventClient) Forward(channel string, payload []byte) error {
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
func (hub MultiEventClient) Start(pipe *EventPipe) {
	for _, s := range hub {
		s.Start(pipe)
	}
}
