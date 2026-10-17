package presets

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

// A sub-page in the menu (InMenu) is a child of its page's item, in the order
// registered — but one MenuVisible refuses —; shown, it is the active one and
// its page is opened.
func TestMenuSubPages(t *testing.T) {
	b := menuBuilder()
	layout := func(*web.EventContext) (r web.PageResponse, err error) { return }
	title := func(s string) func(context.Context) string { return func(context.Context) string { return s } }
	pg := b.PagesRegistrator().New(HttpPage("/files")).Title(title("Files")).Layout(layout)
	a := pg.SubPage("/a").InMenu(true).Title(title("Sub A")).Layout(layout)
	hidden := pg.SubPage("/h").InMenu(true).Title(title("Sub H")).
		MenuVisible(func(*http.Request) bool { return false }).Layout(layout)
	out := pg.SubPage("/u/{user}").Title(title("Sub U")).Layout(layout)
	for _, p := range []*PageBuilder{pg, a, hidden, out} {
		p.Build()
	}
	b.Build(http.NewServeMux())

	menu := func(path string) string {
		r := httptest.NewRequest("GET", path, nil)
		ctx := &web.EventContext{R: r}
		return h.MustString(b.CreateMenus(ctx), r.Context())
	}
	body := menu("/files/a")
	if !strings.Contains(body, "Sub A") || strings.Contains(body, "Sub H") || strings.Contains(body, "Sub U") {
		t.Errorf("the children: %s", body)
	}
	if !strings.Contains(body, `activated: [&#34;/files/a&#34;]`) && !strings.Contains(body, `activated: ["/files/a"]`) {
		t.Errorf("the child not active: %s", body)
	}
	if !strings.Contains(body, `opened: [&#34;/files&#34;]`) && !strings.Contains(body, `opened: ["/files"]`) {
		t.Errorf("its page not opened: %s", body)
	}
}
