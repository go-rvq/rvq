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
	"github.com/google/uuid"
)

type (
	ID      = model.ID
	IDSlice = model.IDSlice
	Schema  = model.Schema
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	uuidType  = reflect.TypeOf(uuid.UUID{})
)

// IdParserFallback holds per-type fallback parsers for record ids, tried as a
// last resort when a primary-key field is neither a basic kind, nor uuid.UUID,
// nor exposes Parse(string) (T, error), nor implements sql.Scanner. Register a
// parser to support a custom id type without changing this package.
var IdParserFallback = map[reflect.Type]func(string) (any, error){}

// parseByMethod parses v using a Parse(string) (T, error) method exposed by t
// (or *t), returning the parsed value. It supports id types such as uuid.UUID
// that carry their own string parser. ok is false when no such method exists.
func parseByMethod(t reflect.Type, v string) (av any, ok bool, err error) {
	for _, rt := range []reflect.Type{t, reflect.PointerTo(t)} {
		m, found := rt.MethodByName("Parse")
		if !found {
			continue
		}
		ft := m.Func.Type()
		// method signature (receiver included): func(recv, string) (T, error)
		if ft.NumIn() != 2 || ft.In(1).Kind() != reflect.String ||
			ft.NumOut() != 2 || !ft.Out(1).Implements(errorType) {
			continue
		}
		var recv reflect.Value
		if rt.Kind() == reflect.Pointer {
			recv = reflect.New(rt.Elem())
		} else {
			recv = reflect.New(rt).Elem()
		}
		out := m.Func.Call([]reflect.Value{recv, reflect.ValueOf(v)})
		if !out[1].IsNil() {
			return nil, true, out[1].Interface().(error)
		}
		res := out[0]
		if res.Kind() == reflect.Pointer {
			res = res.Elem()
		}
		return res.Interface(), true, nil
	}
	return nil, false, nil
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
		// Resolve the decoder's parts (it names them freely) to schema fields.
		byField := make(map[string]string, len(s.PrimaryFields()))
		for fieldName, value := range sd.PrimaryColumnValuesBySlug(v) {
			// the decoder names its parts freely; an unknown one is an error, not
			// an index panic
			f := id.Schema.FieldsByName(fieldName)
			if len(f) == 0 {
				return id, fmt.Errorf("slug %q refers to unknown field %q", v, fieldName)
			}
			byField[f[0].Name()] = value
		}
		// Emit the fields/values in the schema's primary-key order — NOT the map's
		// random iteration order — so the resulting ID and its String() are stable
		// and match MustRecordID().String(). Otherwise a composite slug like
		// "2_en-US" could round-trip to "en-US_2", and re-decoding that would feed
		// "en-US" to the uint ID field (strconv.ParseUint: invalid syntax).
		for _, pf := range s.PrimaryFields() {
			if value, ok := byField[pf.Name()]; ok {
				fields = append(fields, pf)
				parts = append(parts, value)
			}
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
			fieldName = fields[i].Name()
			fieldType = fields[i].Type()
			av        any
		)

		if fieldType == nil {
			// the field does not know its type: fall back to the model's struct
			sf, _ := modelType.Elem().FieldByName(fieldName)
			fieldType = sf.Type
		}
		if fieldType == nil {
			err = fmt.Errorf("field %s of %s has no type", fieldName, modelType)
			return
		}

		switch fieldType.Kind() {
		case reflect.String:
			av = v
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bsize := strconv.IntSize
			switch fieldType.Kind() {
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
			switch fieldType.Kind() {
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
			// uuid.UUID (an [16]byte array) parses from its string form.
			if fieldType == uuidType {
				var u uuid.UUID
				if u, err = uuid.Parse(v); err != nil {
					return
				}
				av = u
				break
			}
			// types exposing Parse(string) (T, error) parse from their string
			// form.
			if parsed, hasParse, e := parseByMethod(fieldType, v); hasParse {
				if e != nil {
					err = e
					return
				}
				av = parsed
				break
			}
			fv := reflect.New(fieldType).Interface()
			if s, _ := fv.(sql.Scanner); s != nil {
				if err = s.Scan(v); err != nil {
					return
				}
				// store the scanned value (deref the pointer), not the *T, so
				// ID.SetTo can Convert it to the field type without panicking.
				av = reflect.ValueOf(s).Elem().Interface()
				break
			}
			// last resort: a globally registered fallback parser.
			if fn := IdParserFallback[fieldType]; fn != nil {
				if av, err = fn(v); err != nil {
					return
				}
				break
			}
			err = errors.New(fmt.Sprintf("Unsupported type: %v of field %s", fieldType, fieldName))
			return
		}
		id.Values = append(id.Values, av)
	}
	return
}
