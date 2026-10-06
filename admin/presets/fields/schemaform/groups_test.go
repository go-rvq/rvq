package schemaform

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

// rowsOf are the fields of s by their rows: "a|b+c" — b and c a group.
func rowsOf(s *Schema) string {
	var out []string
	for _, row := range Rows(s.Fields) {
		var cols []string
		for _, f := range row {
			c := f.Name
			if f.Row != 0 {
				c += fmt.Sprintf("@%g", f.Percent)
			}
			cols = append(cols, c)
		}
		out = append(out, strings.Join(cols, "+"))
	}
	return strings.Join(out, "|")
}

// A group — `{ … }` in the body, the field `$N` — is no field of the form:
// its fields are, in its place, the columns of one row, each at its share
// (`[width=N]` of 5; the ones with none share what is left).
func TestGroupsAreRows(t *testing.T) {
	s, err := Parse(`{ a str; { [width=2] b str; c str; d? int }; e str; { f str; g str } }`)
	if err != nil {
		t.Fatal(err)
	}
	if got := rowsOf(s); got != "a|b@40+c@30+d@30|e|f@50+g@50" {
		t.Errorf("rows: %s", got)
	}
	// a class's too, a parent's rows after its own numbering
	cls := runClass(t, New(), `
class P { { x str; y str } }
class Form { *P; { a str; [width=4] b str }; c str }
return Form`, nil)
	if got := rowsOf(cls); got != "x@50+y@50|a@20+b@80|c" {
		t.Errorf("class rows: %s", got)
	}
	if cls.Fields[0].Row == cls.Fields[2].Row {
		t.Error("the parent's group and the class's are one row")
	}
}

// The widths are whole units of 5, and a row is whole: refused otherwise,
// saying which group and fields.
func TestGroupWidthsRefused(t *testing.T) {
	for src, want := range map[string]string{
		`{ { [width=0] a str; b str } }`:           "from 1 to 5",
		`{ { [width=6] a str } }`:                  "from 1 to 5",
		`{ { [width=3] a str; [width=3] b str } }`: "more than 5",
		`{ { [width=2] a str; [width=2] b str } }`: "not 5",
		`{ { [width=5] a str; b str } }`:           "none left",
		`{ { a str; { b str } } }`:                 "a group inside",
	} {
		if _, err := Parse(src); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v, want %q", src, err, want)
		}
	}
}

// What a group holds is its record's: posted and read back by the record's
// keys, with no name of the group.
func TestGroupValuesAreTheRecords(t *testing.T) {
	s, err := Parse(`{ a str; { b str; c? int } }`)
	if err != nil {
		t.Fatal(err)
	}
	v := s.Decode(url.Values{"Value.a": {"1"}, "Value.b": {"2"}, "Value.c": {"3"}}, "Value")
	got, _ := json.Marshal(v)
	if string(got) != `{"a":"1","b":"2","c":3}` {
		t.Errorf("decoded: %s", got)
	}
}

// A group is drawn as a row of columns, each field bound to the record's
// key; the detail too.
func TestGroupDrawn(t *testing.T) {
	got := render(t, New(), `{ a str; { [width=2] b str; c str } }`)
	for _, want := range []string{
		`data-schemaform-row`,
		`data-schemaform-column='b'`, `--sf-col: 40%`,
		`data-schemaform-column='c'`, `--sf-col: 60%`,
		`v-model='form["Value"].b'`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("no %s in:\n%.3000s", want, got)
		}
	}
	if strings.Contains(got, "$1") {
		t.Error("the group's name is drawn")
	}
}

// `[tabs] { … }` is a group of tabs, `[steps] { … }` of steps (`[steps,
// vertical]` down the page): each field a tab, a step; no widths to share.
func TestGroupTabsAndSteps(t *testing.T) {
	s, err := Parse(`{
	[tabs] { [label="Person", icon="mdi-account"] person: { name str }; [width=9] other? str }
	[steps, vertical] { a str; b str }
	[steps] { c str }
	d str
}`)
	if err != nil {
		t.Fatal(err)
	}
	f := s.Fields
	if g := f[0].Group; g == nil || g.Kind != GroupTabs || f[1].Group != g {
		t.Errorf("tabs: %+v", f[0].Group)
	}
	if g := f[2].Group; g == nil || g.Kind != GroupSteps || !g.Vertical {
		t.Errorf("vertical steps: %+v", g)
	}
	if g := f[4].Group; g == nil || g.Kind != GroupSteps || g.Vertical {
		t.Errorf("steps: %+v", g)
	}
	if f[5].Group != nil || f[5].Row != 0 {
		t.Errorf("d is of no group: %+v", f[5])
	}

	got := render(t, New(), `{ [tabs] { [label="Person", icon="mdi-account"] person: { name str }; other? str }; [steps] { a str; b str } }`)
	for _, want := range []string{
		`data-schemaform-tabs='$1'`, `data-schemaform-tab='person'`, `Person`, `mdi-account`,
		`v-show='locals.tab === 1'`, `data-schemaform-steps='$2'`, `2. B`, `locals.tab++`,
		`v-model='form["Value"].other'`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("no %s in:\n%.4000s", want, got)
		}
	}
}
