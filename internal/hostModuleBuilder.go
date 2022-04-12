package internal

type HostModuleBuilder struct {
	module *HostModule
}

func NewHostModuleBuilder() *HostModuleBuilder {
	module := &HostModule{}
	return &HostModuleBuilder{
		module: module,
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

func (builder *HostModuleBuilder) Configure(action ConfigureAction) *HostModuleBuilder {
	builder.module.configureAction = action
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
	m.componentService = NewComponentService()
	return m
}
