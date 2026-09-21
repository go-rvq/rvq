package schemaform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/web"
	"gopkg.in/yaml.v3"
)

// Record is a decoded record: its fields IN THE ORDER THE SCHEMA DECLARES them.
// A map would sort them alphabetically and rewrite the whole value on every
// save, so the order the schema gives is kept all the way to the file.
type Record []RecordField

// RecordField is one field of a decoded record.
type RecordField struct {
	Name  string
	Value any
}

// MarshalJSON writes the record as an object, in its own order — encoding/json
// would sort a map's keys and rewrite the whole value on every save.
func (r Record) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, f := range r {
		if i > 0 {
			buf.WriteByte(',')
		}
		name, err := json.Marshal(f.Name)
		if err != nil {
			return nil, err
		}
		value, err := json.Marshal(f.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(name)
		buf.WriteByte(':')
		buf.Write(value)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// MarshalYAML writes the record as a mapping, in its own order.
func (r Record) MarshalYAML() (any, error) {
	n := &yaml.Node{Kind: yaml.MappingNode}
	for _, f := range r {
		var key, value yaml.Node
		if err := key.Encode(f.Name); err != nil {
			return nil, err
		}
		if err := value.Encode(f.Value); err != nil {
			return nil, err
		}
		n.Content = append(n.Content, &key, &value)
	}
	return n, nil
}

// Decode reads back the value the form posted for a schema drawn under key.
//
// The browser flattens what the form holds by the same two rules the builder
// binds by: an array indexes (`key[0]`) and a record names (`key.label`), so a
// list of records arrives as `key[0].label`. The schema says which shape those
// names have and what each value IS — an int is an int, a switch is a bool —
// and the result is the value itself: a Record, a list of them, or a list of
// plain values, ready to be written out.
//
// It reads; it does not check. What a field may HOLD is a question for the
// builder, which is the one that knows the values a select offered: see
// Builder.DecodeForm.
func (s *Schema) Decode(values url.Values, key string) any {
	d := &decoder{values: values}
	return d.schema(s, key, "")
}

// DecodeForm reads the posted form back AND checks it: a field that holds one
// of a closed list of values — an enum the schema declared, or a list
// EnumItemsFunc answered with — may only hold a value that list offered, since
// the select that offered it is the only thing that should have written it.
//
// The value comes back whole either way, so a caller may show the errors beside
// the form the user is still editing.
func (b *Builder) DecodeForm(ctx *web.EventContext, schema *Schema, values url.Values, key string) (any, error) {
	d := &decoder{
		values: values,
		items: func(path string, f *Field) []EnumItem {
			return b.itemsOf(ctx, path, f)
		},
	}
	value := d.schema(schema, key, "")
	return value, errors.Join(d.errs...)
}

// decoder walks a schema over what the form posted. It carries two paths: the
// form KEY, which grows with each index (`links[0].href`), and the schema PATH,
// which does not (`links.href`) — a list adds no name of its own, which is the
// path the words and the values of a field are asked by.
type decoder struct {
	values url.Values
	// items are the values a field may hold, when it holds one of a list. Nil
	// when nobody is checking.
	items func(path string, f *Field) []EnumItem
	errs  []error
}

func (d *decoder) schema(s *Schema, key, path string) any {
	if !s.Slice {
		return d.record(s, key, path)
	}

	// A list is as long as the indexes the form carries: the browser posts one
	// per item, in order, and stops.
	items := []any{}
	for i := 0; ; i++ {
		at := key + "[" + strconv.Itoa(i) + "]"
		if !hasKeyUnder(d.values, at) {
			break
		}
		if s.Item != nil {
			items = append(items, d.field(s.Item, at, path))
			continue
		}
		items = append(items, d.record(s, at, path))
	}
	return items
}

func (d *decoder) record(s *Schema, key, path string) Record {
	rec := make(Record, 0, len(s.Fields))
	for _, f := range s.Fields {
		name := f.Name
		if key != "" {
			name = key + "." + f.Name
		}
		rec = append(rec, RecordField{Name: f.Name, Value: d.field(f, name, join(path, f.Name))})
	}
	return rec
}

func (d *decoder) field(f *Field, key, path string) any {
	if f.Schema != nil {
		return d.schema(f.Schema, key, path)
	}
	return d.value(f, key, path)
}

// value reads ONE input. An empty input is nil when the field accepts nil and
// the empty value of its type otherwise — a required field that was left empty
// is written as empty, not dropped, so the value keeps its shape.
func (d *decoder) value(f *Field, key, path string) any {
	raw, present := d.values[key]
	var s string
	if len(raw) > 0 {
		s = raw[0]
	}

	if d.items != nil {
		if items := d.items(path, f); len(items) > 0 {
			d.check(f, path, s, items)
		}
	}

	switch f.Type {
	case "bool":
		// A switch that is off may not post at all.
		return present && (s == "true" || s == "1" || s == "on")
	case "int":
		if v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err == nil {
			return v
		}
	case "uint":
		if v, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64); err == nil {
			return v
		}
	case "float", "decimal":
		if v, err := strconv.ParseFloat(strings.TrimSpace(s), 64); err == nil {
			return v
		}
	}

	if s == "" {
		if f.Nullable {
			return nil
		}
		switch f.Type {
		case "int", "uint":
			return 0
		case "float", "decimal":
			return 0.0
		}
	}
	return s
}

// check refuses a value the list did not offer. An empty value is the field's
// own business: a field that accepts nil may be left empty, and a required one
// that was is reported as empty, not as an impostor.
func (d *decoder) check(f *Field, path, value string, items []EnumItem) {
	where := path
	if where == "" {
		where = f.Name
	}

	if value == "" {
		if !f.Nullable {
			d.errs = append(d.errs, fmt.Errorf("%s: escolha um valor", where))
		}
		return
	}

	names := make([]string, len(items))
	for i, it := range items {
		if it.Name == value {
			return
		}
		names[i] = it.Name
	}
	d.errs = append(d.errs, fmt.Errorf("%s: %q não está entre os valores possíveis (%s)",
		where, value, strings.Join(names, ", ")))
}

// hasKeyUnder reports whether the form carries anything at this path: the value
// itself (a plain item) or a field of it (`[0].label`, `[0][1]`).
func hasKeyUnder(values url.Values, path string) bool {
	if _, ok := values[path]; ok {
		return true
	}
	for k := range values {
		if strings.HasPrefix(k, path+".") || strings.HasPrefix(k, path+"[") {
			return true
		}
	}
	return false
}
