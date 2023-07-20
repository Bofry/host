package app_test

import (
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

	FooEvent app.EventHandler `channel:"foo_topic"`
	BarEvent app.EventHandler `channel:"bar_topic"`

	App *MockApp
	app.ModuleOptionCollection
}{
	ModuleOptionCollection: app.ModuleOptions(
		app.WithDefaultMessageHandler(func(ctx *app.Context, message *app.Message) {
			ctx.Send(message.Format, append([]byte("[default]"), message.Body...))
		}),
		app.WithMessageCodeResolver(func(format app.MessageFormat, payload []byte) string {
			if len(payload) > 4 && payload[3] == '$' {
				return string(payload[:3])
			}
			return ""
		}),
	),
}

type MockApp struct {
	ServiceProvider *MockServiceProvider
	Config          *MockConfig

	Env string
}

func (app *MockApp) Init() {
	app.Env = app.Config.Env
}

func (app *MockApp) Foo(ctx *app.Context, message *app.Message) {
	data := message.Body[4:]
	prefix := fmt.Sprintf("[Foo:%s]", app.Env)
	ctx.Send(message.Format, append([]byte(prefix), data...))
}

func (app *MockApp) Bar(ctx *app.Context, message *app.Message) {
	data := message.Body[4:]
	prefix := fmt.Sprintf("[Bar:%s]", app.Env)
	ctx.Send(message.Format, append([]byte(prefix), data...))
}

func (app *MockApp) FooEvent(ctx *app.Context, event *app.Event) error { return nil }
func (app *MockApp) BarEvent(ctx *app.Context, event *app.Event) error { return nil }
