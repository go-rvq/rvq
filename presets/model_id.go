package presets

import (
	"fmt"
	"net/http"
	"reflect"
)

func (id ID) String() string {
	if id.Value == nil {
		return ""
	}
	s := fmt.Sprint(id.Value)
	if s == "0" {
		s = ""
	}
	return s
}

func (id ID) IsZero() bool {
	if id.Value != nil {
		if z, ok := id.Value.(interface {
			IsZero() bool
		}); ok {
			return z.IsZero()
		}
		if s := id.String(); len(s) > 0 && s != "0" {
			return false
		}
	}
	return true
}

func (id ID) Of(field string) IDOfField {
	return IDOfField{id, field}
}
func (f ID) SetTo(obj interface{}) {
	f.Of("ID").SetTo(obj)
}

func (f IDOfField) SetTo(obj interface{}) {
	field := reflect.ValueOf(obj).Elem().FieldByName(f.Field)
	if field.Kind() == reflect.Pointer {
		field.Set(reflect.New(field.Type().Elem()))
		field = field.Elem()
	}
	field.Set(reflect.ValueOf(f.ID.Value))
}

type IDSlice []ID

func (ids IDSlice) Len() int {
	return len(ids)
}

func (ids IDSlice) Last() ID {
	return ids[len(ids)-1]
}

func (ids IDSlice) LastValue() any {
	return ids.Last().Value
}

func (ids *IDSlice) Pop() ID {
	v := (*ids)[len((*ids))-1]
	*ids = (*ids)[:len(*ids)-1]
	return v
}

func (ids *IDSlice) PopValue() any {
	return ids.Pop().Value
}

func (ids IDSlice) Values() (values []any) {
	values = make([]any, len(ids))
	for i, id := range ids {
		values[i] = id.Value
	}
	return
}

func ParentsModelID(r *http.Request) IDSlice {
	if v := r.Context().Value(ParentsModelIDKey); v != nil {
		return v.(IDSlice)
	}
	return nil
}

func (mb *ModelBuilder) ParseParentsID(r *http.Request) (ids IDSlice, err error) {
	if parents := mb.Parents(); len(parents) > 0 {
		ids = make(IDSlice, len(parents))
		var s string

		for i, p := range parents {
			s = r.PathValue(fmt.Sprintf("parent_%d_id", i))
			if ids[i], err = p.ParseID(s); err != nil {
				return
			}
			ids[i].Model = p
		}
	}
	return
}

type ParentsModelIDResolver func(r *http.Request) (ids IDSlice, err error)

var DefaultParentsModelIDResolver ParentsModelIDResolver = func(r *http.Request) (ids IDSlice, err error) {
	return ParentsModelID(r), nil
}

func ResolveParentsModelID(resolver ParentsModelIDResolver, r *http.Request) (ids IDSlice, err error) {
	if resolver == nil {
		resolver = DefaultParentsModelIDResolver
	}
	return resolver(r)
}
