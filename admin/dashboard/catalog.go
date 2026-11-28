package dashboard

import (
	"cmp"
	"context"
	"fmt"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/packages/helper"
	"net/url"
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
	// the selects of the records of the widgets' pages (catalogScope)
	helper.EnsureMessages(b.p.I18n())
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
		icon := cmp.Or(g.g.Icon, "mdi-view-grid-outline")
		page := b.p.PagesRegistrator().New(presets.HttpPage(CatalogPath(g.ID())).
			TitleFunc(g.title).DescriptionFunc(g.description).MenuIcon(icon)).
			Layout(func(ctx *web.EventContext) (r web.PageResponse, err error) {
				r.PageTitle = g.title(ctx.Context())
				r.Body = b.catalogGroup(ctx, g)
				return
			})
		for _, w := range g.widgets {
			w := w.w
			sub := page.SubPage("/" + w.Name).InMenu(true).Title(w.title)
			sub.Page().MenuIcon(icon)
			sub.Layout(func(ctx *web.EventContext) (r web.PageResponse, err error) {
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
		t := pn.TTitle(c)
		if pn.record != nil {
			// of whom: "Painel (Organização)"
			t += " (" + pn.record.TTitle(c) + ")"
		}
		ts = append(ts, t)
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

// catalogWidget is the page of the widget w in the catalog: at its top, the
// choice of its scope (catalogScope: the record it is loaded for, when it
// needs one) and its preview there; its group, its width, how it loads, the
// panels that take it, its documentation.
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
	body := h.Div().Class("pa-4").Attr("data-dashboard-catalog", w.ID)
	body.AppendChildren(b.catalogScope(ctx, w)...)
	body.AppendChildren(v.VList(facts...).Density(v.DensityCompact).Class("mt-4"))
	if b.docs != nil {
		if href := b.docs(ctx, cmp.Or(w.Doc, w.ID)); href != "" {
			body.AppendChildren(v.VBtn(msgs.Help).Href(href).Variant(v.VariantTonal).PrependIcon("mdi-help-circle-outline").
				Class("my-2").Attr("data-dashboard-catalog-doc", w.ID))
		}
	}
	return body
}

// The query of the page of a widget in the catalog: the scope chosen (""
// global, else the URI name of the model of the record) and the records
// chosen, from the root of its parents down ("<org>/<project>").
const (
	paramScope = "scope"
	paramRec   = "rec"
)

// catalogScope is the top of the page of the widget w in the catalog: of a
// widget of records, the kind of scope (when more than one: Consolidated,
// Organization, Project…), the selects of its record, cascaded by its
// parents (Organization › Project: each in a dialog, the records the user
// lists under the ones chosen above — none they do not reach, and no select
// of a model they may not list —; the first the one chosen in the menu at
// first), then its preview, when each one chosen is read and the widget is
// the user's (else not drawn) — loaded by its panel at the page of the record (Panel.at:
// its record the path's there, read again); a global one, its preview.
func (b *Builder) catalogScope(ctx *web.EventContext, w *Widget) (out h.HTMLComponents) {
	msgs := MessagesOf(ctx.Context())
	sc := b.scopeOf(w)
	type kind struct {
		key, title string
		mb         *presets.ModelBuilder
	}
	var kinds []kind
	if sc.Global {
		kinds = append(kinds, kind{title: msgs.Consolidated})
	}
	for _, mb := range sc.Records {
		// of a model whose records the user may list (from the root of its
		// parents): else not offered
		if root := mb.Root(); root.Permissioner().Lister(ctx.R).Denied() {
			continue
		}
		kinds = append(kinds, kind{key: mb.Info().URIName(), title: mb.TTitle(ctx.Context()), mb: mb})
	}
	noAccess := func() h.HTMLComponent {
		return v.VAlert(h.Text(msgs.NoAccess)).Type("warning").Variant(v.VariantTonal).Attr("data-dashboard-catalog-no-access", true)
	}
	if len(kinds) == 0 || !w.visible(ctx.R) {
		// nothing the user may choose, or the widget not theirs: not drawn
		return append(out, noAccess())
	}
	cur := kinds[0]
	for _, k := range kinds {
		if k.key == ctx.R.FormValue(paramScope) {
			cur = k
		}
	}
	base := b.catalogHref(w.Group, w.Name)
	href := func(key, rec string) string {
		q := url.Values{}
		if key != "" {
			q.Set(paramScope, key)
		}
		if rec != "" {
			q.Set(paramRec, rec)
		}
		if len(q) == 0 {
			return base
		}
		return base + "?" + q.Encode()
	}
	preview := func(c h.HTMLComponent) h.HTMLComponents {
		return h.HTMLComponents{h.Div(h.Text(msgs.Preview)).Class("text-subtitle-1 font-weight-bold mt-2 mb-2"), c}
	}
	if len(kinds) > 1 {
		var btns h.HTMLComponents
		for _, k := range kinds {
			btn := v.VBtn(k.title).Href(href(k.key, "")).Size(v.SizeSmall).Class("me-2").
				Attr("data-dashboard-catalog-scope", cmp.Or(k.key, "global"))
			if k.key == cur.key {
				btn.Variant(v.VariantFlat).Color(v.ColorPrimary)
			} else {
				btn.Variant(v.VariantTonal)
			}
			btns = append(btns, btn)
		}
		out = append(out, h.Div(btns...).Class("mb-3"))
	}
	if cur.mb == nil {
		for _, pn := range b.panelsOf(w) {
			if pn.record == nil {
				ctx, _ = pn.scoped(ctx) // of no record: always
				return append(out, preview(pn.cards(ctx, w.ID))...)
			}
		}
		return append(out, v.VAlert(h.Text(msgs.InNoPanel)).Type("info").Variant(v.VariantTonal))
	}
	// the selects, from the root of the parents of the model down: each the
	// one chosen (when the user reads it, under the ones above), the next
	// only once one is
	chain := append(cur.mb.Parents(), cur.mb)
	given := strings.Split(ctx.R.FormValue(paramRec), "/")
	_, asked := ctx.R.Form[paramRec]
	var (
		ids     []model.ID
		titles  []string
		selects h.HTMLComponents
	)
	for i, mb := range chain {
		titles = append(titles, mb.TTitle(ctx.Context()))
		val := ""
		if i < len(given) {
			val = given[i]
		}
		if val == "" && !asked {
			// at first, the one chosen in the menu
			if id, ok := mb.MenuRecordChosen(ctx, ids); ok {
				val = id.String()
			}
		}
		prefix := ""
		for _, id := range ids {
			prefix += id.String() + "/"
		}
		rs := helper.RecordSelect(mb).Label(mb.TTitle(ctx.Context())).Parents(ids...).Value(val).
			OnChange(fmt.Sprintf("plaid().vars(vars).pushStateURL(%s + ($event ? encodeURIComponent($event) : '')).go()",
				h.JSONString(href(cur.key, "")+map[bool]string{true: "?", false: "&"}[cur.key == ""]+paramRec+"="+url.QueryEscape(prefix))))
		comp := rs.Component(ctx)
		if comp == nil {
			// records the user may not list: nothing to choose, no widget
			out = append(out, v.VRow(selects...).Dense(true))
			return append(out, noAccess())
		}
		chosen, _ := rs.Chosen(ctx)
		selects = append(selects, v.VCol(comp.Attr("data-dashboard-catalog-record", mb.Info().URIName())).Cols(12).Md(6))
		if chosen == "" {
			out = append(out, v.VRow(selects...).Dense(true))
			return append(out, v.VAlert(h.Text(fmt.Sprintf(msgs.ChooseRecord, strings.Join(titles[i:], " › ")))).
				Type("info").Variant(v.VariantTonal).Attr("data-dashboard-catalog-choose", true))
		}
		id, _ := mb.ParseRecordID(chosen)
		ids = append(ids, id)
	}
	out = append(out, v.VRow(selects...).Dense(true))
	for _, pn := range b.panelsOf(w) {
		if pn.record == cur.mb {
			last := ids[len(ids)-1]
			pn = pn.atRecord(last.String(), ids[:len(ids)-1]...)
			ctx.R = WithScope(ctx.R, Scope{Panel: pn.id, Record: last.String()})
			return append(out, preview(pn.cards(ctx, w.ID))...)
		}
	}
	return append(out, v.VAlert(h.Text(msgs.InNoPanel)).Type("info").Variant(v.VariantTonal))
}
