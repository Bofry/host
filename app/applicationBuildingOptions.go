package app

import (
	"io"

	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

const (
	APP_BUILDING_OPT_SESSION_STATE_MANAGER = "WithSessionStateManager"
	APP_BUILDING_OPT_PROTOCOL_RESOLVER     = "WithProtocolResolver"
	APP_BUILDING_OPT_PROTOCOL_EMITTER      = "WithProtocolEmitter"

	APP_BUILDING_OPT_LOGGER_OUTPUT       = "WithLoggerOutput"
	APP_BUILDING_OPT_LOGGER_FLAGS        = "WithLoggerFlags"
	APP_BUILDING_OPT_TRACER_PROVIDER     = "WithTracerProvider"
	APP_BUILDING_OPT_TEXT_MAP_PROPAGATOR = "WithTextMapPropagator"

	APP_BUILDING_OPT_EVENT_CLIENT   = "WithEventClient"
	APP_BUILDING_OPT_MESSAGE_ROUTER = "WithMessageRouter"
	APP_BUILDING_OPT_EVENT_ROUTER   = "WithEventRouter"

	APP_BUILDING_OPT_INVALID_MESSAGE_HANDLER = "WithInvalidMessageHandler"
	APP_BUILDING_OPT_INVALID_EVENT_HANDLER   = "WithInvalidEventHandler"
	APP_BUILDING_OPT_DEFAULT_MESSAGE_HANDLER = "WithDefaultMessageHandler"
	APP_BUILDING_OPT_DEFAULT_EVENT_HANDLER   = "WithDefaultEventHandler"
	APP_BUILDING_OPT_ERROR_HANDLER           = "WithErrorHandler"
)

var (
	_ ApplicationBuildingOption = GenericApplicationBuildingOption{}
)

type GenericApplicationBuildingOption struct {
	_apply    func(*Application) error
	_typename string
}

// apply implements ApplicationBuildingOption.
func (opt GenericApplicationBuildingOption) apply(app *Application) error {
	return opt._apply(app)
}

// typeName implements ApplicationBuildingOption.
func (opt GenericApplicationBuildingOption) typeName() string {
	return opt._typename
}

// ----------------------------------------------

func WithSessionStateManager(manager SessionStateManager) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.sessionStateManager = manager
			return nil
		},
		_typename: APP_BUILDING_OPT_SESSION_STATE_MANAGER,
	}
}

func WithProtocolResolver(resolver ProtocolResolver) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.configureProtocolResolver(resolver)
			return nil
		},
		_typename: APP_BUILDING_OPT_PROTOCOL_RESOLVER,
	}
}

func WithProtocolEmitter(emitter ProtocolEmitter) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.configureProtocolEmitter(emitter)
			return nil
		},
		_typename: APP_BUILDING_OPT_PROTOCOL_EMITTER,
	}
}

func WithInvalidMessageHandler(handler MessageHandler) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.configureInvalidMessageHandler(handler)
			return nil
		},
		_typename: APP_BUILDING_OPT_INVALID_MESSAGE_HANDLER,
	}
}

func WithInvalidEventHandler(handler EventHandler) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.configureInvalidEventHandler(handler)
			return nil
		},
		_typename: APP_BUILDING_OPT_INVALID_EVENT_HANDLER,
	}
}

func WithLoggerOutput(w io.Writer) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			if w == nil {
				w = io.Discard
			}
			ap.logger.SetOutput(w)
			return nil
		},
		_typename: APP_BUILDING_OPT_LOGGER_OUTPUT,
	}
}

func WithLoggerFlags(flags int) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.logger.SetFlags(flags)
			return nil
		},
		_typename: APP_BUILDING_OPT_LOGGER_FLAGS,
	}
}

func WithTracerProvider(tp *trace.SeverityTracerProvider) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.tracerProvider = tp
			return nil
		},
		_typename: APP_BUILDING_OPT_TRACER_PROVIDER,
	}
}

func WithTextMapPropagator(p propagation.TextMapPropagator) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.textMapPropagator = p
			return nil
		},
		_typename: APP_BUILDING_OPT_TEXT_MAP_PROPAGATOR,
	}
}

func WithEventClient(source EventClient) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			if source == nil {
				source = NoopEventClient{}
			}
			ap.eventClient = source
			return nil
		},
		_typename: APP_BUILDING_OPT_EVENT_CLIENT,
	}
}

func WithMessageRouter(router MessageRouter) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.configureMessageRouter(router)
			return nil
		},
		_typename: APP_BUILDING_OPT_MESSAGE_ROUTER,
	}
}

func WithEventRouter(router EventRouter) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.configureEventRouter(router)
			return nil
		},
		_typename: APP_BUILDING_OPT_EVENT_ROUTER,
	}
}

func WithDefaultMessageHandler(handler MessageHandler) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.configureDefaultMessageHandler(handler)
			return nil
		},
		_typename: APP_BUILDING_OPT_DEFAULT_MESSAGE_HANDLER,
	}
}

func WithDefaultEventHandler(handler EventHandler) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.configureDefaultEventHandler(handler)
			return nil
		},
		_typename: APP_BUILDING_OPT_DEFAULT_EVENT_HANDLER,
	}
}

func WithErrorHandler(handler ErrorHandler) ApplicationBuildingOption {
	return GenericApplicationBuildingOption{
		_apply: func(ap *Application) error {
			ap.errorHandler = handler
			return nil
		},
		_typename: APP_BUILDING_OPT_ERROR_HANDLER,
	}
}
