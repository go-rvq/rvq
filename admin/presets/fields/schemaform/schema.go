// Package schemaform builds an edit form out of a gad interface.
//
// A value edited as a form is described by a SCHEMA: the gad interface its
// content satisfies, written without the keyword —
//
//	[]{label str; icon str; href}
//
// which is `interface []{label str; icon str; href}`, a list of records with
// three fields. The form the schema describes is what the user edits; the value
// itself is stored as it always was (YAML, for a LocaleMessage).
//
// The schema says, for each field, WHICH component renders it: the field's type
// names a component registered on the Builder. An untyped field is `str` —
// plain text — and a field that is itself an interface is a form again, drawn
// by the same component one level down.
//
// A list needs no record: `[]str` is a list of plain values, each edited by the
// component of its type, and `[][]str` a list of those. A field may be one too
// (`tags []str`).
package schemaform

import (
	"fmt"
	"strings"

	"github.com/gad-lang/gad/parser"
	"github.com/gad-lang/gad/parser/node"
	"github.com/gad-lang/gad/parser/source"
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
	// Enums are the enums the schema declared, by name. A field typed with one
	// carries it in Field.Enum.
	Enums map[string]*Enum
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
	// Enum is the set of values the field may hold, set when its type names an
	// enum the schema declared.
	Enum *Enum
}

// Required reports whether the field may be left empty: only an optional one
// may, and optional is what the `?` after the name says.
func (f *Field) Required() bool { return !f.Nullable }

// Parse reads a schema — the body of a gad interface, written without the
// keyword — into the form it describes.
//
// The three spellings of a field that is an interface arrive the same way:
//
//	sub interface {x str}     sub: {x str}
//	sub interface[] {x str}   sub: []{x str}
func Parse(src string) (*Schema, error) {
	// `[]str` — a list of plain values, not of records. It is a TYPE, not an
	// interface, so it is read where a type is read: as the type of a field of
	// a throw-away interface.
	if isSliceOfType(src) {
		return parseSliceOfType(src)
	}

	fs := source.NewFileSet()
	f := fs.AddFileData("schema", -1, []byte(withKeyword(src)))

	file, err := parser.NewParser(f, nil).ParseFile()
	if err != nil {
		return nil, fmt.Errorf("schemaform: %w", err)
	}

	var (
		iface *node.InterfaceExpr
		enums = map[string]*Enum{}
	)

	for _, stmt := range file.Stmts {
		switch st := stmt.(type) {
		case *node.EnumStmt:
			if e := enumOf(&st.EnumExpr); e != nil {
				enums[e.Name] = e
			}
		case *node.ExprStmt:
			found, ok := st.Expr.(*node.InterfaceExpr)
			if !ok {
				return nil, fmt.Errorf("schemaform: a schema is one interface and its enums, got %T", st.Expr)
			}
			if iface != nil {
				return nil, fmt.Errorf("schemaform: a schema is ONE interface, got another")
			}
			iface = found
		default:
			return nil, fmt.Errorf("schemaform: a schema is one interface and its enums, got %T", stmt)
		}
	}

	if iface == nil {
		return nil, fmt.Errorf("schemaform: a schema is one interface, and there is none")
	}

	s := schemaOf(iface)
	s.Enums = enums
	bindEnums(s, enums)
	return s, nil
}

// withKeyword writes the `interface` a schema may leave out. A schema that is
// only the interface is written bare — `[]{label str}` — and one that declares
// enums beside it writes the keyword, because the bare form has nowhere to put
// them:
//
//	enum Perm { Read, Write }
//	interface { perm Perm }
func withKeyword(src string) string {
	trimmed := strings.TrimSpace(src)
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		return "interface " + trimmed
	}
	return trimmed
}

// sliceItemSrc is what `[]…` holds, and how many `[]` are before it: `[]str`
// holds `str` at depth 1, `[][]str` at depth 2. It is only a slice of a TYPE
// when what it holds is not an interface body, so `[]{…}` — a list of records —
// answers depth 0.
func sliceItemSrc(src string) (item string, depth int) {
	rest := strings.TrimSpace(src)
	for strings.HasPrefix(rest, "[]") {
		rest = strings.TrimSpace(rest[2:])
		depth++
	}
	if depth == 0 || rest == "" || strings.HasPrefix(rest, "{") {
		return "", 0
	}
	return rest, depth
}

// isSliceOfType reports whether the schema is a list of plain values.
func isSliceOfType(src string) bool {
	_, depth := sliceItemSrc(src)
	return depth > 0
}

// parseSliceOfType reads `[]str` (and `[][]str`, a list of lists) into the form
// it describes: a list whose item is one value of that type. The type is read
// where a type is read — as the type of a field — so a list holds whatever a
// field may hold.
func parseSliceOfType(src string) (*Schema, error) {
	trimmed := strings.TrimSpace(src)

	s, err := Parse("interface { item " + trimmed + " }")
	if err != nil {
		return nil, err
	}
	if len(s.Fields) != 1 || s.Fields[0].Schema == nil {
		return nil, fmt.Errorf("schemaform: %q não descreve uma lista de valores: uma lista guarda UM type", trimmed)
	}
	return s.Fields[0].Schema, nil
}

// sliceSchemaOf reads a slice type — `[]str`, `[][]int` — into a list whose item
// is one value of that type, one schema per `[]` so a list of lists nests. A
// slice of SEVERAL types (`[]<int|str>`) describes no single input, so it has no
// form: it comes back nil and is reported as the type it is written as.
func sliceSchemaOf(t *node.SliceTypeExpr) *Schema {
	if len(t.Types) != 1 {
		return nil
	}

	item := &Field{Type: t.Types[0].String()}
	if nested, ok := t.Types[0].Expr.(*node.InterfaceExpr); ok {
		item.Type, item.Schema = FormType, schemaOf(nested)
	}

	s := &Schema{Slice: true, Item: item}
	for i := 1; i < t.Depth; i++ {
		s = &Schema{Slice: true, Item: &Field{Type: FormType, Schema: s}}
	}
	return s
}

// enumOf reads `enum Name { A, B }`. An anonymous enum names nothing a field
// could refer to, so it is left out.
func enumOf(e *node.EnumExpr) *Enum {
	name, _ := e.NameExpr.(*node.IdentExpr)
	if name == nil {
		return nil
	}

	out := &Enum{Name: name.Name}
	for _, f := range e.Fields {
		if f.Name != nil {
			out.Names = append(out.Names, f.Name.Name)
		}
	}
	return out
}

// bindEnums hands each field the enum its type names, walking into the forms
// inside the form — an enum declared once serves the whole schema.
func bindEnums(s *Schema, enums map[string]*Enum) {
	if s.Item != nil {
		if e, ok := enums[s.Item.Type]; ok {
			s.Item.Enum = e
		}
		if s.Item.Schema != nil {
			s.Item.Schema.Enums = enums
			bindEnums(s.Item.Schema, enums)
		}
	}
	for _, f := range s.Fields {
		if e, ok := enums[f.Type]; ok {
			f.Enum = e
		}
		if f.Schema != nil {
			f.Schema.Enums = enums
			bindEnums(f.Schema, enums)
		}
	}
}

// schemaOf reads one interface. Only its FIELDS are a form: a method, a getter,
// a setter or a prop describes behaviour, and there is nothing to edit in it.
func schemaOf(iface *node.InterfaceExpr) *Schema {
	s := &Schema{Slice: iface.ArrayDepth > 0}

	for _, m := range iface.Members {
		if m.Kind != node.IfaceField || m.Name == nil || m.Name.Ident == nil {
			continue
		}

		f := &Field{
			Name:     m.Name.Ident.Name,
			Type:     DefaultType,
			Nullable: m.Name.Nullable,
		}

		if len(m.Name.Type) > 0 {
			switch t := m.Name.Type[0].Expr.(type) {
			case *node.InterfaceExpr:
				f.Type, f.Schema = FormType, schemaOf(t)
			case *node.SliceTypeExpr:
				// `tags []str` — a list of plain values is a form of its own,
				// one level down, like any other nested schema.
				if sub := sliceSchemaOf(t); sub != nil {
					f.Type, f.Schema = FormType, sub
				} else {
					f.Type = m.Name.Type[0].String()
				}
			default:
				f.Type = m.Name.Type[0].String()
			}
		}

		s.Fields = append(s.Fields, f)
	}

	return s
}

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
