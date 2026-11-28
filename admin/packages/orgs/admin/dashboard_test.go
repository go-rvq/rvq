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
	d.Panel(dashboard.PanelHome).AllGroups()
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

	// the page of a widget: under the organization's page; not the global
	// one (the record of a query)
	w = httptest.NewRecorder()
	b.ServeHTTP(w, as(httptest.NewRequest("GET", page+"/widget?id=orgs/probe", nil)))
	if body := w.Body.String(); w.Code != 200 || !strings.Contains(body, "data-dashboard-widget-page") {
		t.Errorf("the page of the widget under the organization: %d %.300s", w.Code, body)
	}
	w = httptest.NewRecorder()
	b.ServeHTTP(w, as(httptest.NewRequest("GET", "/admin/dashboard/widget?id=orgs/probe&panel=org&record="+org.ID.String(), nil)))
	if strings.Contains(w.Body.String(), "data-dashboard-widget-page") {
		t.Error("the global page of a widget of an organization")
	}
	// the home takes no widget of an organization
	home := d.GetPanel(dashboard.PanelHome)
	for _, w := range home.Widgets() {
		if w.Group == GroupOrg {
			t.Errorf("the home has %s", w.ID)
		}
	}

	// another's organization (the user does not read it): nothing of it,
	// at its page or by the query of an event elsewhere (at the page of
	// theirs, the page's own wins)
	other := &models.Organizacao{Nome: "Beta Ltda", OwnerID: model.UUID(uuid.New())}
	db.Create(other)
	for _, url := range []string{
		"/admin/orgs/" + other.ID.String() + PainelPath + "?panel=org&widget=orgs/probe",
		page + "?panel=org&widget=orgs/probe&record=" + other.ID.String(),
		"/admin?panel=org&widget=orgs/probe&record=" + other.ID.String(),
	} {
		scope = dashboard.Scope{}
		r := multipartestutils.NewMultipartBuilder().PageURL(url).EventFunc(dashboard.EventLoad).BuildEventFuncRequest()
		w = httptest.NewRecorder()
		b.ServeHTTP(w, as(r))
		if scope.Record == other.ID.String() || strings.Contains(url, "/admin?") && strings.Contains(w.Body.String(), "probed") {
			t.Errorf("%s: another's organization loaded (%d)", url, w.Code)
		}
	}
	w = httptest.NewRecorder()
	b.ServeHTTP(w, as(httptest.NewRequest("GET", "/admin/orgs/"+other.ID.String()+PainelPath, nil)))
	if strings.Contains(w.Body.String(), "data-dashboard-widget") {
		t.Error("the page of another's organization shows its widgets")
	}
}
