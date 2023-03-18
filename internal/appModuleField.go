package internal

import (
	"reflect"
	"unsafe"
)

type AppModuleField reflect.Value

func (f AppModuleField) MakeGetter() interface{} {
	rv := reflect.Value(f)
	if rv.IsValid() {
		impl := func(in []reflect.Value) []reflect.Value {
			return []reflect.Value{rv}
		}
		decl := reflect.FuncOf([]reflect.Type{}, []reflect.Type{rv.Type()}, false)

		return reflect.MakeFunc(decl, impl).Interface()
	}
	return nil
}

func (f AppModuleField) As(typ reflect.Type) AppModuleField {
	var (
		rv         = reflect.Value(f)
		rvInstance reflect.Value
	)
	if rv.Type().ConvertibleTo(typ) {
		rvInstance = rv.Convert(typeOfHost)
	} else {
		rvInstance = reflect.NewAt(typ, unsafe.Pointer(rv.Pointer()))
	}
	return AppModuleField(rvInstance)
}

func (f AppModuleField) Value() reflect.Value {
	return reflect.Value(f)
}
