package role

import (
	"strings"
	"testing"

	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// A system role follows its originals as the code changes them: an original
// gained added, one lost taken; one taken by hand not given again, one
// added by hand kept. A role made before the originals were kept (no
// SystemPolicies) is given those it lacks, and loses its former versions'.
func TestSystemRoleOriginalsFollowTheCode(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&Role{}, &perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	ensure := func(sr SystemRole) {
		t.Helper()
		b := New(db)
		b.SystemRoles(sr)
		if err := b.EnsureSystemRoles(); err != nil {
			t.Fatal(err)
		}
	}
	res := func(key string) string {
		t.Helper()
		var r Role
		if err := db.First(&r, "system_key = ?", key).Error; err != nil {
			t.Fatal(err)
		}
		return linesOf(db, r)
	}
	add := func(key, line string) {
		var r Role
		db.First(&r, "system_key = ?", key)
		db.Model(&r).Update("permissions", append(r.Permissions, line))
	}
	drop := func(key, line string) {
		var r Role
		db.First(&r, "system_key = ?", key)
		var kept pq.StringArray
		for _, l := range r.Permissions {
			if l != line {
				kept = append(kept, l)
			}
		}
		db.Model(&r).Update("permissions", kept)
	}

	v1 := SystemRole{Key: "fe", Name: "Front", Policies: []SystemPolicy{Allow("admin:a:*"), Allow("admin:b:*")}}
	ensure(v1)
	if got := res("fe"); got != "admin:a:*;admin:b:*" {
		t.Fatalf("made: %s", got)
	}
	// by hand: c added, a taken
	add("fe", "admin:c:*")
	drop("fe", "admin:a:*")
	// the code: b lost, d gained (a still an original: not given again)
	v2 := SystemRole{Key: "fe", Name: "Front", Policies: []SystemPolicy{Allow("admin:a:*"), Allow("admin:d:*"),
		{Effect: perm.Denied, Actions: []string{"*"}, Resources: []string{"admin:d:*:secret*"}}}}
	ensure(v2)
	if got := res("fe"); got != "admin:c:*;admin:d:*;deny admin:d:*:secret*" {
		t.Errorf("v2: %s", got)
	}
	ensure(v2)
	if got := res("fe"); !strings.Contains(got, "admin:d:*;") || strings.Contains(got, "admin:a:*") {
		t.Errorf("v2 again: %s", got)
	}

	// a role made before the originals were kept: an old version of one
	// (Formerly) and some of them
	old := Role{ID: uuid.New(), Name: "Old", SystemKey: "old",
		Permissions: pq.StringArray{"admin:x:*", "admin:keep:*", "admin:mine:*"}}
	db.Create(&old)
	ensure(SystemRole{Key: "old", Name: "Old", Policies: []SystemPolicy{Allow("admin:keep:*"), Allow("admin:new:*")},
		Formerly: []SystemPolicy{Allow("admin:x:*")}})
	if got := res("old"); got != "admin:keep:*;admin:mine:*;admin:new:*" {
		t.Errorf("an old role: %s", got)
	}
}
