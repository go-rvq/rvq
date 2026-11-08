package models

import (
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// The default organization of a user: one, changed, cleared; another
// user's apart.
func TestUserDefault(t *testing.T) {
	db := testDB(t)
	ana, bia := uuid.New(), uuid.New()
	a := &Organizacao{Nome: "A", OwnerID: model.UUID(ana)}
	b := &Organizacao{Nome: "B", OwnerID: model.UUID(ana)}
	db.Create(a)
	db.Create(b)
	if DefaultOf(db, ana) != uuid.Nil {
		t.Fatal("a default of no choice")
	}
	if err := SetDefault(db, ana, a.ID); err != nil {
		t.Fatal(err)
	}
	if err := SetDefault(db, ana, b.ID); err != nil {
		t.Fatal(err)
	}
	if err := SetDefault(db, bia, a.ID); err != nil {
		t.Fatal(err)
	}
	if DefaultOf(db, ana) != b.ID || DefaultOf(db, bia) != a.ID {
		t.Fatalf("the defaults: %v %v", DefaultOf(db, ana), DefaultOf(db, bia))
	}
	var n int64
	db.Model(&UserDefault{}).Where("user_id = ?", ana).Count(&n)
	if n != 1 {
		t.Errorf("rows of ana: %d", n)
	}
	if err := ClearDefault(db, ana); err != nil || DefaultOf(db, ana) != uuid.Nil || DefaultOf(db, bia) != a.ID {
		t.Errorf("cleared: %v", err)
	}
}

// The defaults of before — an owner's organization marked padrao — become
// the owners' defaults, the column gone.
func TestUserDefaultsOfBefore(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&Organizacao{}); err != nil {
		t.Fatal(err)
	}
	ana := uuid.New()
	a := &Organizacao{Nome: "A", OwnerID: model.UUID(ana)}
	b := &Organizacao{Nome: "B", OwnerID: model.UUID(ana)}
	db.Create(a)
	db.Create(b)
	if err := db.Exec(`ALTER TABLE org_organizacoes ADD COLUMN padrao boolean NOT NULL DEFAULT false`).Error; err != nil {
		t.Fatal(err)
	}
	db.Exec(`UPDATE org_organizacoes SET padrao = true WHERE id = ?`, b.ID)
	if err := AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	if DefaultOf(db, ana) != b.ID {
		t.Errorf("the default of before: %v", DefaultOf(db, ana))
	}
	if hasColumn(db, "padrao") {
		t.Error("the column padrao stays")
	}
	if err := AutoMigrate(db); err != nil {
		t.Errorf("again: %v", err)
	}
}
