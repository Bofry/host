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

func isErrorHandler(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfErrorHandler)
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

func asErrorHandler(rv reflect.Value) ErrorHandler {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfErrorHandler).Interface().(ErrorHandler); ok {
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

func sliceToMap[E comparable, V any](slice []E, setter func(key E) V) map[E]V {
	var m map[E]V = make(map[E]V, len(slice))

	for _, elem := range slice {
		m[elem] = setter(elem)
	}
	return m
}
