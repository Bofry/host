package app

var (
	_ EventClient    = new(MultiEventClient)
	_ EventForwarder = new(MultiEventClient)
)

type MultiEventClient map[string]EventClient

// Close implements EventBroker.
func (hub MultiEventClient) Close() error {
	for _, c := range hub {
		if c != nil {
			_ = c.Close()
		}
	}
	return nil
}

// Stop implements EventBroker.
func (hub MultiEventClient) Stop() {
	for _, c := range hub {
		if c != nil {
			c.Stop()
		}
	}
}

// Forward implements EventBroker.
func (hub MultiEventClient) Forward(channel string, payload []byte) error {
	if hub != nil {
		c, _ := hub[channel]
		if c != nil {
			return c.Forward(channel, payload)
		}

		c, _ = hub[InvalidChannel]
		if c != nil {
			return c.Forward(channel, payload)
		}
	}
	return Nop
}

// Start implements EventBroker.
func (hub MultiEventClient) Start(pipe *EventPipe) {
	for _, c := range hub {
		if c != nil {
			c.Start(pipe)
		}
	}
}
