package gaddocs

import (
	"io/fs"
	"runtime/debug"
	"strings"
	"testing"
)

// The documentation of Gad: its node and the three parts under it, every
// node a document; the page of the node in each language; the copy of the
// module's, its links made for where they went.
func TestDocs(t *testing.T) {
	nodes := Nodes()
	if len(nodes) != 1 || nodes[0].ID != Node {
		t.Fatalf("the top: %+v", nodes)
	}
	var ids []string
	for _, c := range nodes[0].Children {
		ids = append(ids, c.ID)
	}
	if got := strings.Join(ids, ","); got != "gad/gad,gad/gadt,gad/gadx" {
		t.Errorf("the parts: %s", got)
	}
	src := Source()
	if !src.Raw {
		t.Error("not raw")
	}
	for _, f := range []string{"en/gad/README.md", "pt-BR/gad/README.md", "en/gad/gad/tutorial/README.md",
		"en/gad/gadt/README.md", "en/gad/gadx/syntax/README.md"} {
		if _, err := fs.Stat(src.FS, f); err != nil {
			t.Errorf("%s: %v", f, err)
		}
	}
	b, _ := fs.ReadFile(src.FS, "pt-BR/gad/README.md")
	if !strings.Contains(string(b), "Gadx") || !strings.Contains(string(b), "gadx/README.md") {
		t.Errorf("the page in pt-BR: %s", b)
	}
	langs, err := src.Langs()
	if err != nil || strings.Join(langs, ",") != "en,pt-BR" {
		t.Errorf("the languages: %v %v", langs, err)
	}
	// a link of the gadx's documentation to another of it
	b, _ = fs.ReadFile(src.FS, "en/gad/gadx/README.md")
	if !strings.Contains(string(b), "](syntax/README.md)") {
		t.Error("the links of gadx/README.md")
	}
}

// The links to the repository of Gad follow the version of Gad of the program:
// the copy holds RefPlaceholder, FS gives Ref in its place.
func TestRef(t *testing.T) {
	for v, want := range map[string]string{
		"v0.1.0-rc.33.0.20261018132931-08adf6abed4d": "08adf6abed4d",
		"v0.0.0-20261002142357-9e9dae3d3353":         "9e9dae3d3353",
		"v0.1.0":                                     "v0.1.0",
		"v0.1.0-rc.33":                               "v0.1.0-rc.33",
		"":                                           "main",
		"(devel)":                                    "main",
	} {
		if got := refOf(v); got != want {
			t.Errorf("refOf(%q) = %q, want %q", v, got, want)
		}
	}
	// a replace by a directory keeps the version required; by a module, its
	req := "v0.1.0-rc.33.0.20261018132931-08adf6abed4d"
	if v := depVersion(&debug.Module{Version: req, Replace: &debug.Module{Path: "/x", Version: "(devel)"}}); v != req {
		t.Errorf("replaced by a directory: %s", v)
	}
	if v := depVersion(&debug.Module{Version: req, Replace: &debug.Module{Path: "m", Version: "v0.2.0"}}); v != "v0.2.0" {
		t.Errorf("replaced by a module: %s", v)
	}
	src := Source()
	n := 0
	_ = fs.WalkDir(src.FS, "en", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		b, _ := fs.ReadFile(src.FS, p)
		if strings.Contains(string(b), RefPlaceholder) {
			t.Errorf("%s: %s left", p, RefPlaceholder)
		}
		if strings.Contains(string(b), "gad/blob/"+Ref()+"/") {
			n++
		}
		return nil
	})
	if n == 0 {
		t.Errorf("no link to the repository at %s", Ref())
	}
}
