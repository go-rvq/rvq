package reflect_utils

import "reflect"

func Filter(s any, f func(item any) bool) any {
	v := reflect.ValueOf(s)
	out := reflect.MakeSlice(v.Type(), 0, 0)
	for i := 0; i < v.Len(); i++ {
		if f(v.Index(i).Interface()) {
			out = reflect.Append(out, v.Index(i))
		}
	}
	return out.Interface()
}
