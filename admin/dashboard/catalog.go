package dashboard

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// The catalog of the widgets in the menu (Catalog): a group Widgets, in it a
// group of the menu per group of widgets — a page (presets.PageMenuGroup):
// its widgets —, in it a page per widget (its sub-pages): what it is, the
// panels that take it, a preview.

const (
	// MenuGroupWidgets is the group of the menu of the catalog of the
	// widgets.
	MenuGroupWidgets = "widgets"
	// CatalogPrefix begins the name of the group of the menu (and the path of
	// the page) of a group of widgets in the catalog: "widgets-<group>".
	CatalogPrefix = "widgets-"
)

// CatalogPath is the path of the page of the group of widgets group in the
// catalog (its group of the menu's name the path's segment).
func CatalogPath(group string) string { return "/" + CatalogPrefix + group }

// Catalog puts the catalog of the widgets in the group parent of the menu
// (nil: at its root): made once every package registered its groups
// (BeforeBuild of the admin). Before Build.
func (b *Builder) Catalog(parent *presets.MenuGroupBuilder) *Builder {
	b.p.BeforeBuild(func() { b.catalog(parent) })
	return b
}

func (b *Builder) catalog(parent *presets.MenuGroupBuilder) {
	menu := b.p.MenuGroup(MenuGroupWidgets).Icon("mdi-widgets-outline").
		TitleFunc(func(ctx context.Context) string { return MessagesOf(ctx).Widgets }).
		DescriptionFunc(func(ctx context.Context) string { return MessagesOf(ctx).WidgetsDescription })
	if parent != nil {
		parent.Add(presets.GroupItem(MenuGroupWidgets))
	}
	gs := slices.Clone(b.groups)
	slices.SortStableFunc(gs, func(a, c *GroupBuilder) int { return cmp.Compare(a.g.Order, c.g.Order) })
	for _, g := range gs {
		if !g.made {
			continue
		}
		g := g
		page := b.p.PagesRegistrator().New(presets.HttpPage(CatalogPath(g.ID())).
			TitleFunc(g.title).DescriptionFunc(g.description).MenuIcon(cmp.Or(g.g.Icon, "mdi-view-grid-outline"))).
			Layout(func(ctx *web.EventContext) (r web.PageResponse, err error) {
				r.PageTitle = g.title(ctx.Context())
				r.Body = b.catalogGroup(ctx, g)
				return
			})
		for _, w := range g.widgets {
			w := w.w
			page.SubPage("/" + w.Name).InMenu(true).Title(w.title).
				Layout(func(ctx *web.EventContext) (r web.PageResponse, err error) {
					r.PageTitle = w.title(ctx.Context())
					r.Body = b.catalogWidget(ctx, w)
					return
				}).Build()
		}
		page.Build()
		mg := b.p.PageMenuGroup(page.Page())
		menu.Add(presets.GroupItem(mg.Name()))
	}
}

// title is the group's title, its id with none.
func (g *GroupBuilder) title(ctx context.Context) string {
	if g.g.Title != nil {
		return g.g.Title(ctx)
	}
	return g.g.ID
}

func (g *GroupBuilder) description(ctx context.Context) string {
	if g.g.Description != nil {
		return g.g.Description(ctx)
	}
	return ""
}

// title is the widget's title, its name with none.
func (w *Widget) title(ctx context.Context) string {
	if w.Title != nil {
		return w.Title(ctx)
	}
	return w.Name
}

// panelsOf are the panels made that take w.
func (b *Builder) panelsOf(w *Widget) (out []*Panel) {
	for _, pn := range b.Panels() {
		if pn.takes(w) {
			out = append(out, pn)
		}
	}
	return
}

// catalogGroup is the page of the group g in the catalog: its description,
// its widgets — their titles links to their pages, the panels that take
// them.
func (b *Builder) catalogGroup(ctx *web.EventContext, g *GroupBuilder) h.HTMLComponent {
	msgs := MessagesOf(ctx.Context())
	c := ctx.Context()
	var items h.HTMLComponents
	for _, w := range g.widgets {
		items = append(items, v.VListItem().Title(w.w.title(c)).Subtitle(b.panelTitles(c, msgs, w.w)).
			PrependIcon(cmp.Or(g.g.Icon, "mdi-view-grid-outline")).
			Href(b.catalogHref(g.ID(), w.w.Name)).Attr("data-dashboard-catalog-widget", w.ID()))
	}
	if len(g.widgets) == 0 {
		items = append(items, v.VListItem().Title(msgs.NoWidgets))
	}
	return h.Div(
		h.If(g.description(c) != "", h.P(h.Text(g.description(c))).Class("text-body-2 text-medium-emphasis mb-2")),
		v.VList(items...).Lines("two").Attr("data-dashboard-catalog-group", g.ID()),
	).Class("pa-4")
}

// panelTitles says the panels that take w.
func (b *Builder) panelTitles(c context.Context, msgs *Messages, w *Widget) string {
	var ts []string
	for _, pn := range b.panelsOf(w) {
		ts = append(ts, pn.TTitle(c))
	}
	if len(ts) == 0 {
		return msgs.InNoPanel
	}
	return msgs.Panels + ": " + strings.Join(ts, ", ")
}

// catalogHref is the page of the widget name of the group in the catalog.
func (b *Builder) catalogHref(group, name string) string {
	if pg := b.p.PagesRegistrator().GetHttpPage(CatalogPath(group) + "/" + name); pg != nil {
		return pg.FullPath()
	}
	return ""
}

// catalogWidget is the page of the widget w in the catalog: its group, its
// width, how it loads, the panels that take it, its documentation; a
// preview when a panel of no record takes it (else the panel of each record
// shows it, for its record).
func (b *Builder) catalogWidget(ctx *web.EventContext, w *Widget) h.HTMLComponent {
	msgs := MessagesOf(ctx.Context())
	c := ctx.Context()
	g := b.GetGroup(w.Group)
	var facts h.HTMLComponents
	fact := func(label, value string) {
		facts = append(facts, v.VListItem().Title(value).Subtitle(label))
	}
	fact(msgs.Group, g.title(c))
	fact(msgs.Width, fmt.Sprintf(msgs.WidthColumns, strconv.Itoa(int(w.Size))))
	fact(msgs.Panels, b.panelTitles(c, msgs, w))
	if len(w.Live) > 0 {
		fact(msgs.LiveEvery, w.Live[0].String())
	}
	body := h.Div(v.VList(facts...).Density(v.DensityCompact)).Class("pa-4").Attr("data-dashboard-catalog", w.ID)
	if b.docs != nil {
		if href := b.docs(ctx, cmp.Or(w.Doc, w.ID)); href != "" {
			body.AppendChildren(v.VBtn(msgs.Help).Href(href).Variant(v.VariantTonal).PrependIcon("mdi-help-circle-outline").
				Class("my-2").Attr("data-dashboard-catalog-doc", w.ID))
		}
	}
	for _, pn := range b.panelsOf(w) {
		if pn.record != nil {
			continue
		}
		ctx = pn.scoped(ctx)
		return body.AppendChildren(
			h.Div(h.Text(msgs.Preview)).Class("text-subtitle-1 font-weight-bold mt-4 mb-2"),
			pn.cards(ctx, w.ID),
		)
	}
	return body.AppendChildren(v.VAlert(h.Text(msgs.PreviewOfRecord)).Type("info").Variant(v.VariantTonal).Class("mt-4"))
}
