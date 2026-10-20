// Command gen compiles contracts.gad to contracts.gadb (the module of
// bytecode the package embeds). Run by `go generate` in the package's
// directory.
package main

import (
	"fmt"
	"os"

	"github.com/go-rvq/rvq/admin/presets/fields/schemaform/contracts"
)

func main() {
	src, err := os.ReadFile("contracts.gad")
	if err == nil {
		var b []byte
		if b, err = contracts.Compile(src); err == nil {
			err = os.WriteFile("contracts.gadb", b, 0o644)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "gen:", err)
		os.Exit(1)
	}
}
