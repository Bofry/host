package app

import (
	"log"

	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

type AppBase struct {
	application *Application
}

// func (app *AppBase) DumpErrorStackTrace(err error) {
// 	buf := debug.Stack()
// 	app.Logger().Printf("%v\n%s\n", err, string(buf))
// }

func (app *AppBase) EnableStackTrace(enabled bool) {
	app.application.enableErrorStackTrace(enabled)
}

func (app *AppBase) Logger() *log.Logger {
	return app.application.logger
}

func (app *AppBase) TracerProvider() *trace.SeverityTracerProvider {
	return app.application.TracerProvider()
}

func (app *AppBase) TextMapPropagator() propagation.TextMapPropagator {
	return app.application.TextMapPropagator()
}
