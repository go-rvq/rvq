package contracts

import (
	"bytes"
	"testing"
)

// contracts.gadb is contracts.gad compiled: `go generate` when not.
func TestCompiledUpToDate(t *testing.T) {
	b, err := Compile(Source)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, Compiled) {
		t.Fatal("contracts.gadb is not contracts.gad compiled: run `go generate ./admin/presets/fields/schemaform/contracts`")
	}
}
