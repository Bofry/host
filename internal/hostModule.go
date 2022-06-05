package internal

import (
	"context"
	"fmt"

	"github.com/Bofry/config"
)

type HostModule struct {
	appService *AppService

	hostService      HostService
	componentService *ComponentService

	middlewares                  []Middleware
	configureConfigurationAction ConfigureConfigurationAction
	configureAction              ConfigureAction

	// transient
	host Host
}

func (m *HostModule) Init(service InjectionService) {
	// register dependency injection types
	m.appService.RegisterConstructors(service)

	// trigger Init()
	m.hostService.Init(m.getHost(), m.appService.AppContext())
}

func (m *HostModule) LoadConfiguration() {
	m.appService.InitConfig()

	if m.configureConfigurationAction != nil {
		rvConfig := m.appService.AppContext().Field(APP_CONFIG_FIELD)
		service := config.NewConfigurationService(rvConfig.Interface())
		m.configureConfigurationAction(service)
	}

	m.appService.InitHost()
	m.appService.InitApp()
	m.appService.InitServiceProvider()
}

func (m *HostModule) LoadCompoent() {
	m.appService.RegisterComponents(m.componentService)
}

func (m *HostModule) LoadMiddleware() {
	appCtx := m.appService.AppContext()
	for _, v := range m.middlewares {
		v.Init(appCtx)
	}
}

func (m *HostModule) Configure() {
	if m.configureAction != nil {
		rv := m.appService.AppContext().Field(APP_CONFIG_FIELD)
		config := rv.Interface()
		m.configureAction(config)
	}
}

func (m *HostModule) InitComplete() {
	// trigger InitComplete()
	m.hostService.InitComplete(m.getHost(), m.appService.AppContext())
}

func (m *HostModule) Start(ctx context.Context) {
	var (
		host = m.getHost()
	)
	m.componentService.Start()
	host.Start(ctx)
}

func (m *HostModule) Stop(ctx context.Context) error {
	var (
		host = m.getHost()
	)
	// FIXME see if all components are graceful shutdown or don't
	m.componentService.Stop()
	return host.Stop(ctx)
}

func (m *HostModule) getHost() Host {
	if m.host == nil {
		var (
			rvHost          = m.appService.AppContext().Field(APP_HOST_FIELD)
			rvHostInterface = AppContextField(rvHost).As(m.hostService.GetHostType()).Value()
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
