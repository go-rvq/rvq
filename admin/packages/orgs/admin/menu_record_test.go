package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
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

// choiceProject is a model under the organization chosen in the menu (a
// project); choiceItem, a resource of a project.
type choiceProject struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizacaoID uuid.UUID `gorm:"type:uuid;index"`
	Nome          string
	Padrao        bool
}

func (choiceProject) TableName() string { return "choice_projects" }

func (p *choiceProject) String() string { return p.Nome }

func (p *choiceProject) BeforeCreate(*gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type choiceItem struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizacaoID uuid.UUID `gorm:"type:uuid;index"`
	ProjetoID     uuid.UUID `gorm:"type:uuid;index"`
	Nome          string
}

func (choiceItem) TableName() string { return "choice_items" }

func (i *choiceItem) BeforeCreate(*gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

type choiceApp struct {
	b              *presets.Builder
	db             *gorm.DB
	proj, items    *presets.ModelBuilder
	me             uuid.UUID
	org, other     *models.Organizacao
	p1, p2, pOther *choiceProject
}

func newChoiceApp(t *testing.T) *choiceApp {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&choiceProject{}, &choiceItem{}); err != nil {
		t.Fatal(err)
	}
	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.I18n().SupportLanguages(language.BrazilianPortuguese, language.English)
	b.DataOperator(gorm2op.DataOperator(db))
	b.Permission(perm.New().Policies(
		perm.PolicyFor("everyone").WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
	).SubjectsFunc(func(*http.Request) []string { return []string{"everyone"} }))
	Configure(b, db)

	a := &choiceApp{b: b, db: db, me: uuid.New()}
	a.proj = presets.NewModelBuilder(b, &choiceProject{}, presets.ModelWithID("projetos"))
	a.proj.Listing("Nome")
	a.proj.Editing("Nome")
	MountUnderOrg(OrganizacaoModel(b), a.proj, "finance")
	a.items = presets.NewModelBuilder(b, &choiceItem{}, presets.ModelWithID("itens"))
	a.items.Listing("Nome")
	a.items.Editing("Nome")
	a.proj.AddChild(a.items)
	ScopeToOrg(a.items)
	ScopeToParent(a.items, "ProjetoID", "projeto_id")
	// a page of its detail, under it too
	a.proj.Detailing("Nome")
	a.proj.Detailing().PagesRegistrator().New(presets.HttpPage("/relatorios").StringTitle("Relatórios")).
		RecordLayout(func(ctx *web.EventContext, _ model.ID, _ []model.ID, _ func() (any, error)) (r web.PageResponse, err error) {
			return
		}).Build()
	cfg := a.proj.MenuRecordSelectorConfig()
	cfg.Items = []string{"itens"}
	cfg.Pages = []string{"/relatorios"}
	cfg.DefaultConditions = presets.SQLConditions{{Query: "padrao"}}
	a.proj.MenuRecordSelector(true)
	b.Build()

	a.org = &models.Organizacao{Nome: "Alpha", OwnerID: model.UUID(a.me)}
	a.other = &models.Organizacao{Nome: "Beta", OwnerID: model.UUID(a.me)}
	db.Create(a.org)
	db.Create(a.other)
	a.p1 = &choiceProject{OrganizacaoID: a.org.ID, Nome: "Obra 1"}
	a.p2 = &choiceProject{OrganizacaoID: a.org.ID, Nome: "Obra 2"}
	a.pOther = &choiceProject{OrganizacaoID: a.other.ID, Nome: "Obra B"}
	db.Create(a.p1)
	db.Create(a.p2)
	db.Create(a.pOther)
	db.Create(&choiceItem{OrganizacaoID: a.org.ID, ProjetoID: a.p1.ID, Nome: "Item da obra 1"})
	db.Create(&choiceItem{OrganizacaoID: a.org.ID, ProjetoID: a.p2.ID, Nome: "Item da obra 2"})
	return a
}

func (a *choiceApp) mid(t *testing.T, mb *presets.ModelBuilder, id uuid.UUID) model.ID {
	t.Helper()
	m, err := mb.ParseRecordID(id.String())
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// itemsHref is the listing of the items of the project p.
func (a *choiceApp) itemsHref(t *testing.T, p *choiceProject) string {
	return a.items.Info().ListingHref(a.mid(t, OrganizacaoModel(a.b), a.org.ID), a.mid(t, a.proj, p.ID))
}

// menu is the node of the projects in the menu of the page path (cookie, the
// record kept in the browser).
func (a *choiceApp) menu(t *testing.T, path, cookie string) *presets.MenuNode {
	t.Helper()
	r := httptest.NewRequest("GET", path, nil)
	r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: a.me}))
	if cookie != "" {
		r.AddCookie(&http.Cookie{Name: "rvq_menu_projetos", Value: cookie})
	}
	nodes, _, _ := orgMenuRoot(a.b, a.db)(&web.EventContext{R: r})
	if len(nodes) != 1 || len(nodes[0].Children) == 0 {
		t.Fatalf("the menu of the organization: %+v", nodes)
	}
	// the item of the projects (in their group, finance): a group, first the
	// selector, then the listing
	var find func(ns []*presets.MenuNode) *presets.MenuNode
	find = func(ns []*presets.MenuNode) *presets.MenuNode {
		for _, n := range ns {
			if strings.HasPrefix(n.Value, "g:menu-records:projetos:") {
				return n
			}
			if f := find(n.Children); f != nil {
				return f
			}
		}
		return nil
	}
	g := find(nodes[0].Children)
	if g == nil || len(g.Children) != 2 || g.Children[1].Value != a.proj.Info().ListingHref(a.mid(t, OrganizacaoModel(a.b), a.org.ID)) {
		t.Fatalf("the item of the projects: %+v", g)
	}
	return g.Children[0]
}

func chipOf(n *presets.MenuNode) *presets.MenuChip { return n.Props["chip"].(*presets.MenuChip) }

// A record chosen in the menu (presets MenuRecordSelector), of a model under
// the organization: none, the chip to choose; the page's, else the browser's
// (of its organization), else the default; its resources under it.
func TestMenuRecordSelectorUnderOrg(t *testing.T) {
	a := newChoiceApp(t)
	orgPage := "/admin/orgs/" + a.org.ID.String()

	// two projects, no default: none chosen
	n := a.menu(t, orgPage, "")
	// none chosen: the chip "Select" the title too, no clear
	if c := chipOf(n); c.Title != presets.Messages_en_US.MenuRecordSelect || !c.Inline || c.Clear != nil || len(n.Children) != 0 {
		t.Errorf("none chosen: %+v %+v", c, n.Children)
	}
	if c := chipOf(n); c.Event != actions.OpenActionDialog || !strings.Contains(c.Query, presets.ActionChooseMenuRecord) {
		t.Errorf("the chip's action: %+v", c)
	}

	// the default
	a.db.Model(a.p1).Update("padrao", true)
	n = a.menu(t, orgPage, "")
	if c := chipOf(n); c.Title != "Obra 1" || len(n.Children) != 2 || n.Children[0].Value != a.itemsHref(t, a.p1) {
		t.Errorf("the default: %+v %+v", c, n.Children)
	}
	// the page of its detail, after its resources
	if want := a.proj.Info().DetailingHref(a.p1.ID.String(), a.mid(t, OrganizacaoModel(a.b), a.org.ID)) + "/relatorios"; len(n.Children) == 2 &&
		(n.Children[1].Value != want || n.Children[1].Title != "Relatórios") {
		t.Errorf("the page: %+v, want %s", n.Children[1], want)
	}
	// chosen: the chip the title of the item (no label of the model), the
	// clear at its right
	if c := chipOf(n); !c.Inline || c.Clear == nil || c.Clear.Event == "" || n.Title != "Obra 1" {
		t.Errorf("chosen: the chip inline, its clear: %+v %+v", n.Title, c)
	}
	// cleared (the browser's "none"): none, not the default
	n = a.menu(t, orgPage, a.mid(t, OrganizacaoModel(a.b), a.org.ID).String()+"|")
	if c := chipOf(n); c.Title != presets.Messages_en_US.MenuRecordSelect || !c.Inline || c.Clear != nil || len(n.Children) != 0 {
		t.Errorf("cleared: %+v %+v", c, n.Children)
	}
	// the browser's
	n = a.menu(t, orgPage, a.mid(t, OrganizacaoModel(a.b), a.org.ID).String()+"|"+a.p2.ID.String())
	if chipOf(n).Title != "Obra 2" {
		t.Errorf("the browser's: %+v", chipOf(n))
	}
	// the browser's of another organization: not this one's
	n = a.menu(t, orgPage, a.mid(t, OrganizacaoModel(a.b), a.other.ID).String()+"|"+a.pOther.ID.String())
	if chipOf(n).Title != "Obra 1" {
		t.Errorf("another organization's: %+v", chipOf(n))
	}
	// the page's, over the browser's
	n = a.menu(t, a.itemsHref(t, a.p2), a.mid(t, OrganizacaoModel(a.b), a.org.ID).String()+"|"+a.p1.ID.String())
	if chipOf(n).Title != "Obra 2" || n.Children[0].Value != a.itemsHref(t, a.p2) {
		t.Errorf("the page's: %+v %+v", chipOf(n), n.Children)
	}
}

// A resource of a project: its listing, the project's only; a record saved,
// of the project of its path.
func TestScopeToParent(t *testing.T) {
	a := newChoiceApp(t)
	serve := func(r *http.Request) string {
		r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: a.me}))
		w := httptest.NewRecorder()
		a.b.ServeHTTP(w, r)
		return w.Body.String()
	}
	body := serve(multipartestutils.NewMultipartBuilder().PageURL(a.itemsHref(t, a.p1)).
		EventFunc(actions.ReloadList).BuildEventFuncRequest())
	if !strings.Contains(body, "Item da obra 1") || strings.Contains(body, "Item da obra 2") {
		t.Errorf("the items of the project 1: %.600s", body)
	}

	body = serve(multipartestutils.NewMultipartBuilder().PageURL(a.itemsHref(t, a.p2)).
		EventFunc(actions.Create).AddField("Nome", "Novo").BuildEventFuncRequest())
	var got choiceItem
	if err := a.db.Where("nome = ?", "Novo").Take(&got).Error; err != nil {
		t.Fatalf("%v: %.600s", err, body)
	}
	if got.ProjetoID != a.p2.ID || got.OrganizacaoID != a.org.ID {
		t.Errorf("saved of %v / %v, want the project 2 of Alpha", got.ProjetoID, got.OrganizacaoID)
	}
}

// The trail of a page of a resource of a project: every step — the
// organization, its record, the menu group of the projects (finance), the
// projects, the project — and, of a page of the project's detail, the
// project too.
func TestTrailOfAProject(t *testing.T) {
	a := newChoiceApp(t)
	orgMID := a.mid(t, OrganizacaoModel(a.b), a.org.ID)
	// the text of the trail of the page path
	crumbs := func(path string) string {
		t.Helper()
		r := httptest.NewRequest("GET", path, nil)
		r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: a.me}))
		w := httptest.NewRecorder()
		a.b.ServeHTTP(w, r)
		body := w.Body.String()
		i := strings.Index(body, "rvq-presets-breadcrumbs")
		if i < 0 {
			t.Fatalf("%s: no trail (%d): %.300s", path, w.Code, body)
		}
		nav := body[i:]
		if e := strings.Index(nav, "</nav>"); e > 0 {
			nav = nav[:e]
		}
		return regexp.MustCompile(`<[^>]*>|\s+`).ReplaceAllString(nav, " ")
	}
	trail := crumbs(a.itemsHref(t, a.p1))
	for _, want := range []string{"Alpha", "Finance", "Obra 1"} {
		if !strings.Contains(trail, want) {
			t.Errorf("the trail of the items: no %q: %s", want, trail)
		}
	}
	if strings.Index(trail, "Alpha") > strings.Index(trail, "Finance") || strings.Index(trail, "Finance") > strings.Index(trail, "Obra 1") {
		t.Errorf("the trail out of order: %s", trail)
	}
	trail = crumbs(a.proj.Info().DetailingHref(a.p1.ID.String(), orgMID) + "/relatorios")
	if !strings.Contains(trail, "Finance") || !strings.Contains(trail, "Obra 1") {
		t.Errorf("the trail of a page of the project: %s", trail)
	}
}

// The organization in the menu is as any record chosen in the menu: at the
// top, a group — the selector (the chip of the organization chosen, in the
// place of its title, its clear; its resources under it in their menu
// groups), then the listing of the organizations.
func TestOrgInMenu(t *testing.T) {
	a := newChoiceApp(t)
	r := httptest.NewRequest("GET", "/admin/orgs/"+a.org.ID.String(), nil)
	r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: a.me}))
	nodes, opened, _ := orgMenuRoot(a.b, a.db)(&web.EventContext{R: r})
	if len(nodes) != 1 || !strings.HasPrefix(nodes[0].Value, "g:menu-records:") || len(nodes[0].Children) != 2 {
		t.Fatalf("the item of the organization: %+v", nodes)
	}
	sel, list := nodes[0].Children[0], nodes[0].Children[1]
	if c := chipOf(sel); c.Title != "Alpha" || !c.Inline || c.Clear == nil {
		t.Errorf("the organization chosen: %+v", c)
	}
	if list.Value != OrganizacaoModel(a.b).Info().ListingHref() {
		t.Errorf("the listing of the organizations: %+v", list)
	}
	// its resources in their menu group (finance)
	if len(sel.Children) != 1 || !strings.HasSuffix(sel.Children[0].Value, ":finance") ||
		!strings.HasPrefix(sel.Children[0].Children[0].Value, "g:menu-records:projetos:") {
		t.Errorf("the resources of the organization: %+v", sel.Children)
	}
	if !slices.Contains(opened, sel.Value) {
		t.Errorf("the organization chosen not open: %v", opened)
	}
	// a user with no organization: no item
	r = httptest.NewRequest("GET", "/admin", nil)
	r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: uuid.New()}))
	if nodes, _, _ := orgMenuRoot(a.b, a.db)(&web.EventContext{R: r}); len(nodes) != 0 {
		t.Errorf("no organization, an item: %+v", nodes)
	}
}
