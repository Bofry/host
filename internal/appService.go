package internal

import (
	"fmt"
	"log"
	"reflect"
)

type AppService struct {
	ctx    *AppContext
	logger *log.Logger
}

func NewAppService(appCtx *AppContext, logger *log.Logger) *AppService {
	return &AppService{
		ctx:    appCtx,
		logger: logger,
	}
}

func (s *AppService) RegisterConstructors(service InjectionService) error {
	ctx := s.ctx
	var (
		configFieldGetter          = AppContextField(ctx.Field(APP_CONFIG_FIELD)).MakeGetter()
		serviceProviderFieldGetter = AppContextField(ctx.Field(APP_SERVICE_PROVIDER_FIELD)).MakeGetter()
		hostFieldGetter            = AppContextField(ctx.Field(APP_HOST_FIELD)).MakeGetter()
	)

	service.registerConstructors(
		configFieldGetter,
		hostFieldGetter,
		serviceProviderFieldGetter,
	)
	return nil
}

func (s *AppService) RegisterComponents(service *ComponentService) {
	ctx := s.ctx
	var (
		rvApp = ctx.rv
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

func (s *AppService) InitApp() {
	ctx := s.ctx
	var (
		rvConfig = ctx.Field(APP_CONFIG_FIELD)
		rvApp    = ctx.pv
	)
	s.logger.Printf("LOAD App %s", rvApp.Type())

	// get the instance Init() method
	fn := rvApp.MethodByName(APP_COMPONENT_INIT_METHOD)
	if fn.IsValid() {
		if fn.Kind() != reflect.Func {
			panic(fmt.Errorf("invalid func %s() within type %s",
				APP_COMPONENT_INIT_METHOD,
				rvApp.Type().Name()))
		}
		if fn.Type().NumIn() == 0 && fn.Type().NumOut() == 0 {
			fn.Call([]reflect.Value{})
		} else if fn.Type().NumIn() == 1 && fn.Type().NumOut() == 0 &&
			(fn.Type().In(0) == rvConfig.Type()) {
			fn.Call([]reflect.Value{rvConfig})
		} else {
			panic(fmt.Errorf("method type should be func %[1]s.%[2]s() or func %[1]s.%[2]s(conf %[3]s)",
				rvApp.Type().String(),
				APP_COMPONENT_INIT_METHOD,
				rvConfig.Type().String()))
		}
	}
}

func (s *AppService) InitConfig() {
	ctx := s.ctx
	var (
		rvConfig = ctx.Field(APP_CONFIG_FIELD)
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

func (s *AppService) InitHost() {
	ctx := s.ctx
	var (
		rvConfig = ctx.Field(APP_CONFIG_FIELD)
		rvHost   = ctx.Field(APP_HOST_FIELD)
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

func (s *AppService) InitServiceProvider() {
	ctx := s.ctx
	var (
		rvConfig          = ctx.Field(APP_CONFIG_FIELD)
		rvServiceProvider = ctx.Field(APP_SERVICE_PROVIDER_FIELD)
	)
	s.logger.Printf("LOAD ServiceProvider %s", rvServiceProvider.Type())

	// get the instance Init() method
	fn := rvServiceProvider.MethodByName(APP_COMPONENT_INIT_METHOD)
	if fn.IsValid() {
		if fn.Kind() != reflect.Func {
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
				case ctx.pv.Type():
					args = append(args, ctx.pv)
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

func (s *AppService) AppContext() *AppContext {
	return s.ctx
}
