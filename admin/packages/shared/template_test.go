package shared

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func templateTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := AutoMigrateTemplates(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestShareTemplates(t *testing.T) {
	db := templateTestDB(t)
	org := uuid.New()

	// a global template and the per-org defaults
	if _, err := CreateTemplate(db, nil, "Global-View", VerbView); err != nil {
		t.Fatal(err)
	}
	if err := SeedDefaultTemplates(db, &org); err != nil {
		t.Fatal(err)
	}
	// seeding is idempotent
	if err := SeedDefaultTemplates(db, &org); err != nil {
		t.Fatal(err)
	}

	// the org sees its 3 defaults + the global one
	forOrg, err := Templates(db, &org)
	if err != nil {
		t.Fatal(err)
	}
	if len(forOrg) != 4 {
		t.Fatalf("org should see 4 templates (3 defaults + 1 global), got %d", len(forOrg))
	}

	// another org sees only the global one (not the first org's defaults)
	other := uuid.New()
	forOther, _ := Templates(db, &other)
	if len(forOther) != 1 || forOther[0].Name != "Global-View" {
		t.Fatalf("another org should see only the global template, got %v", forOther)
	}

	// the Manager template grants view+edit+delete
	var manager ShareTemplate
	db.Where("scope = ? AND name = ?", org, "Manager").First(&manager)
	acts, err := TemplateActions(db, manager.ID)
	if err != nil {
		t.Fatal(err)
	}
	want := len(VerbView) + len(VerbEdit) + len(VerbDelete)
	if len(acts) != want {
		t.Errorf("Manager actions = %v, want %d verbs", acts, want)
	}
}
