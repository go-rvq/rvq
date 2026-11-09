package shared

import (
	"context"
	"slices"
	"testing"

	"github.com/go-rvq/rvq/admin/packages/shared/sharedtest"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The records shared under a parent reached only as support (what Related
// narrows a supporting model to); none — all — for whom reaches the parent
// itself.
func TestSupportedFor(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:supported_for?mode=memory&cache=shared"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
		t.Cleanup(func() { _ = sqlDB.Close() })
	}
	if err := db.AutoMigrate(&perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	if err := sharedtest.Migrate(db); err != nil {
		t.Fatal(err)
	}
	users := sharedtest.NewUsers()
	for _, id := range []string{"ana", "bia", "caio"} {
		users.Add(&sharedtest.User{ID: id, Name: id})
	}
	s, err := New(Options{DB: db, Users: users})
	if err != nil {
		t.Fatal(err)
	}
	p := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(db)).URIPrefix("/admin")
	pastas := p.Model(&sharedtest.Pasta{}).URIName("pastas")
	docs := presets.NewModelBuilder(p, &sharedtest.Documento{})
	pastas.AddChild(docs)
	viewer := []RoleSeed{{Name: "Viewer", Actions: []string{presets.PermList, presets.PermGet}}}
	if err := s.Install(pastas, ModelOptions{Resource: Resource{Name: "pastas", Table: "sharedtest_pastas",
		OwnerColumn: "owner_id"}, Roles: viewer, ReadScope: true}); err != nil {
		t.Fatal(err)
	}
	if err := s.Install(docs, ModelOptions{Resource: Resource{Name: "docs", Table: "sharedtest_documentos",
		Parent: "pastas", ParentColumn: "pasta_id"}, Roles: viewer, ReadScope: true}); err != nil {
		t.Fatal(err)
	}
	role := func(res string) uuid.UUID {
		rs, _ := RolesOf(db, uuid.Nil, res)
		for _, r := range rs {
			if r.Name == "Viewer" {
				return r.ID
			}
		}
		t.Fatal("no Viewer")
		return uuid.Nil
	}
	pasta := sharedtest.NewPasta(db, "Contratos", "ana")
	d1 := sharedtest.NewDocumento(db, pasta, "Aluguel")
	sharedtest.NewDocumento(db, pasta, "Luz")
	if err := s.grant(db, "docs", d1.ID, UserSubject("caio"), role("docs"), false, "ana", Origin{}); err != nil {
		t.Fatal(err)
	}
	if err := s.grant(db, "pastas", pasta.ID, UserSubject("bia"), role("pastas"), false, "ana", Origin{}); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	// caio reaches the folder as support: the document shared with them
	if ids, ok, err := s.supportedFor(ctx, "caio", pastas, pasta.ID); err != nil || !ok || !slices.Equal(ids, []string{d1.ID}) {
		t.Errorf("caio: %v %v %v", ids, ok, err)
	}
	// bia (the folder shared) and ana (its owner): all
	for _, u := range []string{"bia", "ana"} {
		if _, ok, err := s.supportedFor(ctx, u, pastas, pasta.ID); err != nil || ok {
			t.Errorf("%s: as support (%v)", u, err)
		}
	}
}
