package contracts

import (
	"bytes"
	"os"
	"testing"

	"github.com/gad-lang/gad"
	"github.com/gad-lang/gad/encoder"
	"github.com/gad-lang/gad/stdlib/helper"
)

// schemaform.gadb is contracts.gad compiled (as `gad gomodule` compiles it):
// `go generate` when not.
func TestCompiledUpToDate(t *testing.T) {
	src, err := os.ReadFile("contracts.gad")
	if err != nil {
		t.Fatal(err)
	}
	b, err := encoder.CompileModuleWith(ModuleName, src, gad.SourceKindGad, helper.NewModuleMap(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, Bytecode) {
		t.Fatal("schemaform.gadb is not contracts.gad compiled: run `go generate ./admin/presets/fields/schemaform/contracts`")
	}
}
