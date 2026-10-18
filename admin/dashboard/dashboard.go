// Package dashboard is the admin's dashboard of widgets: small cards, each
// answering one question — "are there comments waiting?", "did the last
// backup work?" —, grouped by context (Group), each loaded on its own (a
// slow or failing widget never holds the others), over one period chosen at
// the top (Period: the last hour to 90 days, or custom).
//
// A widget is there only when its feature is compiled and on (Widget.Feature)
// and the user may see what it shows (Widget.Allowed); a group with no widget
// left is not drawn. Each user hides, orders and collapses them (Prefs).
//
// What a widget shows is a list of Blocks — a Metric, a Chart, a List, a
// Status, a Progress —, drawn by the server, the charts as SVG in the theme's
// colors (no script of their own).
package dashboard

import (
	"context"
	"net/http"
	"slices"
	"sync"
	"time"

	h "github.com/go-rvq/htmlgo"
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
type Group struct {
	ID    string
	Title func(ctx context.Context) string
	Icon  string
	Order int
}

// Widget is a card of the dashboard.
type Widget struct {
	// ID is unique: its place in the preferences, its documentation
	// (docs/dashboard/widgets/<ID>).
	ID    string
	Group string
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
	// Load is what it shows.
	Load func(ctx *web.EventContext, p Period) ([]Block, error)
	// Actions are what its blocks' buttons do (Action.Name), the argument
	// the button's (an id): after one, the widget is loaded again.
	Actions map[string]func(ctx *web.EventContext, arg string) error
	// Link, when given, is where its "open" menu item goes (its listing).
	Link func(ctx *web.EventContext) string
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

// Builder is the dashboard of an admin.
type Builder struct {
	p       *presets.Builder
	db      *gorm.DB
	groups  []*Group
	widgets []*Widget
	cache   sync.Map // cacheKey → cached
	docs    func(ctx *web.EventContext, widgetID string) string
	intro   func(ctx *web.EventContext) h.HTMLComponent
	// Now is the clock (tests set it).
	Now func() time.Time
}

// New is the dashboard of p, its preferences kept in db.
func New(p *presets.Builder, db *gorm.DB) *Builder {
	return &Builder{p: p, db: db, Now: time.Now}
}

// Group adds (or replaces) a group.
func (b *Builder) Group(g *Group) *Builder {
	if i := slices.IndexFunc(b.groups, func(o *Group) bool { return o.ID == g.ID }); i >= 0 {
		b.groups[i] = g
		return b
	}
	b.groups = append(b.groups, g)
	return b
}

// Widget adds (or replaces) widgets.
func (b *Builder) Widget(ws ...*Widget) *Builder {
	for _, w := range ws {
		if w.Size == 0 {
			w.Size = SizeDefault
		}
		if i := slices.IndexFunc(b.widgets, func(o *Widget) bool { return o.ID == w.ID }); i >= 0 {
			b.widgets[i] = w
			continue
		}
		b.widgets = append(b.widgets, w)
	}
	return b
}

// Widgets are the widgets registered.
func (b *Builder) Widgets() []*Widget { return b.widgets }

// Groups are the groups registered.
func (b *Builder) Groups() []*Group { return b.groups }

// widget is the widget of id.
func (b *Builder) widget(id string) *Widget {
	for _, w := range b.widgets {
		if w.ID == id {
			return w
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
