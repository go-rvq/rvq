package dashboard

import (
	"cmp"
	"context"
	"net/http"
	"net/url"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// The page of a widget: the widget alone, the whole width — its full
// version (IsFull: a widget may show more there, more points, more rows) —,
// with the period, its reloading (Live) and the way back. Every widget has
// one; its header links it.

// WidgetPagePath is the path of the page of a widget (under the admin's
// prefix), the widget its query's id.
const WidgetPagePath = "/dashboard/widget"

// paramFull, in the events of a widget's page: its full version.
const paramFull = "full"

// IsFull says whether the widget is loaded for its page (its full version):
// for its Load, which may show more there.
func IsFull(ctx *web.EventContext) bool { return ctx.R.FormValue(paramFull) == "1" }

// WidgetPageURL is the address of the page of the widget id (prefix: the
// admin's), keeping the period of r.
func WidgetPageURL(prefix, id string, r *http.Request) string {
	q := url.Values{"id": {id}}
	for _, k := range []string{"period", "from", "to"} {
		if v := r.FormValue(k); v != "" {
			q.Set(k, v)
		}
	}
	return prefix + WidgetPagePath + "?" + q.Encode()
}

// installWidgetPage registers the page of the widgets.
func (b *Builder) installWidgetPage() {
	b.p.PagesRegistrator().New(
		presets.HttpPage(WidgetPagePath).InMenu(false).
			TitleFunc(func(ctx context.Context) string { return MessagesOf(ctx).Dashboard })).
		Layout(b.widgetPage).
		Build()
}

// widgetPage is the page of the widget of the query's id: not one, or not
// the user's, not found.
func (b *Builder) widgetPage(ctx *web.EventContext) (r web.PageResponse, err error) {
	msgs := MessagesOf(ctx.Context())
	w := b.widget(ctx.R.FormValue("id"))
	if w == nil || !w.visible(ctx.R) {
		return b.p.DefaultNotFoundPageFunc(ctx)
	}
	prefs := b.LoadPrefs(userID(ctx.R))
	period := PeriodOf(ctx.R, b.Now(), cmp.Or(prefs.Period, PeriodDefault))
	title := w.Title(ctx.Context())
	r.PageTitle = title

	// the period, as the dashboard's (a select), the page's id kept
	var periods []map[string]string
	for _, k := range Periods {
		if k != PeriodCustom {
			periods = append(periods, map[string]string{"title": msgs.PeriodLabel(k), "value": k})
		}
	}
	choose := web.Plaid().PushState(true).Query("id", w.ID).Query("period", web.Var("$event"))
	var tools h.HTMLComponents
	if w.Periodic {
		tools = append(tools, v.VSelect().Items(periods).ItemTitle("title").ItemValue("value").
			ModelValue(period.Key).Label(msgs.Period).Density(v.DensityCompact).HideDetails(true).
			Variant(v.VariantOutlined).PrependInnerIcon("mdi-calendar-range").
			Style("min-width: 200px; max-width: 240px").Attr("data-dashboard-period", period.Key).
			Attr("@update:model-value", choose.Go()))
	}
	if len(w.Live) > 0 {
		tools = append(tools, liveSelector(msgs, w))
	}
	tools = append(tools, iconBtn("mdi-refresh", msgs.Refresh, b.fullLoadEvent(ctx.R, w, period).Query("refresh", "1").Go()).
		Attr("data-dashboard-refresh", w.ID))

	r.Body = h.Div(
		h.Div(
			v.VBtn(msgs.BackToDashboard).PrependIcon("mdi-arrow-left").Variant(v.VariantText).Size(v.SizeSmall).
				Href(b.uriPrefix()+"/").Attr("data-dashboard-back", true),
			h.Div(h.Text(title)).Class("text-h6 text-truncate").Title(title),
			v.VSpacer(),
			h.Div(tools...).Class("d-flex align-center ga-2 flex-shrink-0"),
		).Class("d-flex flex-wrap align-center ga-2 mb-3"),
		v.VCard(
			v.VCardText(web.Portal().Name(PortalName(w.ID)).Loader(b.fullLoadEvent(ctx.R, w, period))),
		).Variant(v.VariantOutlined).Attr("data-dashboard-widget", w.ID).Attr("data-dashboard-full", true),
	).Class("dashboard pa-4").Attr("data-dashboard-widget-page", w.ID)
	return
}

// fullLoadEvent is the load of w on its page: its full version.
func (b *Builder) fullLoadEvent(r *http.Request, w *Widget, p Period) *web.VueEventTagBuilder {
	return b.loadEvent(r, w, p).Query(paramFull, "1")
}

// pageBtn is the button of a widget's header that opens its page.
func (b *Builder) pageBtn(ctx *web.EventContext, msgs *Messages, w *Widget) h.HTMLComponent {
	return v.VBtn("").Icon("mdi-arrow-top-right").Variant(v.VariantText).Size(v.SizeXSmall).Density(v.DensityComfortable).
		Href(WidgetPageURL(b.uriPrefix(), w.ID, ctx.R)).
		Attr("title", msgs.OpenPage, "aria-label", msgs.OpenPage, "data-dashboard-page", w.ID)
}

// uriPrefix is the admin's prefix (none without the presets: a test's).
func (b *Builder) uriPrefix() string {
	if b.p == nil {
		return ""
	}
	return b.p.GetURIPrefix()
}
