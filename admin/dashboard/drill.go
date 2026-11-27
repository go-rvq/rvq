package dashboard

import (
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/charts"
	"github.com/go-rvq/rvq/web"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// A point of a chart opened (Chart.Drill): the widget is loaded again with
// its period — a day of 30, by its hours; an hour, by its 5 minutes —, and
// the periods it came from (the trail, ?trail=) are a way back at its top:
// "30 days › 30/10 › 30/10 14:00–15:00", each a link.

// crumb is a period of the trail: its key and, custom, its bounds.
type crumb struct{ key, from, to string }

// trailSep and crumbSep separate the trail's periods and their parts.
const (
	trailSep = "|"
	crumbSep = "~"
)

func crumbOf(p Period) crumb {
	from, to := p.Bounds()
	return crumb{p.Key, from, to}
}

func (c crumb) String() string { return c.key + crumbSep + c.from + crumbSep + c.to }

func parseTrail(s string) (out []crumb) {
	for _, part := range strings.Split(s, trailSep) {
		f := strings.Split(part, crumbSep)
		if len(f) != 3 || f[0] == "" {
			continue
		}
		out = append(out, crumb{f[0], f[1], f[2]})
	}
	return
}

func encodeTrail(cs []crumb) string {
	parts := make([]string, len(cs))
	for i, c := range cs {
		parts[i] = c.String()
	}
	return strings.Join(parts, trailSep)
}

// trail is the request's way back.
func (rc *renderCtx) trail() []crumb { return parseTrail(rc.ctx.R.FormValue("trail")) }

// reload is the event that loads the widget again with the period c, the
// way back trail.
func (rc *renderCtx) reload(c crumb, trail []crumb) *web.VueEventTagBuilder {
	ev := web.Plaid().EventFunc(EventLoad).Query("widget", rc.w.ID).Query("refresh", "1").Query("period", c.key)
	if c.key == PeriodCustom {
		ev.Query("from", c.from).Query("to", c.to)
	}
	if len(trail) > 0 {
		ev.Query("trail", encodeTrail(trail))
	}
	return ev
}

// crumbs is the way back of a widget opened at a point: the periods it came
// from, links, and its own; nothing when it was not opened.
func (rc *renderCtx) crumbs() h.HTMLComponent {
	trail := rc.trail()
	if len(trail) == 0 {
		return nil
	}
	var items h.HTMLComponents
	for i, c := range trail {
		p := NewPeriod(c.key, rc.b.now(), c.from, c.to)
		items = append(items,
			h.A(h.Text(p.Label(rc.msgs))).Href("#").Class("text-primary text-decoration-none").
				Attr("@click.prevent", rc.reload(c, trail[:i]).Go()).Attr("data-dashboard-crumb", i),
			h.Span(" › ").Class("text-medium-emphasis"),
		)
	}
	items = append(items, h.Span(rc.period.Label(rc.msgs)).Class("font-weight-medium").Attr("data-dashboard-crumb-current", true))
	return h.Div(items...).Class("text-caption mb-1").Attr("data-dashboard-trail", true)
}

// vxChart is c as a vx-chart; a point of a period (Drill), clicked, opens
// it.
func (rc *renderCtx) vxChart(c *Chart) h.HTMLComponent {
	series := make([]vx.VXChartSeries, len(c.Series))
	for i, s := range c.Series {
		series[i] = vx.VXChartSeries{Name: s.Name, Color: s.Color, Values: s.Values}
	}
	ch := vx.VXChart().Src(charts.Src(rc.b.p)).Kind(string(c.Kind)).Title(c.Title).
		Labels(c.Labels).Series(series...).Hints(c.Hints).Colors(c.Colors).ResetLabel(rc.msgs.ResetZoom)
	if len(c.Drill) > 0 {
		drill := make([]any, len(c.Drill))
		for i, p := range c.Drill {
			if p == nil {
				continue
			}
			from, to := p.Bounds()
			drill[i] = map[string]string{"from": from, "to": to}
		}
		trail := encodeTrail(append(rc.trail(), crumbOf(rc.period)))
		ch.Drill(drill).Attr("@point", web.Plaid().EventFunc(EventLoad).Query("widget", rc.w.ID).Query("refresh", "1").
			Query("period", PeriodCustom).Query("from", web.Var("$event.from")).Query("to", web.Var("$event.to")).
			Query("trail", trail).Go())
	}
	return ch
}

// now is the clock of the dashboard.
func (b *Builder) now() time.Time {
	if b == nil || b.Now == nil {
		return time.Now()
	}
	return b.Now()
}
