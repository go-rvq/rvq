package schemaform

import (
	"testing"

	"github.com/gad-lang/gad"
)

// runClass runs src — the builder's types in scope (Builtins), and extra
// names — and reads the class it returns as a form.
func runClass(t *testing.T, b *Builder, src string, extra map[string]gad.Object) *Schema {
	t.Helper()
	bi := b.Builtins()
	for name, o := range extra {
		bi.Set(name, o)
	}
	res, err := gad.Compile(gad.NewSymbolTable(bi.NameSet), []byte(src), gad.CompileOptions{})
	if err != nil {
		t.Fatal(err)
	}
	vm := gad.NewVM(bi.Build(), res.Bytecode)
	ret, err := vm.Run()
	if err != nil {
		t.Fatal(err)
	}
	s, err := b.Class(vm, ret.(*gad.Class))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// A class that extends another (`*P`) is a form of the parent's fields, then
// its own, in the order declared; each says whose it is (Owner). The
// builder's types (`text`, `date`) type a class field as they do an
// interface's.
func TestClassExtendsAndBuilderTypes(t *testing.T) {
	s := runClass(t, New(), `
class P { pz str; [label="A"] pa? str }
class Form { *P; z str; a str; m text; d? date; b? int }
return Form`, nil)
	want := []struct{ name, typ, owner string }{
		{"pz", "str", "P"}, {"pa", "str", "P"}, {"z", "str", ""}, {"a", "str", ""},
		{"m", "text", ""}, {"d", "date", ""}, {"b", "int", ""},
	}
	if len(s.Fields) != len(want) {
		t.Fatalf("fields: %d", len(s.Fields))
	}
	for i, w := range want {
		f := s.Fields[i]
		if f.Name != w.name || f.Type != w.typ || f.Owner != w.owner {
			t.Errorf("field %d = %s %s (%s), want %s %s (%s)", i, f.Name, f.Type, f.Owner, w.name, w.typ, w.owner)
		}
	}
	if !s.Fields[1].Nullable || s.Fields[1].Meta.String("label") != "A" {
		t.Error("the parent's field lost its ? or its metadata")
	}
}

// A type an application names by a marker (TypeMarker) — a module's member,
// `contact_form.email` — types a field as the builder's types do.
func TestClassTypeMarker(t *testing.T) {
	b := New().Type("contact_form.email", TextComponentFunc)
	s := runClass(t, b, `class Form { email contact_form.email }; return Form`, map[string]gad.Object{
		"contact_form": gad.Dict{"email": TypeMarker("contact_form.email")},
	})
	if f := s.Fields[0]; f.Name != "email" || f.Type != "contact_form.email" {
		t.Errorf("field = %s %s", f.Name, f.Type)
	}
}

// ClassName names a class extended: the fields it gives are of that name,
// not of the class's own (two classes Form, one extending the other).
func TestClassName(t *testing.T) {
	b := New()
	var parent *gad.Class
	b.ClassName(func(c *gad.Class) string {
		if c == parent {
			return "forms.base"
		}
		return ""
	})
	bi := b.Builtins()
	res, _ := gad.Compile(gad.NewSymbolTable(bi.NameSet), []byte(`class Form { a str }; return Form`), gad.CompileOptions{})
	vm := gad.NewVM(bi.Build(), res.Bytecode)
	ret, err := vm.Run()
	if err != nil {
		t.Fatal(err)
	}
	parent = ret.(*gad.Class)
	s := runClass(t, b, `class Form { *P; b str }; return Form`, map[string]gad.Object{"P": parent})
	if s.Fields[0].Name != "a" || s.Fields[0].Owner != "forms.base" || s.Fields[1].Owner != "" {
		t.Errorf("fields: %s %q, %s %q", s.Fields[0].Name, s.Fields[0].Owner, s.Fields[1].Name, s.Fields[1].Owner)
	}
}
