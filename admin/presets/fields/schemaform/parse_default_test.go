package schemaform

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

// The form is what the schema exports as default: a class, a record — its
// fields, their defaults, its metadata the form's —; a type of a list of the
// records of a class, or of plain values, the type's metadata the form's.
func TestParseDefault(t *testing.T) {
	list, err := Parse(`class Link {
	label str
	href? str
	n int = 3
}

[layout={name: "table", columns: [#label, #href]}]
export type default []Link`)
	if err != nil {
		t.Fatal(err)
	}
	if !list.Slice || list.Item != nil || list.Layout != LayoutTable || len(list.Fields) != 3 {
		t.Fatalf("a list of records: %+v", list)
	}
	if cols := fmt.Sprint(list.LayoutConfig["columns"]); cols != "[label href]" {
		t.Errorf("the columns: %v", list.LayoutConfig)
	}
	if f := list.Fields[1]; f.Name != "href" || !f.Nullable || fmt.Sprint(list.Fields[2].Default) != "3" {
		t.Errorf("the fields: %+v %+v", f, list.Fields[2])
	}
	values := url.Values{"V[0].label": {"a"}, "V[0].n": {"1"}, "V[1].label": {"b"}, "V[1].n": {"2"}}
	if got := string(mustJSON(t, list.Decode(values, "V"))); got != `[{"label":"a","href":null,"n":1},{"label":"b","href":null,"n":2}]` {
		t.Errorf("decoded: %s", got)
	}

	one, err := Parse("export class default { a str }")
	if err != nil || one.Slice || len(one.Fields) != 1 {
		t.Errorf("a record: %+v %v", one, err)
	}
	values2, err := Parse(`[layout="form"] export type default []str`)
	if err != nil || !values2.Slice || values2.Item == nil || values2.Item.Type != "str" {
		t.Errorf("a list of values: %+v %v", values2, err)
	}
	ro, err := Parse("export class default { [readonly] path str; content text }")
	if err != nil || !ro.Fields[0].ReadOnly || ro.Fields[1].ReadOnly {
		t.Errorf("[readonly]: %+v %v", ro, err)
	}
}

// What is not a default export of a class or of a list type is refused,
// saying how to write it: the short forms of before, an interface, a class
// by another name, a record typed by an interface, two defaults.
func TestParseRefuses(t *testing.T) {
	for src, want := range map[string]string{
		`{a str}`:                   "export class default",
		`[]{a str}`:                 "export class default",
		`[]str`:                     "export class default",
		`interface Form {a str}`:    "export class default",
		`class Form {a str}`:        "export class default",
		`export class Form {a str}`: "export class default",
		`export class default { sub interface {x str} }`:            "a record is a class",
		"export class default { a str }\nexport type default []str": "2 times",
		`export default = 1`: "export class default",
	} {
		if _, err := Parse(src); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: err = %v, want %q", src, err, want)
		}
	}
}

// A form's fields in the order its class declares them (gad's `@flat`): a
// parent's (`*Parent`) and a mixin's (`use M`) where they are written.
func TestParseClassFlatOrder(t *testing.T) {
	s, err := Parse(`class Contact { name str; email str }
mixin Notes { notes? str }
export class default {
	zip str
	*Contact
	city str
	use Notes
	state str
}`)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range s.Fields {
		got = append(got, f.Name)
	}
	if want := "zip name email city notes state"; strings.Join(got, " ") != want {
		t.Errorf("the order: %q, want %q", strings.Join(got, " "), want)
	}
	if s.Fields[1].Owner != "Contact" || s.Fields[0].Owner != "" {
		t.Errorf("the owners: %q %q", s.Fields[0].Owner, s.Fields[1].Owner)
	}
}
