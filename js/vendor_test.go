package js

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoVendorSubfolders: a module's zip leaves out what is under a subfolder
// of a folder named vendor (golang.org/x/mod/zip, the vendored packages of
// old) — what is embedded from one is missing in a dependent of rvq.
func TestNoVendorSubfolders(t *testing.T) {
	root := ".."
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		switch d.Name() {
		case "node_modules", ".git", ".tmp":
			return filepath.SkipDir
		}
		rel, _ := filepath.Rel(root, p)
		parts := strings.Split(filepath.ToSlash(rel), "/")
		for i, s := range parts[:len(parts)-1] {
			if s == "vendor" {
				t.Errorf("%s: a subfolder of a folder vendor (%s): left out of the module",
					rel, strings.Join(parts[:i+1], "/"))
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
