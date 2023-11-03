package app_test

import (
	"bytes"
	"fmt"

	"github.com/Bofry/host/app"
)

type (
	MockServiceProvider struct {
		ID string
	}

	MockConfig struct {
		Env string
	}
)

var MockModule = struct {
	Foo app.MessageHandler `protocol:"foo"`
	Bar app.MessageHandler `protocol:"bar"`

	FooEvent app.EventHandler `channel:"foo_topic"   @ExpandEnv:"off"`
	BarEvent app.EventHandler `channel:"bar_topic"   @ExpandEnv:"off"`

	App *MockApp
	app.ModuleOptionCollection
}{
	ModuleOptionCollection: app.ModuleOptions(
		app.WithProtocolResolver(func(format app.MessageFormat, payload []byte) (string, []byte) {
			if len(payload) > 4 && payload[3] == '$' {
				return string(payload[:3]), payload[4:]
			}
			return "", payload
		}),
		app.WithProtocolEmitter(func(format app.MessageFormat, protocol string, body []byte) []byte {
			return bytes.Join(
				[][]byte{
					[]byte(protocol + "$"),
					body,
				}, nil)
		}),
	),
}

type MockApp struct {
	ServiceProvider *MockServiceProvider
	Config          *MockConfig

	// EventClient app.EventClient

	Env string
}

func (ap *MockApp) Init() {
	ap.Env = ap.Config.Env

	// ap.EventClient = app.MultiEventClient{
	// 	"foo_topic":        app.NoopEventClient{},
	// 	"bar_topic":        app.NoopEventClient{},
	// 	app.InvalidChannel: app.NoopEventClient{},
	// }
}

func (app *MockApp) Foo(ctx *app.Context, message *app.Message) {
	data := message.Body
	prefix := fmt.Sprintf("[Foo:%s]", app.Env)
	ctx.Send(message.Format, message.Protocol, append([]byte(prefix), data...))
}

func (app *MockApp) Bar(ctx *app.Context, message *app.Message) {
	data := message.Body
	prefix := fmt.Sprintf("[Bar:%s]", app.Env)
	ctx.Send(message.Format, message.Protocol, append([]byte(prefix), data...))
}

func (app *MockApp) FooEvent(ctx *app.Context, event *app.Event) error { return nil }
func (app *MockApp) BarEvent(ctx *app.Context, event *app.Event) error { return nil }

func (app *MockApp) DefaultMessageHandler(ctx *app.Context, message *app.Message) {
	prefix := fmt.Sprintf("[default:%s]", app.Env)
	err := ctx.Send(message.Format, message.Protocol, append([]byte(prefix), message.Body...))
	if err != nil {
		ctx.Logger().Printf("err on send 1st:: %v", err)
	}
	ctx.CloseSend()
	err = ctx.Send(message.Format, message.Protocol, append([]byte(prefix+"??11"), message.Body...))
	if err != nil {
		ctx.Logger().Printf("err on send 2nd:: %v", err)
	}
	panic("... error occurred ...")
}

func (app *MockApp) DefaultEventHandler(ctx *app.Context, event *app.Event) error {
	return nil
}
