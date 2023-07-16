package app

import (
	"fmt"

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
