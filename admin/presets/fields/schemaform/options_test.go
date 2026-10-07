package schemaform

import (
	"net/url"
	"strings"
	"testing"

	"github.com/gad-lang/gad"
	"github.com/go-rvq/rvq/web"
)

// `[options=…]` gives a field the values it may hold, with their labels, in
// the order written: a key-value array, an array of pairs [value, label], or of
// values alone.
func TestOptions(t *testing.T) {
	want := []EnumItem{{Name: "opt1", Label: "option 1"}, {Name: "opt2", Label: "option 2"}}
	for _, src := range []string{
		`interface { [options=(;opt1="option 1", opt2="option 2")]; value str }`,
		`interface { [options=[["opt1", "option 1"], ["opt2", "option 2"]]]; value str }`,
	} {
		s, err := Parse(src)
		if err != nil {
			t.Fatalf("%s: %v", src, err)
		}
		f := s.Fields[0]
		if f.Type != DefaultType || f.Enum == nil || !f.Enum.Options {
			t.Fatalf("%s: field = %+v", src, f)
		}
		if strings.Join(f.Enum.Names, ",") != "opt1,opt2" || len(f.Enum.Items) != 2 ||
			f.Enum.Items[0] != want[0] || f.Enum.Items[1] != want[1] {
			t.Errorf("%s: enum = %+v", src, f.Enum)
		}
	}

	// values alone are their own labels; a value of another type, as text
	s, err := Parse(`interface { [options=["a", "b"]]; v str; [options=[[1, "one"], [2, "two"]]]; n int }`)
	if err != nil {
		t.Fatal(err)
	}
	items, err := New().itemsOf(&web.EventContext{}, "v", s.Fields[0])
	if err != nil || len(items) != 2 || items[0] != (EnumItem{Name: "a", Label: "a"}) {
		t.Errorf("values alone: %+v %v", items, err)
	}
	if names := strings.Join(s.Fields[1].Enum.Names, ","); names != "1,2" || s.Fields[1].Type != "int" {
		t.Errorf("int options: %s %s", names, s.Fields[1].Type)
	}
}

// The field is a select of its options — whatever its type, a registered one
// included —, each value shown by its label.
func TestOptionsDrawASelect(t *testing.T) {
	got := render(t, New(), `interface { [options=(;opt1="option 1", opt2="option 2")]; value str; [options=[[1, "one"]]]; n int }`)
	for _, want := range []string{`v-select`, `option 1`, `option 2`, `opt1`, `v-model='form["Value"].value'`, `one`} {
		if !strings.Contains(got, want) {
			t.Errorf("o select não traz %s:\n%s", want, got)
		}
	}
	if strings.Contains(got, "v-text-field") {
		t.Errorf("o int com options não virou select:\n%s", got)
	}
	// a list of them: one autocomplete of the options, several chosen as
	// chips, each closable — no list of selects —; its [max=N] disables the
	// ones not chosen when it is full
	got = render(t, New(), `interface { [options=(;a="A", b="B"), max=1]; tags []str }`)
	for _, want := range []string{"v-autocomplete", `"title":"B"`, `:multiple='true'`, `:chips='true'`, `:closable-chips='true'`,
		`v-model='form["Value"].tags'`, `.length >= 1`} {
		if !strings.Contains(got, want) {
			t.Errorf("a lista não virou um autocomplete (%s):\n%s", want, got)
		}
	}
	if strings.Contains(got, "v-select") || strings.Contains(got, "vx-array-sorter") {
		t.Errorf("a lista ainda é uma lista de selects:\n%s", got)
	}
}

// Saving: only an option is accepted; an int keeps being an int.
func TestOptionsDecodeForm(t *testing.T) {
	s, err := Parse(`interface { [options=(;opt1="option 1", opt2="option 2")]; value str; [options=[[1, "one"], [2, "two"]]]; n? int }`)
	if err != nil {
		t.Fatal(err)
	}
	b := New()
	value, err := b.DecodeForm(&web.EventContext{}, s, url.Values{"V.value": {"opt2"}, "V.n": {"2"}}, "V")
	if err != nil {
		t.Fatalf("uma opção foi recusada: %v", err)
	}
	if rec, ok := value.(Record); !ok || rec[0].Value != "opt2" || rec[1].Value != int64(2) {
		t.Errorf("value = %#v", value)
	}
	_, err = b.DecodeForm(&web.EventContext{}, s, url.Values{"V.value": {"option 1"}}, "V")
	if err == nil || !strings.Contains(err.Error(), `"option 1"`) || !strings.Contains(err.Error(), "opt1, opt2") {
		t.Errorf("um rótulo passou por valor: %v", err)
	}
}

// What cannot be options is said, with the field.
func TestOptionsErrors(t *testing.T) {
	for src, want := range map[string]string{
		`interface { [options=3]; v str }`:                        "want a key-value array",
		`interface { [options=[["a"]]]; v str }`:                  "a pair is [value, label]",
		`interface { [options=["a", "a"]]; v str }`:               `the value "a" twice`,
		`interface { [options=[]]; v str }`:                       "no option",
		`interface { [options=(;a="A")]; v interface { x str } }`: "a record has no options",
	} {
		_, err := Parse(src)
		if err == nil || !strings.Contains(err.Error(), want) || !strings.Contains(err.Error(), "v:") {
			t.Errorf("%s: err = %v, want %q", src, err, want)
		}
	}
}

// A class's field takes `[options=…]` as an interface's does — the layouts'
// options are classes —, its default kept; an int's items are numbers, so its
// default (or what was saved) shows chosen.
func TestOptionsOfAClass(t *testing.T) {
	b := New()
	vm, ret, err := b.runProgram(`
class Banner {
    [label="Alignment", options=(;left="On the left", center="Centered")]
    align str = "center"
    [options=[[2, "Two"], [3, "Three"]]]
    columns int = 3
}
return Banner
`)
	if err != nil {
		t.Fatal(err)
	}
	s, err := b.Class(vm, ret.(*gad.Class))
	if err != nil {
		t.Fatal(err)
	}
	align := s.FieldAt("align")
	if align.Enum == nil || strings.Join(align.Enum.Names, ",") != "left,center" ||
		align.Enum.Items[1].Label != "Centered" || align.Default != "center" || align.Meta["label"] != "Alignment" {
		t.Errorf("align: %+v %+v", align, align.Enum)
	}
	columns := s.FieldAt("columns")
	if columns.Enum == nil || columns.Default != int64(3) {
		t.Fatalf("columns: %+v", columns)
	}
	if v := itemValue(columns, "3"); v != int64(3) {
		t.Errorf("an int's item: %#v", v)
	}
	if v := itemValue(align, "left"); v != "left" {
		t.Errorf("a str's item: %#v", v)
	}

	got := render(t, b, `interface { [options=[[2, "Two"], [3, "Three"]]]; columns int }`)
	if !strings.Contains(got, `"value":3`) {
		t.Errorf("the int's items are not numbers:\n%s", got)
	}

	// an error says the field
	_, ret, err = b.runProgram("class C { [options=3] x str }\nreturn C")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Class(vm, ret.(*gad.Class)); err == nil || !strings.Contains(err.Error(), "x:") {
		t.Errorf("err = %v", err)
	}
}
