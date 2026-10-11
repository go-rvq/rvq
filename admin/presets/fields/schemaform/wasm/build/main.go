// Command build builds a check of forms for the browser (wasm.Build):
//
//	go run github.com/go-rvq/rvq/admin/presets/fields/schemaform/wasm/build -o dist -name forms ./cmd
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-rvq/rvq/admin/presets/fields/schemaform/wasm"
)

func main() {
	out := flag.String("o", "dist", "the directory it is built into")
	name := flag.String("name", "forms", "its name: <name>.wasm.gz")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: build [-o dist] [-name forms] <the command's package>")
		os.Exit(2)
	}
	if err := wasm.Build("", flag.Arg(0), *out, *name); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	info, err := os.Stat(*out + "/" + *name + ".wasm.gz")
	if err == nil {
		fmt.Printf("%s/%s.wasm.gz: %d bytes\n", *out, *name, info.Size())
	}
}
