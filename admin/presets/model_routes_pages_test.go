package presets

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

type routesInvoice struct {
	ID   uint
	Name string
}

type routesSettings struct {
	ID   uint
	Name string
}

func nopHandler() http.Handler {
	return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
}

func nopPageFunc(*web.EventContext) (web.PageResponse, error) {
	return web.PageResponse{}, nil
}

// mountedIn reports the pattern a path resolves to, "" when nothing is mounted.
func mountedIn(mux *http.ServeMux, pth string) string {
	_, pattern := mux.Handler(httptest.NewRequest("GET", pth, nil))
	return pattern
}

// A page registered on the listing or on the detailing is a CHILD of that
// page: the listing's hang under the listing route, the detailing's under the
// record route.
func TestModelPagesAreMountedUnderListingAndDetail(t *testing.T) {
	b := New(i18n.New())
	mb := b.Model(&routesSettings{}).URIName("settings")

	mb.Listing().PagesRegistrator().AddHttpPage(HttpPage("/import").Handler(nopHandler()))
	mb.Listing().AddRawPageFunc("/export", nopPageFunc)
	mb.Detailing().PagesRegistrator().AddHttpPage(HttpPage("/log").Handler(nopHandler()))
	mb.Detailing().AddRawPageFunc("/audit", nopPageFunc)

	mux := http.NewServeMux()
	mb.SetupRoutes(mux)

	want := map[string]string{
		"listing page (registrator)": "/settings/import",
		"listing page (raw)":         "/settings/export",
		"detail page (registrator)":  "/settings/1/log",
		"detail page (raw)":          "/settings/1/audit",
	}

	for what, pth := range want {
		if got := mountedIn(mux, pth); got == "" {
			t.Errorf("%s: %s não foi montada", what, pth)
		}
	}
}

// A singleton is a single record and no listing: the menu it shows is the
// record's, so only the detailing pages hang under it. The detailing
// registrator was not mounted there at all until now.
func TestSingletonMountsTheDetailPagesOnly(t *testing.T) {
	b := New(i18n.New())
	mb := b.Model(&routesSettings{}).URIName("settings").Singleton(true)

	mb.Detailing().PagesRegistrator().AddHttpPage(HttpPage("/log").Handler(nopHandler()))
	mb.Detailing().AddRawPageFunc("/audit", nopPageFunc)

	mux := http.NewServeMux()
	mb.SetupRoutes(mux)

	for what, pth := range map[string]string{
		"detail page (registrator)": "/settings/log",
		"detail page (raw)":         "/settings/audit",
	} {
		if got := mountedIn(mux, pth); got == "" {
			t.Errorf("%s: %s não foi montada", what, pth)
		}
	}

}

// setupPanic boots the model and returns what the boot died of, or "".
func setupPanic(t *testing.T, build func(b *Builder) *ModelBuilder) (msg string) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Sprint(r)
		}
	}()

	build(New(i18n.New())).SetupRoutes(http.NewServeMux())
	return ""
}

// A page registered on the listing of a singleton would be mounted nowhere:
// it would look registered and answer 404. That is a mistake in the setup, so
// the boot stops and says which pages and where they belong.
func TestSingletonRefusesListingPages(t *testing.T) {
	msg := setupPanic(t, func(b *Builder) *ModelBuilder {
		mb := b.Model(&routesSettings{}).URIName("settings").Singleton(true)
		mb.Listing().PagesRegistrator().AddHttpPage(HttpPage("/import").Handler(nopHandler()))
		mb.Listing().AddRawPageFunc("/export", nopPageFunc)
		return mb
	})

	if msg == "" {
		t.Fatal("o boot seguiu em frente")
	}
	for _, want := range []string{"singleton", "settings", "/import", "/export", "Detailing()"} {
		if !strings.Contains(msg, want) {
			t.Errorf("o erro não menciona %q: %s", want, msg)
		}
	}
}

// Nothing registered on the listing, nothing to refuse.
func TestSingletonWithNoListingPagesBoots(t *testing.T) {
	if msg := setupPanic(t, func(b *Builder) *ModelBuilder {
		mb := b.Model(&routesSettings{}).URIName("settings").Singleton(true)
		mb.Detailing().PagesRegistrator().AddHttpPage(HttpPage("/log").Handler(nopHandler()))
		return mb
	}); msg != "" {
		t.Errorf("o boot morreu sem motivo: %s", msg)
	}
}

// A collection model mounts its listing pages, so there is nothing to refuse
// there either.
func TestCollectionWithListingPagesBoots(t *testing.T) {
	if msg := setupPanic(t, func(b *Builder) *ModelBuilder {
		mb := b.Model(&routesSettings{}).URIName("settings")
		mb.Listing().PagesRegistrator().AddHttpPage(HttpPage("/import").Handler(nopHandler()))
		return mb
	}); msg != "" {
		t.Errorf("o boot morreu sem motivo: %s", msg)
	}
}

// A model's page is a child of the listing or of the detailing, not an item of
// the side menu: it takes no key in the menu tree. Two models are therefore
// free to have a page of the same path — the key would collide.
func TestModelPagesTakeNoKeyInTheMenuTree(t *testing.T) {
	b := New(i18n.New())

	invoices := b.Model(&routesInvoice{}).URIName("invoices")
	settings := b.Model(&routesSettings{}).URIName("settings")

	invoices.Listing().PagesRegistrator().AddHttpPage(HttpPage("/import").Handler(nopHandler()))
	settings.Listing().PagesRegistrator().AddHttpPage(HttpPage("/import").Handler(nopHandler()))
	invoices.Detailing().PagesRegistrator().AddHttpPage(HttpPage("/log").Handler(nopHandler()))

	for _, key := range []string{"page:/import", "page:/log"} {
		if _, ok := b.MenuItems()[key]; ok {
			t.Errorf("%s entrou na árvore do menu", key)
		}
	}

	// the two models, and nothing else
	if got, want := len(b.MenuTree().items), 2; got != want {
		t.Errorf("itens na raiz = %d, want %d", got, want)
	}
}

// A page of the builder itself IS a side-menu item, and keeps its key.
func TestBuilderPagesStillTakeAKeyInTheMenuTree(t *testing.T) {
	b := New(i18n.New())
	b.PagesRegistrator().AddHttpPage(HttpPage("/report").Handler(nopHandler()))

	if _, ok := b.MenuItems()["page:/report"]; !ok {
		t.Error("a página do builder não entrou na árvore do menu")
	}
}

// A model's page has no entry in the tree, but the group it was given still
// shapes its URL — as it did before the tree existed.
func TestModelPageKeepsItsMenuGroupInTheURL(t *testing.T) {
	b := New(i18n.New())
	mb := b.Model(&routesInvoice{}).URIName("invoices")

	page := HttpPage("/import").Handler(nopHandler()).MenuGroup("tools")
	mb.Listing().PagesRegistrator().AddHttpPage(page)

	mux := http.NewServeMux()
	mb.SetupRoutes(mux)

	if got, want := page.FullPath(), "/invoices/tools/import"; got != want {
		t.Errorf("FullPath() = %q, want %q", got, want)
	}
	if got := mountedIn(mux, "/invoices/tools/import"); got == "" {
		t.Error("a página não foi montada sob o grupo")
	}
}
