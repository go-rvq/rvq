package schemaform

import (
	"strings"
	"testing"

	"github.com/gad-lang/gad"
)

// The `[validation=…]` of a field and of a class: a function or an array of
// them, kept as gad functions; a class's after its parents'; what is no
// function, refused.
func TestValidationMeta(t *testing.T) {
	read := func(src string) (*Schema, error) {
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
	s, err := read(`
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
	for _, src := range []string{
		`[validation=1] class Form { a str }; return Form`,
		`class Form { [validation=["x"]] a str }; return Form`,
	} {
		if _, err := read(src); err == nil || !strings.Contains(err.Error(), "want a function") {
			t.Errorf("%s: %v", src, err)
		}
	}
}
