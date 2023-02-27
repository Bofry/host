package internal

import "log"

type HostModuleBuilder struct {
	module *HostModule
	logger *log.Logger
}

func NewHostModuleBuilder(logger *log.Logger) *HostModuleBuilder {
	module := &HostModule{
		logger: logger,
	}

	return &HostModuleBuilder{
		module: module,
		logger: logger,
	}
}

func (builder *HostModuleBuilder) AppService(service *AppService) *HostModuleBuilder {
	builder.module.appService = service
	return builder
}

func (builder *HostModuleBuilder) ConfigureConfiguration(action ConfigureConfigurationAction) *HostModuleBuilder {
	builder.module.configureConfigurationAction = action
	return builder
}

func (builder *HostModuleBuilder) HostService(service HostService) *HostModuleBuilder {
	builder.module.hostService = service
	return builder
}

func (builder *HostModuleBuilder) Middlewares(middlewares []Middleware) *HostModuleBuilder {
	builder.module.middlewares = middlewares
	return builder
}

func (builder *HostModuleBuilder) Build() *HostModule {
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
