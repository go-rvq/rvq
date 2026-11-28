package dashboard

import (
	"context"
	"encoding/json"
	"html"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

// renderIn draws bs in a dashboard of an admin (its vx-charts), for the
// request url.
func renderIn(t *testing.T, url string, bs ...Block) string {
	t.Helper()
	p := presets.New(i18n.New()).URIPrefix("/admin")
	InstallCharts(p)
	bd := New(p, nil)
	bd.Now = func() time.Time { return now }
	b := bd.Panel(PanelHome).AllGroups()
	ctx := &web.EventContext{R: httptest.NewRequest("GET", url, nil)}
	rc := &renderCtx{b: b, w: &Widget{ID: "g/w"}, ctx: ctx, msgs: Messages_en_US, period: b.eventPeriod(ctx)}
	out := []h.HTMLComponent{rc.crumbs()}
	for _, bl := range bs {
		out = append(out, bl.render(rc))
	}
	return h.MustString(h.Components(out...), context.Background())
}

// attr is the value of the attribute name of the first <vx-chart> of s,
// unescaped.
func attr(t *testing.T, s, name string) string {
	t.Helper()
	m := regexp.MustCompile(`<vx-chart[^>]*\s` + regexp.QuoteMeta(name) + `='([^']*)'`).FindStringSubmatch(s)
	if m == nil {
		t.Fatalf("no %s in %s", name, s)
	}
	return html.UnescapeString(m[1])
}

// A chart is a vx-chart, of each kind — its labels, its series, its hints,
// its colors as JSON, the code of the charts its src —, the SVG of before
// only printed.
func TestChartVX(t *testing.T) {
	for _, k := range []ChartKind{ChartLine, ChartBars, ChartStacked, ChartHBars, ChartDonut} {
		s := renderIn(t, "/?period=7d", &Chart{Kind: k, Title: "Visits", Labels: []string{"a", "", "c"}, Hints: []string{"A", "B", "C"},
			Colors: []string{"success"}, Series: []Series{{Name: "n", Color: "info", Values: []float64{1, 0, 3}}, {Name: "m", Values: nil}}})
		if got := attr(t, s, "kind"); got != string(k) {
			t.Errorf("%s: kind %q", k, got)
		}
		if !strings.Contains(attr(t, s, "src"), "/admin/!assets/extra/echarts.js?v=") {
			t.Errorf("%s: src %s", k, attr(t, s, "src"))
		}
		var labels, hints, colors []string
		var series []map[string]any
		for name, v := range map[string]any{":labels": &labels, ":hints": &hints, ":colors": &colors, ":series": &series} {
			if err := json.Unmarshal([]byte(attr(t, s, name)), v); err != nil {
				t.Fatalf("%s: %s: %v", k, name, err)
			}
		}
		if strings.Join(labels, ",") != "a,,c" || strings.Join(hints, ",") != "A,B,C" || strings.Join(colors, ",") != "success" {
			t.Errorf("%s: labels %q hints %q colors %q", k, labels, hints, colors)
		}
		if len(series) != 2 || series[0]["name"] != "n" || series[0]["color"] != "info" || len(series[1]["values"].([]any)) != 0 {
			t.Errorf("%s: series %v", k, series)
		}
		if !strings.Contains(s, "d-none d-print-block") || !strings.Contains(s, "d-print-none") {
			t.Errorf("%s: the paper's drawing, only printed: %s", k, s)
		}
		if strings.Contains(s, ":drill") {
			t.Errorf("%s: no drill given, none: %s", k, s)
		}
	}
}

// A period: the full label of each point; each point opened (a day by its
// hours, an hour by its 5 minutes, 5 minutes not), its bounds as the
// request says them, its name.
func TestPeriodHintsDrill(t *testing.T) {
	p := NewPeriod(Period7d, now, "", "")
	hs := p.Hints()
	if len(hs) != 7 || hs[6] != "18/10/2026" {
		t.Errorf("days: %q", hs)
	}
	d := p.Drill()
	if len(d) != 7 || d[6].Step != time.Hour || len(d[6].Buckets()) != 24 || d[6].Label(Messages_en_US) != "18/10" {
		t.Fatalf("a day opened: %+v", d[6])
	}
	if from, to := d[6].Bounds(); from != "2026-10-18" || to != "2026-10-18" {
		t.Errorf("a day's bounds %s %s", from, to)
	}
	// its bounds, back: the same period
	from, to := d[6].Bounds()
	if q := NewPeriod(PeriodCustom, now, from, to); !q.From.Equal(d[6].From) || !q.To.Equal(d[6].To) {
		t.Errorf("a day again: %s..%s", q.From, q.To)
	}

	day := NewPeriod(PeriodCustom, now, "2026-10-17", "2026-10-17")
	if hs := day.Hints(); len(hs) != 24 || hs[14] != "17/10 14:00–15:00" {
		t.Errorf("hours: %q", hs)
	}
	hour := day.Drill()[14]
	from, to = hour.Bounds()
	if from != "2026-10-17T14:00" || to != "2026-10-17T15:00" || hour.Label(Messages_en_US) != "17/10 14:00–15:00" {
		t.Errorf("an hour: %s %s %s", from, to, hour.Label(Messages_en_US))
	}
	q := NewPeriod(PeriodCustom, now, from, to)
	if q.Step != 5*time.Minute || len(q.Buckets()) != 12 {
		t.Errorf("an hour again: %s %d", q.Step, len(q.Buckets()))
	}
	if q.Drill() != nil {
		t.Error("5 minutes are not opened")
	}
	if bad := NewPeriod(PeriodCustom, now, "2026-10-17T15:00", "2026-10-17T14:00"); bad.Key != PeriodDefault {
		t.Errorf("an inverted range of times: %s", bad.Key)
	}
	if l := NewPeriod(Period30d, now, "", "").Label(Messages_en_US); l != Messages_en_US.Period30d {
		t.Errorf("a chip's name: %s", l)
	}
	if l := NewPeriod(PeriodCustom, now, "2026-10-01", "2026-10-15").Label(Messages_en_US); l != "01/10–15/10" {
		t.Errorf("a range of days: %s", l)
	}
}

// A chart of a period's points opens them: its drill, its event (the
// widget loaded with the point's period, the way back with the period it
// came from); the widget loaded so says the way back, each period a link.
func TestChartDrillTrail(t *testing.T) {
	p := NewPeriod(Period7d, now, "", "")
	s := renderIn(t, "/?period=7d", &Chart{Kind: ChartLine, Labels: p.Labels(7), Hints: p.Hints(), Drill: p.Drill(),
		Series: []Series{{Name: "n", Values: []float64{1, 2, 3, 4, 5, 6, 7}}}})
	var drill []map[string]string
	if err := json.Unmarshal([]byte(attr(t, s, ":drill")), &drill); err != nil {
		t.Fatal(err)
	}
	if len(drill) != 7 || drill[6]["from"] != "2026-10-18" || drill[6]["to"] != "2026-10-18" {
		t.Errorf("drill %v", drill)
	}
	ev := attr(t, s, "@point")
	for _, want := range []string{EventLoad, `"period", "custom"`, "$event.from", "$event.to", `"trail", "7d~~"`, `"refresh", "1"`} {
		if !strings.Contains(ev, want) {
			t.Errorf("the event has no %s: %s", want, ev)
		}
	}

	// opened: 18/10 of 7 days, then its 14:00
	s = renderIn(t, "/?period=custom&from=2026-10-18T14:00&to=2026-10-18T15:00&trail=7d~~|custom~2026-10-18~2026-10-18",
		&Chart{Kind: ChartLine, Labels: []string{"x"}, Series: []Series{{Name: "n", Values: []float64{1}}}})
	if !strings.Contains(s, "data-dashboard-trail") || strings.Count(s, "data-dashboard-crumb=") != 2 {
		t.Fatalf("the way back: %s", s)
	}
	for _, want := range []string{Messages_en_US.Period7d, ">18/10<", "18/10 14:00–15:00"} {
		if !strings.Contains(s, want) {
			t.Errorf("the way back has no %q: %s", want, s)
		}
	}
	// the day's link: back to it, the way before it (7 days)
	if !strings.Contains(html.UnescapeString(s), `"from", "2026-10-18").query("to", "2026-10-18").query("trail", "7d~~")`) {
		t.Errorf("the day's link: %s", s)
	}
	// a point of the hour opened: its way back the two before and the hour
	p = NewPeriod(PeriodCustom, now, "2026-10-18T14:00", "2026-10-18T15:00")
	s = renderIn(t, "/?period=custom&from=2026-10-18T14:00&to=2026-10-18T15:00&trail=7d~~|custom~2026-10-18~2026-10-18",
		&Chart{Kind: ChartLine, Labels: p.Labels(4), Drill: p.Drill(), Series: []Series{{Name: "n", Values: make([]float64, 12)}}})
	if strings.Contains(s, ":drill") {
		t.Errorf("5 minutes not opened: %s", s)
	}
}
