package schemaform

import (
	"fmt"
	"strings"
	"testing"
)

// The schema PAC carries: a list of records, one field of it untyped.
func TestParseSliceOfRecords(t *testing.T) {
	s, err := Parse("class Item { label str; icon str; href }\nexport type default []Item")
	if err != nil {
		t.Fatal(err)
	}

	if !s.Slice {
		t.Error("[]{…} descreve uma lista")
	}
	if got, want := len(s.Fields), 3; got != want {
		t.Fatalf("fields = %d, want %d", got, want)
	}

	for i, want := range []struct {
		name, typ string
	}{
		{"label", "str"},
		{"icon", "str"},
		{"href", DefaultType}, // untyped is plain text
	} {
		if got := s.Fields[i]; got.Name != want.name || got.Type != want.typ {
			t.Errorf("field %d = %q %q, want %q %q", i, got.Name, got.Type, want.name, want.typ)
		}
	}
}

// A form, not a list.
func TestParseRecord(t *testing.T) {
	s, err := Parse("export class default { title str; count int }")
	if err != nil {
		t.Fatal(err)
	}
	if s.Slice {
		t.Error("{…} descreve um registro, não uma lista")
	}
	if got, want := s.Fields[1].Type, "int"; got != want {
		t.Errorf("type = %q, want %q", got, want)
	}
}

// A field typed by a class is a form again; by a list of a class, a list of
// its records.
func TestParseNestedInterface(t *testing.T) {
	cases := map[string]struct {
		src   string
		slice bool
	}{
		"record": {"class Sub { x str }\nexport class default { a str; sub Sub }", false},
		"list":   {"class Sub { x str }\nexport class default { a str; sub []Sub }", true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s, err := Parse(c.src)
			if err != nil {
				t.Fatal(err)
			}
			sub := s.Fields[1]
			if sub.Name != "sub" {
				t.Fatalf("name = %q", sub.Name)
			}
			if sub.Type != FormType {
				t.Errorf("type = %q, want %q", sub.Type, FormType)
			}
			if sub.Schema == nil {
				t.Fatal("o field não trouxe o schema dele")
			}
			if sub.Schema.Slice != c.slice {
				t.Errorf("slice = %v, want %v", sub.Schema.Slice, c.slice)
			}
			if got, want := sub.Schema.Fields[0].Name, "x"; got != want {
				t.Errorf("nested field = %q, want %q", got, want)
			}
		})
	}
}

// The `?` after a name marks the field nullable.
func TestParseNullable(t *testing.T) {
	s, err := Parse("export class default { a str; b? str }")
	if err != nil {
		t.Fatal(err)
	}
	if s.Fields[0].Nullable {
		t.Error("a não é nullable")
	}
	if !s.Fields[1].Nullable {
		t.Error("b? é nullable")
	}
}

func TestParseRejectsWhatIsNotASchema(t *testing.T) {
	for name, src := range map[string]string{
		"empty":      "",
		"not a form": "42",
		"broken":     "export class default {a str",
		// the form is exported once
		"two defaults": "export class default { a str }\nexport class default { b str }",
		// written in parts, the form is the default — and there is none
		"no default": "class User { name }\nclass Owner { user User }",
		// the short forms of before, an interface
		"short form":   "{a str}",
		"another name": "interface Schema { a str }",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse(src); err == nil {
				t.Error("devia recusar")
			} else if !strings.HasPrefix(err.Error(), "schemaform:") {
				t.Errorf("erro sem contexto: %v", err)
			}
		})
	}
}

// An enum declared beside the interface is a closed set of values, and a field
// typed with it carries it.
func TestParseEnum(t *testing.T) {
	s, err := Parse("enum Perm { Read, Write }\nexport class default { perm Perm; other str }")
	if err != nil {
		t.Fatal(err)
	}

	e, ok := s.Enums["Perm"]
	if !ok {
		t.Fatalf("o enum não foi declarado: %v", s.Enums)
	}
	if got, want := strings.Join(e.Names, ","), "Read,Write"; got != want {
		t.Errorf("valores = %q, want %q", got, want)
	}

	if s.Fields[0].Enum != e {
		t.Error("o field não recebeu o enum do type dele")
	}
	if s.Fields[1].Enum != nil {
		t.Error("um field de outro type não tem enum")
	}
}

// An enum declared once serves the forms inside the form too.
func TestParseEnumReachesNestedForms(t *testing.T) {
	s, err := Parse("enum Perm { Read, Write }\nclass Sub { perm Perm }\nexport class default { sub Sub }")
	if err != nil {
		t.Fatal(err)
	}
	if s.Fields[0].Schema.Fields[0].Enum == nil {
		t.Error("o form aninhado não enxergou o enum")
	}
}

// Empty is only allowed where the schema said it is: the `?` after the name.
func TestFieldRequired(t *testing.T) {
	s, err := Parse("export class default { btnColor? color; label str }")
	if err != nil {
		t.Fatal(err)
	}
	if s.Fields[0].Required() {
		t.Error("btnColor? é opcional")
	}
	if !s.Fields[1].Required() {
		t.Error("um field sem ? é obrigatório")
	}
}

// A list of PLAIN VALUES — `export type default []str`, the shape PAC's service_area uses — is a
// list whose item is one value, not a record: no fields, one Item.
func TestParseSliceOfValues(t *testing.T) {
	s, err := Parse("export type default []str")
	if err != nil {
		t.Fatal(err)
	}
	if !s.Slice {
		t.Error("[]str descreve uma lista")
	}
	if len(s.Fields) != 0 {
		t.Errorf("fields = %d, want 0: o item é o valor", len(s.Fields))
	}
	if s.Item == nil {
		t.Fatal("Item = nil, want o valor da lista")
	}
	if got, want := s.Item.Type, "str"; got != want {
		t.Errorf("item type = %q, want %q", got, want)
	}

	// any registered type, not only str
	for _, typ := range []string{"int", "color", "date"} {
		s, err := Parse("export type default []" + typ)
		if err != nil {
			t.Fatalf("[]%s: %v", typ, err)
		}
		if s.Item == nil || s.Item.Type != typ {
			t.Errorf("[]%s: item = %+v", typ, s.Item)
		}
	}
}

// `export type default [][]str` is a list of lists: one schema per `[]`, down to the value.
func TestParseSliceOfSlices(t *testing.T) {
	s, err := Parse("export type default [][]str")
	if err != nil {
		t.Fatal(err)
	}
	if s.Item == nil || s.Item.Type != FormType || s.Item.Schema == nil {
		t.Fatalf("item = %+v, want a lista de dentro", s.Item)
	}
	inner := s.Item.Schema
	if !inner.Slice || inner.Item == nil || inner.Item.Type != "str" {
		t.Errorf("lista de dentro = %+v", inner)
	}
}

// A field of the form may be a list of values too.
func TestParseFieldSliceOfValues(t *testing.T) {
	s, err := Parse("export class default { title str; tags []str }")
	if err != nil {
		t.Fatal(err)
	}
	tags := s.Fields[1]
	if tags.Type != FormType || tags.Schema == nil {
		t.Fatalf("tags = %+v, want um form de um nível abaixo", tags)
	}
	if !tags.Schema.Slice || tags.Schema.Item == nil || tags.Schema.Item.Type != "str" {
		t.Errorf("tags.Schema = %+v", tags.Schema)
	}
}

// A list holds ONE type: several enveloped types describe no single input.
func TestParseSliceOfSeveralTypesRejected(t *testing.T) {
	if _, err := Parse("export type default []<int|str>"); err == nil {
		t.Fatal("uma lista de vários types não descreve um input")
	} else if !strings.Contains(err.Error(), "ONE type") {
		t.Errorf("err = %v", err)
	}
}

// FieldAt walks the same path the form asks its words by.
func TestFieldAt(t *testing.T) {
	s, err := Parse("class Sub { x int }\nclass Links { href str }\nexport class default { title str; sub Sub; links []Links; tags []str }")
	if err != nil {
		t.Fatal(err)
	}

	for path, want := range map[string]string{
		"title":      "str",
		"sub.x":      "int",
		"links.href": "str", // a list adds no level
		"tags":       "str", // a list of values: the path of the list is its item
	} {
		f := s.FieldAt(path)
		if f == nil {
			t.Errorf("%s: nil", path)
			continue
		}
		if f.Type != want {
			t.Errorf("%s: type = %q, want %q", path, f.Type, want)
		}
	}

	for _, path := range []string{"nope", "sub.nope", "title.x"} {
		if f := s.FieldAt(path); f != nil {
			t.Errorf("%s = %+v, want nil", path, f)
		}
	}
}

// A schema may be written in parts: a named interface is a declaration another
// may use by ITS NAME, and the one named Form is the form.
func TestParseNamedInterfaces(t *testing.T) {
	s, err := Parse("class User { name; id }\nexport class default { owner User; creator User }")
	if err != nil {
		t.Fatal(err)
	}

	if len(s.Fields) != 2 {
		t.Fatalf("o form é o último interface: %+v", s.Fields)
	}
	for _, f := range s.Fields {
		if f.Type != FormType || f.Schema == nil {
			t.Errorf("%s = %q, want um form", f.Name, f.Type)
			continue
		}
		if got, want := len(f.Schema.Fields), 2; got != want {
			t.Errorf("%s: %d fields, want %d", f.Name, got, want)
		}
		if f.Schema.Fields[0].Name != "name" || f.Schema.Fields[1].Name != "id" {
			t.Errorf("%s: %+v", f.Name, f.Schema.Fields)
		}
	}

	// each occurrence gets its own form: editing the words of one is not
	// editing the other's
	if s.Fields[0].Schema == s.Fields[1].Schema {
		t.Error("os dois fields dividem o mesmo schema")
	}
}

// A declared class serves a list — a list of its records — and a nested
// field the same way.
func TestParseNamedInterfaceInAList(t *testing.T) {
	s, err := Parse("class Link { label str; href }\nexport class default { links []Link; main Link }")
	if err != nil {
		t.Fatal(err)
	}

	links := s.Fields[0]
	if links.Schema == nil || !links.Schema.Slice || links.Schema.Item != nil {
		t.Fatalf("links = %+v", links)
	}
	if got := links.Schema.Fields[0].Name; got != "label" {
		t.Errorf("fields[0] = %q", got)
	}
	if s.Fields[1].Type != FormType {
		t.Errorf("main = %q, want um form", s.Fields[1].Type)
	}
}

// Enums declared beside the schema reach the interfaces it names.
func TestParseNamedInterfaceWithEnum(t *testing.T) {
	s, err := Parse("enum Perm { Read, Write }\nclass User { perm Perm }\nexport class default { owner User }")
	if err != nil {
		t.Fatal(err)
	}
	perm := s.Fields[0].Schema.Fields[0]
	if perm.Enum == nil || len(perm.Enum.Names) != 2 {
		t.Errorf("perm = %+v", perm)
	}
}

// An interface that contains itself describes a form without end: it is
// refused where the schema is read — one that names itself, by the schema
// ("contains itself"); two that name each other, by the gad compiler (an
// initialization cycle).
func TestParseNamedInterfaceCycle(t *testing.T) {
	for name, c := range map[string]struct{ src, want string }{
		"direto":   {"class Node { child Node }\nexport class default { root Node }", `the class "Node" contains itself`},
		"indireto": {"class A { b B }\nclass B { a A }\nexport class default { a A }", `contains itself`},
		"na raiz":  {"export class default { self default }", `contains itself`},
	} {
		_, err := Parse(c.src)
		if err == nil {
			t.Errorf("%s: um ciclo passou", name)
			continue
		}
		if !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err = %v, want %q", name, err, c.want)
		}
	}
}

// A type name nothing declared is not a type: the schema is a gad program, and
// the gad compiler refuses it naming the reference — where the schema is read,
// not later, in the form.
func TestParseUnknownTypeNameIsRefused(t *testing.T) {
	_, err := Parse("export class default { owner User }")
	if err == nil {
		t.Fatal("um type que ninguém declarou passou")
	}
	if !strings.Contains(err.Error(), `unresolved reference "User"`) {
		t.Errorf("err = %v", err)
	}
}

// The form is the interface NAMED Form, wherever it is among the declarations —
// which may come in any order, as anywhere in gad.
func TestParseFormNameWinsOverOrder(t *testing.T) {
	for name, src := range map[string]string{
		"o form primeiro":   "export class default { owner User }\nclass User { name; id }",
		"o form por último": "class User { name; id }\nexport class default { owner User }",
		"entre outras":      "class User { name; id }\nexport class default { owner User }\nclass Other { x str }",
		"na mesma linha":    "class User { name; id }\nexport class default { owner User }",
	} {
		s, err := Parse(src)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if len(s.Fields) != 1 || s.Fields[0].Name != "owner" {
			t.Errorf("%s: o form não é o interface Form: %+v", name, s.Fields)
			continue
		}
		if s.Fields[0].Schema == nil || s.Fields[0].Schema.Fields[0].Name != "name" {
			t.Errorf("%s: owner não virou o form de User: %+v", name, s.Fields[0])
		}
	}
}

// A schema written in ONE piece needs no name for anything, and goes on being
// written without one.
func TestParseAnonymousStillIsTheForm(t *testing.T) {
	for _, src := range []string{"export class default { a str }", "class Item { a str }\nexport type default []Item", "export type default []str", "export class default { a str }"} {
		if _, err := Parse(src); err != nil {
			t.Errorf("%q: %v", src, err)
		}
	}
}

// A clone may be changed without touching the schema Parse cached.
func TestSchemaClone(t *testing.T) {
	src := `class Item { a str; b str }
[layout={name: "table", columns: [#a, #b]}]
export type default []Item`
	s, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	c := s.Clone()
	c.Layout = LayoutGrid
	c.LayoutConfig["columns"] = 2
	c.Fields[0].Name = "changed"

	again, _ := Parse(src)
	if again.Fields[0].Name != "a" || again.Layout != LayoutTable || len(again.TableColumns()) != 2 {
		t.Errorf("changing the clone changed the cached schema: %s %v", again.Layout, again.LayoutConfig)
	}
}

// The value a form starts from has every record and list the schema has, so
// what the form binds through is there.
func TestFilled(t *testing.T) {
	s, err := Parse(`class Words { label? str }
class Fields { icon? Words; link? Words }
class Layout { name? str; tags? []str }
export class default { layout? Layout; fields? Fields }`)
	if err != nil {
		t.Fatal(err)
	}
	got := s.Filled(map[string]any{"fields": map[string]any{"icon": map[string]any{"label": "Ícone"}}})
	want := map[string]any{
		"layout": map[string]any{"tags": []any{}},
		"fields": map[string]any{
			"icon": map[string]any{"label": "Ícone"},
			"link": map[string]any{},
		},
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("filled:\n %v\nwant\n %v", got, want)
	}
	if list := (&Schema{Slice: true}).Filled(nil); list == nil {
		t.Error("a list is filled in as an empty list")
	}
}
