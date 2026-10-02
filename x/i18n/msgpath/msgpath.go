// Package msgpath reaches the texts of a struct of messages (the Messages of a
// module, as i18n registers them) by the path of their fields — Title,
// Form.Submit, Help[0][1] — and reads what their i18n tag says of them: how a
// text is edited (type), its label, its hint and the fields its value holds.
//
// The tag is gad metadata (parser.ParseMetadata), its texts in single quotes:
//
//	Title   string        `i18n:"label='Page title', hint='The title of the page.'"`
//	Deleted string        `i18n:"hint='Shown once records were deleted.', fields=(;'%d'='how many')"`
//	Help    h.RawHTML     `i18n:"type=html"`
//	Usage   gadxtpl.Template `i18n:"type=gadx"`
//
// A text is a field whose kind is string (string, h.RawHTML, a template, …).
// A struct, a pointer to one, a slice or an array holds texts: its own path is
// the prefix of theirs; an embedded struct adds no name. A field of type form
// is one text however it is made: its value is YAML. Anything else — a func, a
// number, a map — is not a text, and Walk tells it as skipped.
package msgpath

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gad-lang/gad/parser"
	"github.com/gad-lang/gad/parser/node"
	"github.com/go-rvq/rvq/web/str_utils"
	"gopkg.in/yaml.v3"
)

// Type is how a text is edited.
type Type string

const (
	// Text is plain text, the default.
	Text Type = "text"
	// HTML is HTML.
	HTML Type = "html"
	// Gadx is a Gadx template (gadxtpl).
	Gadx Type = "gadx"
	// Gadt is a Gad template (gadttpl).
	Gadt Type = "gadt"
	// Form is a value made of several (a struct, a slice), one text in YAML,
	// edited by a form (schemaform) whose schema the tag gives.
	Form Type = "form"
)

// Field is one of the fields the value of a text holds — a formatting verb
// (%s, %[2]d), a {name} replaced, a template's datum — and what it is.
type Field struct {
	Name        string
	Description string
}

// Tag is what the i18n tag of a field says.
type Tag struct {
	// Type is how the text is edited; Text when the tag says nothing.
	Type Type
	// Label names the text; empty, it is the field's name humanized (Label).
	Label string
	// Hint says what the text is, where it is shown.
	Hint string
	// Fields are the fields its value holds, in the order of the tag.
	Fields []Field
	// Schema is the schema of a Form's form.
	Schema string
}

// TagKey is the key of the tag in a struct tag.
const TagKey = "i18n"

// ParseTag reads the i18n tag of a struct tag. No i18n tag is a Text with
// nothing else said.
func ParseTag(st reflect.StructTag) (t Tag, err error) {
	t.Type = Text
	src, ok := st.Lookup(TagKey)
	if !ok || strings.TrimSpace(src) == "" {
		return
	}
	kva, err := parser.ParseMetadata(src, 0)
	if err != nil {
		return t, fmt.Errorf("i18n tag %q: %w", src, err)
	}
	for _, e := range kva.Elements {
		kv, ok := e.(*node.KeyValuePairLit)
		if !ok {
			return t, fmt.Errorf("i18n tag %q: %s is not key=value", src, e)
		}
		key, err := name(kv.Key)
		if err != nil {
			return t, fmt.Errorf("i18n tag %q: %w", src, err)
		}
		switch key {
		case "type":
			v, err := name(kv.Value)
			if err != nil {
				return t, fmt.Errorf("i18n tag %q: type: %w", src, err)
			}
			t.Type = Type(v)
		case "label", "hint", "schema":
			v, err := text(kv.Value)
			if err != nil {
				return t, fmt.Errorf("i18n tag %q: %s: %w", src, key, err)
			}
			switch key {
			case "label":
				t.Label = v
			case "hint":
				t.Hint = v
			default:
				t.Schema = v
			}
		case "fields":
			fields, ok := kv.Value.(*node.KeyValueArrayLit)
			if !ok {
				return t, fmt.Errorf("i18n tag %q: fields is not (;name='…', …)", src)
			}
			for _, e := range fields.Elements {
				f, ok := e.(*node.KeyValuePairLit)
				if !ok {
					return t, fmt.Errorf("i18n tag %q: field %s is not name='…'", src, e)
				}
				n, err := name(f.Key)
				if err != nil {
					return t, fmt.Errorf("i18n tag %q: field: %w", src, err)
				}
				d, err := text(f.Value)
				if err != nil {
					return t, fmt.Errorf("i18n tag %q: field %s: %w", src, n, err)
				}
				t.Fields = append(t.Fields, Field{Name: n, Description: d})
			}
		default:
			return t, fmt.Errorf("i18n tag %q: unknown key %s", src, key)
		}
	}
	switch t.Type {
	case Text, HTML, Gadx, Gadt, Form:
	default:
		return t, fmt.Errorf("i18n tag %q: unknown type %s", src, t.Type)
	}
	return
}

// name is an identifier or a string: a key, a type.
func name(e node.Expr) (string, error) {
	switch v := e.(type) {
	case *node.IdentExpr:
		return v.Name, nil
	case *node.StrLit:
		return v.Value(), nil
	case *node.RawStrLit:
		return v.Value(), nil
	}
	return "", fmt.Errorf("%s is not a name", e)
}

// text is a string.
func text(e node.Expr) (string, error) {
	switch v := e.(type) {
	case *node.StrLit:
		return v.Value(), nil
	case *node.RawStrLit:
		return v.Value(), nil
	}
	return "", fmt.Errorf("%s is not a string", e)
}

// Entry is a text of a struct of messages.
type Entry struct {
	// Path is where it is: Title, Form.Submit, Help[0][1].
	Path string
	// Name is the name of its field — of the slice, for an element.
	Name string
	Tag
	// Value is the text; YAML for a Form.
	Value string
}

// Label is the tag's label, or the field's name humanized.
func (e *Entry) Label() string {
	if e.Tag.Label != "" {
		return e.Tag.Label
	}
	return str_utils.HumanizeString(e.Name)
}

// Skipped is a field that holds no text.
type Skipped struct {
	Path string
	Type reflect.Type
}

// Walk is the texts of msgs, a struct or a pointer to one, in the order of its
// fields; and the fields that are not texts.
func Walk(msgs any) (entries []Entry, skipped []Skipped, err error) {
	v := reflect.Indirect(reflect.ValueOf(msgs))
	if v.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("msgpath: %T is not a struct", msgs)
	}
	w := walker{}
	if err = w.walkStruct("", v); err != nil {
		return nil, nil, err
	}
	return w.entries, w.skipped, nil
}

type walker struct {
	entries []Entry
	skipped []Skipped
}

func (w *walker) walkStruct(prefix string, v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		fv := v.Field(i)
		if sf.Anonymous && reflect.Indirect(fv).Kind() == reflect.Struct {
			if fv.Kind() == reflect.Pointer && fv.IsNil() {
				continue
			}
			if err := w.walkStruct(prefix, reflect.Indirect(fv)); err != nil {
				return err
			}
			continue
		}
		tag, err := ParseTag(sf.Tag)
		if err != nil {
			return fmt.Errorf("msgpath: %s: %w", join(prefix, sf.Name), err)
		}
		if err = w.walkValue(join(prefix, sf.Name), sf.Name, tag, fv); err != nil {
			return err
		}
	}
	return nil
}

func (w *walker) walkValue(path, name string, tag Tag, v reflect.Value) error {
	if tag.Type == Form {
		b, err := yaml.Marshal(v.Interface())
		if err != nil {
			return fmt.Errorf("msgpath: %s: %w", path, err)
		}
		w.entries = append(w.entries, Entry{Path: path, Name: name, Tag: tag, Value: string(b)})
		return nil
	}
	switch v.Kind() {
	case reflect.String:
		w.entries = append(w.entries, Entry{Path: path, Name: name, Tag: tag, Value: v.String()})
	case reflect.Struct:
		return w.walkStruct(path, v)
	case reflect.Pointer:
		if !v.IsNil() && v.Elem().Kind() == reflect.Struct {
			return w.walkStruct(path, v.Elem())
		}
		w.skipped = append(w.skipped, Skipped{Path: path, Type: v.Type()})
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := w.walkValue(path+"["+strconv.Itoa(i)+"]", name, tag, v.Index(i)); err != nil {
				return err
			}
		}
	default:
		w.skipped = append(w.skipped, Skipped{Path: path, Type: v.Type()})
	}
	return nil
}

func join(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

// segment is a step of a path: a field's name, or an index.
type segment struct {
	name  string
	index int
}

func parsePath(path string) ([]segment, error) {
	var segs []segment
	for i := 0; i < len(path); {
		switch path[i] {
		case '.':
			i++
		case '[':
			end := strings.IndexByte(path[i:], ']')
			if end < 0 {
				return nil, fmt.Errorf("msgpath: %q: [ not closed", path)
			}
			n, err := strconv.Atoi(path[i+1 : i+end])
			if err != nil || n < 0 {
				return nil, fmt.Errorf("msgpath: %q: bad index %q", path, path[i+1:i+end])
			}
			segs = append(segs, segment{index: n})
			i += end + 1
		default:
			end := strings.IndexAny(path[i:], ".[")
			if end < 0 {
				end = len(path) - i
			}
			segs = append(segs, segment{name: path[i : i+end], index: -1})
			i += end
		}
	}
	if len(segs) == 0 {
		return nil, fmt.Errorf("msgpath: empty path")
	}
	return segs, nil
}

// find is the value at path in v and the tag of its field; grow makes room
// for an index past the end of a slice (Set).
func find(v reflect.Value, path string, grow bool) (reflect.Value, Tag, error) {
	segs, err := parsePath(path)
	if err != nil {
		return reflect.Value{}, Tag{}, err
	}
	tag := Tag{Type: Text}
	for _, s := range segs {
		for v.Kind() == reflect.Pointer {
			if v.IsNil() {
				if !grow {
					return reflect.Value{}, tag, fmt.Errorf("msgpath: %s: nil", path)
				}
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
		if s.index < 0 {
			if v.Kind() != reflect.Struct {
				return reflect.Value{}, tag, fmt.Errorf("msgpath: %s: %s is not a field of a %s", path, s.name, v.Type())
			}
			sf, ok := v.Type().FieldByName(s.name)
			if !ok || !sf.IsExported() {
				return reflect.Value{}, tag, fmt.Errorf("msgpath: %s: no field %s", path, s.name)
			}
			if tag, err = ParseTag(sf.Tag); err != nil {
				return reflect.Value{}, tag, fmt.Errorf("msgpath: %s: %w", path, err)
			}
			v = v.FieldByIndex(sf.Index)
			if tag.Type == Form {
				continue
			}
			continue
		}
		switch v.Kind() {
		case reflect.Slice:
			if s.index >= v.Len() {
				if !grow {
					return reflect.Value{}, tag, fmt.Errorf("msgpath: %s: index %d out of %d", path, s.index, v.Len())
				}
				v.Set(reflect.AppendSlice(v, reflect.MakeSlice(v.Type(), s.index+1-v.Len(), s.index+1-v.Len())))
			}
		case reflect.Array:
			if s.index >= v.Len() {
				return reflect.Value{}, tag, fmt.Errorf("msgpath: %s: index %d out of %d", path, s.index, v.Len())
			}
		default:
			return reflect.Value{}, tag, fmt.Errorf("msgpath: %s: a %s has no index", path, v.Type())
		}
		v = v.Index(s.index)
	}
	return v, tag, nil
}

// Get is the text at path in msgs; YAML for a Form.
func Get(msgs any, path string) (string, error) {
	v, tag, err := find(reflect.ValueOf(msgs), path, false)
	if err != nil {
		return "", err
	}
	if tag.Type == Form {
		b, err := yaml.Marshal(v.Interface())
		return string(b), err
	}
	if v.Kind() != reflect.String {
		return "", fmt.Errorf("msgpath: %s: a %s is not a text", path, v.Type())
	}
	return v.String(), nil
}

// Set sets the text at path in msgs, a pointer to a struct; YAML for a Form.
// An index past the end of a slice grows it.
func Set(msgs any, path, value string) error {
	rv := reflect.ValueOf(msgs)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("msgpath: %T is not a pointer", msgs)
	}
	v, tag, err := find(rv, path, true)
	if err != nil {
		return err
	}
	if tag.Type == Form {
		n := reflect.New(v.Type())
		if err = yaml.Unmarshal([]byte(value), n.Interface()); err != nil {
			return fmt.Errorf("msgpath: %s: %w", path, err)
		}
		v.Set(n.Elem())
		return nil
	}
	if v.Kind() != reflect.String {
		return fmt.Errorf("msgpath: %s: a %s is not a text", path, v.Type())
	}
	v.SetString(value)
	return nil
}

// Clone is a deep copy of msgs, a pointer to a struct: its texts may be Set
// without touching those of msgs. A func, a chan or an interface holding no
// struct is shared.
func Clone[T any](msgs T) T {
	return deepCopy(reflect.ValueOf(msgs)).Interface().(T)
}

func deepCopy(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return v
		}
		n := reflect.New(v.Type().Elem())
		n.Elem().Set(deepCopy(v.Elem()))
		return n
	case reflect.Struct:
		n := reflect.New(v.Type()).Elem()
		n.Set(v)
		for i := 0; i < v.NumField(); i++ {
			if n.Field(i).CanSet() {
				n.Field(i).Set(deepCopy(v.Field(i)))
			}
		}
		return n
	case reflect.Slice:
		if v.IsNil() {
			return v
		}
		n := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			n.Index(i).Set(deepCopy(v.Index(i)))
		}
		return n
	case reflect.Array:
		n := reflect.New(v.Type()).Elem()
		for i := 0; i < v.Len(); i++ {
			n.Index(i).Set(deepCopy(v.Index(i)))
		}
		return n
	case reflect.Map:
		if v.IsNil() {
			return v
		}
		n := reflect.MakeMapWithSize(v.Type(), v.Len())
		it := v.MapRange()
		for it.Next() {
			n.SetMapIndex(it.Key(), deepCopy(it.Value()))
		}
		return n
	}
	return v
}
