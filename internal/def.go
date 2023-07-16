package internal

import (
	"context"
	"io"
	"log"
	"reflect"

	"github.com/Bofry/config"
	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

const (
	APP_HOST_FIELD             string = "Host"
	APP_CONFIG_FIELD           string = "Config"
	APP_SERVICE_PROVIDER_FIELD string = "ServiceProvider"
	APP_COMPONENT_INIT_METHOD  string = "Init"

	LOGGER_PREFIX string = "[host] "
)

var (
	typeOfApp              = reflect.TypeOf(App(nil))
	typeOfHost             = reflect.TypeOf((*Host)(nil)).Elem()
	nopHostModuleSingleton = NopHostModule{}
)

type (
	App interface {
		Init()
		OnInit()
		OnInitComplete()
		OnStart(ctx context.Context)
		OnStop(ctx context.Context)
	}

	AppStaterConfigurator interface {
		ConfigureLogger(logger *log.Logger)
		Logger() *log.Logger
	}

	AppTracingConfigurator interface {
		ConfigureTextMapPropagator()
		ConfigureTracerProvider()
		TextMapPropagator() propagation.TextMapPropagator
		TracerProvider() *trace.SeverityTracerProvider
	}

	Host interface {
		Start(ctx context.Context)
		Stop(ctx context.Context) error
	}

	HostOnErrorEventHandler interface {
		OnError(err error) (disposed bool)
	}

	HostModule interface {
		Init(host Host, app *AppModule)
		InitComplete(host Host, app *AppModule)
		DescribeHostType() reflect.Type
		ConfigureLogger(logflags int, w io.Writer)
	}

	InjectionService interface {
		registerConstructors(constructors ...interface{})
		registerFunctions(functions ...interface{})
		build()
	}

	Middleware interface {
		Init(app *AppModule)
	}

	Runner interface {
		Start()
		Stop()
	}

	Runable interface {
		Runner() Runner
	}

	ConfigurationLoader func(service *config.ConfigurationService)
)
