package app

import (
	"fmt"
	"reflect"

	"github.com/Bofry/host"
	"github.com/Bofry/structproto/reflecting"
)

var _ ModuleBindingOption = ModuleBindingOptionFunc(nil)

type ModuleBindingOptionFunc func(reflect.Value) error

func (f ModuleBindingOptionFunc) apply(v reflect.Value) error {
	return f(v)
}

func BindServiceProvider(v interface{}) ModuleBindingOption {
	return ModuleBindingOptionFunc(func(rv reflect.Value) error {
		if v == nil {
			return nil
		}

		var (
			rvApp reflect.Value = indirectValue(rv)
			rvVal reflect.Value = indirectValue(reflect.ValueOf(v))
		)

		if !rvApp.IsValid() {
			return fmt.Errorf("specified target is invalid")
		}

		target := rvApp.FieldByName(host.APP_SERVICE_PROVIDER_FIELD)
		target = indirectValue(reflecting.AssignZero(target))

		if rvVal.Type().ConvertibleTo(target.Type()) {
			target.Set(rvVal.Convert(target.Type()))
		}
		return nil
	})
}

func BindConfig(v interface{}) ModuleBindingOption {
	return ModuleBindingOptionFunc(func(rv reflect.Value) error {
		if v == nil {
			return nil
		}

		var (
			rvApp reflect.Value = indirectValue(rv)
			rvVal reflect.Value = indirectValue(reflect.ValueOf(v))
		)

		if !rvApp.IsValid() {
			return fmt.Errorf("specified target is invalid")
		}

		target := rvApp.FieldByName(host.APP_CONFIG_FIELD)
		target = indirectValue(reflecting.AssignZero(target))

		if rvVal.Type().ConvertibleTo(target.Type()) {
			target.Set(rvVal.Convert(target.Type()))
		}
		return nil
	})
}
