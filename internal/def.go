package internal

import (
	"context"
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
	typeOfHost     = reflect.TypeOf((*Host)(nil)).Elem()
	stdHostService = &StdHostService{}
)

type (
	Host interface {
		Start(ctx context.Context)
		Stop(ctx context.Context) error
	}

	HostService interface {
		Init(host Host, app *AppContext)
		InitComplete(host Host, app *AppContext)
		GetHostType() reflect.Type
	}

	InjectionService interface {
		registerConstructors(constructors ...interface{})
		registerFunctions(functions ...interface{})
		build()
	}

	Middleware interface {
		Init(appCtx *AppContext)
	}

	Runner interface {
		Start()
		Stop()
	}

	Runable interface {
		Runner() Runner
	}

	OnInitAction                 func(appContext interface{})
	OnInitCompleteAction         func(appContext interface{})
	ConfigureConfigurationAction func(service *config.ConfigurationService)
)
