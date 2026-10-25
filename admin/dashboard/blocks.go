package dashboard

import (
	"fmt"
	"math"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// Block is a part of what a widget shows: a Metric, a Chart, a List, a
// Status, a Progress, a Note, or Columns of them.
type Block interface {
	render(rc *renderCtx) h.HTMLComponent
}

// renderCtx is what a block is drawn with: the widget (its actions), the
// event, the words.
type renderCtx struct {
	b      *Builder
	w      *Widget
	ctx    *web.EventContext
	msgs   *Messages
	period Period
}

// Level is how a row, an item says it: ok, warning, error, info (or none).
type Level string

const (
	LevelNone    Level = ""
	LevelOK      Level = "success"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
	LevelInfo    Level = "info"
)

// Action is a button of a row or an item: the widget's Actions[Name] with
// Arg, the widget loaded again after it.
type Action struct {
	Name  string
	Arg   string
	Label string
	Icon  string
	// Removes: done, the row is gone (read, approved): it leaves at the
	// click — slides, fades, its height closes —, the action runs after; a
	// failure brings it back (the widget drawn again)
	Removes bool
}

// rowLeaveMS is how long a row takes to leave (Action.Removes).
const rowLeaveMS = 280

func (rc *renderCtx) action(a *Action) h.HTMLComponent {
	if a == nil {
		return nil
	}
	ev := web.Plaid().EventFunc(EventAction).
		Query("widget", rc.w.ID).
		Query("action", a.Name).
		Query("arg", a.Arg).
		Query("period", rc.period.Key).
		Query("from", rc.ctx.R.FormValue("from")).
		Query("to", rc.ctx.R.FormValue("to")).
		Go()
	click := ev
	if a.Removes {
		// the row leaves first, the action after it
		click = fmt.Sprintf(`(e) => {
	const row = e.currentTarget.closest("[data-dashboard-row], [data-dashboard-status]");
	if (row) {
		row.style.height = row.offsetHeight + "px";
		void row.offsetHeight;
		row.classList.add("dash-row--leaving");
	}
	e.currentTarget.disabled = true;
	window.setTimeout(() => { %s }, %d);
}`, ev, rowLeaveMS)
	}
	btn := v.VBtn(a.Label).Size(v.SizeXSmall).Variant(v.VariantTonal).Color("primary").
		Attr("@click.stop", click).Attr("data-dashboard-action", a.Name)
	if a.Icon != "" {
		btn.PrependIcon(a.Icon)
	}
	return btn
}

func link(href, text string) h.HTMLComponent {
	if href == "" || text == "" {
		return nil
	}
	return h.A(h.Text(text + " →")).Href(href).Class("text-caption text-primary text-decoration-none d-inline-block mt-2")
}

// Metric is a number of the period, how it changed from the period before
// (Previous, when known), its values by step (Spark).
type Metric struct {
	Value    float64
	Previous *float64
	// Format writes the value (default: a number, thousands grouped).
	Format  func(float64) string
	Spark   []float64
	Caption string
	Link    string
	// DownIsGood: a fall is the good news (failures).
	DownIsGood bool
}

func (m *Metric) render(rc *renderCtx) h.HTMLComponent {
	f := m.Format
	if f == nil {
		f = formatNumber
	}
	var delta h.HTMLComponent
	if m.Previous != nil {
		d, up := 0.0, m.Value >= *m.Previous
		if *m.Previous != 0 {
			d = math.Abs(m.Value-*m.Previous) / *m.Previous * 100
		}
		good := up != m.DownIsGood
		color := "text-success"
		if !good {
			color = "text-error"
		}
		arrow := "▲"
		if !up {
			arrow = "▼"
		}
		txt := fmt.Sprintf("%s %.0f%%", arrow, d)
		if *m.Previous == 0 {
			txt = arrow + " " + rc.msgs.New
		}
		if m.Value == *m.Previous {
			txt, color = "= 0%", "text-medium-emphasis"
		}
		delta = h.Div(h.Text(txt)).Class("text-caption "+color).Title(rc.msgs.VsPrevious).Attr("data-dashboard-delta", true)
	}
	return h.Div(
		h.Div(
			h.Div(h.Text(f(m.Value))).Class("text-h5 font-weight-bold").Attr("data-dashboard-value", true),
			delta,
			h.If(m.Caption != "", h.Div(h.Text(m.Caption)).Class("text-caption text-medium-emphasis")),
		),
		h.If(len(m.Spark) > 1, h.Div(h.RawHTML(sparkSVG(m.Spark))).Class("ms-2 flex-shrink-0")),
	).Class("d-flex align-center justify-space-between").AppendChildren(link(m.Link, rc.msgs.Open))
}

// ChartKind is how a chart draws its series.
type ChartKind string

const (
	ChartLine    ChartKind = "line"
	ChartBars    ChartKind = "bars"
	ChartStacked ChartKind = "stacked"
	// ChartHBars: one series as horizontal bars, a label each (a ranking).
	ChartHBars ChartKind = "hbars"
	// ChartDonut: one series as the parts of a ring, a label each.
	ChartDonut ChartKind = "donut"
)

// Series is a named list of values, a color of the theme (or CSS).
type Series struct {
	Name   string
	Color  string
	Values []float64
}

// Chart draws Series: by step (line, bars, stacked: Labels the steps'), or
// by label (hbars, donut: Labels the parts', one series).
type Chart struct {
	Kind   ChartKind
	Title  string
	Labels []string
	Series []Series
	// Colors are the donut's parts' (default the palette).
	Colors   []string
	Link     string
	LinkText string
}

func (c *Chart) legend() h.HTMLComponent {
	if len(c.Series) < 2 {
		return nil
	}
	var items h.HTMLComponents
	for i, s := range c.Series {
		items = append(items, h.Span("").Children(
			h.Span("").Style(fmt.Sprintf("display:inline-block;width:10px;height:10px;border-radius:2px;margin-right:4px;background:%s", seriesColor(s, i))),
			h.Text(s.Name),
		).Class("text-caption text-medium-emphasis me-3"))
	}
	return h.Div(items...).Class("mb-1")
}

func (c *Chart) render(rc *renderCtx) h.HTMLComponent {
	lt := c.LinkText
	if lt == "" {
		lt = rc.msgs.Open
	}
	empty := true
	for _, s := range c.Series {
		if Total(s.Values) != 0 {
			empty = false
		}
	}
	if empty {
		return h.Components(h.Div(h.Text(rc.msgs.NoData)).Class("text-medium-emphasis text-body-2 py-6 text-center").Attr("data-dashboard-empty", true), link(c.Link, lt))
	}
	switch c.Kind {
	case ChartHBars:
		var rows h.HTMLComponents
		var s Series
		if len(c.Series) > 0 {
			s = c.Series[0]
		}
		mx := maxOf([]Series{s}, false)
		for i, val := range s.Values {
			pct := 0.0
			if mx > 0 {
				pct = val / mx * 100
			}
			rows = append(rows, h.Div(
				h.Div(h.Text(labelAt(c.Labels, i))).Class("text-body-2 text-truncate").Style("width:38%").Title(labelAt(c.Labels, i)),
				h.Div(h.Div().Style(fmt.Sprintf("width:%.1f%%;height:14px;border-radius:3px;background:%s", pct, seriesColor(s, 0)))).Style("flex:1"),
				h.Div(h.Text(formatNumber(val))).Class("text-caption text-medium-emphasis text-end").Style("width:52px"),
			).Class("d-flex align-center ga-2 mb-1").Attr("data-dashboard-bar", true))
		}
		return h.Div(rows...).AppendChildren(link(c.Link, lt))
	case ChartDonut:
		var legend h.HTMLComponents
		var s Series
		if len(c.Series) > 0 {
			s = c.Series[0]
		}
		for i, val := range s.Values {
			col := cssColor(palette[i%len(palette)])
			if i < len(c.Colors) && c.Colors[i] != "" {
				col = cssColor(c.Colors[i])
			}
			legend = append(legend, h.Div(
				h.Span("").Style(fmt.Sprintf("display:inline-block;width:10px;height:10px;border-radius:2px;margin-right:6px;background:%s", col)),
				h.Text(fmt.Sprintf("%s (%s)", labelAt(c.Labels, i), formatNumber(val))),
			).Class("text-caption text-medium-emphasis"))
		}
		return h.Div(
			h.Div(h.RawHTML(donutSVG(c))).Class("flex-shrink-0"),
			h.Div(legend...).Class("ms-4"),
		).Class("d-flex align-center flex-wrap").AppendChildren(link(c.Link, lt))
	case ChartBars, ChartStacked:
		return h.Components(c.legend(), h.Div(h.RawHTML(barsSVG(c))).Class("dashboard-chart"), link(c.Link, lt))
	default:
		return h.Components(c.legend(), h.Div(h.RawHTML(lineSVG(c))).Class("dashboard-chart"), link(c.Link, lt))
	}
}

// Row is a line of a List: its text (a link), under it Sub, at its right
// Right, a dot of its Level, a button (Action).
type Row struct {
	Text   string
	Sub    string
	Right  string
	Link   string
	Level  Level
	Action *Action
}

// List is rows (the top N, the newest), Empty said with none.
type List struct {
	Rows     []Row
	Empty    string
	Link     string
	LinkText string
}

func dot(l Level) h.HTMLComponent {
	if l == LevelNone {
		return nil
	}
	return h.Span("").Class("flex-shrink-0").Style(fmt.Sprintf("display:inline-block;width:8px;height:8px;border-radius:50%%;background:%s;margin-right:8px", cssColor(string(l))))
}

func (l *List) render(rc *renderCtx) h.HTMLComponent {
	lt := l.LinkText
	if lt == "" {
		lt = rc.msgs.SeeAll
	}
	if len(l.Rows) == 0 {
		empty := l.Empty
		if empty == "" {
			empty = rc.msgs.NoData
		}
		return h.Components(h.Div(h.Text(empty)).Class("text-medium-emphasis text-body-2 py-6 text-center").Attr("data-dashboard-empty", true), link(l.Link, lt))
	}
	var rows h.HTMLComponents
	for _, r := range l.Rows {
		var text h.HTMLComponent = h.Text(r.Text)
		if r.Link != "" {
			text = h.A(h.Text(r.Text)).Href(r.Link).Class("text-decoration-none").Style("color:inherit")
		}
		rows = append(rows, h.Div(
			dot(r.Level),
			h.Div(
				h.Div(text).Class("text-body-2 text-truncate").Title(r.Text),
				h.If(r.Sub != "", h.Div(h.Text(r.Sub)).Class("text-caption text-medium-emphasis text-truncate")),
			).Style("min-width:0;flex:1"),
			h.If(r.Right != "", h.Div(h.Text(r.Right)).Class("text-caption text-medium-emphasis ms-2 flex-shrink-0")),
			h.If(r.Action != nil, h.Div(rc.action(r.Action)).Class("ms-2 flex-shrink-0")),
		).Class("d-flex align-center py-1 dash-row").Style("border-bottom:1px solid rgba(var(--v-border-color), var(--v-border-opacity))").Attr("data-dashboard-row", true))
	}
	return h.Div(rowLeaveStyle, h.Div(rows...)).AppendChildren(link(l.Link, lt))
}

// StatusItem is a state: a dot of its Level, a title, a detail, a link, a
// button.
type StatusItem struct {
	Level  Level
	Title  string
	Detail string
	Link   string
	Action *Action
}

// Status is states, said one under the other.
type Status struct {
	Items []StatusItem
	Empty string
}

func (s *Status) render(rc *renderCtx) h.HTMLComponent {
	if len(s.Items) == 0 {
		empty := s.Empty
		if empty == "" {
			empty = rc.msgs.NoData
		}
		return h.Div(h.Text(empty)).Class("text-medium-emphasis text-body-2 py-6 text-center").Attr("data-dashboard-empty", true)
	}
	var items h.HTMLComponents
	for _, it := range s.Items {
		var title h.HTMLComponent = h.Text(it.Title)
		if it.Link != "" {
			title = h.A(h.Text(it.Title)).Href(it.Link).Class("text-decoration-none").Style("color:inherit")
		}
		items = append(items, h.Div(
			h.Div(dot(it.Level)).Class("pt-2"),
			h.Div(
				h.Div(title).Class("text-body-2 font-weight-bold"),
				h.If(it.Detail != "", h.Div(h.Text(it.Detail)).Class("text-caption text-medium-emphasis")),
				h.If(it.Action != nil, h.Div(rc.action(it.Action)).Class("mt-1")),
			),
		).Class("d-flex mb-2 dash-row").Attr("data-dashboard-status", string(it.Level)))
	}
	return h.Div(rowLeaveStyle, h.Div(items...))
}

// Bar is a share done of a total (a language's texts translated).
type Bar struct {
	Label       string
	Done, Total float64
	Link        string
}

// Progress is bars, one under the other.
type Progress struct {
	Bars     []Bar
	Link     string
	LinkText string
}

func (p *Progress) render(rc *renderCtx) h.HTMLComponent {
	var rows h.HTMLComponents
	for _, b := range p.Bars {
		pct := 0.0
		if b.Total > 0 {
			pct = b.Done / b.Total * 100
		}
		color := "primary"
		switch {
		case pct >= 100:
			color = "success"
		case pct < 80:
			color = "warning"
		}
		rows = append(rows, h.Div(
			h.Div(h.Text(b.Label)).Class("text-body-2 text-truncate").Style("width:30%"),
			v.VProgressLinear().ModelValue(pct).Color(color).Height(10).Rounded(true).Style("flex:1"),
			h.Div(h.Text(fmt.Sprintf("%s/%s", formatNumber(b.Done), formatNumber(b.Total)))).Class("text-caption text-medium-emphasis text-end").Style("width:80px"),
		).Class("d-flex align-center ga-2 mb-2").Attr("data-dashboard-progress", b.Label))
	}
	lt := p.LinkText
	if lt == "" {
		lt = rc.msgs.Open
	}
	return h.Div(rows...).AppendChildren(link(p.Link, lt))
}

// Note is a line of text (a hint, where the data comes from).
type Note struct {
	Text  string
	Level Level
}

func (n *Note) render(rc *renderCtx) h.HTMLComponent {
	class := "text-caption text-medium-emphasis mt-1"
	if n.Level != LevelNone {
		class = "text-caption mt-1 text-" + string(n.Level)
	}
	return h.Div(h.Text(n.Text)).Class(class)
}

// Columns are blocks side by side (one under the other on a phone).
type Columns struct{ Blocks []Block }

func (c *Columns) render(rc *renderCtx) h.HTMLComponent {
	var cols h.HTMLComponents
	for _, b := range c.Blocks {
		cols = append(cols, h.Div(b.render(rc)).Style("flex:1 1 200px;min-width:0"))
	}
	return h.Div(cols...).Class("d-flex flex-wrap ga-4")
}

// Heading is a small title inside a widget (between blocks).
type Heading struct{ Text string }

func (hd *Heading) render(rc *renderCtx) h.HTMLComponent {
	return h.Div(h.Text(hd.Text)).Class("text-caption font-weight-bold text-uppercase text-medium-emphasis mt-2 mb-1")
}

// Stack are blocks one under the other (a column of Columns).
type Stack struct{ Blocks []Block }

func (s *Stack) render(rc *renderCtx) h.HTMLComponent {
	var out h.HTMLComponents
	for _, b := range s.Blocks {
		if b != nil {
			out = append(out, b.render(rc))
		}
	}
	return h.Div(out...)
}

// rowLeaveStyle is how a row leaves (Action.Removes): to the side, fading,
// its height closing.
var rowLeaveStyle = h.Tag("component").Attr(":is", "'style'").Children(h.RawHTML(fmt.Sprintf(`
.dash-row { transition: opacity %[1]dms ease, transform %[1]dms ease, height %[1]dms ease, padding %[1]dms ease, margin %[1]dms ease; }
.dash-row--leaving { opacity: 0; transform: translateX(32px); height: 0 !important; padding-top: 0 !important; padding-bottom: 0 !important; margin: 0 !important; border-color: transparent !important; overflow: hidden; }
`, rowLeaveMS)))
