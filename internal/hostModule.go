package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/Bofry/config"
)

type HostModule struct {
	appService *AppService

	hostService      HostService
	componentService *ComponentService

	middlewares                  []Middleware
	configureConfigurationAction ConfigureConfigurationAction

	logger *log.Logger

	host Host
}

func (m *HostModule) Init(service InjectionService) {
	// register dependency injection types
	m.appService.RegisterConstructors(service)

	m.appService.App().ConfigureLogger(m.logger)
	m.appService.App().OnInit()

	// pass logger to HostService
	m.hostService.ConfigureLogger(m.logger.Flags(), m.logger.Writer())

	// trigger Init()
	m.hostService.Init(m.getHost(), m.appService.App())
}

func (m *HostModule) LoadConfiguration() {
	m.appService.InitConfig()

	if m.configureConfigurationAction != nil {
		rvConfig := m.appService.App().Field(APP_CONFIG_FIELD)
		service := config.NewConfigurationService(rvConfig.Interface())
		m.configureConfigurationAction(service)
	}

	m.appService.InitApp()
	m.appService.InitHost()
	m.appService.InitServiceProvider()
}

func (m *HostModule) LoadComponent() {
	m.appService.RegisterComponents(m.componentService)
}

func (m *HostModule) LoadMiddleware() {
	app := m.appService.App()
	for _, v := range m.middlewares {
		m.logger.Printf("load middleware %T", v)
		v.Init(app)
	}
}

func (m *HostModule) InitComplete() {
	// trigger InitComplete()
	m.hostService.InitComplete(m.getHost(), m.appService.App())
	m.appService.App().OnInitComplete()
}

func (m *HostModule) Start(ctx context.Context) {
	var (
		host = m.getHost()
	)
	m.componentService.Start()
	host.Start(ctx)
	m.appService.App().OnStart(ctx)
}

func (m *HostModule) Stop(ctx context.Context) error {
	var (
		host = m.getHost()
	)
	m.appService.App().OnStop(ctx)
	m.componentService.Stop()
	return host.Stop(ctx)
}

func (m *HostModule) getHost() Host {
	if m.host == nil {
		var (
			rvHost          = m.appService.App().Field(APP_HOST_FIELD)
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
