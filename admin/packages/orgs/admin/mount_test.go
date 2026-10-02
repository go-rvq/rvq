package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// sampleModel is a domain model mounted under the organization in the tests.
type sampleModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizacaoID uuid.UUID `gorm:"type:uuid;index"`
	Nome          string
}

func (sampleModel) TableName() string { return "sample_models" }

func (s *sampleModel) BeforeCreate(*gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func newMountTestAdmin(t *testing.T) (*presets.Builder, *gorm.DB, *presets.ModelBuilder) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&sampleModel{}); err != nil {
		t.Fatal(err)
	}
	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.I18n().SupportLanguages(language.BrazilianPortuguese, language.English)
	b.DataOperator(gorm2op.DataOperator(db))
	b.Permission(perm.New().Policies(
		perm.PolicyFor("everyone").WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
	).SubjectsFunc(func(*http.Request) []string { return []string{"everyone"} }))

	Configure(b, db)

	orgMb := OrganizacaoModel(b)
	child := presets.NewModelBuilder(b, &sampleModel{}, presets.ModelWithID("samples"))
	child.Listing("Nome")
	child.Editing("Nome")
	MountUnderOrg(orgMb, child, "finance")

	b.Build()
	return b, db, child
}

func TestMountUnderOrgPathAndScope(t *testing.T) {
	b, db, child := newMountTestAdmin(t)

	// the nested route sits under /admin/orgs/{parent_0_id}/finance/samples
	tmpl := child.Info().ListingHref()
	if !strings.Contains(tmpl, "/admin/orgs/") || !strings.Contains(tmpl, "/finance/samples") {
		t.Fatalf("child listing href = %q, want /admin/orgs/{id}/finance/samples", tmpl)
	}
	if !strings.Contains(tmpl, "{parent_0_id}") {
		t.Fatalf("child listing href = %q, want a {parent_0_id} segment", tmpl)
	}

	// two real organizations owned by the current user, one sample each
	me := uuid.New()
	orgA := &models.Organizacao{Nome: "Alpha SA", OwnerID: me}
	orgB := &models.Organizacao{Nome: "Beta Ltda", OwnerID: me}
	db.Create(orgA)
	db.Create(orgB)
	db.Create(&sampleModel{OrganizacaoID: orgA.ID, Nome: "Alpha Item"})
	db.Create(&sampleModel{OrganizacaoID: orgB.ID, Nome: "Beta Item"})

	hrefA := strings.Replace(tmpl, "{parent_0_id}", orgA.ID.String(), 1)
	r := multipartestutils.NewMultipartBuilder().
		PageURL(hrefA).
		EventFunc(actions.ReloadList).
		BuildEventFuncRequest()
	r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: me}))
	w := httptest.NewRecorder()
	b.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("listing under org A = %d, want 200 (%s)", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, "Alpha Item") {
		t.Errorf("org A listing should show its own Alpha Item")
	}
	if strings.Contains(body, "Beta Item") {
		t.Errorf("org A listing must not show org B's Beta Item")
	}
}
