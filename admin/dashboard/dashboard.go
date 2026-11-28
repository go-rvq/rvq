// Package dashboard is the admin's dashboards of widgets: small cards, each
// answering one question — "are there comments waiting?", "did the last
// backup work?" —, each loaded on its own (a slow or failing widget never
// holds the others), over one period chosen at the top (Period: the last
// hour to 90 days, or custom).
//
// One Builder mounts them all. Any package registers its groups of widgets
// and their widgets — d.Group(GroupX).Title(…).Widget(WidgetY, func(w
// *WidgetBuilder) {…}) —, and adds widgets to another's group
// (d.GetGroup(id).Widget(…)); a widget's id is "<group>/<name>", each name
// a public constant of the package it comes from. A Panel is a dashboard:
// the groups and widgets it takes (Group, AllGroups, Widget, Without,
// WithoutGroup), drawn as the admin's home page (Home), as a page of the
// detail of a model's records (Record: its widgets loaded for the record of
// the page, ScopeOf) or in another page (Body: the profile's). Panels too are
// referred to before made (GetPanel); a group or a panel never made panics
// once the admin is built.
//
// A widget is there only when its feature is compiled and on (Widget.Feature)
// and the user may see what it shows (Widget.Allowed); a group with no widget
// left is not drawn. Each user customizes each panel for themself (Prefs):
// none, everything shown; drags the widgets (to any section) and the
// sections (among their siblings), resizes the widgets (columns of the row),
// removes, disables, collapses.
//
// What a widget shows is a list of Blocks — a Metric, a Chart, a List, a
// Status, a Progress —, drawn by the server, the charts by vx-chart. The
// catalog of the widgets (Catalog) puts in the menu a page per group and per
// widget.
package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/features"
	"gorm.io/gorm"
)

// Size is a widget's width, in columns of a grid of 12 (on a phone: 12).
type Size int

const (
	Quarter     Size = 3
	Third       Size = 4
	Half        Size = 6
	TwoThirds   Size = 8
	Full        Size = 12
	SizeDefault      = Third
)

// Group is a context of widgets: a title, an icon, its place among the
// others (Order, the smaller first). An empty Title: no header (the KPIs).
// Made by Builder.Group (GroupBuilder).
type Group struct {
	ID          string
	Title       func(ctx context.Context) string
	Description func(ctx context.Context) string
	Icon        string
	Order       int
	// Untitled: its section on a panel has no heading (a strip of numbers)
	// but while customizing; its title is the catalog's, the help's
	Untitled bool
	// Scope is what its widgets are loaded for (GroupBuilder.Global,
	// Records): none declared, global
	Scope Scopes
}

// Scopes are the scopes a widget is loaded for, as the URI gives them: Global
// — a panel of no record (the home, the profile's) —, and Records — a panel
// of a record of one of these models (Panel.Record: its page). A panel takes
// only the widgets of its scope.
type Scopes struct {
	Global  bool
	Records []*presets.ModelBuilder
}

// accepts says the scopes take a panel of a record of mb (nil: of none).
func (s Scopes) accepts(mb *presets.ModelBuilder) bool {
	if mb == nil {
		return s.Global
	}
	return slices.Contains(s.Records, mb)
}

// scopes is the declared (Global, Records) as Scopes: none declared, global.
func scopes(global bool, records []*presets.ModelBuilder) Scopes {
	return Scopes{Global: global || len(records) == 0, Records: records}
}

// Widget is a card of the dashboard, of a group: made by
// GroupBuilder.Widget (WidgetBuilder).
type Widget struct {
	// ID is its group's and its name, "<group>/<name>": unique, its place in
	// the preferences.
	ID    string
	Group string
	Name  string
	Title func(ctx context.Context) string
	Size  Size
	// Order is its place in its group (the smaller first), before the user's.
	Order int
	// Feature, when given, must be on (compiled and enabled).
	Feature *features.Feature
	// Allowed, when given, says whether the request's user may see it — the
	// permission of what it shows (ListAllowed).
	Allowed func(r *http.Request) bool
	// Periodic: it follows the period chosen (its header says it).
	Periodic bool
	// TTL is how long its blocks are kept for a user and a period (default
	// a minute); its refresh button skips it.
	TTL time.Duration
	// Scope, when given, is what it is loaded for, in place of its group's
	// (WidgetBuilder.Global, Records).
	Scope *Scopes
	// Load is what it shows (ScopeOf: the panel and the record it is of).
	Load func(ctx *web.EventContext, p Period) ([]Block, error)
	// Actions are what its blocks' buttons do (Action.Name), the argument
	// the button's (an id): after one, the widget is loaded again.
	Actions map[string]func(ctx *web.EventContext, arg string) error
	// Link, when given, is where its "open" menu item goes (its listing).
	Link func(ctx *web.EventContext) string
	// Doc, when given, is the document of the widget in place of its ID's
	// (Docs): the one of its kind, for the widgets of the records (Dynamic),
	// whose IDs are made of the record's.
	Doc string
	// Live, when given, are the intervals it reloads itself at, the first
	// by default: its header has a selector of them (the choice kept in the
	// browser), and it reloads while the page is shown. Its blocks are kept
	// (TTL) at most the shortest.
	Live []time.Duration
}

// ListAllowed is Allowed by the permission to list the model mb.
func ListAllowed(mb *presets.ModelBuilder) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		return mb != nil && mb.Permissioner().ReqLister(r).IsAllowed() == nil
	}
}

// PageAllowed is Allowed by the permission of the page pb (a page of the
// admin's, not a model's).
func PageAllowed(pb *presets.PageBuilder) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		return pb != nil && pb.Page().Verifier(r).IsAllowed() == nil
	}
}

// AnyAllowed is allowed when one of fs is.
func AnyAllowed(fs ...func(r *http.Request) bool) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		for _, f := range fs {
			if f != nil && f(r) {
				return true
			}
		}
		return false
	}
}

// Builder is the dashboards of an admin: their groups and widgets —
// registered by any package (Group, GetGroup) — and their panels (Panel,
// GetPanel): the home page, the dashboard of a record (an organization, a
// project), the profile's. Every panel is customized by each user.
type Builder struct {
	p      *presets.Builder
	db     *gorm.DB
	groups []*GroupBuilder
	panels []*Panel
	cache  sync.Map // cacheKey → cached
	docs   func(ctx *web.EventContext, widgetID string) string
	// Now is the clock (tests set it).
	Now func() time.Time
}

// New is the dashboards of p, their preferences kept in db.
func New(p *presets.Builder, db *gorm.DB) *Builder {
	return &Builder{p: p, db: db, Now: time.Now}
}

// Group makes the group id, its owner's (a package): a group made twice is
// an error — two owners. Its widgets are added by GroupBuilder.Widget; other
// packages add theirs by GetGroup.
func (b *Builder) Group(id string) *GroupBuilder {
	g := b.GetGroup(id)
	if g.made {
		panic(fmt.Sprintf("dashboard: the group %q made twice", id))
	}
	g.made = true
	return g
}

// GetGroup is the group id, to add widgets to (GroupBuilder.Widget) or to put
// in a panel: made by its owner (Group) before or after — Install says one
// never made.
func (b *Builder) GetGroup(id string) *GroupBuilder {
	for _, g := range b.groups {
		if g.g.ID == id {
			return g
		}
	}
	g := &GroupBuilder{b: b, g: &Group{ID: id, Scope: scopes(false, nil)}}
	b.groups = append(b.groups, g)
	return g
}

// Groups are the groups made.
func (b *Builder) Groups() []*Group {
	var out []*Group
	for _, g := range b.groups {
		if g.made {
			out = append(out, g.g)
		}
	}
	return out
}

// Widgets are the widgets registered, of every group (not the dynamic ones).
func (b *Builder) Widgets() []*Widget {
	var out []*Widget
	for _, g := range b.groups {
		for _, w := range g.widgets {
			out = append(out, w.w)
		}
	}
	return out
}

// GroupBuilder is a group of widgets: its title, icon, order; its widgets.
type GroupBuilder struct {
	b       *Builder
	g       *Group
	made    bool
	widgets []*WidgetBuilder
	dynamic func() []*WidgetBuilder
	global  bool
	records []*presets.ModelBuilder
}

// ID is the group's id.
func (g *GroupBuilder) ID() string { return g.g.ID }

// Get is the group.
func (g *GroupBuilder) Get() *Group { return g.g }

func (g *GroupBuilder) Title(f func(ctx context.Context) string) *GroupBuilder {
	g.g.Title = f
	return g
}

func (g *GroupBuilder) Description(f func(ctx context.Context) string) *GroupBuilder {
	g.g.Description = f
	return g
}

func (g *GroupBuilder) Icon(v string) *GroupBuilder {
	g.g.Icon = v
	return g
}

// Records says its widgets are loaded for a record of one of mbs: only on
// the panels of their records (Panel.Record), never on the home nor on the
// profile's — unless Global too.
func (g *GroupBuilder) Records(mbs ...*presets.ModelBuilder) *GroupBuilder {
	g.records = append(g.records, mbs...)
	g.g.Scope = scopes(g.global, g.records)
	return g
}

// Global says its widgets are loaded for a panel of no record too (with
// Records: what they show of every record the user reaches, consolidated).
func (g *GroupBuilder) Global() *GroupBuilder {
	g.global = true
	g.g.Scope = scopes(g.global, g.records)
	return g
}

// Untitled draws its section with no heading (but while customizing).
func (g *GroupBuilder) Untitled() *GroupBuilder {
	g.g.Untitled = true
	return g
}

func (g *GroupBuilder) Order(v int) *GroupBuilder {
	g.g.Order = v
	return g
}

// Widget registers the widget name of the group (its ID "<group>/<name>"),
// made by configure; a name registered twice is an error.
func (g *GroupBuilder) Widget(name string, configure func(w *WidgetBuilder)) *GroupBuilder {
	if g.GetWidget(name) != nil {
		panic(fmt.Sprintf("dashboard: the widget %q of the group %q registered twice", name, g.g.ID))
	}
	w := NewWidget(g.g.ID, name)
	if configure != nil {
		configure(w)
	}
	g.widgets = append(g.widgets, w)
	return g
}

// GetWidget is the widget name of the group, nil none.
func (g *GroupBuilder) GetWidget(name string) *WidgetBuilder {
	for _, w := range g.widgets {
		if w.w.Name == name {
			return w
		}
	}
	return nil
}

// Widgets are the widgets of the group: the registered, then the dynamic.
func (g *GroupBuilder) Widgets() []*WidgetBuilder {
	if g.dynamic == nil {
		return g.widgets
	}
	out := slices.Clone(g.widgets)
	for _, w := range g.dynamic() {
		if w.w.Group == g.g.ID {
			out = append(out, w)
		}
	}
	return out
}

// Dynamic gives the widgets of the group known at run time besides the
// registered ones — those of the records of an app (a process's), as many as
// they are now (NewWidget, of this group): f is asked each time a panel needs
// them.
func (g *GroupBuilder) Dynamic(f func() []*WidgetBuilder) *GroupBuilder {
	g.dynamic = f
	return g
}

// WidgetBuilder is a widget: what of it is set.
type WidgetBuilder struct {
	w       *Widget
	global  bool
	records []*presets.ModelBuilder
}

// Records says it is loaded for a record of one of mbs, in place of its
// group's scope (GroupBuilder.Records).
func (w *WidgetBuilder) Records(mbs ...*presets.ModelBuilder) *WidgetBuilder {
	w.records = append(w.records, mbs...)
	sc := scopes(w.global, w.records)
	w.w.Scope = &sc
	return w
}

// Global says it is loaded for a panel of no record, in place of its
// group's scope (GroupBuilder.Global).
func (w *WidgetBuilder) Global() *WidgetBuilder {
	w.global = true
	sc := scopes(w.global, w.records)
	w.w.Scope = &sc
	return w
}

// NewWidget is the widget name of the group group, not registered (a
// dynamic one: GroupBuilder.Dynamic).
func NewWidget(group, name string) *WidgetBuilder {
	return &WidgetBuilder{w: &Widget{ID: group + "/" + name, Group: group, Name: name, Size: SizeDefault}}
}

// ID is its id, "<group>/<name>".
func (w *WidgetBuilder) ID() string { return w.w.ID }

// Get is the widget.
func (w *WidgetBuilder) Get() *Widget { return w.w }

func (w *WidgetBuilder) Title(f func(ctx context.Context) string) *WidgetBuilder {
	w.w.Title = f
	return w
}

func (w *WidgetBuilder) Size(v Size) *WidgetBuilder {
	w.w.Size = v
	return w
}

func (w *WidgetBuilder) Order(v int) *WidgetBuilder {
	w.w.Order = v
	return w
}

func (w *WidgetBuilder) Feature(v *features.Feature) *WidgetBuilder {
	w.w.Feature = v
	return w
}

func (w *WidgetBuilder) Allowed(f func(r *http.Request) bool) *WidgetBuilder {
	w.w.Allowed = f
	return w
}

// Periodic: it follows the period chosen.
func (w *WidgetBuilder) Periodic() *WidgetBuilder {
	w.w.Periodic = true
	return w
}

func (w *WidgetBuilder) TTL(v time.Duration) *WidgetBuilder {
	w.w.TTL = v
	return w
}

func (w *WidgetBuilder) Load(f func(ctx *web.EventContext, p Period) ([]Block, error)) *WidgetBuilder {
	w.w.Load = f
	return w
}

// Action is what the button of its blocks named name does.
func (w *WidgetBuilder) Action(name string, f func(ctx *web.EventContext, arg string) error) *WidgetBuilder {
	if w.w.Actions == nil {
		w.w.Actions = map[string]func(ctx *web.EventContext, arg string) error{}
	}
	w.w.Actions[name] = f
	return w
}

func (w *WidgetBuilder) Link(f func(ctx *web.EventContext) string) *WidgetBuilder {
	w.w.Link = f
	return w
}

func (w *WidgetBuilder) Doc(v string) *WidgetBuilder {
	w.w.Doc = v
	return w
}

// Live are the intervals it reloads itself at (the first by default).
func (w *WidgetBuilder) Live(ds ...time.Duration) *WidgetBuilder {
	w.w.Live = ds
	return w
}

// widget is the widget of id, of any group (the dynamic too).
func (b *Builder) widget(id string) *Widget {
	group, name, ok := strings.Cut(id, "/")
	if !ok {
		return nil
	}
	for _, g := range b.groups {
		if g.g.ID != group {
			continue
		}
		for _, w := range g.Widgets() {
			if w.w.Name == name {
				return w.w
			}
		}
	}
	return nil
}

// visible says whether w is there for r: its feature on, allowed.
func (w *Widget) visible(r *http.Request) bool {
	if w.Feature != nil && !w.Feature.Enabled() {
		return false
	}
	return w.Allowed == nil || w.Allowed(r)
}

// DB is the dashboard's database.
func (b *Builder) DB() *gorm.DB { return b.db }
