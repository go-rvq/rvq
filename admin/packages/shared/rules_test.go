package shared_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/go-rvq/rvq/admin/packages/shared"
	"github.com/go-rvq/rvq/admin/packages/shared/sharedtest"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
)

// A role of rules: its Actions on the record and all under it, each rule's
// actions on its resource under the record.
func TestRoleRules(t *testing.T) {
	f := newFixture(t)
	if err := f.db.AutoMigrate(&perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	p := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(f.db)).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	pastas := p.Model(&sharedtest.Pasta{}).URIName("pastas")
	view := []string{presets.PermList, presets.PermGet}
	if err := f.s.Install(pastas, shared.ModelOptions{
		Resource: shared.Resource{Name: "admin_pastas", Table: "sharedtest_pastas", OwnerColumn: "owner_id"},
		Roles: []shared.RoleSeed{{Name: "Redator", Actions: view,
			Rules: []shared.Rule{{On: "documentos:*", Do: []string{presets.PermCreate, presets.PermUpdate}}}}},
	}); err != nil {
		t.Fatal(err)
	}
	roles, _ := shared.RolesOf(f.db, uuid.Nil, "admin_pastas")
	f.roles["admin_pastas"] = map[string]uuid.UUID{}
	for _, r := range roles {
		f.roles["admin_pastas"][r.Name] = r.ID
	}
	pasta := sharedtest.NewPasta(f.db, "Contratos", "ana")
	f.share("admin_pastas", pasta.ID, shared.UserSubject("bia"), "Redator", "ana")
	var ps []perm.DefaultDBPolicy
	f.db.Where("subject = ?", shared.UserSubject("bia")).Order("refer_id").Find(&ps)
	rec := "admin:pasta:<" + pasta.ID + ">:"
	if len(ps) != 2 || ps[0].Resources[0] != rec+"*" || !slices.Equal([]string(ps[0].Actions), view) ||
		ps[1].Resources[0] != rec+"documentos:*" ||
		!slices.Equal([]string(ps[1].Actions), []string{presets.PermCreate, presets.PermUpdate}) {
		t.Fatalf("the policies: %+v", ps)
	}
}

// More: an action on a resource the role from has not — on it nor on all.
func TestMoreByRules(t *testing.T) {
	view := &shared.ShareTemplate{Actions: []string{"@list", "@get"}}
	editDocs := &shared.ShareTemplate{Actions: []string{"@list", "@get"}, Rules: []string{"documentos:*|@edit"}}
	editAll := &shared.ShareTemplate{Actions: []string{"@list", "@get", "@edit"}}
	if !shared.More(view, editDocs) || shared.More(editDocs, view) {
		t.Error("view → edit the documents")
	}
	if shared.More(editAll, editDocs) || !shared.More(editDocs, editAll) {
		t.Error("edit all covers edit the documents")
	}
	if !editDocs.Allows("@get") || editDocs.Allows("@edit") {
		t.Error("Allows is of the record and all under it")
	}
	all := &shared.ShareTemplate{Actions: []string{"*"}}
	if !all.Allows("!share") || shared.More(all, editAll) || !shared.More(editAll, all) {
		t.Error("* is any action")
	}
}

// SeedRoles makes the roles the code's (by name), RenameRoles renames the
// ones of before — the accesses keep them.
func TestSeedAndRenameRoles(t *testing.T) {
	f := newFixture(t)
	if err := shared.SeedRoles(f.db, uuid.Nil, "x", shared.RoleSeed{Name: "Viewer", Actions: []string{"@get"}}); err != nil {
		t.Fatal(err)
	}
	roles, _ := shared.RolesOf(f.db, uuid.Nil, "x")
	id := roles[0].ID
	if err := shared.RenameRoles(f.db, uuid.Nil, "x", map[string]string{"Viewer": "Leitor"}); err != nil {
		t.Fatal(err)
	}
	if err := shared.SeedRoles(f.db, uuid.Nil, "x", shared.RoleSeed{Name: "Leitor", Actions: []string{"@list", "@get"},
		Rules: []shared.Rule{{On: "a:*", Do: []string{"@edit"}}}}); err != nil {
		t.Fatal(err)
	}
	if err := shared.SeedRoles(f.db, uuid.Nil, "x", shared.RoleSeed{Name: "Velho"}); err != nil {
		t.Fatal(err)
	}
	if err := shared.RemoveRoles(f.db, uuid.Nil, "x", "Velho", "Inexistente"); err != nil {
		t.Fatal(err)
	}
	roles, _ = shared.RolesOf(f.db, uuid.Nil, "x")
	if len(roles) != 2 || roles[0].ID != id || roles[0].Name != "Leitor" || len(roles[0].Actions) != 2 ||
		len(roles[0].Rules) != 1 || !roles[1].None {
		t.Fatalf("the roles: %+v", roles)
	}
}

// Grant: at once, in the history; ConfirmRequest: the owner confirming a
// request by its id (not another).
func TestGrantAndConfirmRequest(t *testing.T) {
	f := newFixture(t)
	p := sharedtest.NewPasta(f.db, "Contratos", "ana")
	if err := f.s.Grant(f.ctx, "pastas", p.ID, shared.UserSubject("bia"), f.role("pastas", "Viewer"), "ana", origin); err != nil {
		t.Fatal(err)
	}
	if f.roleName("bia", "pastas", p.ID) != "Viewer" {
		t.Error("granted")
	}
	req, err := f.s.Share(f.ctx, "pastas", p.ID, shared.UserSubject("caio"), f.role("pastas", "Viewer"), "ana", origin)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.s.ConfirmRequest(f.ctx, req.ID, "bia", origin); !errors.Is(err, shared.ErrNotOwner) {
		t.Errorf("another confirmed: %v", err)
	}
	if err := f.s.ConfirmRequest(f.ctx, req.ID, "ana", origin); err != nil {
		t.Fatal(err)
	}
	if invs, _ := f.s.Invites(f.ctx, "caio"); len(invs) != 1 {
		t.Errorf("the invite: %v", invs)
	}
}
