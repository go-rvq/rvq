package shared_test

import (
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

// nota is a model of the folders not shared by itself: a supporting one.
type nota struct {
	ID      string `gorm:"primaryKey"`
	PastaID string
	Texto   string
	Segredo string
}

func (nota) TableName() string { return "sharedtest_notas" }

// anexo is a model of the folders declared nowhere: not read as support.
type anexo struct {
	ID      string `gorm:"primaryKey"`
	PastaID string
}

func (anexo) TableName() string { return "sharedtest_anexos" }

// A role on a record under another: its actions there, the listing of its
// model under the parent; its parent read — the parent's detail and
// listing, its fields declared, its supporting models declared (Support:
// their fields and children declared; not the shared ones, not the ones
// declared nowhere) —, unless NoParentRead.
func TestInstallParentRead(t *testing.T) {
	f := newFixture(t)
	if err := f.db.AutoMigrate(&perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	p := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(f.db)).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	pastas := p.Model(&sharedtest.Pasta{}).URIName("pastas")
	docs := presets.NewModelBuilder(p, &sharedtest.Documento{})
	notas := presets.NewModelBuilder(p, &nota{})
	pastas.AddChild(docs)
	pastas.AddChild(notas)
	pastas.AddChild(presets.NewModelBuilder(p, &anexo{}))
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
	f.s.Support(pastas, shared.SupportOptions{Fields: []string{"Nome"}})
	f.s.Support(notas, shared.SupportOptions{Fields: []string{"Texto"}, Children: []string{"revisoes"}})
	pasta := sharedtest.NewPasta(f.db, "Contratos", "ana")
	doc := sharedtest.NewDocumento(f.db, pasta, "Aluguel")
	resources := func(user string) (out []string) {
		var ps []perm.DefaultDBPolicy
		f.db.Where("subject = ?", shared.UserSubject(user)).Find(&ps)
		for _, p := range ps {
			out = append(out, p.Resources...)
		}
		slices.Sort(out)
		return
	}
	f.share("admin_docs", doc.ID, shared.UserSubject("caio"), "Viewer", "ana")
	pr := "admin:pasta:<" + pasta.ID + ">:"
	want := []string{
		"admin:pasta:#Nome:*", "admin:pasta:@list",
		pr + "#Nome:*", pr + "@get",
		pr + "documentos:#*", pr + "documentos:<" + doc.ID + ">:*", pr + "documentos:@list",
		pr + "nota:#Texto:*", pr + "nota:<*>:#Texto:*", pr + "nota:<*>:@get", pr + "nota:<*>:revisoes:*",
		pr + "nota:@get", pr + "nota:@list",
	}
	slices.Sort(want)
	if got := resources("caio"); !slices.Equal(got, want) {
		t.Errorf("the policies of caio:\n got %q\nwant %q", got, want)
	}
	// the patterns, by the matcher of the permissions
	m := &perm.PathMatcher{}
	matches := func(needle string) bool {
		ok, err := m.Matches(nil, want, needle)
		return err == nil && ok
	}
	for needle, ok := range map[string]bool{
		pr + "nota:<n1>:#Texto:@get":             true,
		pr + "nota:#Texto:@list":                 true,
		pr + "nota:<n1>:#Segredo:@get":           false,
		pr + "nota:#Segredo:@list":               false,
		pr + "nota:<n1>:revisoes:<r1>:@get":      true,
		pr + "nota:<n1>:outros:<o1>:@get":        false,
		pr + "anexo:@list":                       false,
		pr + "#Nome:@get":                        true,
		pr + "#Owner:@get":                       false,
		pr + "documentos:<outro>:@get":           false,
		pr + "documentos:<" + doc.ID + ">:@edit": true,
	} {
		if matches(needle) != ok {
			t.Errorf("%s: matched %v", needle, !ok)
		}
	}
	// installed again (SyncPolicies), the same
	if err := f.s.SyncPolicies(); err != nil {
		t.Fatal(err)
	}
	if got := resources("caio"); !slices.Equal(got, want) {
		t.Errorf("synced again: %q", got)
	}
	// an access of a document deleted for good: synced again all the same,
	// the others kept
	gone := sharedtest.NewDocumento(f.db, pasta, "Apagado")
	f.share("admin_docs", gone.ID, shared.UserSubject("bia"), "Viewer", "ana")
	if err := f.db.Exec("DELETE FROM sharedtest_documentos WHERE id = ?", gone.ID).Error; err != nil {
		t.Fatal(err)
	}
	if err := f.s.SyncPolicies(); err != nil {
		t.Fatalf("an access of a record gone: %v", err)
	}
	if got := resources("caio"); !slices.Equal(got, want) {
		t.Errorf("synced with a record gone: %q", got)
	}
	f.share("admin_docs", doc.ID, shared.UserSubject("dani"), "Viewer", "ana", shared.NoParentRead())
	if got := resources("dani"); len(got) != 3 {
		t.Errorf("NoParentRead: %q", got)
	}
}

// None of a group on a record under another: denied to the members it
// decides for — not to one with an access of their own there —, again as
// members join and leave and as their own access changes.
func TestInstallGroupNone(t *testing.T) {
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
	f.users.AddGroup("estag", "Estagiários")
	f.users.Join("bia", "estag")
	f.users.Join("caio", "estag")
	pasta := sharedtest.NewPasta(f.db, "Contratos", "ana")
	doc := sharedtest.NewDocumento(f.db, pasta, "Aluguel")
	f.share("admin_docs", doc.ID, shared.UserSubject("bia"), "Viewer", "ana")
	f.share("admin_docs", doc.ID, shared.GroupSubject("estag"), shared.RoleNone, "ana")

	denied := func() (out []string) {
		var ps []perm.DefaultDBPolicy
		f.db.Where("effect = ?", perm.Denied).Find(&ps)
		for _, p := range ps {
			out = append(out, p.Subject)
		}
		slices.Sort(out)
		return
	}
	want := func(step string, subs ...string) {
		t.Helper()
		var w []string
		for _, u := range subs {
			w = append(w, shared.UserSubject(u))
		}
		if got := denied(); !slices.Equal(got, w) {
			t.Errorf("%s: denied %q, want %q", step, got, w)
		}
	}
	want("bia's own access before the group's None", "caio")
	f.users.Leave("caio", "estag")
	if err := f.s.MemberChanged(f.ctx, "estag", "caio", false, "ana", origin); err != nil {
		t.Fatal(err)
	}
	want("caio left")
	f.users.Join("dani", "estag")
	if err := f.s.MemberChanged(f.ctx, "estag", "dani", true, "ana", origin); err != nil {
		t.Fatal(err)
	}
	want("dani joined", "dani")
	if err := f.s.Revoke(f.ctx, "admin_docs", doc.ID, shared.UserSubject("bia"), "ana", origin); err != nil {
		t.Fatal(err)
	}
	want("bia's own access revoked", "bia", "dani")
}
