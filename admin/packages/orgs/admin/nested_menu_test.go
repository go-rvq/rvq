package admin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func renderComp(t *testing.T, comp h.HTMLComponent, ctx *web.EventContext) string {
	t.Helper()
	var buf bytes.Buffer
	if err := h.Fprint(&buf, comp, ctx.Context()); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// nestedMenuAdmin mounts a sample resource under the org with the given
// permission builder, and returns the builder + a real organization id.
func nestedMenuAdmin(t *testing.T, pb *perm.Builder) (*presets.Builder, *gorm.DB, string) {
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
	b.Permission(pb)
	Configure(b, db)
	child := presets.NewModelBuilder(b, &sampleModel{}, presets.ModelWithID("samples"))
	child.Listing("Nome")
	MountUnderOrg(OrganizacaoModel(b), child, "finance")
	b.Build()

	org := &models.Organizacao{Nome: "Alpha", OwnerID: uuid.New()}
	if err := db.Create(org).Error; err != nil {
		t.Fatal(err)
	}
	return b, db, org.ID.String()
}

func TestNestedMenuContent(t *testing.T) {
	// permissive permission: the sample resource is listed
	allow := perm.New().Policies(
		perm.PolicyFor("everyone").WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
	).SubjectsFunc(func(*http.Request) []string { return []string{"everyone"} })

	b, db, orgID := nestedMenuAdmin(t, allow)
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/admin", nil)}

	html := renderComp(t, NestedMenuContent(b, db, ctx, orgID), ctx)
	if !strings.Contains(html, "/admin/orgs/"+orgID+"/finance/samples") {
		t.Errorf("permitted menu should link to the org's finance/samples; got:\n%s", html)
	}

	// no org selected -> empty
	if strings.Contains(renderComp(t, NestedMenuContent(b, db, ctx, ""), ctx), "samples") {
		t.Errorf("no org selected should render no resources")
	}
}

func TestNestedMenuContentHidesForbidden(t *testing.T) {
	// deny-by-default permission (no allow policy): the user may not list samples
	deny := perm.New().SubjectsFunc(func(*http.Request) []string { return []string{"visitor"} })

	b, db, orgID := nestedMenuAdmin(t, deny)
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/admin", nil)}

	html := renderComp(t, NestedMenuContent(b, db, ctx, orgID), ctx)
	if strings.Contains(html, "samples") {
		t.Errorf("a forbidden resource must not appear in the nested menu; got:\n%s", html)
	}
}
