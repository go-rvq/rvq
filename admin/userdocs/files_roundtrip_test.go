package userdocs_test

import (
	"testing"

	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/admin/userdocs/gaddocs"
)

// The documents written read back as they are — the ones yaml.v3 writes in a
// form it does not parse (a tab, spaces at the end of a line of a block: Gad's
// documentation has them) too, as JSON.
func TestFormatFilesRoundTrip(t *testing.T) {
	files, err := gaddocs.Source().FilesOf("en")
	if err != nil {
		t.Fatal(err)
	}
	s, err := userdocs.FormatFiles(files)
	if err != nil {
		t.Fatal(err)
	}
	back, err := userdocs.ParseFiles(s)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(back) != len(files) {
		t.Fatalf("files: %d, want %d", len(back), len(files))
	}
	for i := range files {
		if back[i] != files[i] {
			t.Errorf("%s: not as written", files[i].Path)
		}
	}
}
