package schemaform

import (
	"regexp"
	"strings"
	"testing"

	"github.com/gad-lang/gad"
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
