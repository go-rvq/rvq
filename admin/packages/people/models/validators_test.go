package models

import (
	"testing"

	validators "github.com/go-rvq/rvq/admin/packages/validators/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := SeedDocumentValidators(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestDocumentValidatorsSeeded(t *testing.T) {
	db := testDB(t)
	for _, name := range []string{CPFValidatorName, CNPJValidatorName} {
		if _, err := validators.Get(db, name); err != nil {
			t.Errorf("validator %q not seeded: %v", name, err)
		}
	}
	// idempotent re-seed
	if err := SeedDocumentValidators(db); err != nil {
		t.Fatalf("re-seed: %v", err)
	}
	var n int64
	db.Model(&validators.Validator{}).Count(&n)
	if n != 2 {
		t.Errorf("expected 2 validators after re-seed, got %d", n)
	}
}

func TestCPFViaEngine(t *testing.T) {
	db := testDB(t)
	// valid CPFs pass
	for _, d := range []string{"529.982.247-25", "52998224725"} {
		if err := validators.Check(db, CPFValidatorName, d, "pt-BR"); err != nil {
			t.Errorf("CPF %q should be valid: %v", d, err)
		}
	}
	// invalid CPF returns the translated message
	if err := validators.Check(db, CPFValidatorName, "123", "en-US"); err == nil || err.Error() != "Invalid CPF." {
		t.Errorf("invalid CPF (en-US) = %v, want \"Invalid CPF.\"", err)
	}
	if err := validators.Check(db, CPFValidatorName, "111.111.111-11", "pt-BR"); err == nil || err.Error() != "CPF inválido." {
		t.Errorf("repeated CPF (pt-BR) = %v, want \"CPF inválido.\"", err)
	}
}

func TestCNPJViaEngine(t *testing.T) {
	db := testDB(t)
	if err := validators.Check(db, CNPJValidatorName, "11.222.333/0001-81", "pt-BR"); err != nil {
		t.Errorf("valid CNPJ: %v", err)
	}
	if err := validators.Check(db, CNPJValidatorName, "11.222.333/0001-80", "en-US"); err == nil || err.Error() != "Invalid CNPJ." {
		t.Errorf("invalid CNPJ = %v, want \"Invalid CNPJ.\"", err)
	}
}

func TestValidatorNameForDocumentType(t *testing.T) {
	if ValidatorNameForDocumentType(DocumentCPF) != CPFValidatorName ||
		ValidatorNameForDocumentType(DocumentCNPJ) != CNPJValidatorName ||
		ValidatorNameForDocumentType(DocumentOther) != "" {
		t.Errorf("document type → validator name mapping wrong")
	}
}
