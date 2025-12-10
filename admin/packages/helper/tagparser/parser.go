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

	builtints := gad.NewBuiltins()
	builtints.Set("NEW", gad.Uint(presets.NEW))
	builtints.Set("EDIT", gad.Uint(presets.EDIT))
	builtints.Set("DETAIL", gad.Uint(presets.DETAIL))
	builtints.Set("LIST", gad.Uint(presets.LIST))
	builtints.Set("WRITE", gad.Uint(presets.WRITE))
	builtints.Set("ALL", gad.Uint(presets.ALL))

	src := "return (;" + s + ")"

	var bc *gad.Bytecode

	if bc, err = gad.Compile([]byte(src), gad.CompileOptions{
		ScannerOptions: parser.ScannerOptions{
			Mode: parser.ScanCharAsString,
		},
		CompilerOptions: gad.CompilerOptions{
			SymbolTable: gad.NewSymbolTable(builtints),
		},
	}); err != nil {
		return
	}

	var ret gad.Object
	vm = gad.NewVM(bc)

	if ret, err = vm.Run(); err != nil {
		return
	}

	kva = ret.(gad.KeyValueArray)
	return
}
