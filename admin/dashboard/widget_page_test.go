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

// The panel of a group (the page of an area): its widgets alone, in its
// order, the period at its top, no tools of the dashboard's.
func TestGroupPanel(t *testing.T) {
	b := New(presets.New(i18n.New()).URIPrefix("/admin"), nil)
	b.Now = func() time.Time { return now }
	load := func(*web.EventContext, Period) ([]Block, error) { return nil, nil }
	title := func(context.Context) string { return "w" }
	b.Widget(&Widget{ID: "fin-b", Group: "fin", Order: 2, Title: title, Load: load},
		&Widget{ID: "fin-a", Group: "fin", Order: 1, Title: title, Load: load},
		&Widget{ID: "other", Group: "sys", Title: title, Load: load})
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/admin/orgs/1/fin?period=7d", nil)}
	body := h.MustString(b.GroupPanel(ctx, "fin"), context.Background())
	a, bb := strings.Index(body, `data-dashboard-widget='fin-a'`), strings.Index(body, `data-dashboard-widget='fin-b'`)
	if a < 0 || bb < a || strings.Contains(body, `data-dashboard-widget='other'`) {
		t.Errorf("the widgets of fin, in its order:\n%s", body)
	}
	if !strings.Contains(body, `data-dashboard-period='7d'`) || strings.Contains(body, "data-dashboard-customize") ||
		!strings.Contains(body, `data-dashboard-group-panel='fin'`) {
		t.Errorf("the period, no tools:\n%s", body)
	}
}

// The blocks of a widget are kept by the page they are of: the same widget
// on another page (the panel of an organization's group) loads its own.
func TestCacheByPage(t *testing.T) {
	b := New(presets.New(i18n.New()).URIPrefix("/admin"), nil)
	b.Now = func() time.Time { return now }
	var loads []string
	b.Widget(&Widget{ID: "w", TTL: time.Hour, Title: func(context.Context) string { return "w" },
		Load: func(ctx *web.EventContext, _ Period) ([]Block, error) {
			loads = append(loads, ctx.R.URL.Path)
			return nil, nil
		}})
	p := NewPeriod(Period7d, now, "", "")
	for _, path := range []string{"/admin", "/admin/orgs/1/finance", "/admin", "/admin/orgs/1/finance"} {
		b.Render(&web.EventContext{R: httptest.NewRequest("POST", path, nil)}, b.widget("w"), p, false)
	}
	if strings.Join(loads, " ") != "/admin /admin/orgs/1/finance" {
		t.Errorf("the loads: %v", loads)
	}
}
