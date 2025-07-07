package zeroer

import "reflect"

type Zeroer interface {
	IsZero() bool
}

var ZeroerType = reflect.TypeOf((*Zeroer)(nil)).Elem()

func IsNilValueOk(value reflect.Value) (is, ok bool) {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.UnsafePointer, reflect.Interface,
		reflect.Slice:
		is = value.IsNil()
		ok = true
	case reflect.Pointer:
		if !value.IsNil() {
			return IsNilValueOk(value.Elem())
		}
		is = true
		ok = true
	}
	return
}

func IsNilValue(value reflect.Value) (is bool) {
	is, _ = IsNilValueOk(value)
	return
}

func IsNil(value any) bool {
	if value == nil {
		return true
	}
	return IsNilValue(reflect.ValueOf(value))
}

func IsZero(value any) (ok bool) {
	if value == nil {
		return true
	}

	if value, ok_ := value.(reflect.Value); ok_ {
	do:
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
			if value.Type().Implements(ZeroerType) {
				return value.Interface().(Zeroer).IsZero()
			}
			value = value.Elem()
			goto do
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
	return IsZero(reflect.ValueOf(value))
}
