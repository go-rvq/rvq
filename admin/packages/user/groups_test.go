package user

import (
	"context"
	"fmt"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func groupsDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := stackDB(t)
	if err := AutoMigrateGroups(db); err != nil {
		t.Fatal(err)
	}
	return db
}

// The groups of a user in their order (joined last, reordered); a member
// joining and leaving told to the hooks; the generation changed.
func TestGroupsMembers(t *testing.T) {
	db := groupsDB(t)
	ana := uuid.New()
	a, _ := CreateGroup(db, "A", "")
	b, _ := CreateGroup(db, "B", "")
	if _, err := CreateGroup(db, " ", ""); err == nil {
		t.Error("a group with no name")
	}
	var told []string
	OnMemberChange(func(_ *gorm.DB, g, u uuid.UUID, added bool, by string) error {
		told = append(told, fmt.Sprintf("%v %v %s", g == b.ID, added, by))
		return nil
	})
	gen := GroupGeneration()
	for _, g := range []*Group{b, a, b} {
		if err := AddMember(db, g.ID, ana, "root"); err != nil {
			t.Fatal(err)
		}
	}
	names := func() string {
		gs, err := GroupsOf(db, ana)
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, g := range gs {
			out = append(out, g.Name)
		}
		return strings.Join(out, ",")
	}
	if got := names(); got != "B,A" {
		t.Errorf("joined: %s", got)
	}
	if err := SetGroupPositions(db, ana, []uuid.UUID{a.ID, b.ID}); err != nil || names() != "A,B" {
		t.Errorf("reordered: %s %v", names(), err)
	}
	if err := RemoveMember(db, b.ID, ana, "root"); err != nil || names() != "A" {
		t.Errorf("left: %s %v", names(), err)
	}
	if GroupGeneration() == gen {
		t.Error("the generation")
	}
	if !slices.Equal(told, []string{"true true root", "false true root", "true false root"}) {
		t.Errorf("told %v", told)
	}
}

// The stack of a user: the user's roles, then each group's in the user's
// order of their groups, each group's in its order; the lines of the
// database for the groups' subjects.
func TestStackWithGroups(t *testing.T) {
	db := groupsDB(t)
	mk := func(name string, perms ...string) *role.Role {
		r := &role.Role{ID: uuid.New(), Name: name, Permissions: pq.StringArray(perms)}
		if err := db.Create(r).Error; err != nil {
			t.Fatal(err)
		}
		return r
	}
	own, g1, g2a, g2b := mk("Own", "admin:a:list"), mk("G1", "admin:b:list"), mk("G2a", "admin:c:list"), mk("G2b", "admin:d:list")
	u := &stackUser{ID: uuid.New(), Name: "Ana", Roles: []*role.Role{own}}
	u.Account = "ana"
	if err := db.Create(u).Error; err != nil {
		t.Fatal(err)
	}
	x, _ := CreateGroup(db, "X", "")
	y, _ := CreateGroup(db, "Y", "")
	_ = SetGroupRoles(db, x.ID, []uuid.UUID{g1.ID})
	_ = SetGroupRoles(db, y.ID, []uuid.UUID{g2b.ID, g2a.ID})
	_ = AddMember(db, y.ID, u.ID, "")
	_ = AddMember(db, x.ID, u.ID, "")
	b := &Builder{db: db, lb: login.New(i18n.New())}
	s := b.Stack(perm.New(), u)
	var got []string
	for _, l := range s.Lines {
		got = append(got, l.Origin)
	}
	want := []string{"role: Own", "role: G2b", "role: G2a", "role: G1"}
	if !slices.Equal(got, want) {
		t.Errorf("origins %v, want %v", got, want)
	}
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), login.UserKey, u))
	_, v1, _ := b.StackSource(perm.New()).StackOf(r)
	_ = SetGroupRoles(db, x.ID, nil)
	if _, v2, _ := b.StackSource(perm.New()).StackOf(r); v1 == v2 {
		t.Error("a group changed: the same version")
	}
}

// The default language of a user; the e-mail of a user.
func TestPreferenceAndEmail(t *testing.T) {
	db := groupsDB(t)
	id := uuid.New()
	if LanguageOf(db, id) != "" {
		t.Error("none chosen")
	}
	_ = SetLanguage(db, id, "pt-BR")
	_ = SetLanguage(db, id, "en-US")
	if LanguageOf(db, id) != "en-US" {
		t.Errorf("chosen: %q", LanguageOf(db, id))
	}
	u := &stackUser{}
	u.Account = "ana@x.test"
	if EmailOf(u) != "ana@x.test" {
		t.Error("the account as the address")
	}
	u.Account = "admin"
	if EmailOf(u) != "" {
		t.Error("an account that is no address")
	}
}
