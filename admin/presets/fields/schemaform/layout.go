package schemaform

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gad-lang/gad"
)

// Layout is a way a schema may be drawn — `[layout="grid"]` — and what it may
// be told about how: its Config, a gad class written as code
//
//	class Config { columns int = 4 }
//
// whose fields are what `layout={name: "grid", columns: 3}` may say, typed and
// with the default a field has when it is not said. A layout is registered
// once (RegisterLayout); its Config is run and read like a schema, so what it
// declares is what is checked, defaulted — and, by an application, offered to
// edit (Config).
type Layout struct {
	Name string
	// Records says the layout draws a list of records, and only that.
	Records bool
	// Config are the fields of its Config class: what it may be told. nil when
	// it may be told nothing.
	Config *Schema
	// Check refuses what the Config's types alone cannot — a grid of 40
	// columns. It is given the config already typed and defaulted.
	Check func(s *Schema, config Meta) error
}

// The field metadata a Config may use:
//
//   - `fields=true` — the field names fields of the list the layout draws: a
//     `Field`, or a list of them; each must be one.
//   - `empty_as_all` — a list of fields left empty means every field, in the
//     schema's order.
//   - `sorted` — the order of the list is the order things are drawn in.
const (
	MetaFields     = "fields"
	MetaEmptyAsAll = "empty_as_all"
	MetaSorted     = "sorted"
)

// FieldRefType is the type of a Config field that names a field of the list:
// `columns []Field`.
const FieldRefType = "Field"

var (
	layoutsMu   sync.RWMutex
	layouts     = map[string]*Layout{}
	layoutNames []string
)

// RegisterLayout registers the layout name, with its Config written as code
// ("" when it takes none) and a check of its own (or nil). records says it
// draws a list of records. A name registered again is replaced.
func RegisterLayout(name string, records bool, config string, check func(*Schema, Meta) error) error {
	l := &Layout{Name: name, Records: records, Check: check}
	if strings.TrimSpace(config) != "" {
		s, err := ParseConfig(config)
		if err != nil {
			return fmt.Errorf("schemaform: layout %q: %w", name, err)
		}
		l.Config = s
	}

	layoutsMu.Lock()
	defer layoutsMu.Unlock()
	if _, ok := layouts[name]; !ok {
		layoutNames = append(layoutNames, name)
	}
	layouts[name] = l
	return nil
}

// MustRegisterLayout is RegisterLayout, for a layout the program cannot be
// without: it panics on a Config that does not read.
func MustRegisterLayout(name string, records bool, config string, check func(*Schema, Meta) error) {
	if err := RegisterLayout(name, records, config, check); err != nil {
		panic(err)
	}
}

// LookupLayout is the layout registered as name, or nil.
func LookupLayout(name string) *Layout {
	layoutsMu.RLock()
	defer layoutsMu.RUnlock()
	return layouts[name]
}

// LayoutNames are the registered layouts, in the order they were registered.
func LayoutNames() []string {
	layoutsMu.RLock()
	defer layoutsMu.RUnlock()
	return append([]string(nil), layoutNames...)
}

// ParseConfig reads a layout's Config: the code declares `class Config {…}`,
// and its fields are read as a schema's are — types, metadata — with the
// default each has. `Field` is in scope, as the type of a field that names a
// field of the list (FieldRefType).
func ParseConfig(src string) (*Schema, error) {
	b := New().Type(FieldRefType, TextComponentFunc)
	vm, ret, err := b.runProgram(src + "\nreturn Config\n")
	if err != nil {
		return nil, err
	}
	c, ok := ret.(*gad.Class)
	if !ok {
		return nil, fmt.Errorf("schemaform: Config is not a class, it is %s", ret.Type().Name())
	}
	r := &reader{vm: vm, enums: map[string]*Enum{}, visiting: map[*gad.Interface]bool{}}
	return r.class(c)
}

// read types what cfg says against the Config, fills in the defaults of what
// it does not say, and checks it: for the schema s the layout draws.
func (l *Layout) read(s *Schema, cfg Meta) (Meta, error) {
	if l.Config == nil {
		if len(cfg) > 0 {
			return nil, fmt.Errorf("it takes no config, and was given %v", keysOf(cfg))
		}
		return nil, nil
	}

	out := Meta{}
	for k, v := range cfg {
		f := l.Config.field(k)
		if f == nil {
			return nil, fmt.Errorf("it has no %q (it has %v)", k, l.Config.fieldNames())
		}
		tv, err := configValue(f, v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", k, err)
		}
		out[k] = tv
	}
	for _, f := range l.Config.Fields {
		if _, ok := out[f.Name]; !ok && f.Default != nil {
			tv, err := configValue(f, f.Default)
			if err != nil {
				return nil, fmt.Errorf("%s: the default: %w", f.Name, err)
			}
			out[f.Name] = tv
		}
	}

	// what names fields of the list must name them
	for _, f := range l.Config.Fields {
		if !metaBool(f.Meta, MetaFields) {
			continue
		}
		for _, name := range configStrings(out[f.Name]) {
			if s.field(name) == nil {
				return nil, fmt.Errorf("%s: %q names no field", f.Name, name)
			}
		}
	}

	if l.Check != nil {
		if err := l.Check(s, out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// configValue is v typed as the Config field f says: a number, a text, a yes
// or no, or a list of them.
func configValue(f *Field, v any) (any, error) {
	if f.Schema != nil && f.Schema.Slice && f.Schema.Item != nil {
		list, ok := v.([]any)
		if !ok {
			if ss, isStrings := v.([]string); isStrings {
				for _, s := range ss {
					list = append(list, s)
				}
			} else {
				return nil, fmt.Errorf("want a list, got %T", v)
			}
		}
		out := make([]any, len(list))
		for i, e := range list {
			tv, err := configValue(f.Schema.Item, e)
			if err != nil {
				return nil, fmt.Errorf("[%d]: %w", i, err)
			}
			out[i] = tv
		}
		return out, nil
	}

	switch f.Type {
	case "int", "uint":
		switch n := v.(type) {
		case int64:
			return int(n), nil
		case int:
			return n, nil
		case float64:
			if n == float64(int(n)) {
				return int(n), nil
			}
		}
		return nil, fmt.Errorf("want a whole number, got %v", v)
	case "bool":
		if b, ok := v.(bool); ok {
			return b, nil
		}
		return nil, fmt.Errorf("want true or false, got %v", v)
	case "float", "decimal":
		switch n := v.(type) {
		case float64:
			return n, nil
		case int64:
			return float64(n), nil
		case int:
			return float64(n), nil
		}
		return nil, fmt.Errorf("want a number, got %v", v)
	default: // str, Field, text…
		if s, ok := v.(string); ok {
			return s, nil
		}
		return nil, fmt.Errorf("want a text, got %T", v)
	}
}

// configStrings are the texts a config value holds: itself, or its items.
func configStrings(v any) []string {
	switch t := v.(type) {
	case string:
		return []string{t}
	case []any:
		out := make([]string, 0, len(t))
		for _, e := range t {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// metaBool is a flag of the metadata: `[sorted]` or `[sorted=true]`.
func metaBool(m Meta, key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	b, isBool := v.(bool)
	return !isBool || b
}

func keysOf(m Meta) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// LayoutFields are the fields of the list the config's key names (a Config
// field with `fields=true`), in the order it names them — every field, in the
// schema's order, when it names none and the Config says `empty_as_all`.
func (s *Schema) LayoutFields(key string) []*Field {
	names := configStrings(s.LayoutConfig[key])
	if len(names) == 0 {
		if l := LookupLayout(s.Layout); l != nil && l.Config != nil {
			if f := l.Config.field(key); f != nil && metaBool(f.Meta, MetaEmptyAsAll) {
				return s.Fields
			}
		}
		return nil
	}
	out := make([]*Field, 0, len(names))
	for _, name := range names {
		if f := s.field(name); f != nil {
			out = append(out, f)
		}
	}
	return out
}

// LayoutInt is a whole number of the layout's config, 0 when it has none.
func (s *Schema) LayoutInt(key string) int {
	n, _ := s.LayoutConfig[key].(int)
	return n
}

// field is the schema's field name, or nil.
func (s *Schema) field(name string) *Field {
	for _, f := range s.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (s *Schema) fieldNames() []string {
	out := make([]string, len(s.Fields))
	for i, f := range s.Fields {
		out[i] = f.Name
	}
	return out
}

// The layouts there are by default. A table's `columns` are the fields it
// shows, in order — every field when none; a grid's `columns` are how many
// cards it puts side by side, 1 to MaxGridColumns.
func init() {
	MustRegisterLayout(LayoutForm, false, "", nil)
	MustRegisterLayout(LayoutTable, true,
		"class Config { ["+MetaFields+"=true, "+MetaEmptyAsAll+", "+MetaSorted+"] columns []"+FieldRefType+" }", nil)
	MustRegisterLayout(LayoutGrid, true,
		fmt.Sprintf("class Config { columns int = %d }", DefaultGridColumns),
		func(_ *Schema, config Meta) error {
			if n, _ := config["columns"].(int); n < 1 || n > MaxGridColumns {
				return fmt.Errorf("columns %d: a grid has 1 to %d columns", n, MaxGridColumns)
			}
			return nil
		})
}
