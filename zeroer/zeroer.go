package zeroer

import "reflect"

type Zeroer interface {
	IsZero() bool
}

func IsZero(value interface{}) bool {
	return IsZeroValue(reflect.ValueOf(value))
}

func IsNilValueH(value reflect.Value) (handled, ok bool) {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface,
		reflect.Slice:
		return true, value.IsNil()
	}
	return
}

func IsNilValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface,
		reflect.Slice:
		return value.IsNil()
	}
	return false
}

func IsZeroValue(value reflect.Value) (ok bool) {
	if !value.IsValid() {
		return true
	}

	if IsNilValue(value) {
		return
	}

	switch t := value.Interface().(type) {
	case Zeroer:
		return t.IsZero()
	}

	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return true
		}
		if value.Type().Implements(reflect.TypeOf((*Zeroer)(nil)).Elem()) {
			return value.Interface().(Zeroer).IsZero()
		}
		return IsZeroValue(value.Elem())
	case reflect.Slice:
		return value.Len() == 0
	case reflect.Struct:
		if z, ok := value.Interface().(Zeroer); ok {
			return z.IsZero()
		}
	case reflect.Func:
		return false
	}

	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
