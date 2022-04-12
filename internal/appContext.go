package internal

import (
	"fmt"
	"reflect"
)

type AppContext struct {
	rv reflect.Value
	pv reflect.Value
}

func NewAppContext(target interface{}) *AppContext {
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

	return &AppContext{
		rv: rv,
		pv: rv.Addr(),
	}
}

func (ctx *AppContext) Field(name string) reflect.Value {
	var rv = ctx.rv
	rvfield := rv.FieldByName(name)
	if rvfield.Kind() != reflect.Ptr {
		panic(fmt.Errorf("specified appContext field '%s' should be of type *%s", name, rvfield.Type().String()))
	}
	if rvfield.IsNil() {
		rvfield.Set(reflect.New(rvfield.Type().Elem()))
	}
	return rvfield
}

func (ctx *AppContext) Host() reflect.Value {
	return ctx.Field(APP_HOST_FIELD)
}

func (ctx *AppContext) Config() reflect.Value {
	return ctx.Field(APP_CONFIG_FIELD)
}

func (ctx *AppContext) ServiceProvider() reflect.Value {
	return ctx.Field(APP_SERVICE_PROVIDER_FIELD)
}
