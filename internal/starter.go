package internal

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var _ InjectionService = new(Starter)

type Starter struct {
	app *fx.App

	logger *log.Logger

	constructors []interface{}
	functions    []interface{}

	appServiceBuilder *AppServiceBuilder
}

func NewStarter(app interface{}) *Starter {
	var (
		hostLogger        = createStarterLogger()
		appModule         = NewAppModule(app)
		appServiceBuilder = NewAppServiceBuilder(hostLogger)
	)

	appServiceBuilder.AppModule(appModule)

	return &Starter{
		logger:            hostLogger,
		appServiceBuilder: appServiceBuilder,
	}
}

func (s *Starter) Middlewares(middlewares ...Middleware) *Starter {
	s.appServiceBuilder.Middlewares(middlewares)
	return s
}

func (s *Starter) ConfigureConfiguration(action ConfigureConfigurationAction) *Starter {
	s.appServiceBuilder.ConfigureConfigurationAction(action)
	return s
}

func (s *Starter) Start(ctx context.Context) error {
	s.build()
	if s.app == nil {
		panic(fmt.Errorf("Starter does not be initialized"))
	}
	return s.app.Start(ctx)
}

func (s *Starter) Stop(ctx context.Context) error {
	if s.app == nil {
		panic(fmt.Errorf("Starter did not call Start() yet"))
	}
	return s.app.Stop(ctx)
}

func (s *Starter) Run() {
	s.build()
	if s.app == nil {
		panic(fmt.Errorf("Starter does not be initialized"))
	}
	s.app.Run()
}

func (s *Starter) registerConstructors(constructors ...interface{}) {
	s.constructors = append(s.constructors, constructors...)
}

func (s *Starter) registerFunctions(functions ...interface{}) {
	s.functions = append(s.functions, functions...)
}

func (s *Starter) build() {
	if s.app == nil {
		// build and initialize AppService
		service := s.appServiceBuilder.Build()
		{
			service.Init(s)
			service.LoadConfiguration()
			service.LoadComponent()
			service.LoadMiddleware()
			service.InitComplete()
		}

		// register service hook
		hook := s.makeServiceHook(service)
		s.registerFunctions(hook)

		// build fx.App
		s.app = fx.New(
			fx.Provide(s.constructors...),
			fx.Invoke(s.functions...),
			fx.WithLogger(
				func() fxevent.Logger {
					return &StarterLogger{
						Flags: s.logger.Flags(),
						Logger: &fxevent.ConsoleLogger{
							W: s.logger.Writer(),
						},
					}
				},
			),
		)
	}
}

func (s *Starter) makeServiceHook(module *AppService) interface{} {
	return func(lc fx.Lifecycle) {
		lc.Append(
			fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						s.logger.Println("STARTING")
						module.Start(ctx)
						s.logger.Println("RUNNING")
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					s.logger.Println("STOPPING")
					err := module.Stop(ctx)
					s.logger.Println("SHUTDOWN")
					return err
				},
			},
		)
	}
}
