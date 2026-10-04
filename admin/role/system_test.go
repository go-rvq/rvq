package role

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"
	"testing"

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

func resourcesOf(db *gorm.DB, r Role) string {
	var ps []perm.DefaultDBPolicy
	db.Where("refer_id = ?", r.ID.String()).Find(&ps)
	var res []string
	for _, p := range ps {
		res = append(res, strings.Join(p.Resources, ",")+"="+p.Effect+"/"+p.Subject)
	}
	sort.Strings(res)
	return strings.Join(res, ";")
}

// The system roles: made on boot with their originals; a role of their name
// adopted — its permissions kept, the originals added —; a fixed
// one with none; made once.
func TestEnsureSystemRoles(t *testing.T) {
	db, b, _ := systemApp(t, func(db *gorm.DB) {
		// the editor there, customized; the administrator there, bare
		ed := Role{ID: uuid.New(), Name: "Editor"}
		db.Create(&ed)
		db.Create(&perm.DefaultDBPolicy{ReferID: ed.ID.String(), Subject: "Editor", Effect: perm.Allowed,
			Actions: pq.StringArray{"*"}, Resources: pq.StringArray{"admin:custom:*"}})
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
	if r := get("Editor"); r.SystemKey != "editor" || resourcesOf(db, r) != "admin:custom:*=allow/Editor;admin:{posts,pages}:*=allow/Editor" {
		t.Errorf("the editor adopted, its permissions kept, the originals added: %+v %s", r, resourcesOf(db, r))
	}
	if r := get("Administrador"); r.SystemKey != "admin" || resourcesOf(db, r) != "" {
		t.Errorf("the administrator, fixed: %+v %s", r, resourcesOf(db, r))
	}
	if r := get("Autor"); r.SystemKey != "author" || resourcesOf(db, r) != "admin:posts:*@{list,get,create,edit}=allow/Autor" {
		t.Errorf("the author made with its originals: %+v %s", r, resourcesOf(db, r))
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
	if r := get("Editor"); strings.Count(resourcesOf(db, r), ";") != 1 {
		t.Errorf("the originals added twice: %s", resourcesOf(db, r))
	}
	// adopted again (its key lost): nothing added twice
	db.Model(&Role{}).Where("name = ?", "Editor").Update("system_key", "")
	if err := b.EnsureSystemRoles(); err != nil {
		t.Fatal(err)
	}
	if r := get("Editor"); r.SystemKey != "editor" || strings.Count(resourcesOf(db, r), ";") != 1 {
		t.Errorf("adopted again: %+v %s", r, resourcesOf(db, r))
	}

	// reset: the editor's originals; the fixed one refuses
	ed := get("Editor")
	if err := b.ResetPermissions(&ed); err != nil {
		t.Fatal(err)
	}
	if got := resourcesOf(db, ed); got != "admin:{posts,pages}:*=allow/Editor" {
		t.Errorf("reset: %s", got)
	}
	admin := get("Administrador")
	if err := b.ResetPermissions(&admin); err != ErrSystemRoleFixed {
		t.Errorf("reset of the fixed: %v", err)
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
	db.Where("refer_id = ?", ed.ID.String()).Delete(&perm.DefaultDBPolicy{})
	w = event(eventResetPermissions, [][2]string{{"id", ed.ID.String()}}, nil)
	if got := resourcesOf(db, ed); got != "admin:{posts,pages}:*=allow/Editor" {
		t.Errorf("the reset action: %s %.300s", got, w.Body.String())
	}
}
