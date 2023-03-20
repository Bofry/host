package host

import "github.com/Bofry/host/internal"

const (
	APP_HOST_FIELD             = internal.APP_HOST_FIELD
	APP_CONFIG_FIELD           = internal.APP_CONFIG_FIELD
	APP_SERVICE_PROVIDER_FIELD = internal.APP_SERVICE_PROVIDER_FIELD
	APP_COMPONENT_INIT_METHOD  = internal.APP_COMPONENT_INIT_METHOD
)

// interface
type (
	Host       = internal.Host
	HostModule = internal.HostModule
	Middleware = internal.Middleware
	Runner     = internal.Runner
	Component  = internal.Runner
	Runable    = internal.Runable
)

// func
type (
	ConfigureConfigurationAction = internal.ConfigureConfigurationAction
)

// struct
type (
	App                    = internal.App
	AppStaterConfigurator  = internal.AppStaterConfigurator
	AppTracingConfigurator = internal.AppTracingConfigurator
	AppModule              = internal.AppModule
	Starter                = internal.Starter
)
