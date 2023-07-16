package app

var (
	_ MessageSource = NoopMessageSrouce{}
	_ MessageSender = NoopMessageSrouce{}
)

type NoopMessageSrouce struct{}

// Close implements MessageSource.
func (NoopMessageSrouce) Close() error {
	return nil
}

// Stop implements MessageSource.
func (NoopMessageSrouce) Stop() {}

// Start implements MessageSource.
func (NoopMessageSrouce) Start(chan *Message, chan error) {}

// Send implements MessageSource.
func (NoopMessageSrouce) Send(format MessageFormat, payload []byte) error {
	return nil
}
