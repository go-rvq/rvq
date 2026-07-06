package tagparser

import (
	"github.com/gad-lang/gad"
	"github.com/gad-lang/gad/parser"
	"github.com/go-rvq/rvq/admin/presets"
)

func Parse(s string) (vm *gad.VM, kva gad.KeyValueArray, err error) {
	if len(s) == 0 {
		return
	}

	builtins := gad.NewBuiltins()
	builtins.Set("NEW", gad.Uint(presets.NEW))
	builtins.Set("EDIT", gad.Uint(presets.EDIT))
	builtins.Set("DETAIL", gad.Uint(presets.DETAIL))
	builtins.Set("LIST", gad.Uint(presets.LIST))
	builtins.Set("WRITE", gad.Uint(presets.WRITE))
	builtins.Set("FORM", gad.Uint(presets.FORM))
	builtins.Set("ALL", gad.Uint(presets.ALL))
	builtins.Set("FIELD_TYPES", gad.Dict{
		"text":       gad.Str("text"),
		"inlineText": gad.Str("inlineText"),
	})
	staticBuiltins := builtins.Build()
	symbols := gad.NewSymbolTable(staticBuiltins.Builtins().NameSet)

	src := "return (;" + s + ")"

	var bc *gad.Bytecode

	if _, bc, err = gad.Compile(symbols, []byte(src), gad.CompileOptions{
		ScannerOptions: parser.ScannerOptions{
			Mode: parser.ScanCharAsString,
		},
	}); err != nil {
		return
	}

	var ret gad.Object
	vm = gad.NewVM(staticBuiltins, bc)

	if ret, err = vm.Run(); err != nil {
		return
	}

	kva = ret.(gad.KeyValueArray)
	return
}
