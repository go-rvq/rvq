package schemaform

import (
	"fmt"
	"strings"
)

const (
	// DefaultType is what an untyped field is: plain text.
	DefaultType = "str"
	// FormType is the type of a field that is itself an interface — a form
	// inside the form, or a list of them. It is also the type an application
	// gives a value it wants edited this way.
	FormType = "form"
)

// Schema is one form: the fields it edits, and whether the value it edits is a
// LIST of them rather than one.
type Schema struct {
	Slice  bool
	Fields []*Field
	// Item is what a list of PLAIN VALUES holds — `[]str`, a list of lines,
	// against `[]{…}`, a list of records. It is set (and Slice with it) only
	// for that shape, and then there are no Fields: the item IS the value.
	Item *Field
	// Enums are the enums the schema's fields hold, by name. A field typed with
	// one carries it in Field.Enum.
	Enums map[string]*Enum
	// Meta is the `[k=v, …]` block written before the interface — what the
	// schema says about how it is DRAWN, not about the value.
	Meta Meta
	// Layout is how a list of records is drawn: LayoutForm (each record a form,
	// one under the other) or LayoutTable (one row per record). Written as
	// `[layout="table"]`.
	Layout string
	// Columns are the fields a table shows, in the order it shows them —
	// `[columns=[#label, #href]]`. Empty: every field, in the schema's order. A
	// field left out keeps its value; it is only not drawn.
	Columns []string
}

// The layouts a list of records may be drawn in.
const (
	LayoutForm  = "form"
	LayoutTable = "table"
)

// Meta is a `[k=v, …]` block read into Go values: a string, an int64, a
// float64, a bool, a []any, or a Meta nested.
type Meta map[string]any

// String is the value of key when it is a string, or "".
func (m Meta) String(key string) string {
	s, _ := m[key].(string)
	return s
}

// Strings is the value of key when it is a list of strings — symbols included,
// since `#label` is the string "label".
func (m Meta) Strings(key string) ([]string, error) {
	raw, ok := m[key]
	if !ok {
		return nil, nil
	}
	list, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("%s: want a list, got %T", key, raw)
	}
	out := make([]string, len(list))
	for i, v := range list {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%s[%d]: want a name, got %T", key, i, v)
		}
		out[i] = s
	}
	return out, nil
}

// Enum is a closed set of values a field may hold, declared beside the
// interface:
//
//	enum Perm { Read, Write }
//	{ perm Perm }
type Enum struct {
	Name string
	// Names are its members, in the order declared. The NAME is what the form
	// stores — `Read`, not the 1 behind it — because the value is written out
	// as YAML and read back by name.
	Names []string
}

// Field is one entry of a form.
type Field struct {
	Name string
	// Type names the component that renders it — DefaultType when the schema
	// left it out, FormType when the field is an interface of its own.
	Type string
	// Nullable is the `?` after the name: the field accepts nil. A field that
	// is NOT nullable is required — it may not be left empty.
	Nullable bool
	// Schema is the field's own form, set when Type is FormType.
	Schema *Schema
	// Enum is the set of values the field may hold, set when its type is an
	// enum.
	Enum *Enum
	// Meta is the `[k=v, …]` block written before the field.
	Meta Meta
}

// Required reports whether the field may be left empty: only an optional one
// may, and optional is what the `?` after the name says.
func (f *Field) Required() bool { return !f.Nullable }

// FieldAt is the field at PATH — the names from the root down, the same path
// the form asks its words by. A list adds no name of its own, so the field of
// the records in `links` is `links.href`, and the item of a list of plain
// values is the list's own path.
//
// It is nil when the path names nothing.
func (s *Schema) FieldAt(path string) *Field {
	if path == "" {
		return s.Item
	}

	cur := s
	var found *Field

	for _, name := range strings.Split(path, ".") {
		if cur == nil {
			return nil
		}
		found = nil
		for _, f := range cur.Fields {
			if f.Name == name {
				found = f
				break
			}
		}
		if found == nil {
			return nil
		}
		cur = found.Schema
		// A list of plain values ends the path: its item has no name, so the
		// path of the list is the path of what it holds.
		if cur != nil && cur.Item != nil {
			found = cur.Item
			cur = cur.Item.Schema
		}
	}
	return found
}

// TableColumns are the fields a table layout draws, in order: the ones Columns
// names, or every field when it names none.
func (s *Schema) TableColumns() []*Field {
	if len(s.Columns) == 0 {
		return s.Fields
	}
	out := make([]*Field, 0, len(s.Columns))
	for _, name := range s.Columns {
		for _, f := range s.Fields {
			if f.Name == name {
				out = append(out, f)
				break
			}
		}
	}
	return out
}
