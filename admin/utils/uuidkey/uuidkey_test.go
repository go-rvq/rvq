package uuidkey

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type keyed struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title string
}

// A model with a BeforeCreate of its own — which would shadow an embedded one.
type hooked struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Hooked bool
}

func (h *hooked) BeforeCreate(*gorm.DB) error { h.Hooked = true; return nil }

// Composite: a UUID and a locale, as pages and posts.
type localized struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	LocaleCode string    `gorm:"primaryKey"`
}

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := Register(db); err != nil {
		t.Fatal(err)
	}
	if err := Register(db); err != nil { // idempotent
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&keyed{}, &hooked{}, &localized{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestNewIsV7(t *testing.T) {
	a, b := New(), New()
	if a.Version() != 7 {
		t.Errorf("version %d, want 7", a.Version())
	}
	if a.String() >= b.String() {
		t.Errorf("keys made one after the other must sort: %s >= %s", a, b)
	}
}

func TestCreateFillsZeroKeys(t *testing.T) {
	db := testDB(t)

	one := &keyed{Title: "a"}
	if err := db.Create(one).Error; err != nil {
		t.Fatal(err)
	}
	if one.ID == uuid.Nil || one.ID.Version() != 7 {
		t.Errorf("no v7 key: %s", one.ID)
	}

	many := []*keyed{{Title: "b"}, {Title: "c"}}
	if err := db.Create(&many).Error; err != nil {
		t.Fatal(err)
	}
	if many[0].ID == uuid.Nil || many[1].ID == uuid.Nil || many[0].ID == many[1].ID {
		t.Errorf("batch keys: %s %s", many[0].ID, many[1].ID)
	}

	// a key given is kept
	given := New()
	if err := db.Create(&keyed{ID: given, Title: "d"}).Error; err != nil {
		t.Fatal(err)
	}
	var got keyed
	db.First(&got, "id = ?", given)
	if got.Title != "d" {
		t.Error("a given key was replaced")
	}

	// the model's own hook still runs, and the key is filled anyway
	h := &hooked{}
	if err := db.Create(h).Error; err != nil {
		t.Fatal(err)
	}
	if !h.Hooked || h.ID == uuid.Nil {
		t.Errorf("hooked: ran=%v id=%s", h.Hooked, h.ID)
	}

	// composite key: the UUID part is filled, the locale kept; a translation
	// that shares the id keeps it
	l := &localized{LocaleCode: "en-US"}
	if err := db.Create(l).Error; err != nil {
		t.Fatal(err)
	}
	tr := &localized{ID: l.ID, LocaleCode: "pt-BR"}
	if err := db.Create(tr).Error; err != nil {
		t.Fatal(err)
	}
	if l.ID == uuid.Nil || tr.ID != l.ID {
		t.Errorf("localized: %s / %s", l.ID, tr.ID)
	}
}

func TestShortPath(t *testing.T) {
	id := uuid.MustParse("019a8f3c-1111-7222-8333-444455556666")
	p := ShortPath(id)
	if p != "01/9a8f3c-1111-7222-8333-444455556666" {
		t.Errorf("ShortPath = %s", p)
	}
	if !strings.HasPrefix(p, id.String()[:2]+"/") {
		t.Error("the first level is the first two characters")
	}
}
