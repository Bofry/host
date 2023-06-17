package internal

import (
	"reflect"
	"unsafe"
)

type ReflectHelper reflect.Value

func (f ReflectHelper) MakeGetter() interface{} {
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

func (f ReflectHelper) As(typ reflect.Type) ReflectHelper {
	var (
		rv         = reflect.Value(f)
		rvInstance reflect.Value
	)
	if rv.Type().ConvertibleTo(typ) {
		rvInstance = rv.Convert(typ)
	} else {
		rvInstance = reflect.NewAt(typ, unsafe.Pointer(rv.Pointer()))
	}
	return ReflectHelper(rvInstance)
}

func (f ReflectHelper) Value() reflect.Value {
	return reflect.Value(f)
}
