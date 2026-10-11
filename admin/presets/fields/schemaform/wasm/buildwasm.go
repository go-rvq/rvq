//go:build !(js && wasm)

package wasm

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Build builds the command pkg (its Serve) for the browser into dir — the
// build formscore, small: <name>.wasm.gz, and Go's wasm_exec.js —, run in
// the module of workdir ("" the current directory).
func Build(workdir, pkg, dir, name string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	wasm := filepath.Join(abs, name+".wasm")
	cmd := exec.Command("go", "build", "-tags", "formscore", "-trimpath", "-ldflags=-s -w", "-o", wasm, pkg)
	cmd.Dir, cmd.Stdout, cmd.Stderr = workdir, os.Stdout, os.Stderr
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wasm: build %s: %w", pkg, err)
	}
	b, err := os.ReadFile(wasm)
	if err != nil {
		return err
	}
	var gz bytes.Buffer
	zw, _ := gzip.NewWriterLevel(&gz, gzip.BestCompression)
	if _, err := zw.Write(b); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	if err := os.WriteFile(wasm+".gz", gz.Bytes(), 0o644); err != nil {
		return err
	}
	if err := os.Remove(wasm); err != nil {
		return err
	}
	out, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return fmt.Errorf("wasm: GOROOT: %w", err)
	}
	root := strings.TrimSpace(string(out))
	for _, cand := range []string{filepath.Join(root, "lib", "wasm", "wasm_exec.js"), filepath.Join(root, "misc", "wasm", "wasm_exec.js")} {
		if js, err := os.ReadFile(cand); err == nil {
			return os.WriteFile(filepath.Join(abs, "wasm_exec.js"), js, 0o644)
		}
	}
	return fmt.Errorf("wasm: no wasm_exec.js in %s", root)
}
