package gaddocs

import (
	"io/fs"
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
