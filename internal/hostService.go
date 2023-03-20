package internal

import (
	"context"
	"fmt"
	"io"
	"log"
)

type HostService struct {
	hostModule HostModule

	logger *log.Logger

	host Host
}

func (m *HostService) ConfigureLogger(logflags int, w io.Writer) {
	m.hostModule.ConfigureLogger(logflags, w)
}

func (m *HostService) Init(app *AppModule) {
	m.hostModule.Init(m.host, app)
}

func (m *HostService) InitComplete(app *AppModule) {
	m.hostModule.InitComplete(m.host, app)
}

func (m *HostService) Start(ctx context.Context) {
	m.host.Start(ctx)
}

func (m *HostService) Stop(ctx context.Context) error {
	return m.host.Stop(ctx)
}

func (m *HostService) registerHost(app *AppModule) Host {
	if m.host == nil {
		var (
			rvHost          = app.Field(APP_HOST_FIELD)
			rvHostInterface = AppModuleField(rvHost).As(m.hostModule.DescribeHostType()).Value()
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
