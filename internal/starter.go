package internal

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

var _ InjectionService = new(Starter)

type Starter struct {
	app *fx.App

	constructors []interface{}
	functions    []interface{}

	hostModuleBuilder *HostModuleBuilder
}

func NewStarter(app interface{}) *Starter {
	var (
		appContext        = NewAppContext(app)
		appService        = NewAppService(appContext)
		hostModuleBuilder = NewHostModuleBuilder()
	)

	hostModuleBuilder.AppService(appService)
	hostModuleBuilder.HostService(stdHostService)

	return &Starter{
		hostModuleBuilder: hostModuleBuilder,
	}
}

func (s *Starter) Middlewares(middlewares ...Middleware) *Starter {
	s.hostModuleBuilder.Middlewares(middlewares)
	return s
}

func (s *Starter) ConfigureConfiguration(action ConfigureConfigurationAction) *Starter {
	s.hostModuleBuilder.ConfigureConfiguration(action)
	return s
}

func (s *Starter) Configure(action ConfigureAction) *Starter {
	s.hostModuleBuilder.Configure(action)
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
		// build and initialize HostModule
		module := s.hostModuleBuilder.Build()
		{
			module.Init(s)
			module.LoadConfiguration()
			module.LoadCompoent()
			module.LoadMiddleware()
			module.Configure()
			module.InitComplete()
		}

		// register service hook
		hook := s.makeServiceHook(module)
		s.registerFunctions(hook)

		// build fx.App
		s.app = fx.New(
			fx.Provide(s.constructors...),
			fx.Invoke(s.functions...),
		)
	}
}

func (s *Starter) makeServiceHook(module *HostModule) interface{} {
	return func(lc fx.Lifecycle) {
		lc.Append(
			fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						module.Start(ctx)
						logger.Println("Started")
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Println("Shutdown")
					return module.Stop(ctx)
				},
			},
		)
	}
}
