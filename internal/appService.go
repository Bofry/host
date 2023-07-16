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

	componentService *ComponentService

	middlewares                  []Middleware
	configureConfigurationAction ConfigurationLoader

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
}

func (s *AppService) InitComplete() {
	s.hostService.InitComplete(s.appModule)
	s.appModule.app().OnInitComplete()
}

func (s *AppService) LoadConfiguration() {
	s.initConfig()

	if s.configureConfigurationAction != nil {
		rvConfig := s.appModule.Config()
		service := config.NewConfigurationService(rvConfig.Interface())
		s.configureConfigurationAction(service)
	}

	s.initApp()
	s.hostService.Init(s.appModule)
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
		configFieldGetter          = ReflectHelper(app.Config()).MakeGetter()
		serviceProviderFieldGetter = ReflectHelper(app.ServiceProvider()).MakeGetter()
		hostFieldGetter            = ReflectHelper(app.Host()).MakeGetter()
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
	s.appModule.appTracingConfigurator().ConfigureTextMapPropagator()
	s.appModule.appTracingConfigurator().ConfigureTracerProvider()
}

func (s *AppService) initConfig() {
	app := s.appModule
	var (
		rvConfig = app.Config()
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
		rvConfig = app.Config()
		rvHost   = app.Host()
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
		rvConfig          = app.Config()
		rvServiceProvider = app.ServiceProvider()
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
