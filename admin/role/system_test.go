package role

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testSystemRoles = []SystemRole{
	{Key: "admin", Name: "Administrador", Fixed: true},
	{Key: "editor", Name: "Editor", Policies: []SystemPolicy{Allow("admin:{posts,pages}:*")}},
	{Key: "author", Name: "Autor", Policies: []SystemPolicy{Allow("admin:posts:*@{list,get,create,edit}")}},
}

// systemApp is an admin of roles with the system roles, its database holding
// before what before writes (roles there already).
func systemApp(t *testing.T, before func(db *gorm.DB)) (*gorm.DB, *Builder, http.Handler) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	b := New(db)
	if err := db.AutoMigrate(&Role{}, &perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	if before != nil {
		before(db)
	}
	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.DataOperator(gorm2op.DataOperator(db))
	p.Permission(perm.New().AllowAll().SubjectsFunc(func(*http.Request) []string { return []string{"x"} }))
	b.SystemRoles(testSystemRoles...)
	if err := b.Install(p); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	p.Build(mux)
	return db, b, mux
}

// linesOf are the lines of the role r, as kept, in order.
func linesOf(db *gorm.DB, r Role) string {
	var got Role
	db.First(&got, "id = ?", r.ID)
	return strings.Join(got.Permissions, ";")
}

// The system roles: made on boot with their originals; a role of their name
// adopted — its permissions kept, the originals added —; a fixed
// one with none; made once.
func TestEnsureSystemRoles(t *testing.T) {
	db, b, _ := systemApp(t, func(db *gorm.DB) {
		// the editor there, customized; the administrator there, bare
		db.Create(&Role{ID: uuid.New(), Name: "Editor", Permissions: pq.StringArray{"admin:custom:*"}})
		db.Create(&Role{ID: uuid.New(), Name: "Administrador"})
		db.Create(&Role{ID: uuid.New(), Name: "Mine"})
	})
	get := func(name string) Role {
		var r Role
		if err := db.First(&r, "name = ?", name).Error; err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		return r
	}
	if r := get("Editor"); r.SystemKey != "editor" || linesOf(db, r) != "admin:custom:*;admin:{posts,pages}:*" {
		t.Errorf("the editor adopted, its permissions kept, the originals added: %+v %s", r, linesOf(db, r))
	}
	if r := get("Administrador"); r.SystemKey != "admin" || linesOf(db, r) != "" {
		t.Errorf("the administrator, fixed: %+v %s", r, linesOf(db, r))
	}
	if r := get("Autor"); r.SystemKey != "author" || linesOf(db, r) != "admin:posts:*@{list,get,create,edit}" {
		t.Errorf("the author made with its originals: %+v %s", r, linesOf(db, r))
	}
	if r := get("Mine"); r.SystemKey != "" {
		t.Errorf("a role of the users marked: %+v", r)
	}

	// on boot again: nothing made twice — nor added
	if err := b.EnsureSystemRoles(); err != nil {
		t.Fatal(err)
	}
	var n int64
	db.Model(&Role{}).Count(&n)
	if n != 4 {
		t.Errorf("%d roles", n)
	}
	if r := get("Editor"); strings.Count(linesOf(db, r), ";") != 1 {
		t.Errorf("the originals added twice: %s", linesOf(db, r))
	}
	// adopted again (its key lost): nothing added twice
	db.Model(&Role{}).Where("name = ?", "Editor").Update("system_key", "")
	if err := b.EnsureSystemRoles(); err != nil {
		t.Fatal(err)
	}
	if r := get("Editor"); r.SystemKey != "editor" || strings.Count(linesOf(db, r), ";") != 1 {
		t.Errorf("adopted again: %+v %s", r, linesOf(db, r))
	}

	// reset: the editor's originals; the fixed one refuses
	ed := get("Editor")
	if err := b.ResetPermissions(&ed); err != nil {
		t.Fatal(err)
	}
	if got := linesOf(db, ed); got != "admin:{posts,pages}:*" {
		t.Errorf("reset: %s", got)
	}
	admin := get("Administrador")
	if err := b.ResetPermissions(&admin); err != ErrSystemRoleFixed {
		t.Errorf("reset of the fixed: %v", err)
	}
}

// An original changed (Formerly): a role made before loses the former one on
// boot — kept when it was changed by its users —, and gets the new.
func TestSystemRoleFormerly(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	b := New(db)
	if err := db.AutoMigrate(&Role{}, &perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	old := SystemRole{Key: "dev", Name: "Dev", Policies: []SystemPolicy{Allow("admin:/files:*"), Allow("admin:docs:*")}}
	b.SystemRoles(old)
	if err := b.EnsureSystemRoles(); err != nil {
		t.Fatal(err)
	}
	var r Role
	db.First(&r, "system_key = ?", "dev")
	// one of its users' too
	db.Model(&r).Update("permissions", append(r.Permissions, "admin:mine:*"))

	b2 := New(db)
	b2.SystemRoles(SystemRole{Key: "dev", Name: "Dev",
		Policies: []SystemPolicy{Allow("admin:/files:!{edit,git}"), Allow("admin:docs:*")},
		Formerly: []SystemPolicy{Allow("admin:/files:*")}})
	for i := 0; i < 2; i++ { // and again: nothing more
		if err := b2.EnsureSystemRoles(); err != nil {
			t.Fatal(err)
		}
	}
	if got := linesOf(db, r); got != "admin:docs:*;admin:mine:*;admin:/files:!{edit,git}" {
		t.Errorf("after the change: %s", got)
	}
}

// In the admin: a system role is not deleted nor renamed; its permissions
// are reset by its action; a role of the users is deleted.
func TestSystemRolesAdmin(t *testing.T) {
	db, _, h := systemApp(t, func(db *gorm.DB) { db.Create(&Role{ID: uuid.New(), Name: "Mine"}) })
	var ed, mine Role
	db.First(&ed, "name = ?", "Editor")
	db.First(&mine, "name = ?", "Mine")
	event := func(ev string, query, form [][2]string) *httptest.ResponseRecorder {
		mb := multipartestutils.NewMultipartBuilder().PageURL("/admin/roles").EventFunc(ev)
		for _, q := range query {
			mb = mb.Query(q[0], q[1])
		}
		for _, f := range form {
			mb = mb.AddField(f[0], f[1])
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, mb.BuildEventFuncRequest())
		return w
	}

	w := event(actions.DoDelete, [][2]string{{presets.ParamID, ed.ID.String()}}, nil)
	if db.First(&Role{}, "id = ?", ed.ID).Error != nil {
		t.Errorf("a system role deleted: %.300s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "not deleted") {
		t.Errorf("no reason: %.300s", w.Body.String())
	}
	w = event(actions.DoDelete, [][2]string{{presets.ParamID, mine.ID.String()}}, nil)
	if db.First(&Role{}, "id = ?", mine.ID).Error == nil {
		t.Errorf("a role of the users not deleted: %.400s", w.Body.String())
	}

	// the form of the role, its stamp: the update is the validator's to refuse
	w = event(actions.Edit, [][2]string{{presets.ParamID, ed.ID.String()}}, nil)
	stamp := regexp.MustCompile(`__formSign\\?":\s*\\?"([^"\\]+)`).FindStringSubmatch(w.Body.String())
	if stamp == nil {
		t.Fatalf("the edit form carries no stamp: %.300s", w.Body.String())
	}
	w = event(actions.Update, [][2]string{{presets.ParamID, ed.ID.String()}},
		[][2]string{{"Name", "Chief"}, {presets.RecordStampFormKey, stamp[1]}})
	var kept Role
	db.First(&kept, "id = ?", ed.ID)
	if kept.Name != "Editor" {
		t.Errorf("a system role renamed: %q", kept.Name)
	}
	if !strings.Contains(w.Body.String(), "keeps its name") {
		t.Errorf("no reason: %.300s", w.Body.String())
	}

	// changed, then reset by its action
	db.Model(&ed).Update("permissions", pq.StringArray{"admin:other:*"})
	w = event(eventResetPermissions, [][2]string{{"id", ed.ID.String()}}, nil)
	if got := linesOf(db, ed); got != "admin:{posts,pages}:*" {
		t.Errorf("the reset action: %s %.300s", got, w.Body.String())
	}
}

// The lines as the form posts them — a list of texts, each item its
// __value —: kept as written, in their order, the blank ones out; a prefix
// of no resource refused.
func TestRolePermissionsAsTheFormPosts(t *testing.T) {
	db, _, h := systemApp(t, nil)
	post := func(name string, lines ...string) {
		mb := multipartestutils.NewMultipartBuilder().PageURL("/admin/roles").EventFunc(actions.Create).
			AddField("Name", name).AddField("Permissions.__present", "1")
		for i, l := range lines {
			mb = mb.AddField(fmt.Sprintf("Permissions[%d].__value", i), l).AddField(fmt.Sprintf("Permissions[%d].$id", i), fmt.Sprint(i+1))
		}
		h.ServeHTTP(httptest.NewRecorder(), mb.BuildEventFuncRequest())
	}
	post("Lists", "admin:x:*", " ", "deny admin:x:@delete", "admin:users:{*, deny @edit}")
	var r Role
	db.First(&r, "name = ?", "Lists")
	if got := strings.Join(r.Permissions, "|"); got != "admin:x:*|deny admin:x:@delete|admin:users:{*, deny @edit}" {
		t.Errorf("the lines: %q", got)
	}
	post("Bad", "admin:x:*", "deny")
	if db.First(&Role{}, "name = ?", "Bad").Error == nil {
		t.Error("a role of a prefix alone made")
	}
}

// The roles of before the lines: their policies made lines — in the order
// they were made, "deny " before a deny's —, and their originals kept in the
// form of before made lines; the policies gone; once.
func TestMigrateLines(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&Role{}, &perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	old := Role{ID: uuid.New(), Name: "Old", SystemKey: "old", SystemPolicies: "allow|*|admin:a:*\ndeny|*|admin:a:secret*"}
	db.Create(&old)
	for i, p := range []perm.DefaultDBPolicy{
		{Effect: perm.Allowed, Actions: pq.StringArray{"*"}, Resources: pq.StringArray{"admin:a:*"}},
		{Effect: perm.Denied, Actions: pq.StringArray{"*"}, Resources: pq.StringArray{"admin:a:secret*"}},
		{Effect: perm.Allowed, Actions: pq.StringArray{"@get"}, Resources: pq.StringArray{"admin:b:*"}},
	} {
		p.ID = uuid.New()
		p.ReferID, p.Subject = old.ID.String(), "Old"
		p.CreatedAt = time.Now().Add(time.Duration(i) * time.Second)
		if err := db.Create(&p).Error; err != nil {
			t.Fatal(err)
		}
	}
	lined := Role{ID: uuid.New(), Name: "Lined", Permissions: pq.StringArray{"admin:c:*"}}
	db.Create(&lined)
	for range 2 {
		if err := MigrateLines(db); err != nil {
			t.Fatal(err)
		}
	}
	var r Role
	db.First(&r, "id = ?", old.ID)
	if got := strings.Join(r.Permissions, "|"); got != "admin:a:*|deny admin:a:secret*|admin:b:*@get" {
		t.Errorf("the lines: %s", got)
	}
	if r.SystemPolicies != "admin:a:*\ndeny admin:a:secret*" {
		t.Errorf("the originals: %q", r.SystemPolicies)
	}
	var n int64
	db.Model(&perm.DefaultDBPolicy{}).Count(&n)
	if n != 0 {
		t.Errorf("%d policies left", n)
	}
	if linesOf(db, lined) != "admin:c:*" {
		t.Errorf("a role of lines changed: %s", linesOf(db, lined))
	}
}

// The stack of a user of roles: the code's lines first (of its subjects),
// then each role's, in the order of the roles, each of its origin.
func TestStack(t *testing.T) {
	pb := perm.New().Policies(perm.PolicyFor("Editor").WhoAre(perm.Allowed).ToDo(perm.Anything).On("admin:code:*"))
	s := Stack(pb, []string{"Editor"}, []*Role{
		{Name: "Editor", Permissions: pq.StringArray{"admin:users:{list,new}"}},
		{Name: "Locked", Permissions: pq.StringArray{"deny admin:users:new"}},
	}).Resolve()
	var got []string
	for _, l := range s.Lines {
		got = append(got, fmt.Sprintf("%d %s %s %d", l.Index, l.Text(), l.Origin, l.ReplacedBy))
	}
	want := "1 admin:code:* code: Editor 0|2 admin:users:list role: Editor 0|3 admin:users:new role: Editor 4|4 deny admin:users:new role: Locked 0"
	if strings.Join(got, "|") != want {
		t.Errorf("the stack:\n%s\nwant\n%s", strings.Join(got, "|"), want)
	}
}
