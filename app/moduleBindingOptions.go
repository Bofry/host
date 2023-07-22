package app

import (
	"fmt"
	"reflect"

	"github.com/Bofry/host"
	"github.com/Bofry/structproto/reflecting"
)

var (
	_ ModuleBindingOption = ModuleBindingOptionFunc(nil)
)

type ModuleBindingOptionFunc func(reflect.Value, TargetValueRole) error

func (f ModuleBindingOptionFunc) apply(v reflect.Value, role TargetValueRole) error {
	return f(v, role)
}

func BindServiceProvider(v interface{}) ModuleBindingOption {
	return ModuleBindingOptionFunc(func(rv reflect.Value, role TargetValueRole) error {
		if role != APP {
			return nil
		}
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
		if rvApp.Kind() != reflect.Struct {
			return fmt.Errorf("specified target is not a struct")
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
	return ModuleBindingOptionFunc(func(rv reflect.Value, role TargetValueRole) error {
		if role != APP {
			return nil
		}
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
		if rvApp.Kind() != reflect.Struct {
			return fmt.Errorf("specified target is not a struct")
		}

		target := rvApp.FieldByName(host.APP_CONFIG_FIELD)
		target = indirectValue(reflecting.AssignZero(target))

		if rvVal.Type().ConvertibleTo(target.Type()) {
			target.Set(rvVal.Convert(target.Type()))
		}
		return nil
	})
}

func BindEventClient(v interface{}) ModuleBindingOption {
	eventClient, ok := v.(EventClient)
	if !ok {
		panic("specified value cannot convert to EventClient")
	}

	return ModuleBindingOptionFunc(func(rv reflect.Value, role TargetValueRole) error {
		if role != MODULE_OPTIONS {
			return nil
		}
		if v == nil {
			return nil
		}

		var (
			rvOpt reflect.Value = indirectValue(rv)
		)

		if !rvOpt.IsValid() {
			return fmt.Errorf("specified target is invalid")
		}
		if rvOpt.Kind() != reflect.Slice {
			return fmt.Errorf("specified target is not a slice")
		}

		rvOpt.Set(reflect.Append(rvOpt,
			reflect.ValueOf(WithEventClient(eventClient))))

		return nil
	})
}
