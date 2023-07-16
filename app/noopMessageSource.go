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

// Receive implements MessageSource.
func (NoopMessageSrouce) Receive(chan *Message) {}

// Send implements MessageSource.
func (NoopMessageSrouce) Send(format MessageFormat, payload []byte) error {
	return nil
}
