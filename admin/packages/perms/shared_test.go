package perms

import (
	"testing"

	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func sharedTestDB(t *testing.T) *gorm.DB {
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

func TestSharedGrantListRevoke(t *testing.T) {
	db := sharedTestDB(t)
	share := uuid.New()

	// an ad-hoc grant (no share) and two grants belonging to the share
	if _, err := Grant(db, "adhoc", "ana", "res:x", VerbView); err != nil {
		t.Fatal(err)
	}
	if _, err := GrantShared(db, "share:ana", "ana", "res:org", VerbView, share); err != nil {
		t.Fatal(err)
	}
	if _, err := GrantShared(db, "share:bruno", "bruno", "res:org", VerbEdit, share); err != nil {
		t.Fatal(err)
	}

	// IsShared distinguishes the two kinds
	if shared, _ := IsShared(db, "adhoc"); shared {
		t.Error("ad-hoc grant must not be marked shared")
	}
	if shared, _ := IsShared(db, "share:ana"); !shared {
		t.Error("share grant must be marked shared")
	}

	// ListShared returns only the share's policies, ordered by subject
	got, err := ListShared(db, share)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Subject != "ana" || got[1].Subject != "bruno" {
		t.Fatalf("ListShared = %+v, want ana then bruno", got)
	}
	for _, p := range got {
		if p.SharedID == nil || *p.SharedID != share {
			t.Errorf("policy %s missing shared id", p.ReferID)
		}
	}

	// RevokeShared removes only the share's policies, leaving the ad-hoc one
	if err := RevokeShared(db, share); err != nil {
		t.Fatal(err)
	}
	if remaining, _ := ListShared(db, share); len(remaining) != 0 {
		t.Errorf("share should be empty after revoke, got %d", len(remaining))
	}
	var adhoc int64
	db.Model(&perm.DefaultDBPolicy{}).Where("refer_id = ?", "adhoc").Count(&adhoc)
	if adhoc != 1 {
		t.Errorf("ad-hoc grant must survive share revoke, got %d", adhoc)
	}
}
