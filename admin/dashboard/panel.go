package dashboard

import (
	"cmp"
	"sort"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// cards are the cards of the widgets ids of the panel — those visible for
// the request — loaded as on the panel (its period the user's); the menu of
// one adds it to the user's layout of the panel while it is not there.
func (b *Panel) cards(ctx *web.EventContext, ids ...string) h.HTMLComponent {
	msgs := MessagesOf(ctx.Context())
	prefs := b.LoadPrefs(userID(ctx.R))
	period := PeriodOf(ctx.R, b.Now(), cmp.Or(prefs.Period, PeriodDefault))
	layout := b.layout(&prefs)
	var cols h.HTMLComponents
	for _, id := range ids {
		w := b.widget(id)
		if w == nil || !w.visible(ctx.R) {
			continue
		}
		onDashboard := sectionOf(layout, w.ID) != nil
		cols = append(cols, v.VCol(b.panelCard(ctx, msgs, w, period, onDashboard)).Cols(12).Md(int(w.Size)))
	}
	return h.Div(v.VRow(cols...).Dense(true)).Class("dashboard").Attr("data-dashboard-panel", true)
}

func (b *Panel) panelCard(ctx *web.EventContext, msgs *Messages, w *Widget, period Period, onDashboard bool) h.HTMLComponent {
	title := w.Title(ctx.Context())
	items := h.HTMLComponents{}
	if !onDashboard {
		items = append(items, v.VListItem(v.VListItemTitle(h.Text(msgs.AddToDashboard))).PrependIcon("mdi-view-dashboard-edit-outline").
			Attr("@click", b.prefEvent("add", w.ID)).Attr("data-dashboard-panel-add", w.ID))
	} else {
		items = append(items, v.VListItem(v.VListItemTitle(h.Text(msgs.OnDashboard))).PrependIcon("mdi-check").Disabled(true))
	}
	var tools h.HTMLComponents
	if len(w.Live) > 0 {
		tools = append(tools, liveSelector(msgs, w))
	}
	tools = append(tools, b.pageBtn(ctx, msgs, w),
		iconBtn("mdi-refresh", msgs.Refresh, b.loadEvent(ctx.R, w, period).Query("refresh", "1").Go()).Attr("data-dashboard-refresh", w.ID),
		v.VMenu(
			web.Slot(iconBtn("mdi-dots-vertical", title, "").Attr("v-bind", "props")).Name("activator").Scope("{ props }"),
			v.VList(items...).Density(v.DensityCompact),
		),
	)
	return v.VCard(
		h.Div(
			h.Div(h.Text(title)).Class("text-subtitle-2 text-truncate").Title(title),
			v.VSpacer(),
			h.Div(tools...).Class("d-flex align-center flex-shrink-0"),
		).Class("d-flex align-center px-4 pt-3"),
		v.VCardText(web.Portal().Name(PortalName(w.ID)).Loader(b.loadEvent(ctx.R, w, period))).Class("pt-2"),
	).Variant(v.VariantOutlined).Class("h-100").Attr("data-dashboard-widget", w.ID)
}

// GroupPage draws the widgets of the group group of the panel — those
// visible for the request, in the group's order — in a page of their own
// (the page of a group of the menu, presets.PageMenuGroup: an area's), the
// period chosen at its top (in the URL), in the panel's scope; nothing of a
// record the user may not read.
func (b *Panel) GroupPage(ctx *web.EventContext, group string) h.HTMLComponent {
	ctx, ok := b.scoped(ctx)
	if !ok {
		return nil
	}
	msgs := MessagesOf(ctx.Context())
	prefs := b.LoadPrefs(userID(ctx.R))
	period := PeriodOf(ctx.R, b.Now(), cmp.Or(prefs.Period, PeriodDefault))
	var ws []*Widget
	for _, w := range b.all() {
		if w.Group == group {
			ws = append(ws, w)
		}
	}
	sort.SliceStable(ws, func(i, j int) bool { return ws[i].Order < ws[j].Order })
	ids := make([]string, len(ws))
	for i, w := range ws {
		ids[i] = w.ID
	}
	return h.Div(
		b.headerWith(ctx, msgs, &prefs, period, false, false),
		h.Div(b.cards(ctx, ids...)).Class("mt-3"),
	).Class("dashboard pa-4").Attr("data-dashboard-group-panel", group)
}
