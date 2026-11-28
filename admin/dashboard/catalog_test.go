package dashboard

import (
	"html"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

// The catalog in the menu: the group Widgets in the parent given, a page per
// group of widgets (its widgets, the panels that take them), a page per
// widget (what it is, a preview); groups registered after Catalog too.
func TestCatalog(t *testing.T) {
	p := presets.New(i18n.New()).URIPrefix("/admin")
	d := New(p, nil)
	if err := d.Install(); err != nil {
		t.Fatal(err)
	}
	load := func(*web.EventContext, Period) ([]Block, error) { return []Block{&Note{Text: "loaded"}}, nil }
	d.Catalog(p.MenuGroup("ui").Title("UI"))
	d.Group("fin").Title(titled("Financeiro")).Icon("mdi-finance").Description(titled("o dinheiro")).
		Widget("due", func(w *WidgetBuilder) { w.Title(titled("A vencer")).Size(Half).Load(load) })
	d.Group("sys").Title(titled("Sistema")).Widget("cpu", func(w *WidgetBuilder) { w.Title(titled("CPU")).Load(load) })
	d.Panel(PanelHome).Title(titled("Início")).AllGroups().WithoutGroup(d.GetGroup("sys"))
	p.Build()

	get := func(url string) (string, int) {
		t.Helper()
		w := httptest.NewRecorder()
		p.ServeHTTP(w, httptest.NewRequest("GET", url, nil))
		return html.UnescapeString(strings.NewReplacer("\\u003c", "<", "\\u003e", ">", "\\u0026", "&", `\"`, `"`).Replace(w.Body.String())), w.Code
	}
	body, code := get("/admin/ui/widgets/widgets-fin")
	if code != 200 || !strings.Contains(body, "data-dashboard-catalog-widget='fin/due'") || !strings.Contains(body, "o dinheiro") ||
		!strings.Contains(body, "/admin/ui/widgets/widgets-fin/due") || !strings.Contains(body, "Dashboards: Início") {
		t.Fatalf("the page of the group: %d %.500s", code, body)
	}
	body, code = get("/admin/ui/widgets/widgets-fin/due")
	if code != 200 || !strings.Contains(body, "data-dashboard-catalog='fin/due'") || !strings.Contains(body, "6 of 12 columns") ||
		!strings.Contains(body, "data-dashboard-widget='fin/due'") {
		t.Fatalf("the page of the widget: %d %.500s", code, body)
	}
	for _, want := range []string{`"groupHref":"/admin/ui/widgets"`, `"groupHref":"/admin/ui/widgets/widgets-fin"`,
		`"groupHref":"/admin/ui/widgets/widgets-sys"`, `"href":"/admin/ui/widgets/widgets-sys/cpu"`,
		// each group's icon, its widgets' its group's
		`"groupHref":"/admin/ui/widgets/widgets-fin","prependIcon":"mdi-finance"`, `"href":"/admin/ui/widgets/widgets-fin/due","prependIcon":"mdi-finance"`} {
		if !strings.Contains(body, want) {
			t.Errorf("the menu has no %s", want)
		}
	}
	// on no panel: said, no preview
	body, _ = get("/admin/ui/widgets/widgets-sys/cpu")
	if !strings.Contains(body, "On no dashboard") || strings.Contains(body, "data-dashboard-widget='sys/cpu'") {
		t.Errorf("a widget of no panel: %.500s", body)
	}
}
