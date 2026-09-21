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
package schemaform

import (
	"fmt"

	"github.com/gad-lang/gad/parser"
	"github.com/gad-lang/gad/parser/node"
	"github.com/gad-lang/gad/parser/source"
)

const (
	// DefaultType is what an untyped field is: plain text.
	DefaultType = "str"
	// FormType is the type of a field that is itself an interface — a form
	// inside the form, or a list of them.
	FormType = "yaml_form"
)

// Schema is one form: the fields it edits, and whether the value it edits is a
// LIST of them rather than one.
type Schema struct {
	Slice  bool
	Fields []*Field
}

// Field is one entry of a form.
type Field struct {
	Name string
	// Type names the component that renders it — DefaultType when the schema
	// left it out, FormType when the field is an interface of its own.
	Type string
	// Nullable is the `?` after the name: the field accepts nil.
	Nullable bool
	// Schema is the field's own form, set when Type is FormType.
	Schema *Schema
}

// Parse reads a schema — the body of a gad interface, written without the
// keyword — into the form it describes.
//
// The three spellings of a field that is an interface arrive the same way:
//
//	sub interface {x str}     sub: {x str}
//	sub interface[] {x str}   sub: []{x str}
func Parse(src string) (*Schema, error) {
	fs := source.NewFileSet()
	f := fs.AddFileData("schema", -1, []byte("interface "+src))

	file, err := parser.NewParser(f, nil).ParseFile()
	if err != nil {
		return nil, fmt.Errorf("schemaform: %w", err)
	}
	if len(file.Stmts) != 1 {
		return nil, fmt.Errorf("schemaform: a schema is one interface, got %d statements", len(file.Stmts))
	}

	stmt, ok := file.Stmts[0].(*node.ExprStmt)
	if !ok {
		return nil, fmt.Errorf("schemaform: a schema is one interface, got %T", file.Stmts[0])
	}
	iface, ok := stmt.Expr.(*node.InterfaceExpr)
	if !ok {
		return nil, fmt.Errorf("schemaform: a schema is one interface, got %T", stmt.Expr)
	}

	return schemaOf(iface), nil
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
			if nested, ok := m.Name.Type[0].Expr.(*node.InterfaceExpr); ok {
				f.Type, f.Schema = FormType, schemaOf(nested)
			} else {
				f.Type = m.Name.Type[0].String()
			}
		}

		s.Fields = append(s.Fields, f)
	}

	return s
}
