package shared

import (
	"testing"

	"github.com/go-rvq/rvq/admin/packages/perms"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestShareLifecycle(t *testing.T) {
	db := testDB(t)
	const resource = "presets:orgs:o1:*"

	// share with two subjects
	sid, err := Create(db, resource, []string{"ana", "bruno"}, VerbView)
	if err != nil {
		t.Fatal(err)
	}
	subs, err := Subjects(db, sid)
	if err != nil {
		t.Fatal(err)
	}
	if len(subs) != 2 {
		t.Fatalf("share should have 2 subjects, got %d", len(subs))
	}
	for _, p := range subs {
		if p.SharedID == nil || *p.SharedID != sid {
			t.Errorf("policy %s not stamped with the share id", p.ReferID)
		}
		if len(p.Resources) != 1 || p.Resources[0] != resource {
			t.Errorf("policy resource = %v, want %s", p.Resources, resource)
		}
	}

	// add a third subject with edit rights
	if err := AddSubjects(db, sid, resource, []string{"clara"}, append(append([]string{}, VerbView...), VerbEdit...)); err != nil {
		t.Fatal(err)
	}
	if subs, _ = Subjects(db, sid); len(subs) != 3 {
		t.Fatalf("share should have 3 subjects after add, got %d", len(subs))
	}

	// revoke a single subject
	if err := RevokeSubject(db, sid, "bruno"); err != nil {
		t.Fatal(err)
	}
	if subs, _ = Subjects(db, sid); len(subs) != 2 {
		t.Fatalf("after removing bruno the share should have 2 subjects, got %d", len(subs))
	}

	// revoke the whole share
	if err := Revoke(db, sid); err != nil {
		t.Fatal(err)
	}
	if subs, _ = Subjects(db, sid); len(subs) != 0 {
		t.Fatalf("after revoke the share should be empty, got %d", len(subs))
	}

	// an unrelated ad-hoc grant is untouched by share operations
	if _, err := perms.Grant(db, "adhoc", "dora", resource, VerbView); err != nil {
		t.Fatal(err)
	}
	sid2, _ := Create(db, resource, []string{"ana"}, VerbView)
	if err := Revoke(db, sid2); err != nil {
		t.Fatal(err)
	}
	var adhoc int64
	db.Model(&perm.DefaultDBPolicy{}).Where("refer_id = ?", "adhoc").Count(&adhoc)
	if adhoc != 1 {
		t.Errorf("ad-hoc grant must survive share revoke, got %d", adhoc)
	}
}
