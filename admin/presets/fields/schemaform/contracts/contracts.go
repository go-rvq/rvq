// Package contracts holds the contracts of the field types of the schemaform
// (contracts.gad), compiled once to a module of bytecode (contracts.gadb,
// `go generate`): a schema's program imports it — linked, not compiled again
// — as the module "schemaform" (Module).
package contracts

import (
	_ "embed"

	"github.com/gad-lang/gad"
	"github.com/gad-lang/gad/encoder"
)

//go:generate go run ./gen

// Name is the module the contracts are imported as.
const Name = "schemaform"

// Source is the contracts, as written.
//
//go:embed contracts.gad
var Source []byte

// Compiled is Source compiled (Compile): what Module links.
//
//go:embed contracts.gadb
var Compiled []byte

// Compile is src compiled as the module Name, with gad's builtins.
func Compile(src []byte) ([]byte, error) { return encoder.CompileModule(Name, src) }

// Module is the contracts as a module to import: linked into the program
// that imports it, run by its VM.
func Module() gad.Importable { return encoder.BytecodeModule(Compiled) }
