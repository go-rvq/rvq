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

// The page of a widget of records in the catalog: at its top, the kind of
// scope (Consolidated, Organization, Sample), the selects of its record
// cascaded by its parents — the organization the one chosen in the menu at
// first, the sample under it —, its preview loaded at the page of the
// record (its events there, at=1); a record the user does not read, not
// chosen (no preview).
func TestCatalogRecords(t *testing.T) {
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
	// "limited" (the header X-Subject) may not list the samples
	b.Permission(perm.New().Policies(
		perm.PolicyFor("everyone").WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
		perm.PolicyFor("limited").WhoAre(perm.Denied).ToDo(perm.Anything).On("*:samples:*"),
		perm.PolicyFor("limited").WhoAre(perm.Denied).ToDo(perm.Anything).On("*:samples"),
	).SubjectsFunc(func(r *http.Request) []string {
		if s := r.Header.Get("X-Subject"); s != "" {
			return []string{"everyone", s}
		}
		return []string{"everyone"}
	}))
	ob := New(db)
	b.Use(ob)
	orgMB := OrganizacaoModel(b)
	child := presets.NewModelBuilder(b, &sampleModel{}, presets.ModelWithID("samples"))
	child.Listing("Nome")
	MountUnderOrg(orgMB, child, "finance")
	d := dashboard.New(b, db)
	if err := d.Install(); err != nil {
		t.Fatal(err)
	}
	ui := b.MenuGroup("ui")
	d.Catalog(ui)
	ob.Dashboard(d)
	load := func(ctx *web.EventContext, _ dashboard.Period) ([]dashboard.Block, error) {
		return []dashboard.Block{&dashboard.Note{Text: "of " + dashboard.ScopeOf(ctx).Record}}, nil
	}
	fin := d.Group("fin").Global().Records(orgMB, child).Widget("saldo", func(w *dashboard.WidgetBuilder) {
		w.Title(func(context.Context) string { return "Saldo" }).Load(load)
	}).Widget("secret", func(w *dashboard.WidgetBuilder) {
		w.Title(func(context.Context) string { return "Secret" }).Load(load).
			Allowed(func(*http.Request) bool { return false })
	})
	d.GetPanel(PanelOrg).Group(fin)
	d.Panel("sample").Group(fin).Record(child, "/painel")
	d.Panel(dashboard.PanelHome).AllGroups()
	b.Build()

	me := uuid.New()
	org := &models.Organizacao{Nome: "Alpha SA", OwnerID: model.UUID(me)}
	other := &models.Organizacao{Nome: "Beta Ltda", OwnerID: model.UUID(uuid.New())}
	db.Create(org)
	db.Create(other)
	s := &sampleModel{OrganizacaoID: org.ID, Nome: "Obra"}
	db.Create(s)
	get := func(url string) string {
		t.Helper()
		r := httptest.NewRequest("GET", url, nil)
		r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: me}))
		w := httptest.NewRecorder()
		b.ServeHTTP(w, r)
		return html.UnescapeString(strings.NewReplacer("\\u003c", "<", "\\u003e", ">", "\\u0026", "&", `\"`, `"`).Replace(w.Body.String()))
	}
	page := "/admin/ui/widgets/widgets-fin/saldo"

	// the consolidated: its preview at the catalog's page
	body := get(page)
	for _, want := range []string{"data-dashboard-catalog-scope='global'", "data-dashboard-catalog-scope='orgs'",
		"data-dashboard-catalog-scope='samples'", "data-dashboard-widget='fin/saldo'"} {
		if !strings.Contains(body, want) {
			t.Errorf("consolidated: no %s", want)
		}
	}
	// of an organization: its select, the one of the menu at first (the only
	// one the user reads), the preview's events at its page
	body = get(page + "?scope=orgs")
	orgPage := "/admin/orgs/" + org.ID.String() + PainelPath
	if !strings.Contains(body, "data-dashboard-catalog-record='orgs'") || !strings.Contains(body, `.url("`+orgPage+`")`) ||
		!strings.Contains(body, `.query("at", "1")`) {
		t.Errorf("of an organization: %.2000s", body)
	}
	// another's: not chosen, no preview
	body = get(page + "?scope=orgs&rec=" + other.ID.String())
	if strings.Contains(body, "data-dashboard-widget='fin/saldo'") || !strings.Contains(body, "data-dashboard-catalog-choose") {
		t.Error("another's organization chosen")
	}
	// of a sample: the organization, then the sample; the preview at the
	// sample's page under the organization
	body = get(page + "?scope=samples&rec=" + org.ID.String())
	if !strings.Contains(body, "data-dashboard-catalog-record='samples'") || !strings.Contains(body, "data-dashboard-catalog-choose") ||
		strings.Contains(body, "data-dashboard-widget='fin/saldo'") {
		t.Error("of a sample: the organization chosen, the sample asked")
	}
	body = get(page + "?scope=samples&rec=" + org.ID.String() + "/" + s.ID.String())
	samplePage := "/admin/orgs/" + org.ID.String() + "/finance/samples/" + s.ID.String() + "/painel"
	if !strings.Contains(body, `.url("`+samplePage+`")`) {
		t.Errorf("of a sample: its page %s not in %.2000s", samplePage, body)
	}
	// in its dialog
	if !strings.Contains(body, "dialog") {
		t.Error("the select of a record not in its dialog")
	}
	// one who may not list the samples: no select of them, no widget
	r0 := httptest.NewRequest("GET", page+"?scope=samples&rec="+org.ID.String()+"/"+s.ID.String(), nil)
	r0.Header.Set("X-Subject", "limited")
	if !child.Permissioner().Lister(r0, model.ID{}).Denied() && !child.Permissioner().Lister(r0).Denied() {
		t.Fatal("the policy does not deny the listing of the samples")
	}
	r0 = r0.WithContext(context.WithValue(r0.Context(), login.UserKey, &fakeUser{id: me}))
	w0 := httptest.NewRecorder()
	b.ServeHTTP(w0, r0)
	if lb := html.UnescapeString(w0.Body.String()); strings.Contains(lb, "data-dashboard-catalog-record='samples'") ||
		strings.Contains(lb, "data-dashboard-widget='fin/saldo'") || !strings.Contains(lb, "data-dashboard-catalog-no-access") {
		t.Error("one who may not list the samples chooses one, or sees the widget")
	}

	// a widget not the user's: not drawn, nothing to choose
	for _, q := range []string{"", "?scope=orgs&rec=" + org.ID.String()} {
		if body := get("/admin/ui/widgets/widgets-fin/secret" + q); strings.Contains(body, "data-dashboard-widget='fin/secret'") ||
			strings.Contains(body, "data-dashboard-catalog-record") || !strings.Contains(body, "data-dashboard-catalog-no-access") {
			t.Errorf("a widget not the user's%s: drawn", q)
		}
	}

	// its load at the sample's page, of it
	r := multipartestutils.NewMultipartBuilder().PageURL(samplePage + "?panel=sample&widget=fin/saldo&at=1").
		EventFunc(dashboard.EventLoad).BuildEventFuncRequest()
	r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: me}))
	w := httptest.NewRecorder()
	b.ServeHTTP(w, r)
	if !strings.Contains(w.Body.String(), "of "+s.ID.String()) {
		t.Errorf("the load at the sample's page: %d %.300s", w.Code, w.Body.String())
	}
}
