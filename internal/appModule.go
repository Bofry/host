package internal

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/Bofry/trace"
)

var (
	_ App                    = AppProxy{}
	_ AppStaterConfigurator  = AppStaterConfiguratorProxy{}
	_ AppTracingConfigurator = AppTracingConfiguratorProxy{}
)

type AppModule struct {
	target interface{}
	rv     reflect.Value
	pv     reflect.Value
}

func NewAppModule(target interface{}) *AppModule {
	var rv reflect.Value
	switch target.(type) {
	case reflect.Value:
		rv = target.(reflect.Value)
	default:
		rv = reflect.ValueOf(target)
	}

	if !rv.IsValid() {
		panic("host: specified argument 'target' is invalid")
	}

	rv = reflect.Indirect(rv)

	return &AppModule{
		target: target,
		rv:     rv,
		pv:     rv.Addr(),
	}
}

func (module *AppModule) Field(name string) reflect.Value {
	var rv = module.rv
	rvfield := rv.FieldByName(name)
	if rvfield.Kind() != reflect.Ptr {
		panic(fmt.Errorf("specified App field '%s' should be of type *%s", name, rvfield.Type().String()))
	}
	if rvfield.IsNil() {
		rvfield.Set(reflect.New(rvfield.Type().Elem()))
	}
	return rvfield
}

func (module *AppModule) Host() reflect.Value {
	return module.Field(APP_HOST_FIELD)
}

func (module *AppModule) Config() reflect.Value {
	return module.Field(APP_CONFIG_FIELD)
}

func (module *AppModule) ServiceProvider() reflect.Value {
	return module.Field(APP_SERVICE_PROVIDER_FIELD)
}

func (module *AppModule) TracerProvider() *trace.SeverityTracerProvider {
	return module.appTracingConfigurator().TracerProvider()
}

func (module *AppModule) app() App {
	return AppProxy{module: module}
}

func (module *AppModule) appStaterConfigurator() AppStaterConfigurator {
	return AppStaterConfiguratorProxy{module: module}
}

func (module *AppModule) appTracingConfigurator() AppTracingConfigurator {
	return AppTracingConfiguratorProxy{module: module}
}

type AppProxy struct {
	module *AppModule
}

func (proxy AppProxy) app() App {
	if v, ok := proxy.module.target.(App); ok {
		return v
	}
	return nil
}

// Init implements App
func (proxy AppProxy) Init() {
	if app := proxy.app(); app != nil {
		app.Init()
	}
}

// OnInit implements App
func (proxy AppProxy) OnInit() {
	if app := proxy.app(); app != nil {
		app.OnInit()
	}
}

// OnInitComplete implements App
func (proxy AppProxy) OnInitComplete() {
	if app := proxy.app(); app != nil {
		app.OnInitComplete()
	}
}

// OnStart implements App
func (proxy AppProxy) OnStart(ctx context.Context) {
	if app := proxy.app(); app != nil {
		app.OnStart(ctx)
	}
}

// OnStop implements App
func (proxy AppProxy) OnStop(ctx context.Context) {
	if app := proxy.app(); app != nil {
		app.OnStop(ctx)
	}
}

type AppStaterConfiguratorProxy struct {
	module *AppModule
}

func (proxy AppStaterConfiguratorProxy) app() AppStaterConfigurator {
	if v, ok := proxy.module.target.(AppStaterConfigurator); ok {
		return v
	}
	return nil
}

// ConfigureLogger implements AppStaterConfigurator
func (proxy AppStaterConfiguratorProxy) ConfigureLogger(logger *log.Logger) {
	if app := proxy.app(); app != nil {
		app.ConfigureLogger(logger)
	}
}

// ConfigureLogger implements AppStaterConfigurator
func (proxy AppStaterConfiguratorProxy) Logger() *log.Logger {
	if app := proxy.app(); app != nil {
		return app.Logger()
	}
	return log.Default()
}

type AppTracingConfiguratorProxy struct {
	module *AppModule
}

func (proxy AppTracingConfiguratorProxy) app() AppTracingConfigurator {
	if v, ok := proxy.module.target.(AppTracingConfigurator); ok {
		return v
	}
	return nil
}

// ConfigureTracerProvider implements tracing.AppTracingConfigurator
func (proxy AppTracingConfiguratorProxy) ConfigureTracerProvider() {
	if app := proxy.app(); app != nil {
		app.ConfigureTracerProvider()
	}
}

func (proxy AppTracingConfiguratorProxy) TracerProvider() *trace.SeverityTracerProvider {
	if app := proxy.app(); app != nil {
		return app.TracerProvider()
	}
	return nil
}
