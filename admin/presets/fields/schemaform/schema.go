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
	// Choice says the record holds ONE of its fields, not all: a field typed by
	// a union of classes (`layout listLayout|gridLayout`), each class a field —
	// its name the class's, its schema the class's record. The value is an
	// object of one key, the class chosen, holding what that class holds:
	// `{gridLayout: {columns: 3}}`.
	Choice bool
	// Layout is how the schema is drawn, one of the registered layouts (see
	// Layout): LayoutForm, the list (each record a card of its form, one under the other),
	// LayoutTable (one row per record) or LayoutGrid (each record a card).
	// Written `[layout="table"]`, or with its config
	// `[layout={name: "grid", columns: 3}]`.
	Layout string
	// LayoutConfig is what the layout was told, typed by its Config and with
	// its defaults: a table's `columns` (the fields it shows), a grid's
	// `columns` (the cards side by side).
	LayoutConfig Meta
}

// The layouts a list of records may be drawn in.
const (
	LayoutForm  = "form"
	LayoutTable = "table"
	LayoutGrid  = "grid"
)

// DefaultGridColumns is how many cards a grid puts side by side when the schema
// does not say (its Config's default), and MaxGridColumns the most it may (the
// width of the grid it is drawn in).
const (
	DefaultGridColumns = 4
	MaxGridColumns     = 12
)

// Meta is a `[k=v, …]` block read into Go values: a string, an int64, a
// float64, a bool, a []any, or a Meta nested.
type Meta map[string]any

// String is the value of key when it is a string, or "".
func (m Meta) String(key string) string {
	s, _ := m[key].(string)
	return s
}

// Int is the value of key when it is a whole number; ok is false when the key is
// absent, and err says a value of another kind.
func (m Meta) Int(key string) (n int, ok bool, err error) {
	raw, ok := m[key]
	if !ok {
		return 0, false, nil
	}
	switch v := raw.(type) {
	case int64:
		return int(v), true, nil
	case int:
		return v, true, nil
	}
	return 0, true, fmt.Errorf("%s: want a whole number, got %T", key, raw)
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
	// Items are the members' labels and hints when the schema itself gives
	// them — a choice of classes read as their names (Builder.ChoiceAsName)
	// takes each class's metadata, a field's `[options=…]` its labels; nil
	// otherwise.
	Items []EnumItem
	// Options says the enum is the field's own `[options=…]`, not a type: it
	// is drawn as a select whatever the field's type is.
	Options bool
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
	// Default is the value a class field declares (`columns int = 4`), read
	// into Go as metadata is; nil when it has none.
	Default any
	// ReadOnly is a field declared `get name Type` — a getter: shown, never
	// edited. The form draws it read-only, and Schema.KeepReadOnly puts back
	// the value stored over whatever a post said of it.
	ReadOnly bool
	// Owner is the class that declared it, when the schema is a class
	// extending others: the parent's name for a field it took from one
	// (`*ContactForm`), "" for its own.
	Owner string
	// Row is the group the field is a column of — `{ … }` in the body, the
	// field `$N`, whose fields are the record's own —, 1… in its record; 0
	// for a field alone in its row.
	Row int
	// Percent is the field's share of its group's row: its `[width=N]`, N
	// of WidthUnits, or its part of what the others left.
	Percent float64
	// Group is the group the field is of — its kind: a row, tabs, steps —;
	// nil for a field of no group.
	Group *Group
}

// Clone is a copy of the schema that may be changed — a schema Parse returns
// is shared (cached), and never to be changed. The fields, and the schemas
// inside them, are copied; the enums and the metadata, which nothing changes,
// are shared.
func (s *Schema) Clone() *Schema {
	if s == nil {
		return nil
	}
	c := *s
	if s.LayoutConfig != nil {
		c.LayoutConfig = make(Meta, len(s.LayoutConfig))
		for k, v := range s.LayoutConfig {
			c.LayoutConfig[k] = v
		}
	}
	c.Item = s.Item.clone()
	c.Fields = make([]*Field, len(s.Fields))
	for i, f := range s.Fields {
		c.Fields[i] = f.clone()
	}
	return &c
}

func (f *Field) clone() *Field {
	if f == nil {
		return nil
	}
	c := *f
	c.Schema = f.Schema.Clone()
	return &c
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

// TableColumns are the fields a table layout draws, in order: the ones its
// `columns` name, or every field when they name none.
func (s *Schema) TableColumns() []*Field {
	if cols := s.LayoutFields("columns"); len(cols) > 0 {
		return cols
	}
	return s.Fields
}

// GridColumns is how many cards a grid layout puts side by side.
func (s *Schema) GridColumns() int {
	if n := s.LayoutInt("columns"); n > 0 {
		return n
	}
	return DefaultGridColumns
}

// WithLayout is a copy of the schema drawn in the layout name, told cfg —
// checked as the metadata would be. The schema itself is not changed (Parse
// shares it).
func (s *Schema) WithLayout(name string, cfg Meta) (*Schema, error) {
	c := s.Clone()
	if err := c.setLayout(name, cfg); err != nil {
		return nil, err
	}
	return c, nil
}

// Filled is v with the shape the schema describes filled in where it is
// missing: a record inside the record an empty record, a list an empty list —
// all the way down. A form binds its inputs through those (`value.layout.name`),
// and a record that is not there breaks the whole form in the browser; what v
// holds is kept as it is.
func (s *Schema) Filled(v any) any {
	if s == nil {
		return v
	}
	if s.Slice {
		list, _ := v.([]any)
		if list == nil {
			list = []any{}
		}
		for i, item := range list {
			switch {
			case s.Item == nil:
				list[i] = s.record(item)
			case s.Item.Schema != nil:
				list[i] = s.Item.Schema.Filled(item)
			}
		}
		return list
	}
	if s.Choice {
		return s.chosen(v)
	}
	return s.record(v)
}

// chosen is v as a choice: the class it holds filled in, and no other.
func (s *Schema) chosen(v any) map[string]any {
	m, _ := v.(map[string]any)
	if m == nil {
		m = map[string]any{}
	}
	for _, f := range s.Fields {
		if val, ok := m[f.Name]; ok && f.Schema != nil {
			m[f.Name] = f.Schema.Filled(val)
		}
	}
	return m
}

// record is v as a record of the schema's fields, the records and lists inside
// it filled in, and a field it lacks given the default its class declares.
func (s *Schema) record(v any) map[string]any {
	m, _ := v.(map[string]any)
	if m == nil {
		m = map[string]any{}
	}
	for _, f := range s.Fields {
		switch {
		case f.Schema != nil:
			m[f.Name] = f.Schema.Filled(m[f.Name])
		case f.Default != nil:
			if _, ok := m[f.Name]; !ok {
				m[f.Name] = f.Default
			}
		}
	}
	return m
}
