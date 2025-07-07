package presets

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/admin/model"
)

type (
	ID      = model.ID
	IDSlice = model.IDSlice
	Schema  = model.Schema
)

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
			if ids[i], err = p.ParseRecordID(s); err != nil {
				return
			}
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

func ParseRecordID(s Schema, v string) (id ID, err error) {
	if v == "" {
		return
	}
	id.Schema = s

	var (
		fields model.Fields
		parts  []string
	)

	if sd, _ := s.Model().(SlugDecoder); sd != nil {
		for fieldName, value := range sd.PrimaryColumnValuesBySlug(v) {
			f := id.Schema.FieldsByName(fieldName)[0]
			fields = append(fields, f)
			parts = append(parts, value)
		}
	} else {
		parts = strings.Split(v, "_")
		fields = s.PrimaryFields()
	}

	if len(fields) != len(parts) {
		err = fmt.Errorf("expected %d slug parts, got %d", len(fields), len(parts))
		return
	}

	id.Fields = fields

	modelType := reflect.TypeOf(s.Model())

	for i, v := range parts {
		var (
			field, _ = modelType.Elem().FieldByName(fields[i].Name())
			av       any
		)

		switch field.Type.Kind() {
		case reflect.String:
			av = v
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bsize := strconv.IntSize
			switch field.Type.Kind() {
			case reflect.Uint8:
				bsize = 8
			case reflect.Uint16:
				bsize = 16
			case reflect.Uint32:
				bsize = 32
			case reflect.Uint64:
				bsize = 64
			}
			var i uint64
			if i, err = strconv.ParseUint(v, 10, bsize); err != nil {
				return
			}
			switch bsize {
			case 8:
				av = uint8(i)
			case 16:
				av = uint16(i)
			case 32:
				av = uint32(i)
			case 64:
				av = i
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bsize := strconv.IntSize
			switch field.Type.Kind() {
			case reflect.Int8:
				bsize = 8
			case reflect.Int16:
				bsize = 16
			case reflect.Int32:
				bsize = 32
			case reflect.Int64:
				bsize = 64
			}

			var i int64

			if i, err = strconv.ParseInt(v, 10, bsize); err != nil {
				return
			}

			switch bsize {
			case 8:
				av = int8(i)
			case 16:
				av = int16(i)
			case 32:
				av = int32(i)
			case 64:
				av = i
			}
		default:
			fv := reflect.New(field.Type).Interface()
			if s, _ := fv.(sql.Scanner); s != nil {
				if err = s.Scan(v); err != nil {
					return
				}
				av = s
			} else {
				err = errors.New(fmt.Sprintf("Unsupported type: %v of field %s", field.Type, field.Name))
				return
			}
		}
		id.Values = append(id.Values, av)
	}
	return
}
