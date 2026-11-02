package dashboard

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

// The page of a widget: the widget alone, loaded in its full version (the
// query full=1, IsFull for its Load), with the way back; another id, not
// found.
func TestWidgetPage(t *testing.T) {
	b := New(presets.New(i18n.New()).URIPrefix("/admin"), nil)
	b.Now = func() time.Time { return now }
	var full bool
	b.Widget(&Widget{ID: "mem", Title: func(context.Context) string { return "Memory" }, Periodic: true,
		Live: []time.Duration{2 * time.Second},
		Load: func(ctx *web.EventContext, _ Period) ([]Block, error) { full = IsFull(ctx); return nil, nil }})

	ctx := &web.EventContext{R: httptest.NewRequest("GET", WidgetPagePath+"?id=mem&period=7d", nil), W: web.WrapResponseWriter(httptest.NewRecorder())}
	r, err := b.widgetPage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	body := h.MustString(r.Body, context.Background())
	for _, want := range []string{`data-dashboard-widget-page='mem'`, `data-dashboard-back`, PortalName("mem"), `full`, `data-dashboard-live='mem'`, `data-dashboard-period='7d'`} {
		if !strings.Contains(body, want) {
			t.Errorf("no %q:\n%s", want, body)
		}
	}
	if r.PageTitle != "Memory" {
		t.Errorf("title %q", r.PageTitle)
	}

	// its load, as the page asks it: the full version
	ctx.R = httptest.NewRequest("POST", "/?widget=mem&period=7d&full=1", nil)
	b.Render(ctx, b.widget("mem"), NewPeriod(Period7d, now, "", ""), true)
	if !full {
		t.Error("the load of the page not full")
	}

	ctx.R = httptest.NewRequest("GET", WidgetPagePath+"?id=none", nil)
	r, _ = b.widgetPage(ctx)
	if strings.Contains(h.MustString(r.Body, context.Background()), "data-dashboard-widget-page") {
		t.Error("an unknown widget has a page")
	}
}
