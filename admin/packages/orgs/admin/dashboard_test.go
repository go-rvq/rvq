package admin

import (
	"context"
	"html"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/dashboard"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// The panel of an organization: a page of its detail (PainelPath), its
// group's widgets — another package's too (GetGroup) — loaded for the
// organization of the page (dashboard.ScopeOf).
func TestOrgPanel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.I18n().SupportLanguages(language.BrazilianPortuguese, language.English)
	b.DataOperator(gorm2op.DataOperator(db))
	b.Permission(perm.New().Policies(
		perm.PolicyFor("everyone").WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
	).SubjectsFunc(func(*http.Request) []string { return []string{"everyone"} }))
	ob := New(db)
	b.Use(ob)
	d := dashboard.New(b, db)
	if err := d.Install(); err != nil {
		t.Fatal(err)
	}
	// another package's widget in the group, before the group is made
	var scope dashboard.Scope
	d.GetGroup(GroupOrg).Widget("probe", func(w *dashboard.WidgetBuilder) {
		w.Title(func(context.Context) string { return "Probe" }).
			Load(func(ctx *web.EventContext, _ dashboard.Period) ([]dashboard.Block, error) {
				scope = dashboard.ScopeOf(ctx)
				return []dashboard.Block{&dashboard.Note{Text: "probed"}}, nil
			})
	})
	pn := ob.Dashboard(d)
	if pn.ID() != PanelOrg || pn.RecordPage() == nil {
		t.Fatal("the panel of the organizations")
	}
	b.Build()

	me := uuid.New()
	org := &models.Organizacao{Nome: "Alpha SA", OwnerID: model.UUID(me)}
	db.Create(org)
	as := func(r *http.Request) *http.Request {
		return r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: me}))
	}
	page := "/admin/orgs/" + org.ID.String() + PainelPath
	w := httptest.NewRecorder()
	b.ServeHTTP(w, as(httptest.NewRequest("GET", page, nil)))
	// the page's body is in a template, its JSON escaped
	body := html.UnescapeString(strings.NewReplacer("\\u003c", "<", "\\u003e", ">", "\\u0026", "&").Replace(w.Body.String()))
	if w.Code != 200 || !strings.Contains(body, "data-dashboard='org'") ||
		!strings.Contains(body, "data-dashboard-widget='orgs/probe'") {
		t.Fatalf("the page of the panel: %d %.600s", w.Code, body)
	}

	r := multipartestutils.NewMultipartBuilder().PageURL(page + "?panel=org&widget=orgs/probe").
		EventFunc(dashboard.EventLoad).BuildEventFuncRequest()
	w = httptest.NewRecorder()
	b.ServeHTTP(w, as(r))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "probed") {
		t.Fatalf("the load of the widget: %d %.400s", w.Code, w.Body.String())
	}
	if scope.Panel != PanelOrg || scope.Record != org.ID.String() {
		t.Errorf("the scope of the widget: %+v", scope)
	}
}
