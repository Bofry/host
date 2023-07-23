package app_test

import (
	"fmt"
	"testing"

	"github.com/Bofry/host/app"
)

func TestStandardProtocolRegistry(t *testing.T) {
	registry := app.NewProtocolMessageRegistry()

	message1 := app.Message{
		Format: app.TEXT_MESSAGE,
		Body:   []byte("ping"),
	}
	registry.Add(message1, MockStandardProtocol(0))
	message2 := app.Message{
		Format: app.BINARY_MESSAGE,
		Body:   []byte("ping"),
	}
	registry.Add(message2, MockStandardProtocol(1))

	fmt.Printf("%+v\n", registry)

	protocol := registry.Get(app.Message{
		Format: app.BINARY_MESSAGE,
		Body:   []byte("ping"),
	})

	fmt.Printf("%+v\n", protocol)

	registry.Visit(func(m app.Message, h app.StandardProtocol) {
		fmt.Printf("%+v, %+v\n", m, h)
	})
}

var _ app.StandardProtocol = MockStandardProtocol(0)

type MockStandardProtocol int

// ConfigureProtocol implements app.StandardMessageHandler.
func (MockStandardProtocol) ConfigureProtocol(registry *app.StandardProtocolRegistry) {
	panic("unimplemented")
}

// ReplyMessage implements app.StandardMessageHandler.
func (MockStandardProtocol) ReplyMessage(format app.MessageFormat, sender app.MessageSender) {
	panic("unimplemented")
}
