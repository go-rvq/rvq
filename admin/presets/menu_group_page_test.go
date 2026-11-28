package presets_test

import (
	"context"
	"html"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type gpOrg struct {
	ID   uint
	Name string
}
type gpWallet struct {
	ID      uint
	GpOrgID uint
}
type gpCard struct {
	ID      uint
	GpOrgID uint
}
type gpReport struct {
	ID      uint
	GpOrgID uint
}
type gpBank struct{ ID uint }

// A group of the menu that is a page (PageMenuGroup): its name, title and
// icon the page's, none set otherwise.
func TestPageMenuGroup(t *testing.T) {
	b := presets.New(i18n.New()).URIPrefix("/admin")
	page := presets.HttpPage("/finance").TitleFunc(func(context.Context) string { return "Finanças" }).MenuIcon("mdi-finance")
	b.PagesRegistrator().AddHttpPage(page)
	g := b.PageMenuGroup(page)
	if g.Name() != "finance" || g.TTitle(context.Background()) != "Finanças" || g.GetIcon() != "mdi-finance" ||
		g.GetPage() != page || b.MenuGroup("finance") != g || b.PageMenuGroup(page) != g || page.IsInMenu() {
		t.Fatalf("the group of the page: %q %q %q", g.Name(), g.TTitle(context.Background()), g.GetIcon())
	}
	panics := func(what string, f func()) {
		t.Helper()
		defer func() {
			if recover() == nil {
				t.Errorf("%s: no panic", what)
			}
		}()
		f()
	}
	panics("Title", func() { g.Title("x") })
	panics("TitleFunc", func() { g.TitleFunc(func(context.Context) string { return "x" }) })
	panics("Icon", func() { g.Icon("mdi-x") })
	panics("URIName", func() { g.URIName("x") })
	for _, p := range []string{"/a/b", "/{id}", "/"} {
		p := p
		panics("the path "+p, func() { b.PageMenuGroup(presets.HttpPage(p)) })
	}
	// a group with a title of its own: not of a page
	b.MenuGroup("titled").Title("Titled")
	panics("a group titled", func() { b.PageMenuGroup(presets.HttpPage("/titled")) })
}

// The page of each group of the menu: of a group with no page, its index —
// its items the user sees, their icons, titles (links) and descriptions —;
// of one that is a page, the page. In the menu and in the breadcrumbs, the
// group a link to it.
func TestMenuGroupPages(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&gpOrg{}, &gpWallet{}, &gpCard{}, &gpReport{}, &gpBank{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&gpOrg{ID: 1, Name: "Acme"})
	b := presets.New(i18n.New()).URIPrefix("/admin").DataOperator(gorm2op.DataOperator(db))

	// the root: a group of a model of its own, its index
	b.Model(&gpBank{}, presets.ModelWithID("bancos")).URIName("bancos").DescriptionFunc(func(context.Context) string { return "os bancos do país" })
	b.MenuGroup("cadastros").Title("Cadastros").Add(presets.ModelItem("bancos"))

	// a record (the organization) with its resources in groups
	org := b.Model(&gpOrg{}).URIName("orgs")
	org.MenuRecordSelector(true)
	child := func(v any, uri, group, desc string) *presets.ModelBuilder {
		c := presets.NewModelBuilder(b, v).URIName(uri)
		c.DescriptionFunc(func(context.Context) string { return desc })
		org.AddChild(c)
		c.SetMenuGroupName(group)
		return c
	}
	child(&gpWallet{}, "carteiras", "fin", "as contas e os caixas")
	child(&gpCard{}, "cartoes", "fin", "os cartões de crédito")
	child(&gpReport{}, "relatorios", "painel", "os relatórios")
	b.MenuGroup("fin").Title("Financeiro")
	// a group that is a page of the organization's
	page := presets.HttpPage("/painel").TitleFunc(func(context.Context) string { return "Painel" }).MenuIcon("mdi-view-dashboard")
	org.Detailing().PagesRegistrator().New(page).
		RecordLayout(func(ctx *web.EventContext, mid model.ID, _ []model.ID, _ func() (any, error)) (r web.PageResponse, err error) {
			// with the index of another group of the organization's
			r.Body = h.Div(h.Text("o painel de "+mid.String()), b.MenuGroup("fin").RecordIndex(ctx, org)).Attr("data-painel", true)
			return
		}).Build()
	b.PageMenuGroup(page)
	b.Build()

	get := func(url string) (string, int) {
		t.Helper()
		w := httptest.NewRecorder()
		b.ServeHTTP(w, httptest.NewRequest("GET", url, nil))
		// the page's body is in a template, its JSON escaped
		return html.UnescapeString(strings.NewReplacer("\\u003c", "<", "\\u003e", ">", "\\u0026", "&", `\"`, `"`).Replace(w.Body.String())), w.Code
	}

	// the index of a group of the organization's resources
	body, code := get("/admin/orgs/1/fin")
	if code != 200 || !strings.Contains(body, `data-menu-group-index='fin'`) {
		t.Fatalf("the index of fin: %d %.400s", code, body)
	}
	for _, want := range []string{`"href":"/admin/orgs/1/fin/carteiras"`, `"href":"/admin/orgs/1/fin/cartoes"`,
		`"subtitle":"as contas e os caixas"`, `"subtitle":"os cartões de crédito"`} {
		if !strings.Contains(body, want) {
			t.Errorf("the index of fin has no %s", want)
		}
	}
	if strings.Contains(body, `"href":"/admin/orgs/1/painel/relatorios"`+`,"prependIcon"`) && strings.Contains(body, `data-menu-group-index='painel'`) {
		t.Error("another group's items in the index")
	}

	// a group that is a page: the page
	if body, code := get("/admin/orgs/1/painel"); code != 200 || !strings.Contains(body, "o painel de 1") ||
		!strings.Contains(body, `data-menu-group-index='fin'`) || !strings.Contains(body, `"href":"/admin/orgs/1/fin/carteiras"`) {
		t.Errorf("the page of the group painel, the index of fin in it (RecordIndex): %d %.300s", code, body)
	}

	// the index of a group of the root
	body, code = get("/admin/cadastros")
	if code != 200 || !strings.Contains(body, `"href":"/admin/cadastros/bancos"`) || !strings.Contains(body, `"subtitle":"os bancos do país"`) {
		t.Errorf("the index of cadastros: %d %.400s", code, body)
	}

	// in the menu: the groups' titles links to their pages; a listing of a
	// resource, its group a link in the breadcrumbs
	body, _ = get("/admin/orgs/1/fin/carteiras")
	for _, want := range []string{`"groupHref":"/admin/orgs/1/fin"`, `"groupHref":"/admin/orgs/1/painel"`,
		`"groupHref":"/admin/cadastros"`, "rvq-menu-group-link", "href='/admin/orgs/1/fin'"} {
		if !strings.Contains(body, want) {
			t.Errorf("the page of a resource has no %s", want)
		}
	}
	// the menu's nodes carry no descriptions (only an index's)
	if strings.Contains(body, `"subtitle":"as contas e os caixas"`) {
		t.Error("descriptions in the side menu")
	}
}

// A group that is a page in another group: the page at the group's path
// (under the other's), its sub-pages in the menu (InMenu) the group's
// children, at their paths under it; made before the build (BeforeBuild).
func TestNestedPageMenuGroup(t *testing.T) {
	b := presets.New(i18n.New()).URIPrefix("/admin")
	title := func(s string) func(context.Context) string { return func(context.Context) string { return s } }
	layout := func(s string) func(*web.EventContext) (web.PageResponse, error) {
		return func(*web.EventContext) (r web.PageResponse, err error) {
			r.Body = h.Div(h.Text(s))
			return
		}
	}
	b.BeforeBuild(func() {
		page := b.PagesRegistrator().New(presets.HttpPage("/widgets-fin").TitleFunc(title("Financeiro"))).Layout(layout("a página do grupo"))
		page.SubPage("/pending").InMenu(true).Title(title("Pendências")).Layout(layout("a página do widget")).Build()
		page.Build()
		g := b.PageMenuGroup(page.Page())
		b.MenuGroup("widgets").Title("Widgets").Add(presets.GroupItem(g.Name()))
		b.MenuGroup("ui").Title("UI").Add(presets.GroupItem("widgets"))
	})
	b.Build()

	get := func(url string) (string, int) {
		t.Helper()
		w := httptest.NewRecorder()
		b.ServeHTTP(w, httptest.NewRequest("GET", url, nil))
		return html.UnescapeString(strings.NewReplacer("\\u003c", "<", "\\u003e", ">", "\\u0026", "&", `\"`, `"`).Replace(w.Body.String())), w.Code
	}
	if body, code := get("/admin/ui/widgets/widgets-fin"); code != 200 || !strings.Contains(body, "a página do grupo") {
		t.Fatalf("the page of the group: %d %.300s", code, body)
	}
	body, code := get("/admin/ui/widgets/widgets-fin/pending")
	if code != 200 || !strings.Contains(body, "a página do widget") {
		t.Fatalf("the sub-page: %d %.300s", code, body)
	}
	for _, want := range []string{`"groupHref":"/admin/ui/widgets/widgets-fin"`, `"href":"/admin/ui/widgets/widgets-fin/pending"`,
		`"groupHref":"/admin/ui/widgets"`, "Pendências"} {
		if !strings.Contains(body, want) {
			t.Errorf("the menu has no %s", want)
		}
	}
}
