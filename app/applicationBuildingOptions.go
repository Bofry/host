package app

import (
	"io"

	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

var _ ApplicationBuildingOption = ApplicationBuildingOptionFunc(nil)

type ApplicationBuildingOptionFunc func(*Application) error

func (f ApplicationBuildingOptionFunc) apply(ap *Application) error {
	return f(ap)
}

func WithMessageCodeResolver(resolver MessageCodeResolver) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		ap.configureMessageCodeResolver(resolver)
		return nil
	})
}

func WithInvalidMessageHandler(handler MessageHandler) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		ap.configureInvalidMessageHandler(handler)
		return nil
	})
}

func WithInvalidEventHandler(handler EventHandler) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		ap.configureInvalidEventHandler(handler)
		return nil
	})
}

func WithLoggerOuput(w io.Writer) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		if w == nil {
			w = io.Discard
		}
		ap.logger.SetOutput(w)
		return nil
	})
}

func WithLoggerFlags(flags int) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		ap.logger.SetFlags(flags)
		return nil
	})
}

func WithTracerProvider(tp *trace.SeverityTracerProvider) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		if tp == nil {
			tp = createNoopTracerProvider()
		}
		ap.tracerProvider = tp
		return nil
	})
}

func WithTextMapPropagator(p propagation.TextMapPropagator) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		if p == nil {
			p = createNoopTextMapPropagator()
		}
		ap.textMapPropagator = p
		return nil
	})
}

func WithMessageSource(source MessageSource) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		if source == nil {
			source = NoopMessageSrouce{}
		}
		ap.messageSource = source
		return nil
	})
}

func WithEventSource(source EventSource) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		if source == nil {
			source = NoopEventSource{}
		}
		ap.eventSource = source
		return nil
	})
}

func WithMessageRouter(router MessageRouter) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		ap.configureMessageRouter(router)
		return nil
	})
}

func WithEventRouter(router EventRouter) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		ap.configureEventRouter(router)
		return nil
	})
}

func WithMessageHandler(handler MessageHandler) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		ap.messageHandler = handler
		return nil
	})
}

func WithEventHandler(handler EventHandler) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		ap.eventHandler = handler
		return nil
	})
}

func WithErrorHandler(handler ErrorHandler) ApplicationBuildingOption {
	return ApplicationBuildingOptionFunc(func(ap *Application) error {
		ap.errorHandler = handler
		return nil
	})
}
