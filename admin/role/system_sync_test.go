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
		return resourcesOf(db, r)
	}
	add := func(key, resource string) {
		var r Role
		db.First(&r, "system_key = ?", key)
		db.Create(&perm.DefaultDBPolicy{ReferID: r.ID.String(), Subject: r.Name, Effect: perm.Allowed,
			Actions: pq.StringArray{"*"}, Resources: pq.StringArray{resource}})
	}
	drop := func(key, resource string) {
		var r Role
		db.First(&r, "system_key = ?", key)
		db.Where("refer_id = ? AND resources = ?", r.ID.String(), pq.StringArray{resource}).Delete(&perm.DefaultDBPolicy{})
	}

	v1 := SystemRole{Key: "fe", Name: "Front", Policies: []SystemPolicy{Allow("admin:a:*"), Allow("admin:b:*")}}
	ensure(v1)
	if got := res("fe"); got != "admin:a:*=allow/Front;admin:b:*=allow/Front" {
		t.Fatalf("made: %s", got)
	}
	// by hand: c added, a taken
	add("fe", "admin:c:*")
	drop("fe", "admin:a:*")
	// the code: b lost, d gained (a still an original: not given again)
	v2 := SystemRole{Key: "fe", Name: "Front", Policies: []SystemPolicy{Allow("admin:a:*"), Allow("admin:d:*"),
		{Effect: perm.Denied, Actions: []string{"*"}, Resources: []string{"admin:d:*:secret*"}}}}
	ensure(v2)
	if got := res("fe"); got != "admin:c:*=allow/Front;admin:d:*:secret*=deny/Front;admin:d:*=allow/Front" {
		t.Errorf("v2: %s", got)
	}
	ensure(v2)
	if got := res("fe"); !strings.Contains(got, "admin:d:*=allow") || strings.Contains(got, "admin:a:*") {
		t.Errorf("v2 again: %s", got)
	}

	// a role made before the originals were kept: an old version of one
	// (Formerly) and some of them
	old := Role{ID: uuid.New(), Name: "Old", SystemKey: "old"}
	db.Create(&old)
	for _, r := range []string{"admin:x:*", "admin:keep:*", "admin:mine:*"} {
		db.Create(&perm.DefaultDBPolicy{ReferID: old.ID.String(), Subject: "Old", Effect: perm.Allowed,
			Actions: pq.StringArray{"*"}, Resources: pq.StringArray{r}})
	}
	ensure(SystemRole{Key: "old", Name: "Old", Policies: []SystemPolicy{Allow("admin:keep:*"), Allow("admin:new:*")},
		Formerly: []SystemPolicy{Allow("admin:x:*")}})
	if got := res("old"); got != "admin:keep:*=allow/Old;admin:mine:*=allow/Old;admin:new:*=allow/Old" {
		t.Errorf("an old role: %s", got)
	}
}
