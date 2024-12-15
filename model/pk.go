package model

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/qor5/web/v3/zeroer"
)

type ID struct {
	Values []any
	Fields Fields
	Schema Schema
}

func (s ID) Related(schema Schema, fieldName ...string) ID {
	s.Schema = schema
	fields := make(Fields, len(s.Fields))
	values := make([]any, len(s.Values))

	for i, fieldName := range fieldName {
		fields[i] = schema.FieldsByName(fieldName)[0]
		values[i] = s.Values[i]
	}
	s.Fields = fields
	s.Values = values
	return s
}

func (s ID) Value() any {
	if len(s.Values) != 1 {
		panic("model.ID.Value(): expected one value")
	}
	return s.Values[0]
}

func (id ID) String() string {
	if id.IsZero() {
		return ""
	}
	return strings.Join(id.StringValues(), "_")
}

func (id ID) StringValues() (s []string) {
	s = make([]string, len(id.Values))
	for i, v := range id.Values {
		s[i] = fmt.Sprintf("%v", v)
	}
	return
}

func (id ID) IsZero() bool {
	for _, value := range id.Values {
		if !zeroer.IsZero(value) {
			return false
		}
	}
	return true
}

func (id ID) SetTo(obj interface{}) {
	ov := reflect.ValueOf(obj).Elem()
	for i, f := range id.Fields {
		field, value := ov.FieldByName(f.Name()), reflect.ValueOf(id.Values[i])
		if field.Kind() == reflect.Pointer {
			field.Set(reflect.New(field.Type().Elem()))
			field = field.Elem()
			if value.Kind() == reflect.Ptr {
				if !value.IsNil() {
					field.Set(value.Elem().Convert(field.Type()))
				}
			} else {
				field.Set(value.Elem().Convert(field.Type()))
			}
		} else {
			field.Set(value.Convert(field.Type()))
		}
	}
}

func (id ID) GetValue(fieldName string) interface{} {
	for i, field := range id.Fields {
		if field.Name() == fieldName {
			return id.Values[i]
		}
	}
	return nil
}

func (id ID) Without(fieldName ...string) ID {
	var (
		fields Fields
		values []any
	)

	for i, field := range id.Fields {
		if !slices.Contains(fieldName, field.Name()) {
			fields = append(fields, field)
			values = append(values, id.Values[i])
		}
	}

	id.Fields = fields
	id.Values = values
	return id
}

func (id ID) WithField(fieldName string, value any) ID {
	for i, field := range id.Fields {
		if field.Name() == fieldName {
			id.Values[i] = value
			return id
		}
	}

	id.Fields = append(id.Fields, id.Schema.FieldsByName(fieldName)[0])
	id.Values = append(id.Values, value)
	return id
}

type IDSlice []ID

func (s IDSlice) Len() int {
	return len(s)
}

func (s IDSlice) Last() ID {
	return s[len(s)-1]
}

func (s IDSlice) LastValues() []any {
	return s.Last().Values
}

func (s *IDSlice) Pop() ID {
	v := (*s)[len((*s))-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *IDSlice) PopValues() []any {
	return s.Pop().Values
}

func (s IDSlice) Values() (values [][]any) {
	values = make([][]any, len(s))
	for i, id := range s {
		values[i] = id.Values
	}
	return
}
