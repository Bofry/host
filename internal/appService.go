package internal

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/Bofry/config"
)

type AppService struct {
	appModule   *AppModule
	hostService *HostService
	hostModule  HostModule

	componentService *ComponentService

	middlewares                  []Middleware
	configureConfigurationAction ConfigureConfigurationAction

	logger *log.Logger

	host1 Host
}

func (s *AppService) Init(service InjectionService) {
	// register dependency injection types
	s.registerConstructors(service)

	s.appModule.appStaterConfigurator().ConfigureLogger(s.logger)
	s.appModule.app().OnInit()

	// pass logger to HostService
	s.hostService.ConfigureLogger(s.logger.Flags(), s.logger.Writer())

	// trigger Init()
	s.hostService.Init(s.appModule)
}

func (m *AppService) InitComplete() {
	// trigger InitComplete()
	m.hostService.InitComplete(m.appModule)
	m.appModule.app().OnInitComplete()
}

func (s *AppService) LoadConfiguration() {
	s.initConfig()

	if s.configureConfigurationAction != nil {
		rvConfig := s.appModule.Field(APP_CONFIG_FIELD)
		service := config.NewConfigurationService(rvConfig.Interface())
		s.configureConfigurationAction(service)
	}

	s.initApp()
	s.initHost()
	s.initServiceProvider()
}

func (s *AppService) LoadComponent() {
	var (
		service = s.componentService
		app     = s.appModule
		rvApp   = app.rv
	)
	if rvApp.IsValid() {
		count := rvApp.NumField()
		for i := 0; i < count; i++ {
			rvField := rvApp.Field(i)
			if !rvField.CanInterface() {
				continue
			}

			if !rvField.IsZero() {
				v := rvField.Interface()
				switch v.(type) {
				case Runable:
					service.RegisterComponent(v.(Runable))
					s.logger.Printf("LOAD Component %T", v)
				}
			}
		}
	}
}

func (s *AppService) LoadMiddleware() {
	var (
		app         = s.appModule
		middlewares = s.middlewares
	)
	for _, v := range middlewares {
		s.logger.Printf("load middleware %T", v)
		v.Init(app)
	}
}

func (s *AppService) Start(ctx context.Context) {
	s.componentService.Start()
	s.hostService.Start(ctx)
	s.appModule.app().OnStart(ctx)
}

func (s *AppService) Stop(ctx context.Context) error {
	s.appModule.app().OnStop(ctx)
	s.componentService.Stop()
	return s.hostService.Stop(ctx)
}

func (s *AppService) registerConstructors(service InjectionService) error {
	app := s.appModule
	var (
		configFieldGetter          = AppModuleField(app.Field(APP_CONFIG_FIELD)).MakeGetter()
		serviceProviderFieldGetter = AppModuleField(app.Field(APP_SERVICE_PROVIDER_FIELD)).MakeGetter()
		hostFieldGetter            = AppModuleField(app.Field(APP_HOST_FIELD)).MakeGetter()
	)

	service.registerConstructors(
		configFieldGetter,
		hostFieldGetter,
		serviceProviderFieldGetter,
	)
	return nil
}

func (s *AppService) initApp() {
	app := s.appModule
	var (
		rvApp = app.pv
	)
	s.logger.Printf("LOAD App %s", rvApp.Type())

	s.appModule.app().Init()
	s.appModule.appTracingConfigurator().ConfigureTracerProvider()
}

func (s *AppService) initConfig() {
	app := s.appModule
	var (
		rvConfig = app.Field(APP_CONFIG_FIELD)
	)
	s.logger.Printf("LOAD Configuration %s", rvConfig.Type())

	// get the instance Init() method
	fn := rvConfig.MethodByName(APP_COMPONENT_INIT_METHOD)
	if fn.IsValid() {
		if fn.Kind() != reflect.Func {
			panic(fmt.Errorf("cannot find func %s() within type %s",
				APP_COMPONENT_INIT_METHOD,
				rvConfig.Type().String()))
		}
		if fn.Type().NumIn() != 0 || fn.Type().NumOut() != 0 {
			panic(fmt.Errorf("method type should be func %s.%s()",
				rvConfig.Type().String(),
				APP_COMPONENT_INIT_METHOD))
		}

		fn.Call([]reflect.Value(nil))
	}
}

func (s *AppService) initHost() {
	app := s.appModule
	var (
		rvConfig = app.Field(APP_CONFIG_FIELD)
		rvHost   = app.Field(APP_HOST_FIELD)
	)
	s.logger.Printf("LOAD Host %s", rvHost.Type())

	// get the instance Init() method
	fn := rvHost.MethodByName(APP_COMPONENT_INIT_METHOD)
	if fn.IsValid() {
		if fn.Kind() != reflect.Func {
			panic(fmt.Errorf("invalid func %s.%s(conf %s) within type %s",
				rvHost.Type().String(),
				APP_COMPONENT_INIT_METHOD,
				rvConfig.Type().String(),
				rvHost.Type().Name()))
		}
		if fn.Type().NumIn() != 1 || fn.Type().NumOut() != 0 ||
			(fn.Type().In(0) != rvConfig.Type()) {
			panic(fmt.Errorf("method type should be func %s.%s(conf %s)",
				rvHost.Type().String(),
				APP_COMPONENT_INIT_METHOD,
				rvConfig.Type().String()))
		}

		fn.Call([]reflect.Value{rvConfig})
	}
}

func (s *AppService) initServiceProvider() {
	app := s.appModule
	var (
		rvConfig          = app.Field(APP_CONFIG_FIELD)
		rvServiceProvider = app.Field(APP_SERVICE_PROVIDER_FIELD)
	)
	s.logger.Printf("LOAD ServiceProvider %s", rvServiceProvider.Type())

	// get the instance Init() method
	fn := rvServiceProvider.MethodByName(APP_COMPONENT_INIT_METHOD)
	if fn.IsValid() {
		if fn.Kind() != reflect.Func || fn.Type().NumOut() != 0 {
			panic(fmt.Errorf("invalid func %s(...) within type %s",
				APP_COMPONENT_INIT_METHOD,
				rvServiceProvider.Type().String()))
		}

		var args []reflect.Value
		if fn.Type().NumIn() > 0 {
			count := fn.Type().NumIn()
			for i := 0; i < count; i++ {
				paramType := fn.Type().In(i)
				switch paramType {
				case rvConfig.Type():
					args = append(args, rvConfig)
				case app.pv.Type():
					args = append(args, app.pv)
				default:
					panic(fmt.Errorf("unsupported type '%s' in %s.%s(...)",
						paramType.String(),
						rvServiceProvider.Type().String(),
						APP_COMPONENT_INIT_METHOD))
				}
			}
		}
		fn.Call(args)
	}
}

func (m *AppService) getHost() Host {
	if m.hostService == nil {
		var (
			rvHost          = m.appModule.Field(APP_HOST_FIELD)
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
		m.host1 = host
	}
	return m.hostService
}
