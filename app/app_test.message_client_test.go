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

	onClose func(app.MessageClient)

	stopped bool
}

// Close implements app.MessageClient.
func (c *MockMessageClient) Close() error {
	if c.onClose != nil {
		c.onClose(c)
	}
	return nil
}

// RegisterCloseHandler implements app.MessageClient.
func (c *MockMessageClient) RegisterCloseHandler(listener func(app.MessageClient)) {
	c.onClose = listener
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
