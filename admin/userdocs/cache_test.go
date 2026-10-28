package userdocs

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The files of a package are kept as of its last change: got while its time
// is the same, not when it changed; each package in each language apart.
func TestParsedFiles(t *testing.T) {
	v1, err := FormatFiles([]File{{Path: "README.md", Content: "one"}})
	if err != nil {
		t.Fatal(err)
	}
	v2, _ := FormatFiles([]File{{Path: "README.md", Content: "two"}})
	t1 := time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Microsecond)
	var c parsedFiles
	if _, ok := c.get("en", "pkg", t1); ok {
		t.Fatal("kept before parsed")
	}
	if f, err := c.parse("en", "pkg", t1, v1); err != nil || f["README.md"] != "one" {
		t.Fatalf("the files: %v %v", f, err)
	}
	if f, ok := c.get("en", "pkg", t1.In(time.Local)); !ok || f["README.md"] != "one" {
		t.Errorf("the same time (another zone): %v %v", f, ok)
	}
	if _, ok := c.get("en", "pkg", t2); ok {
		t.Error("a changed package got as it was")
	}
	if _, ok := c.get("pt-BR", "pkg", t1); ok {
		t.Error("another language got")
	}
	c.parse("en", "pkg", t2, v2)
	if f, ok := c.get("en", "pkg", t2); !ok || f["README.md"] != "two" {
		t.Errorf("after the change: %v %v", f, ok)
	}
	if _, err := c.parse("en", "bad", t1, "{"); err == nil {
		t.Error("a value that is not YAML: no error")
	}
}

// A document is read again after its package is edited — its time changed —,
// and from the kept files while it is not.
func TestContentSeesEdits(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	src := testSource("# Things\n")
	if _, err := Sync(db, []*Source{src}, []string{"en"}); err != nil {
		t.Fatal(err)
	}
	b := &Builder{db: db}
	if c, ok, err := b.content("en", src, "things/README.md"); err != nil || !ok || c != "# Things\n" {
		t.Fatalf("the document: %q %v %v", c, ok, err)
	}
	var row UserDocPackage
	db.Take(&row, "locale_code = ? AND id = ?", "en", src.Package)
	files, _ := ParseFiles(row.Value)
	for i := range files {
		if files[i].Path == "things/README.md" {
			files[i].Content = "# Things edited\n"
		}
	}
	edited, _ := FormatFiles(files)
	if err := db.Model(&row).Update("value", edited).Error; err != nil {
		t.Fatal(err)
	}
	if c, _, _ := b.content("en", src, "things/README.md"); c != "# Things edited\n" {
		t.Errorf("after the edit: %q", c)
	}
	// kept: the value itself is not what is read
	db.Exec(`UPDATE user_doc_packages SET value = '{' WHERE id = ?`, src.Package)
	if c, _, err := b.content("en", src, "things/README.md"); err != nil || c != "# Things edited\n" {
		t.Errorf("not from the kept files: %q %v", c, err)
	}
}
