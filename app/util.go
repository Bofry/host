package app

import (
	"fmt"
	"reflect"

	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

func createNoopTracerProvider() *trace.SeverityTracerProvider {
	tp, err := trace.NoopProvider()
	if err != nil {
		panic(fmt.Sprintf("cannot create NoopProvider: %v", err))
	}
	return tp
}

func createNoopTextMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator()
}

func indirectValue(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	return rv
}

func isMessageHandler(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageHandler)
	}
	return false
}

func isEventHandler(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfEventHandler)
	}
	return false
}

func asMessageHandler(rv reflect.Value) MessageHandler {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageHandler).Interface().(MessageHandler); ok {
			return v
		}
	}
	return nil
}

func asEventHandler(rv reflect.Value) EventHandler {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfEventHandler).Interface().(EventHandler); ok {
			return v
		}
	}
	return nil
}

func asEventClient(rv reflect.Value) EventClient {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfEventClient).Interface().(EventClient); ok {
			return v
		}
	}
	return nil
}
