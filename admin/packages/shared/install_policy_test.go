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

// A model under another: None on one of its records denies the subject it
// and what is under it — whatever its parent's role allows —; revoked, the
// deny goes.
func TestInstallNoneUnderAParent(t *testing.T) {
	f := newFixture(t)
	if err := f.db.AutoMigrate(&perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	p := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(f.db)).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	pastas := p.Model(&sharedtest.Pasta{}).URIName("pastas")
	docs := presets.NewModelBuilder(p, &sharedtest.Documento{})
	pastas.AddChild(docs)
	viewer := []shared.RoleSeed{{Name: "Viewer", Actions: []string{presets.PermList, presets.PermGet}}}
	if err := f.s.Install(pastas, shared.ModelOptions{
		Resource: shared.Resource{Name: "admin_pastas", Table: "sharedtest_pastas", OwnerColumn: "owner_id"},
		Roles:    viewer, ReadScope: true,
	}); err != nil {
		t.Fatal(err)
	}
	if err := f.s.Install(docs, shared.ModelOptions{
		Resource: shared.Resource{Name: "admin_docs", Table: "sharedtest_documentos", Parent: "admin_pastas",
			ParentColumn: "pasta_id"},
		Roles: viewer, ReadScope: true,
	}); err != nil {
		t.Fatal(err)
	}
	for _, res := range []string{"admin_pastas", "admin_docs"} {
		roles, _ := shared.RolesOf(f.db, uuid.Nil, res)
		f.roles[res] = map[string]uuid.UUID{}
		for _, r := range roles {
			f.roles[res][r.Name] = r.ID
		}
	}
	pasta := sharedtest.NewPasta(f.db, "Contratos", "ana")
	doc := sharedtest.NewDocumento(f.db, pasta, "Aluguel")
	f.share("admin_pastas", pasta.ID, shared.UserSubject("bia"), "Viewer", "ana")
	f.share("admin_docs", doc.ID, shared.UserSubject("bia"), shared.RoleNone, "ana")

	var denies []perm.DefaultDBPolicy
	f.db.Where("subject = ? AND effect = ?", shared.UserSubject("bia"), perm.Denied).Find(&denies)
	want := "admin:pasta:<" + pasta.ID + ">:documentos:<" + doc.ID + ">:*"
	if len(denies) != 1 || len(denies[0].Resources) != 1 || denies[0].Resources[0] != want {
		t.Fatalf("the deny of None: %+v, want %s", denies, want)
	}
	if err := f.s.Revoke(f.ctx, "admin_docs", doc.ID, shared.UserSubject("bia"), "ana", origin); err != nil {
		t.Fatal(err)
	}
	var n int64
	f.db.Model(&perm.DefaultDBPolicy{}).Where("subject = ? AND effect = ?", shared.UserSubject("bia"), perm.Denied).Count(&n)
	if n != 0 {
		t.Errorf("revoked, the deny stays: %d", n)
	}
}
