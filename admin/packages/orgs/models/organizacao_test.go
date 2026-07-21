package models

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestOrganizacaoCRUDAndValidation(t *testing.T) {
	db := testDB(t)
	owner := uuid.New()

	// name required
	if err := db.Create(&Organizacao{OwnerID: owner}).Error; !errors.Is(err, ErrNomeObrigatorio) {
		t.Fatalf("empty name should fail with ErrNomeObrigatorio, got %v", err)
	}
	// owner required
	if err := db.Create(&Organizacao{Nome: "Acme"}).Error; !errors.Is(err, ErrProprietarioObrigatorio) {
		t.Fatalf("missing owner should fail with ErrProprietarioObrigatorio, got %v", err)
	}

	// valid create assigns a UUID id
	o := &Organizacao{Nome: "Acme", Descricao: "Test", OwnerID: owner}
	if err := db.Create(o).Error; err != nil {
		t.Fatal(err)
	}
	if o.ID == uuid.Nil {
		t.Error("id should be assigned on create")
	}
	if o.String() != "Acme" {
		t.Errorf("String() = %q, want Acme", o.String())
	}

	// read back
	var got Organizacao
	if err := db.First(&got, "id = ?", o.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.Nome != "Acme" || got.OwnerID != owner {
		t.Errorf("read back = %+v", got)
	}
}
