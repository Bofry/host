package internal

import (
	"context"
	"io"
	"log"
	"reflect"

	"github.com/Bofry/config"
)

const (
	APP_HOST_FIELD             string = "Host"
	APP_CONFIG_FIELD           string = "Config"
	APP_SERVICE_PROVIDER_FIELD string = "ServiceProvider"
	APP_COMPONENT_INIT_METHOD  string = "Init"

	LOGGER_PREFIX string = "[host] "
)

var (
	typeOfApp      = reflect.TypeOf(App(nil))
	typeOfHost     = reflect.TypeOf((*Host)(nil)).Elem()
	stdHostService = &StdHostService{}
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
	}

	Host interface {
		Start(ctx context.Context)
		Stop(ctx context.Context) error
	}

	HostService interface {
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

	ConfigureConfigurationAction func(service *config.ConfigurationService)
)
