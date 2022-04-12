package internal

import (
	"reflect"
	"unsafe"
)

type AppContextField reflect.Value

func (f AppContextField) MakeGetter() interface{} {
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

func (f AppContextField) As(typ reflect.Type) AppContextField {
	var (
		rv         = reflect.Value(f)
		rvInstance reflect.Value
	)
	if rv.Type().ConvertibleTo(typ) {
		rvInstance = rv.Convert(typeOfHost)
	} else {
		rvInstance = reflect.NewAt(typ, unsafe.Pointer(rv.Pointer()))
	}
	return AppContextField(rvInstance)
}

func (f AppContextField) Value() reflect.Value {
	return reflect.Value(f)
}
