package dashboard

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

var now = time.Date(2026, 10, 18, 14, 37, 0, 0, time.UTC)

func TestPeriods(t *testing.T) {
	for _, c := range []struct {
		key     string
		step    time.Duration
		buckets int
	}{
		{Period1h, 5 * time.Minute, 12},
		{Period12h, time.Hour, 12},
		{Period24h, time.Hour, 24},
		{Period7d, 24 * time.Hour, 7},
		{Period30d, 24 * time.Hour, 30},
		{Period90d, 24 * time.Hour, 90},
	} {
		p := NewPeriod(c.key, now, "", "")
		if p.Key != c.key || p.Step != c.step || len(p.Buckets()) != c.buckets {
			t.Errorf("%s: key %s step %s buckets %d", c.key, p.Key, p.Step, len(p.Buckets()))
		}
		if !p.To.After(now) || p.To.Sub(now) > p.Step {
			t.Errorf("%s: the last bucket must hold now: to %s", c.key, p.To)
		}
		prev := p.Previous()
		if !prev.To.Equal(p.From) || prev.To.Sub(prev.From) != p.To.Sub(p.From) {
			t.Errorf("%s: previous %s..%s", c.key, prev.From, prev.To)
		}
	}
	p := NewPeriod(PeriodCustom, now, "2026-10-01", "2026-10-03")
	if p.Key != PeriodCustom || len(p.Buckets()) != 3 {
		t.Errorf("custom: %s %d", p.Key, len(p.Buckets()))
	}
	if p := NewPeriod(PeriodCustom, now, "2026-10-03", "2026-10-01"); p.Key != PeriodDefault {
		t.Errorf("an inverted custom range must be the default, got %s", p.Key)
	}
	if p := NewPeriod("bogus", now, "", ""); p.Key != PeriodDefault {
		t.Errorf("unknown: %s", p.Key)
	}
	if l := NewPeriod(Period1h, now, "", "").Labels(4); l[0] == "" || !strings.Contains(l[0], ":") {
		t.Errorf("hour labels: %q", l)
	}
}

// The layout of a user: the groups' sections first (the old order and the
// widgets hidden before, removed), every widget in one; removed and added
// back to its group's; dragged — what is not known or repeated left out,
// what the user did not send kept —; the sections of the user's own made,
// changed, moved and deleted (their widgets back to their groups); a widget
// registered later in its group's.
func TestLayout(t *testing.T) {
	b := &Builder{}
	b.Group(&Group{ID: "g1", Order: 1}).Group(&Group{ID: "g2", Order: 2})
	b.Widget(&Widget{ID: "a", Group: "g1", Order: 1}, &Widget{ID: "b", Group: "g1", Order: 2},
		&Widget{ID: "c", Group: "g1", Order: 3}, &Widget{ID: "d", Group: "g2", Order: 1})
	str := func(ss []Section) string {
		var out []string
		for _, s := range ss {
			out = append(out, s.ID+":"+strings.Join(s.Widgets, ","))
		}
		return strings.Join(out, " ")
	}

	p := Prefs{Order: []string{"c"}, Hidden: []string{"b"}}
	if got := str(b.layout(&p)); got != "g1:c,a g2:d" {
		t.Fatalf("first layout %q", got)
	}
	b.add(&p, "b")
	if got := str(b.layout(&p)); got != "g1:c,a,b g2:d" || p.removed("b") {
		t.Errorf("b added back %q", got)
	}
	b.remove(&p, "a")
	if got := str(b.layout(&p)); got != "g1:c,b g2:d" || !p.removed("a") {
		t.Errorf("a removed %q", got)
	}

	// dragged: b to g2 before d, a removed and x unknown left out, c twice
	b.arrange(&p, []Section{{ID: "g2", Widgets: []string{"b", "d", "x", "a"}}, {ID: "g1", Widgets: []string{"c", "c"}}})
	if got := str(b.layout(&p)); got != "g2:b,d g1:c" {
		t.Errorf("arranged %q", got)
	}
	// a section not sent and its widgets: kept
	b.arrange(&p, []Section{{ID: "g1", Widgets: []string{"c"}}})
	if got := str(b.layout(&p)); got != "g1:c g2:b,d" {
		t.Errorf("arranged, g2 not sent %q", got)
	}

	id := b.saveSection(&p, "", " Mine ", "mdi-star", "")
	b.arrange(&p, []Section{{ID: id, Widgets: []string{"d"}}, {ID: "g1", Widgets: []string{"c"}}, {ID: "g2", Widgets: []string{"b"}}})
	if got := str(b.layout(&p)); got != id+":d g1:c g2:b" {
		t.Errorf("own section %q", got)
	}
	if s := section(p.Sections, id); s == nil || !s.Custom || s.Title != "Mine" || s.Icon != "mdi-star" {
		t.Errorf("own section %+v", s)
	}
	b.saveSection(&p, id, "Ours", "not an icon", "")
	if s := section(p.Sections, id); s.Title != "Ours" || s.Icon != "" {
		t.Errorf("own section changed %+v", s)
	}
	b.moveSection(&p, id, 1)
	if got := str(b.layout(&p)); got != "g1:c "+id+":d g2:b" {
		t.Errorf("moved %q", got)
	}
	b.moveWidget(&p, "c", 1)
	b.deleteSection(&p, "g1") // a group's: not deleted
	b.deleteSection(&p, id)
	if got := str(b.layout(&p)); got != "g1:c g2:b,d" {
		t.Errorf("own deleted %q", got)
	}

	p.Disabled = toggle(p.Disabled, "c", true)
	if !p.disabled("c") || str(b.layout(&p)) != "g1:c g2:b,d" {
		t.Errorf("disabled stays in its place")
	}

	// a tree: g2 in a section of the user's own, in it another; no loop
	top := b.saveSection(&p, "", "Top", "", "")
	b.saveSection(&p, "g2", "ignored", "", top)
	in := b.saveSection(&p, "", "In", "", top)
	b.saveSection(&p, top, "Top", "", in) // under itself: not
	ss := b.layout(&p)
	if section(ss, "g2").Parent != top || section(ss, in).Parent != top || section(ss, top).Parent != "" || section(ss, "g2").Title != "" {
		t.Fatalf("tree %+v", ss)
	}
	if kids := children(ss, top); len(kids) != 2 || kids[0].ID != "g2" || kids[1].ID != in {
		t.Errorf("children of top %+v", kids)
	}
	b.moveSection(&p, in, -1) // among its siblings
	if kids := children(b.layout(&p), top); kids[0].ID != in {
		t.Errorf("moved among siblings %+v", kids)
	}
	b.deleteSection(&p, top) // its sections to its parent (the top)
	ss = b.layout(&p)
	if section(ss, top) != nil || section(ss, "g2").Parent != "" || section(ss, in).Parent != "" {
		t.Errorf("deleted, children up %+v", ss)
	}
	b.deleteSection(&p, in)

	// into a section (its +): a removed one back there, another out of
	// where it was, a disabled one enabled; no such section, nothing
	b.remove(&p, "b")
	if !b.moveTo(&p, "g1", "b", "d", "c") || b.moveTo(&p, "nope", "c") {
		t.Fatal("moveTo")
	}
	if got := str(b.layout(&p)); got != "g1:b,d,c g2:" || p.removed("b") || p.disabled("c") {
		t.Errorf("moved to g1 %q", got)
	}
	if s := sectionOf(b.layout(&p), "d"); s == nil || s.ID != "g1" {
		t.Errorf("sectionOf d %+v", s)
	}
	b.moveTo(&p, "g2", "b", "d")
	if got := str(b.layout(&p)); got != "g1:c g2:b,d" {
		t.Errorf("moved back %q", got)
	}

	// a widget registered later: last in its group's
	b.Widget(&Widget{ID: "e", Group: "g1", Order: 0})
	if got := str(b.layout(&p)); got != "g1:c,e g2:b,d" {
		t.Errorf("new widget %q", got)
	}
}

func render(t *testing.T, bs ...Block) string {
	t.Helper()
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/?period=7d", nil)}
	rc := &renderCtx{w: &Widget{ID: "w"}, ctx: ctx, msgs: Messages_en_US, period: NewPeriod(Period7d, now, "", "")}
	var out []h.HTMLComponent
	for _, b := range bs {
		out = append(out, b.render(rc))
	}
	return h.MustString(h.Components(out...), context.Background())
}

func TestBlocks(t *testing.T) {
	prev := 50.0
	s := render(t, &Metric{Value: 75, Previous: &prev, Spark: []float64{1, 3, 2}})
	for _, want := range []string{"75", "▲ 50%", "<svg", "text-success"} {
		if !strings.Contains(s, want) {
			t.Errorf("metric: missing %q in %s", want, s)
		}
	}
	if s := render(t, &Metric{Value: 3, Previous: &prev, DownIsGood: true}); !strings.Contains(s, "▼") || !strings.Contains(s, "text-success") {
		t.Errorf("a fall that is good: %s", s)
	}
	labels := []string{"a", "b", "c"}
	for _, k := range []ChartKind{ChartLine, ChartBars, ChartStacked, ChartDonut, ChartHBars} {
		s := render(t, &Chart{Kind: k, Labels: labels, Series: []Series{{Name: "n", Values: []float64{1, 2, 3}}, {Name: "m", Values: []float64{3, 2, 1}}}})
		if k == ChartHBars {
			if strings.Count(s, "data-dashboard-bar") != 3 {
				t.Errorf("hbars: %s", s)
			}
		} else if !strings.Contains(s, "<svg") {
			t.Errorf("%s: no svg: %s", k, s)
		}
	}
	if s := render(t, &Chart{Kind: ChartLine, Series: []Series{{Values: []float64{0, 0}}}}); !strings.Contains(s, Messages_en_US.NoData) {
		t.Errorf("an empty chart says so: %s", s)
	}
	s = render(t, &List{Rows: []Row{{Text: "<b>x</b>", Link: "/x", Right: "3", Action: &Action{Name: "approve", Arg: "1", Label: "Approve"}}}})
	if strings.Contains(s, "<b>x</b>") || !strings.Contains(s, "Approve") || !strings.Contains(s, "dashboard_action") {
		t.Errorf("list: %s", s)
	}
	s = render(t, &Status{Items: []StatusItem{{Level: LevelError, Title: "Backup", Detail: "failed"}}}, &Progress{Bars: []Bar{{Label: "pt", Done: 5, Total: 10}}})
	if !strings.Contains(s, "data-dashboard-status='error'") || !strings.Contains(s, "5/10") {
		t.Errorf("status/progress: %s", s)
	}
}

func TestFormat(t *testing.T) {
	if formatNumber(1234567) != "1.234.567" || formatNumber(-1000) != "-1.000" || formatNumber(1.5) != "1.5" {
		t.Error(formatNumber(1234567), formatNumber(-1000), formatNumber(1.5))
	}
	if FormatBytes(999) != "999 B" || FormatBytes(1_400_000_000) != "1,4 GB" {
		t.Error(FormatBytes(999), FormatBytes(1_400_000_000))
	}
}

func TestCache(t *testing.T) {
	b := &Builder{Now: func() time.Time { return now }}
	n := 0
	w := &Widget{ID: "w", Load: func(*web.EventContext, Period) ([]Block, error) { n++; return nil, nil }}
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
	p := NewPeriod(Period7d, now, "", "")
	b.blocks(ctx, w, p, false)
	b.blocks(ctx, w, p, false)
	if n != 1 {
		t.Errorf("cached: loads %d", n)
	}
	b.blocks(ctx, w, p, true)
	if n != 2 {
		t.Errorf("refresh skips the cache: loads %d", n)
	}
	b.Now = func() time.Time { return now.Add(2 * time.Minute) }
	b.blocks(ctx, w, p, false)
	if n != 3 {
		t.Errorf("ttl: loads %d", n)
	}
	b.Invalidate("w")
	b.blocks(ctx, w, p, false)
	if n != 4 {
		t.Errorf("invalidate: loads %d", n)
	}
}
