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

func TestPrefsMove(t *testing.T) {
	p := Prefs{Order: []string{"x"}}
	p.move([]string{"a", "b", "c"}, "c", -1)
	if strings.Join(p.Order, ",") != "x,a,c,b" {
		t.Errorf("order %v", p.Order)
	}
	p.move([]string{"a", "c", "b"}, "a", -1)
	if strings.Join(p.Order, ",") != "x,a,c,b" {
		t.Errorf("first up must not move: %v", p.Order)
	}
	p.Hidden = toggle(p.Hidden, "a", true)
	p.Hidden = toggle(p.Hidden, "a", true)
	if len(p.Hidden) != 1 || !p.hidden("a") {
		t.Errorf("hidden %v", p.Hidden)
	}
	p.Hidden = toggle(p.Hidden, "a", false)
	if p.hidden("a") {
		t.Errorf("hidden %v", p.Hidden)
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
