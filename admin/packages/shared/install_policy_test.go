package shared_test

import (
	"testing"

	"github.com/go-rvq/rvq/admin/packages/shared"
	"github.com/go-rvq/rvq/admin/packages/shared/sharedtest"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
)

// A model shared in the admin with its listing scoped (ReadScope): who is
// given a record lists the model too — only that, and its columns — and
// loses them with the access; installed again, the accesses of before get
// them.
func TestInstallGrantsTheScopedListing(t *testing.T) {
	f := newFixture(t)
	if err := f.db.AutoMigrate(&perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	const res = "admin_pastas"
	install := func() {
		t.Helper()
		p := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(f.db)).URIPrefix("/admin")
		p.Permission(perm.New().AllowAll())
		mb := p.Model(&sharedtest.Pasta{}).URIName("pastas")
		if err := f.s.Install(mb, shared.ModelOptions{
			Resource:  shared.Resource{Name: res, Table: "sharedtest_pastas", OwnerColumn: "owner_id"},
			Roles:     []shared.RoleSeed{{Name: "Viewer", Actions: []string{presets.PermList, presets.PermGet}}},
			ReadScope: true,
		}); err != nil {
			t.Fatal(err)
		}
	}
	install()
	roles, err := shared.RolesOf(f.db, uuid.Nil, res)
	if err != nil {
		t.Fatal(err)
	}
	f.roles[res] = map[string]uuid.UUID{}
	for _, r := range roles {
		f.roles[res][r.Name] = r.ID
	}

	pasta := sharedtest.NewPasta(f.db, "Contratos", "ana")
	f.share(res, pasta.ID, shared.UserSubject("bia"), "Viewer", "ana")

	resources := func() (out []string) {
		var ps []perm.DefaultDBPolicy
		f.db.Where("subject = ?", shared.UserSubject("bia")).Order("refer_id").Find(&ps)
		for _, p := range ps {
			out = append(out, p.Resources...)
		}
		return
	}
	got := resources()
	if len(got) != 3 || got[0] != "admin:pasta:<"+pasta.ID+">:*" || got[1] != "admin:pasta:@list" ||
		got[2] != "admin:pasta:#*" {
		t.Fatalf("the policies of bia: %q", got)
	}

	// installed again (a new version, a new boot): the accesses of before
	// get their policies as they are made now
	f.db.Where("subject = ?", shared.UserSubject("bia")).Delete(&perm.DefaultDBPolicy{})
	install()
	if got := resources(); len(got) != 3 {
		t.Fatalf("installed again: %q", got)
	}

	if err := f.s.Revoke(f.ctx, res, pasta.ID, shared.UserSubject("bia"), "ana", origin); err != nil {
		t.Fatal(err)
	}
	if got := resources(); len(got) != 0 {
		t.Errorf("revoked, still: %q", got)
	}
}
