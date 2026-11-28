package dashboard

import (
	"context"
	"net/http/httptest"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"slices"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

// panics says f panics.
func panics(f func()) (ok bool) {
	defer func() { ok = recover() != nil }()
	f()
	return
}

// The widgets of a panel: of the groups it takes (those added after too, by
// another package: GetGroup), the ones it takes alone, but those left out;
// AllGroups, every group's but the ones left out.
func TestPanelWidgets(t *testing.T) {
	b := New(nil, nil)
	ids := func(pn *Panel) (out []string) {
		for _, w := range pn.all() {
			out = append(out, w.ID)
		}
		return
	}
	// a panel and a group referred to before made
	org := b.GetPanel("org")
	fin := b.GetGroup("fin").Widget("x", nil)
	b.Group("sys").Widget("cpu", nil).Widget("mem", nil)
	sys := b.GetGroup("sys")
	org.Group(fin).Widget(sys.GetWidget("mem"))
	home := b.Panel(PanelHome).AllGroups().WithoutGroup(sys)
	home.Widget(sys.GetWidget("cpu"))
	b.Group("fin").Widget("y", nil)
	b.Panel("org")
	if got := ids(org); !slices.Equal(got, []string{"fin/x", "fin/y", "sys/mem"}) {
		t.Errorf("org: %v", got)
	}
	if got := ids(home); !slices.Equal(got, []string{"fin/x", "fin/y", "sys/cpu"}) {
		t.Errorf("home: %v", got)
	}
	org.Without(fin.GetWidget("y"))
	if got := ids(org); !slices.Equal(got, []string{"fin/x", "sys/mem"}) {
		t.Errorf("org without fin/y: %v", got)
	}
	if len(b.Panels()) != 2 || len(b.Groups()) != 2 || len(b.Widgets()) != 4 {
		t.Errorf("made: %d panels, %d groups, %d widgets", len(b.Panels()), len(b.Groups()), len(b.Widgets()))
	}
	if b.widget("fin/y") == nil || b.widget("fin") != nil || b.widget("none/y") != nil {
		t.Error("the widget by its id")
	}
}

// Made twice, a widget registered twice: an error; referred to and never
// made: an error of check.
func TestPanelErrors(t *testing.T) {
	b := New(nil, nil)
	b.Group("g").Widget("w", nil)
	b.Panel("p")
	if !panics(func() { b.Group("g") }) || !panics(func() { b.Panel("p") }) ||
		!panics(func() { b.GetGroup("g").Widget("w", nil) }) {
		t.Error("made twice, no panic")
	}
	if panics(b.check) {
		t.Error("everything made: a panic")
	}
	b.GetGroup("never")
	if !panics(b.check) {
		t.Error("a group never made: no panic")
	}
	b = New(nil, nil)
	b.GetPanel("never")
	if !panics(b.check) {
		t.Error("a panel never made: no panic")
	}
}

type scopeRec struct{ ID uint }

// The scope of a panel's request: its id and, of a panel of a record, the
// record's (the page's {id}, else the query's); the panel of an event its
// query's, the home's by default.
func TestScope(t *testing.T) {
	b := New(nil, nil)
	home := b.Panel(PanelHome)
	org := b.Panel("org")
	org.record = presets.NewModelBuilder(presets.New(i18n.New()), &scopeRec{})

	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/admin/orgs/7/painel", nil)}
	ctx.R.SetPathValue("id", "7")
	if s := ScopeOf(org.scoped(ctx)); s.Panel != "org" || s.Record != "7" {
		t.Errorf("the page of the record: %+v", s)
	}
	ctx = &web.EventContext{R: httptest.NewRequest("POST", "/admin/dashboard/widget?panel=org&record=9", nil)}
	if s := ScopeOf(org.scoped(ctx)); s.Record != "9" {
		t.Errorf("the record of the query: %+v", s)
	}
	if c := org.in(ctx); c == org || c.recordParam != "9" {
		t.Error("out of the record's page, its events carry no record")
	}
	ctx = &web.EventContext{R: httptest.NewRequest("POST", "/admin?record=9", nil)}
	if s := ScopeOf(home.scoped(ctx)); s.Panel != PanelHome || s.Record != "" {
		t.Errorf("the home: %+v", s)
	}
	for url, want := range map[string]*Panel{"/": home, "/?panel=org": org, "/?panel=none": nil} {
		if got := b.panelOf(httptest.NewRequest("POST", url, nil)); got != want {
			t.Errorf("%s: the panel %v", url, got)
		}
	}
}

// The preferences of before (a widget's name alone) are the panel's
// widgets' of that name, when one only.
func TestNormalize(t *testing.T) {
	b := New(nil, nil)
	b.Group("a").Widget("w", nil).Widget("same", nil)
	b.Group("b").Widget("same", nil)
	pn := b.Panel(PanelHome).AllGroups()
	p := Prefs{Sections: []Section{{ID: "a", Widgets: []string{"w", "same"}}}, Removed: []string{"w"}, Order: []string{"a/w"}}
	pn.normalize(&p)
	if !slices.Equal(p.Sections[0].Widgets, []string{"a/w", "same"}) || !slices.Equal(p.Removed, []string{"a/w"}) ||
		!slices.Equal(p.Order, []string{"a/w"}) {
		t.Errorf("normalized: %+v", p)
	}
}

// Customizing: each section a handle to drag it among its siblings, a place
// to drop on; the sections dragged (their order in the page) reordered, each
// kept where it was — never put in another; a widget dragged to another
// group's section, there. The script saves to the panel of the page.
func TestDragSections(t *testing.T) {
	bd := New(presets.New(i18n.New()).URIPrefix("/admin"), nil)
	bd.Now = func() time.Time { return now }
	load := func(*web.EventContext, Period) ([]Block, error) { return nil, nil }
	w := func(wb *WidgetBuilder) { wb.Title(titled("w")).Load(load) }
	bd.Group("a").Order(1).Widget("x", w)
	bd.Group("b").Order(2).Widget("y", w)
	bd.Group("c").Order(3).Widget("z", w)
	pn := bd.Panel("org").AllGroups()

	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/admin/orgs/1/painel?customize=1", nil), W: web.WrapResponseWriter(httptest.NewRecorder())}
	r, err := pn.Page(ctx)
	if err != nil {
		t.Fatal(err)
	}
	body := hs(r.Body)
	for _, want := range []string{`data-dashboard-section-drag='a'`, `data-dashboard-section-drag='c'`, `secStart`, `secOver`, `secDrop`} {
		if !strings.Contains(body, want) {
			t.Errorf("customizing: no %s", want)
		}
	}
	if s := dashScript(); !strings.Contains(s, `query("panel", root.getAttribute("data-dashboard")`) {
		t.Error("the script saves to no panel")
	}

	// b in a (a section in another, by its edit)
	var p Prefs
	pn.saveSection(&p, "b", "", "", "a")
	// dragged: c before a; b, its parent's widget z
	pn.arrange(&p, []Section{{ID: "c"}, {ID: "a", Widgets: []string{"a/x"}}, {ID: "b", Widgets: []string{"b/y", "c/z"}}})
	got := map[string]Section{}
	var order []string
	for _, s := range pn.layout(&p) {
		got[s.ID] = s
		if s.Parent == "" {
			order = append(order, s.ID)
		}
	}
	if !slices.Equal(order, []string{"c", "a"}) || got["b"].Parent != "a" || got["c"].Parent != "" ||
		!slices.Equal(got["b"].Widgets, []string{"b/y", "c/z"}) || len(got["c"].Widgets) != 0 {
		t.Errorf("arranged: %v %+v", order, got)
	}
}

func hs(c h.HTMLComponent) string { return h.MustString(c, context.Background()) }

// A widget resized: its width the user's, in columns of the row (MinSize to
// Full; its own forgotten); drawn so, a handle on its right edge while
// customizing.
func TestResize(t *testing.T) {
	b, pn := newTest()
	b.Group("g").Widget("w", func(w *WidgetBuilder) { w.Size(Third).Title(titled("w")) })
	w := pn.widget("g/w")
	var p Prefs
	p.resize(w, TwoThirds)
	if p.sizeOf(w) != TwoThirds {
		t.Errorf("resized: %d", p.sizeOf(w))
	}
	p.Sizes[w.ID] = 40
	if p.sizeOf(w) != Third {
		t.Errorf("out of the grid: %d", p.sizeOf(w))
	}
	p.resize(w, Third)
	if _, ok := p.Sizes[w.ID]; ok {
		t.Error("its own size kept")
	}
	p.resize(w, Half)
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/?customize=1", nil)}
	ss := pn.layout(&p)
	body := hs(pn.section(ctx, Messages_en_US, ss, ss[0], []*Widget{w}, nil, &p, NewPeriod(Period7d, now, "", ""), true, true, true))
	for _, want := range []string{`:md='6'`, `data-dashboard-size='6'`, `data-dashboard-resize='g/w'`, "resizeStart"} {
		if !strings.Contains(body, want) {
			t.Errorf("no %s", want)
		}
	}
}

// A group Untitled: its section with no heading on the panel, but while
// customizing; its title still the catalog's.
func TestUntitled(t *testing.T) {
	b, pn := newTest()
	b.Group("k").Title(titled("Indicators")).Untitled().Widget("n", func(w *WidgetBuilder) { w.Title(titled("w")) })
	var p Prefs
	ss := pn.layout(&p)
	draw := func(customize bool) string {
		ctx := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
		return hs(pn.section(ctx, Messages_en_US, ss, ss[0], pn.Widgets(), nil, &p, NewPeriod(Period7d, now, "", ""), customize, true, true))
	}
	if strings.Contains(draw(false), "Indicators") || !strings.Contains(draw(true), "Indicators") {
		t.Error("the heading of an untitled group")
	}
}
