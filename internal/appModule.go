package internal

import (
	"context"
	"fmt"
	"log"
	"reflect"
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

func (module *AppModule) app() App {
	if v, ok := module.target.(App); ok {
		return v
	}
	return nil
}

func (module *AppModule) ConfigureLogger(logger *log.Logger) {
	if configurator, ok := module.target.(AppStaterConfigurator); ok {
		configurator.ConfigureLogger(logger)
	}
}

func (module *AppModule) init() {
	if app := module.app(); app != nil {
		app.Init()
	}
}

func (module *AppModule) onInit() {
	if app := module.app(); app != nil {
		app.OnInit()
	}
}

func (module *AppModule) onInitComplete() {
	if app := module.app(); app != nil {
		app.OnInitComplete()
	}
}

func (module *AppModule) onStart(ctx context.Context) {
	if app := module.app(); app != nil {
		app.OnStart(ctx)
	}
}

func (module *AppModule) onStop(ctx context.Context) {
	if app := module.app(); app != nil {
		app.OnStop(ctx)
	}
}
