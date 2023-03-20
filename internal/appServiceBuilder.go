package internal

import "log"

type AppServiceBuilder struct {
	appService  *AppService
	hostService *HostService
	logger      *log.Logger
}

func NewAppServiceBuilder(logger *log.Logger) *AppServiceBuilder {
	appService := &AppService{
		logger: logger,
	}

	hostService := &HostService{
		hostModule: stdHostModuleSingleton,
		logger:     logger,
	}

	return &AppServiceBuilder{
		appService:  appService,
		hostService: hostService,
		logger:      logger,
	}
}

func (builder *AppServiceBuilder) ConfigureConfigurationAction(action ConfigureConfigurationAction) *AppServiceBuilder {
	builder.appService.configureConfigurationAction = action
	return builder
}

func (builder *AppServiceBuilder) AppModule(app *AppModule) *AppServiceBuilder {
	builder.appService.appModule = app
	return builder
}

func (builder *AppServiceBuilder) HostModule(host HostModule) *AppServiceBuilder {
	builder.hostService.hostModule = host
	return builder
}

func (builder *AppServiceBuilder) Middlewares(middlewares []Middleware) *AppServiceBuilder {
	builder.appService.middlewares = middlewares
	return builder
}

func (builder *AppServiceBuilder) BuildHostService() *HostService {
	s := builder.hostService
	if builder.appService.appModule == nil {
		panic("missing AppModule")
	}
	s.registerHost(builder.appService.appModule)
	return s
}

func (builder *AppServiceBuilder) Build() *AppService {
	s := builder.appService
	if s.appModule == nil {
		panic("missing AppModule")
	}
	if s.hostService == nil {
		s.hostService = builder.BuildHostService()
	}
	s.componentService = NewComponentService(builder.logger)
	return s
}
