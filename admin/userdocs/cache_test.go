package userdocs

import (
	"reflect"
	"testing"
)

// The files of a value are parsed once while the value is the same; a
// changed value (an edit, a sync) is parsed again at once; each package in
// each language is kept apart.
func TestParsedFiles(t *testing.T) {
	v1, err := FormatFiles([]File{{Path: "README.md", Content: "one"}})
	if err != nil {
		t.Fatal(err)
	}
	v2, _ := FormatFiles([]File{{Path: "README.md", Content: "two"}})
	var c parsedFiles
	a, err := c.of("en", "pkg", v1)
	if err != nil || a["README.md"] != "one" {
		t.Fatalf("the files: %v %v", a, err)
	}
	// the same value: the same map, not parsed again
	if b, _ := c.of("en", "pkg", v1); reflect.ValueOf(b).Pointer() != reflect.ValueOf(a).Pointer() {
		t.Error("the same value parsed again")
	}
	if b, _ := c.of("en", "pkg", v2); b["README.md"] != "two" {
		t.Errorf("a changed value not seen: %v", b)
	}
	if b, _ := c.of("pt-BR", "pkg", v1); b["README.md"] != "one" {
		t.Errorf("another language: %v", b)
	}
	if b, _ := c.of("en", "pkg", v2); b["README.md"] != "two" {
		t.Errorf("the languages mixed: %v", b)
	}
	if _, err := c.of("en", "bad", "{"); err == nil {
		t.Error("a value that is not YAML: no error")
	}
}
