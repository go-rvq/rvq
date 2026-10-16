package schemaform

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/gad-lang/gad"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// readClass is the schema of the class src returns.
func readClass(src string) (*Schema, error) {
	b := New()
	bi := b.Builtins()
	res, err := gad.Compile(gad.NewSymbolTable(bi.NameSet), []byte(src), gad.CompileOptions{})
	if err != nil {
		return nil, err
	}
	vm := gad.NewVM(bi.Build(), res.Bytecode)
	ret, err := vm.Run()
	if err != nil {
		return nil, err
	}
	return b.Class(vm, ret.(*gad.Class))
}

// The `[validation=…]` of a field and of a class: a function or an array of
// them, kept as gad functions; a class's after its parents'; what is no
// function, refused.
func TestValidationMeta(t *testing.T) {
	s, err := readClass(`
		a := func(v) => nil
		b := func(v) => nil
		[validation=a] class Base { x int }
		[validation=[a, b]] class Form {
			*Base
			[validation=a] one str
			[validation=[a, b]] two str
			[validation=b] group class { y int }
		}
		return Form`)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Validations) != 3 {
		t.Errorf("the class's, after its parent's: %d", len(s.Validations))
	}
	n := map[string]int{}
	for _, f := range s.Fields {
		n[f.Name] = len(f.Validations)
	}
	if n["one"] != 1 || n["two"] != 2 || n["group"] != 1 || n["x"] != 0 {
		t.Errorf("the fields': %v", n)
	}
	for src, want := range map[string]string{
		`[validation=1] class Form { a str }; return Form`:     "want a function",
		`class Form { [validation=["x"]] a str }; return Form`: "want a function",
	} {
		if _, err := readClass(src); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v", src, err)
		}
	}
}

// A regular expression is a validation too — a literal, `/…/` —: anchored,
// the whole value matches it.
func TestValidationRegexp(t *testing.T) {
	s, err := readClass("class Form { [validation=[/\\d{5}/, /[A-Z]+/]] zip str }; return Form")
	if err != nil {
		t.Fatal(err)
	}
	vs := s.Fields[0].Validations
	if len(vs) != 2 {
		t.Fatalf("validations: %v", vs)
	}
	re := func(o gad.Object) *regexp.Regexp {
		r, ok := o.(*gad.Regexp)
		if !ok {
			t.Fatalf("not a regexp: %T", o)
		}
		return (*regexp.Regexp)(r)
	}
	if !re(vs[0]).MatchString("12345") || re(vs[0]).MatchString("123456") {
		t.Error("the literal, anchored")
	}
	if !re(vs[1]).MatchString("ABC") || re(vs[1]).MatchString("ABC1") {
		t.Error("the second, anchored")
	}
}

// The `[round=(this, parents, value; places=0) => …]` of a float or a
// decimal: called with the value read — places its `[places=N]` —, what it
// answers kept; of another type, or no function, refused.
func TestRoundMeta(t *testing.T) {
	s, err := readClass(`
		up := func(this, parents, value; places=0) {
			m := 1
			for i := 0; i < places; i++ { m *= 10 }
			return float(int(value * m + 0.999999)) / m
		}
		class Form {
			[places=2, round=up] price decimal
			[round=up] whole float
			[round=func(this, parents, value; places=0) => value + this.extra] sum float
			extra float
			plain float
		}
		return Form`)
	if err != nil {
		t.Fatal(err)
	}
	rec, issues := (&Decoder{Schema: s}).Decode(map[string]any{
		"price": "10.011", "whole": "2.1", "sum": "1", "extra": "0.5", "plain": "1.234"})
	if len(issues) != 0 {
		t.Fatalf("issues: %v", issues)
	}
	got := map[string]any{}
	for _, f := range rec {
		got[f.Name] = f.Value
	}
	for name, want := range map[string]float64{"price": 10.02, "whole": 3, "sum": 1.5, "plain": 1.234} {
		if got[name] != want {
			t.Errorf("%s: %v, want %v", name, got[name], want)
		}
	}
	for src, want := range map[string]string{
		`class Form { [round=1] a float }; return Form`:                                 "a function",
		`class Form { [round=func(this, parents, value) => value] a int }; return Form`: "of a float or a decimal",
	} {
		if _, err := readClass(src); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v", src, err)
		}
	}
}

// The `[roundDisplay=(this, parents, value; places=0) => …]` of a float or a
// decimal: its rounding as it is shown — the detail, a field by its path, a
// read-only form; this the record —, the value kept as it is when read.
func TestRoundDisplayMeta(t *testing.T) {
	b := New()
	bi := b.Builtins()
	res, err := gad.Compile(gad.NewSymbolTable(bi.NameSet), []byte(`
		up := func(this, parents, value; places=0) {
			m := 1
			for i := 0; i < places; i++ { m *= 10 }
			return float(int(value * m + 0.999999)) / m + this.bonus
		}
		class Form {
			[places=2, roundDisplay=up] price decimal
			bonus float
		}
		return Form`), gad.CompileOptions{})
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
	// read: as it is
	rec, issues := (&Decoder{Schema: s}).Decode(map[string]any{"price": "10.011", "bonus": "1"})
	if len(issues) != 0 {
		t.Fatalf("issues: %v", issues)
	}
	if rec[0].Value != 10.011 {
		t.Errorf("read: %v", rec[0].Value)
	}
	// shown: rounded, this the record (its bonus)
	value := map[string]any{"price": 10.011, "bonus": 1}
	html := func(comp h.HTMLComponent) string {
		out, err := h.Marshal(comp, context.Background())
		if err != nil {
			t.Fatal(err)
		}
		return string(out)
	}
	if got := html(b.DetailComponent(&web.EventContext{}, s, value)); !strings.Contains(got, "11.02") {
		t.Errorf("detail: no 11.02:\n%s", got)
	}
	show, _ := b.FieldDisplay(&web.EventContext{}, s, value)
	if got := html(show("price")); !strings.Contains(got, "11.02") {
		t.Errorf("by its path: %s", got)
	}
	form := html(b.ComponentFunc(s)(&presets.FieldContext{ToComponentOptions: &presets.ToComponentOptions{},
		Name: "Value", FormKey: "Value", Mode: presets.FieldModeStack{presets.DETAIL}, ValueOverride: value}, &web.EventContext{}))
	if !strings.Contains(form, "11.02") {
		t.Errorf("read-only form: no 11.02")
	}
	// with no record (DisplayText): this empty — the function fails (logged),
	// the value as it is, with its places
	if got := DisplayText(context.Background(), s.Fields[0], 10.011); got != "10.01" {
		t.Errorf("DisplayText: %q", got)
	}
	if _, err := readClass(`class Form { [roundDisplay=1] a float }; return Form`); err == nil {
		t.Error("a roundDisplay of no function")
	}
}
