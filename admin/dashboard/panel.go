package dashboard

import (
	"cmp"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// Panel draws the widgets ids — those visible for the request, known (Widget
// or Dynamic) — on a page of their own: a record's dashboard. Each card loads
// as on the dashboard (its period the user's), and its menu adds it to the
// user's dashboard while it is not there.
func (b *Builder) Panel(ctx *web.EventContext, ids ...string) h.HTMLComponent {
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

func (b *Builder) panelCard(ctx *web.EventContext, msgs *Messages, w *Widget, period Period, onDashboard bool) h.HTMLComponent {
	title := w.Title(ctx.Context())
	items := h.HTMLComponents{}
	if !onDashboard {
		items = append(items, v.VListItem(v.VListItemTitle(h.Text(msgs.AddToDashboard))).PrependIcon("mdi-view-dashboard-edit-outline").
			Attr("@click", prefEvent("add", w.ID)).Attr("data-dashboard-panel-add", w.ID))
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
