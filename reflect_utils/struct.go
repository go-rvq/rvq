package reflect_utils

import (
	"go/ast"
	"reflect"
	"strings"
)

func GetStruct(t reflect.Type) interface{} {
	if t.Kind() == reflect.Struct {
		return reflect.New(t).Interface()
	}
	return GetStruct(t.Elem())
}

type (
	IndexableStructField struct {
		reflect.StructField
		Index []int
		Names []string
	}

	IndexableStructFields []*IndexableStructField
)

func (f *IndexableStructField) String() string {
	return strings.Join(f.Names, "/")
}

func (f IndexableStructFields) String() string {
	var names []string
	for _, field := range f {
		names = append(names, field.String())
	}
	return "[" + strings.Join(names, ", ") + "]"
}

func (f IndexableStructFields) Len() int {
	return len(f)
}

func (f IndexableStructFields) Get(name string) *IndexableStructField {
	for _, field := range f {
		if field.Name == name {
			return field
		}
	}
	return nil
}

func (f IndexableStructFields) Uniquefy() (r IndexableStructFields) {
	names := make(map[string]any)
	for _, field := range f {
		if _, ok := names[field.Name]; ok {
			panic("duplicated field name: " + field.Name)
		}
		names[field.Name] = nil
		r = append(r, field)
	}
	return
}

func (f IndexableStructFields) Only(names ...string) (s IndexableStructFields) {
	for _, f := range f {
		for _, name := range names {
			if f.Name == name {
				s = append(s, f)
				break
			}
		}
	}
	return
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func UniqueFieldsOfReflectType(ityp reflect.Type) (dotFields, result IndexableStructFields) {
	dotFields, result = FieldsOfReflectType(ityp)
	result = result.Uniquefy()
	return
}

func FieldsOfReflectType(ityp reflect.Type) (dotFields, fields IndexableStructFields) {
	FieldsOfReflectTypeCB(
		ityp,
		func(field *IndexableStructField) {
			dotFields = append(dotFields, field)
		},
		func(field *IndexableStructField) {
			fields = append(fields, field)
		},
	)
	return
}

func FieldsOfReflectTypeCB(ityp reflect.Type, appendDotField, appendField func(*IndexableStructField)) {
	if appendDotField == nil {
		appendDotField = func(dotField *IndexableStructField) {}
	}

	if appendField == nil {
		appendField = func(dotField *IndexableStructField) {}
	}

	ityp = indirectType(ityp)

	var walk func(typ reflect.Type, path []int, name []string)

	walk = func(typ reflect.Type, path []int, name []string) {
		typ = indirectType(typ)
		if typ.Kind() != reflect.Struct || (path != nil && typ == ityp) {
			return
		}

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			path := append(append([]int{}, path...), i)
			name := append(append([]string{}, name...), field.Name)

			if field.Anonymous {
				walk(field.Type, path, name)
			} else if field.Name == "_" {
				appendDotField(&IndexableStructField{field, path, name})
			} else if ast.IsExported(field.Name) {
				appendField(&IndexableStructField{field, path, name})
				if field.Anonymous {
					walk(field.Type, path, name)
				}
			}
		}
	}

	walk(ityp, nil, nil)
}
