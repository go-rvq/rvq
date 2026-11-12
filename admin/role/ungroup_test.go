package role

import (
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The lines that give to resources inside groups: by their unique names,
// once; what is given to a group, and the roles of no such line, as they are.
func TestMigrateUngroup(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&Role{}); err != nil {
		t.Fatal(err)
	}
	moved := Role{ID: uuid.New(), Name: "Moved", Permissions: pq.StringArray{
		"admin:content/:posts:*", "deny admin:site/:/report:@delete", "admin:content/:*"}}
	kept := Role{ID: uuid.New(), Name: "Kept", Permissions: pq.StringArray{"admin:{content/:*, posts:@list}"}}
	// placeholders: out — the line left of none, too
	empty := Role{ID: uuid.New(), Name: "Empty", Permissions: pq.StringArray{"admin:{users:, posts:@get}", "deny admin:"}}
	db.Create(&empty)
	db.Create(&moved)
	db.Create(&kept)
	if err := MigrateUngroup(db); err != nil {
		t.Fatal(err)
	}
	if got, want := linesOf(db, moved), "admin:posts:*;deny admin:/report:@delete;admin:content/:*"; got != want {
		t.Errorf("moved: %q, want %q", got, want)
	}
	if got, want := linesOf(db, kept), "admin:{content/:*, posts:@list}"; got != want {
		t.Errorf("kept: %q, want %q", got, want)
	}
	if got, want := linesOf(db, empty), "admin:posts:@get"; got != want {
		t.Errorf("placeholders: %q, want %q", got, want)
	}
	// once done, nothing more
	if err := MigrateUngroup(db); err != nil {
		t.Fatal(err)
	}
	if got, want := linesOf(db, moved), "admin:posts:*;deny admin:/report:@delete;admin:content/:*"; got != want {
		t.Errorf("again: %q", got)
	}
}
