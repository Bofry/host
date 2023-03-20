package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/Bofry/config"
)

type HostService struct {
	appService *AppService

	hostService      HostModule
	componentService *ComponentService

	middlewares                  []Middleware
	configureConfigurationAction ConfigureConfigurationAction

	logger *log.Logger

	host Host
}

func (m *HostService) Init(service InjectionService) {
	// register dependency injection types
	m.appService.RegisterConstructors(service)

	m.appService.AppModule().appStaterConfigurator().ConfigureLogger(m.logger)
	m.appService.AppModule().app().OnInit()

	// pass logger to HostService
	m.hostService.ConfigureLogger(m.logger.Flags(), m.logger.Writer())

	// trigger Init()
	m.hostService.Init(m.getHost(), m.appService.AppModule())
}

func (m *HostService) LoadConfiguration() {
	m.appService.InitConfig()

	if m.configureConfigurationAction != nil {
		rvConfig := m.appService.AppModule().Field(APP_CONFIG_FIELD)
		service := config.NewConfigurationService(rvConfig.Interface())
		m.configureConfigurationAction(service)
	}

	m.appService.InitApp()
	m.appService.InitHost()
	m.appService.InitServiceProvider()
}

func (m *HostService) LoadComponent() {
	m.appService.RegisterComponents(m.componentService)
}

func (m *HostService) LoadMiddleware() {
	app := m.appService.AppModule()
	for _, v := range m.middlewares {
		m.logger.Printf("load middleware %T", v)
		v.Init(app)
	}
}

func (m *HostService) InitComplete() {
	// trigger InitComplete()
	m.hostService.InitComplete(m.getHost(), m.appService.AppModule())
	m.appService.AppModule().app().OnInitComplete()
}

func (m *HostService) Start(ctx context.Context) {
	var (
		host = m.getHost()
	)
	m.componentService.Start()
	host.Start(ctx)
	m.appService.AppModule().app().OnStart(ctx)
}

func (m *HostService) Stop(ctx context.Context) error {
	var (
		host = m.getHost()
	)
	m.appService.AppModule().app().OnStop(ctx)
	m.componentService.Stop()
	return host.Stop(ctx)
}

func (m *HostService) getHost() Host {
	if m.host == nil {
		var (
			rvHost          = m.appService.AppModule().Field(APP_HOST_FIELD)
			rvHostInterface = AppModuleField(rvHost).As(m.hostService.DescribeHostType()).Value()
			host            Host
		)
		// check if rvHost can convert to Host interface
		host, ok := rvHostInterface.Interface().(Host)
		if !ok {
			panic(fmt.Errorf("specified field 'Host' type '%s' cannot convert to '%s' interface",
				rvHost.Type().String(),
				typeOfHost.String()))
		}
		m.host = host
	}
	return m.host
}
