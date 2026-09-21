package schemaform

import (
	"net/url"
	"strconv"
	"strings"

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
func (s *Schema) Decode(values url.Values, key string) any {
	return decodeSchema(s, values, key)
}

func decodeSchema(s *Schema, values url.Values, prefix string) any {
	if !s.Slice {
		return decodeRecord(s, values, prefix)
	}

	// A list is as long as the indexes the form carries: the browser posts one
	// per item, in order, and stops.
	items := []any{}
	for i := 0; ; i++ {
		at := prefix + "[" + strconv.Itoa(i) + "]"
		if !hasKeyUnder(values, at) {
			break
		}
		if s.Item != nil {
			items = append(items, decodeField(s.Item, values, at))
			continue
		}
		items = append(items, decodeRecord(s, values, at))
	}
	return items
}

func decodeRecord(s *Schema, values url.Values, prefix string) Record {
	rec := make(Record, 0, len(s.Fields))
	for _, f := range s.Fields {
		name := f.Name
		if prefix != "" {
			name = prefix + "." + f.Name
		}
		rec = append(rec, RecordField{Name: f.Name, Value: decodeField(f, values, name)})
	}
	return rec
}

func decodeField(f *Field, values url.Values, path string) any {
	if f.Schema != nil {
		return decodeSchema(f.Schema, values, path)
	}
	return decodeValue(f, values, path)
}

// decodeValue reads ONE input. An empty input is nil when the field accepts nil
// and the empty value of its type otherwise — a required field that was left
// empty is written as empty, not dropped, so the value keeps its shape.
func decodeValue(f *Field, values url.Values, path string) any {
	raw, present := values[path]
	var s string
	if len(raw) > 0 {
		s = raw[0]
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
		if f.Type == "int" || f.Type == "uint" {
			return 0
		}
		if f.Type == "float" || f.Type == "decimal" {
			return 0.0
		}
	}
	return s
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
