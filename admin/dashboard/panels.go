package dashboard

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"slices"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// PanelHome is the panel of the admin's home page (Panel.Home).
const PanelHome = "home"

// Panel is a dashboard: the widgets of the groups it takes (Group,
// AllGroups) and the ones it takes alone (Widget), but those left out
// (Without, WithoutGroup); drawn as the admin's home page (Home), as a page
// of the detail of a model's records (Record) or in a page of its own
// (Body). Each user customizes each panel (Prefs): their layout the same for
// every record of it; none, everything shown.
type Panel struct {
	*Builder
	id           string
	made         bool
	title        func(ctx context.Context) string
	allGroups    bool
	groupIDs     []string
	widgetIDs    []string
	without      []string
	withoutGroup []string
	intro        func(ctx *web.EventContext) h.HTMLComponent
	// record: the model whose records' detail it is a page of (Record)
	record *presets.ModelBuilder
	path   string
	// recordParam: the record of the request out of a page of the record
	// (the page of a widget: its query's), given to its events (event)
	recordParam string
}

// Panel makes the panel id, its owner's: made twice is an error.
func (b *Builder) Panel(id string) *Panel {
	pn := b.GetPanel(id)
	if pn.made {
		panic(fmt.Sprintf("dashboard: the panel %q made twice", id))
	}
	pn.made = true
	return pn
}

// GetPanel is the panel id, to put widgets in: made by its owner (Panel)
// before or after — Install says one never made.
func (b *Builder) GetPanel(id string) *Panel {
	for _, pn := range b.panels {
		if pn.id == id {
			return pn
		}
	}
	pn := &Panel{Builder: b, id: id}
	b.panels = append(b.panels, pn)
	return pn
}

// Panels are the panels made.
func (b *Builder) Panels() []*Panel {
	var out []*Panel
	for _, pn := range b.panels {
		if pn.made {
			out = append(out, pn)
		}
	}
	return out
}

// ID is the panel's id.
func (pn *Panel) ID() string { return pn.id }

// Title is the panel's title (its page's; the home's: the dashboard's).
func (pn *Panel) Title(f func(ctx context.Context) string) *Panel {
	pn.title = f
	return pn
}

// TTitle is the title of the panel.
func (pn *Panel) TTitle(ctx context.Context) string {
	if pn.title != nil {
		return pn.title(ctx)
	}
	return MessagesOf(ctx).Dashboard
}

// AllGroups: every group's widgets (the home).
func (pn *Panel) AllGroups() *Panel {
	pn.allGroups = true
	return pn
}

// Group takes the widgets of g — the ones registered after too, by any
// package.
func (pn *Panel) Group(g *GroupBuilder) *Panel {
	if !slices.Contains(pn.groupIDs, g.ID()) {
		pn.groupIDs = append(pn.groupIDs, g.ID())
	}
	return pn
}

// Widget takes the widgets ws alone (not their whole groups).
func (pn *Panel) Widget(ws ...*WidgetBuilder) *Panel {
	for _, w := range ws {
		if !slices.Contains(pn.widgetIDs, w.ID()) {
			pn.widgetIDs = append(pn.widgetIDs, w.ID())
		}
	}
	return pn
}

// Without leaves the widgets ws out (of a group it takes).
func (pn *Panel) Without(ws ...*WidgetBuilder) *Panel {
	for _, w := range ws {
		pn.without = append(pn.without, w.ID())
	}
	return pn
}

// WithoutGroup leaves the group g out (of AllGroups).
func (pn *Panel) WithoutGroup(g *GroupBuilder) *Panel {
	pn.withoutGroup = append(pn.withoutGroup, g.ID())
	return pn
}

// Intro, when given, is drawn under its top (a greeting).
func (pn *Panel) Intro(f func(ctx *web.EventContext) h.HTMLComponent) *Panel {
	pn.intro = f
	return pn
}

// Home makes the panel the admin's home page.
func (pn *Panel) Home() *Panel {
	pn.p.HomePageFunc(pn.Page)
	return pn
}

// Record makes the panel a page of the detail of the records of mb, at path
// ("/painel"): its scope the record of the page (ScopeOf). Its page is
// returned by RecordPage (for the menu: a group's page, the selector's).
func (pn *Panel) Record(mb *presets.ModelBuilder, path string) *Panel {
	pn.record, pn.path = mb, path
	mb.Detailing().PagesRegistrator().New(presets.HttpPage(path).TitleFunc(pn.TTitle).MenuIcon("mdi-view-dashboard-outline")).
		RecordLayout(func(ctx *web.EventContext, _ model.ID, _ []model.ID, _ func() (any, error)) (r web.PageResponse, err error) {
			return pn.Page(ctx)
		}).Build()
	return pn
}

// RecordPage is the page of the panel of a record (Record), nil none.
func (pn *Panel) RecordPage() *presets.HttpPageBuilder {
	if pn.record == nil {
		return nil
	}
	return pn.record.Detailing().PagesRegistrator().GetHttpPage(pn.path)
}

// takes says the panel takes the widget w.
func (pn *Panel) takes(w *Widget) bool {
	if slices.Contains(pn.without, w.ID) {
		return false
	}
	if slices.Contains(pn.widgetIDs, w.ID) || slices.Contains(pn.groupIDs, w.Group) {
		return true
	}
	return pn.allGroups && !slices.Contains(pn.withoutGroup, w.Group)
}

// all are the widgets of the panel (the dynamic too), of the groups made.
func (pn *Panel) all() []*Widget {
	var out []*Widget
	for _, g := range pn.groups {
		if !g.made {
			continue
		}
		for _, w := range g.Widgets() {
			if pn.takes(w.w) {
				out = append(out, w.w)
			}
		}
	}
	return out
}

// Widgets are the widgets of the panel (the dynamic too).
func (pn *Panel) Widgets() []*Widget { return pn.all() }

// RecordModel is the model whose records' detail the panel is a page of
// (Record), nil none.
func (pn *Panel) RecordModel() *presets.ModelBuilder { return pn.record }

// widget is the widget id of the panel, nil none.
func (pn *Panel) widget(id string) *Widget {
	for _, w := range pn.all() {
		if w.ID == id {
			return w
		}
	}
	return nil
}

// sortedGroups are the groups of the panel's widgets, by Order.
func (pn *Panel) sortedGroups() []*Group {
	var gs []*Group
	for _, w := range pn.all() {
		if g := pn.Builder.groupOf(w.Group); g != nil && !slices.Contains(gs, g) {
			gs = append(gs, g)
		}
	}
	slices.SortStableFunc(gs, func(a, c *Group) int { return cmp.Compare(a.Order, c.Order) })
	return gs
}

// groupOf is the group id of the panel, nil none.
func (pn *Panel) groupOf(id string) *Group {
	for _, g := range pn.sortedGroups() {
		if g.ID == id {
			return g
		}
	}
	return nil
}

// groupOf is the group id made, nil none.
func (b *Builder) groupOf(id string) *Group {
	for _, g := range b.groups {
		if g.made && g.g.ID == id {
			return g.g
		}
	}
	return nil
}

// Scope is what a widget is loaded for: its panel and, of a panel of a
// record (Record), the record's id.
type Scope struct {
	Panel  string
	Record string
}

type scopeKey struct{}

// ScopeOf is the scope of the request of a widget (its Load): the panel and
// the record it is of.
func ScopeOf(ctx *web.EventContext) Scope {
	s, _ := ctx.R.Context().Value(scopeKey{}).(Scope)
	return s
}

// WithScope is r in the scope s (what a panel's request is in: ScopeOf) —
// a test's request of a widget.
func WithScope(r *http.Request, s Scope) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), scopeKey{}, s))
}

// scoped is ctx in the scope of the panel: its record the page's (the
// {id} of its path — the events are at the page's URL).
func (pn *Panel) scoped(ctx *web.EventContext) *web.EventContext {
	s := Scope{Panel: pn.id}
	if pn.record != nil {
		s.Record = ctx.R.PathValue(presets.ParamID)
		if s.Record == "" {
			s.Record = ctx.R.FormValue(paramRecord)
		}
	}
	ctx.R = WithScope(ctx.R, s)
	return ctx
}

// paramRecord is the record of a panel's request out of a page of the record
// (the page of a widget).
const paramRecord = "record"

// in is the panel for the request of ctx: a copy whose events carry its
// record when the request is not at a page of it (a widget's page).
func (pn *Panel) in(ctx *web.EventContext) *Panel {
	s := ScopeOf(ctx)
	if pn.record == nil || s.Record == "" || ctx.R.PathValue(presets.ParamID) != "" {
		return pn
	}
	c := *pn
	c.recordParam = s.Record
	return &c
}

// paramPanel is the panel of an event (its query).
const paramPanel = "panel"

// panelOf is the panel of the request (its query's), the home's by default;
// nil one not made.
func (b *Builder) panelOf(r *http.Request) *Panel {
	id := cmp.Or(r.FormValue(paramPanel), PanelHome)
	for _, pn := range b.panels {
		if pn.id == id && pn.made {
			return pn
		}
	}
	return nil
}

// check panics on a group or a panel referred to (GetGroup, GetPanel) and
// never made: once every package registered its (Install's AfterBuild).
func (b *Builder) check() {
	for _, g := range b.groups {
		if !g.made {
			panic(fmt.Sprintf("dashboard: the group %q is referred to (GetGroup) but never made (Group)", g.g.ID))
		}
	}
	for _, pn := range b.panels {
		if !pn.made {
			panic(fmt.Sprintf("dashboard: the panel %q is referred to (GetPanel) but never made (Panel)", pn.id))
		}
	}
}
