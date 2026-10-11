package schemaform

import (
	"net/url"
	"testing"
)

// A schema whose form is a class: its fields, their defaults, its
// metadata the form's; `[list]`, a list of its records — the same as an
// interface's `[]{…}`.
func TestParseClass(t *testing.T) {
	iface, err := Parse(`[layout={name: "table", columns: [#label, #href]}]
interface Form []{label str; href? str; n int}`)
	if err != nil {
		t.Fatal(err)
	}
	cls, err := Parse(`[list, layout={name: "table", columns: [#label, #href]}]
class Form {
	label str
	href? str
	n int = 3
}`)
	if err != nil {
		t.Fatal(err)
	}
	if !cls.Slice || cls.Layout != iface.Layout || len(cls.Fields) != len(iface.Fields) {
		t.Fatalf("class %+v, interface %+v", cls, iface)
	}
	for i, f := range cls.Fields {
		if f.Name != iface.Fields[i].Name || f.Type != iface.Fields[i].Type || f.Nullable != iface.Fields[i].Nullable {
			t.Errorf("field %d: %+v, interface's %+v", i, f, iface.Fields[i])
		}
	}
	if _, ok := cls.Meta[MetaList]; ok {
		t.Error("[list] left in the form's metadata")
	}
	values := url.Values{"V[0].label": {"a"}, "V[0].n": {"1"}, "V[1].label": {"b"}, "V[1].n": {"2"}}
	ci, ii := cls.Decode(values, "V"), iface.Decode(values, "V")
	if string(mustJSON(t, ci)) != string(mustJSON(t, ii)) {
		t.Errorf("decoded: class %s, interface %s", mustJSON(t, ci), mustJSON(t, ii))
	}
	// a record, its own metadata; export written as a module's
	one, err := Parse("export class Form { a str }")
	if err != nil || one.Slice || len(one.Fields) != 1 {
		t.Errorf("a record: %+v %v", one, err)
	}
}
