package internal

import "log"

type HostServiceBuilder struct {
	module *HostService
	logger *log.Logger
}

func NewHostModuleBuilder(logger *log.Logger) *HostServiceBuilder {
	module := &HostService{
		logger: logger,
	}

	return &HostServiceBuilder{
		module: module,
		logger: logger,
	}
}

func (builder *HostServiceBuilder) AppService(service *AppService) *HostServiceBuilder {
	builder.module.appService = service
	return builder
}

func (builder *HostServiceBuilder) ConfigureConfiguration(action ConfigureConfigurationAction) *HostServiceBuilder {
	builder.module.configureConfigurationAction = action
	return builder
}

func (builder *HostServiceBuilder) HostService(service HostModule) *HostServiceBuilder {
	builder.module.hostService = service
	return builder
}

func (builder *HostServiceBuilder) Middlewares(middlewares []Middleware) *HostServiceBuilder {
	builder.module.middlewares = middlewares
	return builder
}

func (builder *HostServiceBuilder) Build() *HostService {
	m := builder.module
	if m.appService == nil {
		panic("missing AppService")
	}
	if m.hostService == nil {
		panic("missing HostService")
	}
	m.componentService = NewComponentService(builder.logger)
	return m
}
