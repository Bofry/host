package app_test

import (
	"fmt"

	"github.com/Bofry/host/app"
)

var (
	_ app.MessageClient = new(MockMessageClient)
)

type MockMessageClient struct {
	In  chan []byte
	Out chan []byte

	onCloseDelegate []func(app.MessageClient)

	stopped bool

	*app.MessageClientInfo
}

// Close implements app.MessageClient.
func (c *MockMessageClient) Close() error {
	restricted := app.NewRestrictedMessageClient(c)

	for _, onClose := range c.onCloseDelegate {
		onClose(restricted)
	}
	return nil
}

// RegisterCloseHandler implements app.MessageClient.
func (c *MockMessageClient) RegisterCloseHandler(listener func(app.MessageClient)) {
	c.onCloseDelegate = append(c.onCloseDelegate, listener)
}

// Send implements app.MessageClient.
func (c *MockMessageClient) Send(format app.MessageFormat, payload []byte) error {
	switch format {
	case app.TEXT_MESSAGE, app.BINARY_MESSAGE:
		c.Out <- payload
	default:
		return fmt.Errorf("unsupported MessageFormat '%v'", format)
	}
	return nil
}

// Start implements app.MessageClient.
func (c *MockMessageClient) Start(pipe *app.MessagePipe) {
	if c.stopped {
		panic("the client stopped")
	}

	go func() {
		for !c.stopped {
			select {
			case v, ok := <-c.In:
				if ok {
					pipe.Forward(c, &app.Message{
						Format: app.BINARY_MESSAGE,
						Body:   v,
					})
				}
			}
		}
	}()
}

// Stop implements app.MessageClient.
func (c *MockMessageClient) Stop() {
	if !c.stopped {
		c.stopped = true
		close(c.In)
	}
}
