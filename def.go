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
	Host        = internal.Host
	HostService = internal.HostService
	Middleware  = internal.Middleware
	Runner      = internal.Runner
	Component   = internal.Runner
	Runable     = internal.Runable
)

// func
type (
	ConfigureAction              = internal.ConfigureAction
	ConfigureConfigurationAction = internal.ConfigureConfigurationAction
)

// struct
type (
	AppContext = internal.AppContext
	Starter    = internal.Starter
)
