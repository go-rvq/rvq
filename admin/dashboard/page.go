package dashboard

import (
	"cmp"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

// The events of the dashboard (global: any page may load a widget).
const (
	// EventLoad loads a widget: ?widget=…&period=…[&from=…&to=…][&refresh=1].
	EventLoad = "dashboard_widget"
	// EventAction runs a widget's action: ?widget=…&action=…&arg=….
	EventAction = "dashboard_action"
	// EventPref changes the user's choices: ?op=…&arg=….
	EventPref = "dashboard_pref"
)

// PortalName is the portal of the widget id's body.
func PortalName(id string) string { return "dashboard-w-" + id }

// Install creates the preferences' table, registers the words and the events,
// and makes the dashboard the admin's home page.
func (b *Builder) Install() error {
	if err := b.Migrate(); err != nil {
		return fmt.Errorf("dashboard: migrate: %w", err)
	}
	b.p.I18n().
		RegisterForModule(language.English, I18nDashboardKey, Messages_en_US).
		RegisterForModule(language.BrazilianPortuguese, I18nDashboardKey, Messages_pt_BR)
	wb := b.p.GetWebBuilder()
	wb.RegisterEventFunc(EventLoad, b.eventLoad)
	wb.RegisterEventFunc(EventAction, b.eventAction)
	wb.RegisterEventFunc(EventPref, b.eventPref)
	b.p.HomePageFunc(b.Page)
	return nil
}

// Docs, when given, is the link of a widget's documentation (its menu's
// "about this widget").
func (b *Builder) Docs(f func(ctx *web.EventContext, widgetID string) string) *Builder {
	b.docs = f
	return b
}

// Intro, when given, is drawn under the title (a greeting).
func (b *Builder) Intro(f func(ctx *web.EventContext) h.HTMLComponent) *Builder {
	b.intro = f
	return b
}

func userID(r *http.Request) uuid.UUID {
	u := user.GetCurrentUser(r)
	if u == nil || u.Anonymous() {
		return uuid.Nil
	}
	return u.GetID()
}

// sortedGroups are the groups by Order.
func (b *Builder) sortedGroups() []*Group {
	gs := slices.Clone(b.groups)
	slices.SortStableFunc(gs, func(a, c *Group) int { return cmp.Compare(a.Order, c.Order) })
	return gs
}

// groupWidgets are the widgets of group g visible for r, in the user's order
// (then their Order); hidden ones only with withHidden.
func (b *Builder) groupWidgets(r *http.Request, g string, prefs *Prefs, withHidden bool) []*Widget {
	var out []*Widget
	for _, w := range b.widgets {
		if w.Group != g || !w.visible(r) {
			continue
		}
		if !withHidden && prefs.hidden(w.ID) {
			continue
		}
		out = append(out, w)
	}
	pos := func(id string) int {
		if i := slices.Index(prefs.Order, id); i >= 0 {
			return i
		}
		return len(prefs.Order)
	}
	slices.SortStableFunc(out, func(a, c *Widget) int {
		if x := cmp.Compare(pos(a.ID), pos(c.ID)); x != 0 {
			return x
		}
		return cmp.Compare(a.Order, c.Order)
	})
	return out
}

// periodQuery are the query values of p (for the events of the widgets).
func periodQuery(ev *web.VueEventTagBuilder, r *http.Request, p Period) *web.VueEventTagBuilder {
	ev.Query("period", p.Key)
	if p.Key == PeriodCustom {
		ev.Query("from", r.FormValue("from")).Query("to", r.FormValue("to"))
	}
	return ev
}

// Page is the dashboard (the admin's home page).
func (b *Builder) Page(ctx *web.EventContext) (r web.PageResponse, err error) {
	msgs := MessagesOf(ctx.Context())
	prefs := b.LoadPrefs(userID(ctx.R))
	def := cmp.Or(prefs.Period, PeriodDefault)
	period := PeriodOf(ctx.R, b.Now(), def)
	customize := ctx.R.FormValue("customize") == "1"

	r.PageTitle = msgs.Dashboard
	body := h.Div().Class("dashboard pa-4").Attr("data-dashboard", true)
	body.AppendChildren(b.header(ctx, msgs, &prefs, period, customize))
	if b.intro != nil {
		body.AppendChildren(b.intro(ctx))
	}

	any := false
	var hidden []*Widget
	for _, g := range b.sortedGroups() {
		ws := b.groupWidgets(ctx.R, g.ID, &prefs, customize)
		if len(ws) == 0 {
			continue
		}
		any = true
		body.AppendChildren(b.group(ctx, msgs, g, ws, &prefs, period, customize))
		for _, w := range ws {
			if prefs.hidden(w.ID) {
				hidden = append(hidden, w)
			}
		}
	}
	if !any {
		body.AppendChildren(v.VAlert(h.Text(msgs.NoWidgets)).Type("info").Variant(v.VariantTonal).Class("mt-4"))
	}
	if customize && len(hidden) > 0 {
		var chips h.HTMLComponents
		for _, w := range hidden {
			chips = append(chips, v.VChip(h.Text(w.Title(ctx.Context()))).PrependIcon("mdi-eye-outline").Class("me-2 mb-2").
				Attr("@click", prefEvent("show", w.ID)).Attr("data-dashboard-show", w.ID))
		}
		body.AppendChildren(h.Div(
			h.Div(h.Text(msgs.HiddenWidgets)).Class("text-subtitle-2 mb-2"),
			h.Div(chips...),
		).Class("mt-6"))
	}
	r.Body = body
	return
}

func prefEvent(op, arg string) string {
	return web.Plaid().EventFunc(EventPref).Query("op", op).Query("arg", arg).Go()
}

// pushPeriod is the event that opens the dashboard with key (custom: the
// locals' from and to), keeping the customize mode.
func pushPeriod(key string, customize bool) *web.VueEventTagBuilder {
	ev := web.Plaid().PushState(true).Query("period", key)
	if key == PeriodCustom {
		ev.Query("from", web.Var("locals.from")).Query("to", web.Var("locals.to"))
	}
	if customize {
		ev.Query("customize", "1")
	}
	return ev
}

func (b *Builder) header(ctx *web.EventContext, msgs *Messages, prefs *Prefs, period Period, customize bool) h.HTMLComponent {
	// the period: a select of the periods — the custom one shows its dates
	var periods []map[string]string
	for _, k := range Periods {
		periods = append(periods, map[string]string{"title": msgs.PeriodLabel(k), "value": k})
	}
	choose := web.Plaid().PushState(true).Query("period", web.Var("$event"))
	if customize {
		choose.Query("customize", "1")
	}
	selector := v.VSelect().Items(periods).ItemTitle("title").ItemValue("value").
		Attr(":model-value", fmt.Sprintf("locals.custom ? %q : %q", PeriodCustom, period.Key)).
		Label(msgs.Period).Density(v.DensityCompact).HideDetails(true).Variant(v.VariantOutlined).
		PrependInnerIcon("mdi-calendar-range").
		Style("min-width: 200px; max-width: 240px").
		Attr("data-dashboard-period", period.Key).
		Attr("aria-label", msgs.Period).
		Attr("@update:model-value", fmt.Sprintf(`(k) => { if (k === %q) { locals.custom = true } else { locals.custom = false; (($event) => %s)(k) } }`,
			PeriodCustom, choose.Go()))
	from, to := ctx.R.FormValue("from"), ctx.R.FormValue("to")
	if period.Key != PeriodCustom {
		from, to = period.From.Format("2006-01-02"), period.To.Add(-time.Second).Format("2006-01-02")
	}
	custom := h.Div(
		v.VTextField().Type("date").Label(msgs.From).Density(v.DensityCompact).HideDetails(true).Attr("v-model", "locals.from").Style("max-width:180px"),
		v.VTextField().Type("date").Label(msgs.To).Density(v.DensityCompact).HideDetails(true).Attr("v-model", "locals.to").Style("max-width:180px"),
		v.VBtn(msgs.Apply).Color("primary").Variant(v.VariantFlat).Attr("@click", pushPeriod(PeriodCustom, customize).Go()),
	).Class("d-flex flex-wrap align-center ga-2 mt-2").Attr("v-if", "locals.custom")

	toggle := "1"
	label, icon := msgs.Customize, "mdi-view-dashboard-edit-outline"
	if customize {
		toggle, label, icon = "", msgs.Done, "mdi-check"
	}
	cust := web.Plaid().PushState(true).Query("period", period.Key)
	if period.Key == PeriodCustom {
		cust.Query("from", from).Query("to", to)
	}
	if toggle != "" {
		cust.Query("customize", toggle)
	}
	tools := h.Div(
		v.VBtn(label).PrependIcon(icon).Variant(v.VariantTonal).Size(v.SizeSmall).Attr("@click", cust.Go()).Attr("data-dashboard-customize", true),
	).Class("d-flex align-center ga-2")
	if customize {
		var items []map[string]string
		for _, k := range Periods {
			if k != PeriodCustom {
				items = append(items, map[string]string{"title": msgs.PeriodLabel(k), "value": k})
			}
		}
		tools.AppendChildren(
			v.VSelect().Items(items).ItemTitle("title").ItemValue("value").ModelValue(cmp.Or(prefs.Period, PeriodDefault)).
				Label(msgs.DefaultPeriod).Density(v.DensityCompact).HideDetails(true).Style("max-width:200px").
				Attr("@update:model-value", web.Plaid().EventFunc(EventPref).Query("op", "period").Query("arg", web.Var("$event")).Go()),
			v.VBtn(msgs.Reset).PrependIcon("mdi-restore").Variant(v.VariantText).Size(v.SizeSmall).Attr("@click", prefEvent("reset", "")),
		)
	}
	return web.Scope(
		h.Div(
			// the title is the page's (its bar, PageTitle): not again here
			selector,
			v.VSpacer(),
			tools,
		).Class("d-flex flex-wrap align-center ga-2"),
		custom,
	).Slot("{ locals }").LocalsInit(fmt.Sprintf("{custom: %t, from: %q, to: %q}", period.Key == PeriodCustom, from, to))
}

func (b *Builder) group(ctx *web.EventContext, msgs *Messages, g *Group, ws []*Widget, prefs *Prefs, period Period, customize bool) h.HTMLComponent {
	collapsed := prefs.collapsed(g.ID)
	box := h.Div().Class("mt-6").Attr("data-dashboard-group", g.ID)
	if g.Title != nil {
		op, ic, tip := "collapse", "mdi-chevron-up", msgs.Collapse
		if collapsed {
			op, ic, tip = "expand", "mdi-chevron-down", msgs.Expand
		}
		box.AppendChildren(h.Div(
			h.If(g.Icon != "", v.VIcon(g.Icon).Size(v.SizeSmall).Class("me-2 text-medium-emphasis")),
			h.Div(h.Text(g.Title(ctx.Context()))).Class("text-subtitle-1 font-weight-bold"),
			v.VSpacer(),
			v.VBtn("").Icon(ic).Variant(v.VariantText).Size(v.SizeSmall).Attr("title", tip, "aria-label", tip).
				Attr("@click", prefEvent(op, g.ID)).Attr("data-dashboard-collapse", g.ID),
		).Class("d-flex align-center mb-2"))
	} else {
		collapsed = false
	}
	if collapsed {
		return box
	}
	ids := make([]string, len(ws))
	for i, w := range ws {
		ids[i] = w.ID
	}
	var cols h.HTMLComponents
	for i, w := range ws {
		cols = append(cols, v.VCol(b.card(ctx, msgs, w, prefs, period, customize, i == 0, i == len(ws)-1)).
			Cols(12).Sm(min(12, int(w.Size)*2)).Md(int(w.Size)))
	}
	return box.AppendChildren(v.VRow(cols...).Dense(true))
}

func iconBtn(icon, tip, click string) *v.VBtnBuilder {
	return v.VBtn("").Icon(icon).Variant(v.VariantText).Size(v.SizeXSmall).Density(v.DensityComfortable).
		Attr("title", tip, "aria-label", tip, "@click.stop", click)
}

func (b *Builder) loadEvent(r *http.Request, w *Widget, p Period) *web.VueEventTagBuilder {
	return periodQuery(web.Plaid().EventFunc(EventLoad).Query("widget", w.ID), r, p)
}

func (b *Builder) card(ctx *web.EventContext, msgs *Messages, w *Widget, prefs *Prefs, period Period, customize, first, last bool) h.HTMLComponent {
	title := w.Title(ctx.Context())
	hidden := prefs.hidden(w.ID)
	var tools h.HTMLComponents
	if customize {
		if !first {
			tools = append(tools, iconBtn("mdi-arrow-up", msgs.MoveUp, prefEvent("up", w.ID)).Attr("data-dashboard-up", w.ID))
		}
		if !last {
			tools = append(tools, iconBtn("mdi-arrow-down", msgs.MoveDown, prefEvent("down", w.ID)).Attr("data-dashboard-down", w.ID))
		}
		if hidden {
			tools = append(tools, iconBtn("mdi-eye-outline", msgs.Show, prefEvent("show", w.ID)).Attr("data-dashboard-show", w.ID))
		} else {
			tools = append(tools, iconBtn("mdi-eye-off-outline", msgs.Hide, prefEvent("hide", w.ID)).Attr("data-dashboard-hide", w.ID))
		}
	} else {
		tools = append(tools, iconBtn("mdi-refresh", msgs.Refresh, b.loadEvent(ctx.R, w, period).Query("refresh", "1").Go()).Attr("data-dashboard-refresh", w.ID))
		var items h.HTMLComponents
		if w.Link != nil {
			if href := w.Link(ctx); href != "" {
				items = append(items, v.VListItem(v.VListItemTitle(h.Text(msgs.Open))).PrependIcon("mdi-open-in-new").Href(href))
			}
		}
		if b.docs != nil {
			if href := b.docs(ctx, w.ID); href != "" {
				items = append(items, v.VListItem(v.VListItemTitle(h.Text(msgs.Help))).PrependIcon("mdi-help-circle-outline").Href(href).Attr("target", "_blank"))
			}
		}
		items = append(items, v.VListItem(v.VListItemTitle(h.Text(msgs.Hide))).PrependIcon("mdi-eye-off-outline").Attr("@click", prefEvent("hide", w.ID)))
		tools = append(tools, v.VMenu(
			web.Slot(iconBtn("mdi-dots-vertical", title, "").Attr("v-bind", "props")).Name("activator").Scope("{ props }"),
			v.VList(items...).Density(v.DensityCompact),
		))
	}

	var body h.HTMLComponent
	if customize {
		body = h.Div(h.Text(msgs.PeriodLabel(period.Key))).Class("text-caption text-medium-emphasis py-4 text-center")
	} else {
		// a portal with children draws them, never its Loader: it starts
		// empty, the widget's body when its event answers
		body = web.Portal().Name(PortalName(w.ID)).Loader(b.loadEvent(ctx.R, w, period))
	}
	card := v.VCard(
		h.Div(
			h.Div(h.Text(title)).Class("text-subtitle-2 text-truncate").Title(title),
			h.If(w.Periodic, v.VIcon("mdi-calendar-range").Size(v.SizeXSmall).Class("ms-1 text-medium-emphasis").Attr("title", msgs.FollowsPeriod)),
			v.VSpacer(),
			h.Div(tools...).Class("d-flex align-center flex-shrink-0"),
		).Class("d-flex align-center px-4 pt-3"),
		v.VCardText(body).Class("pt-2"),
	).Variant(v.VariantOutlined).Class("h-100").Attr("data-dashboard-widget", w.ID)
	if hidden {
		card.Style("opacity:.45")
	}
	return card
}

// cacheKey is a user's widget over a period.
func cacheKey(u uuid.UUID, w string, p Period) string {
	return strings.Join([]string{w, u.String(), p.Key, p.From.Format(time.RFC3339), p.To.Format(time.RFC3339)}, "|")
}

type cached struct {
	blocks []Block
	at     time.Time
}

// blocks are w's blocks over p for the request's user: the cache's while its
// TTL lasts, unless refresh.
func (b *Builder) blocks(ctx *web.EventContext, w *Widget, p Period, refresh bool) ([]Block, time.Time, error) {
	key := cacheKey(userID(ctx.R), w.ID, p)
	ttl := cmp.Or(w.TTL, time.Minute)
	now := b.Now()
	if !refresh {
		if c, ok := b.cache.Load(key); ok && now.Sub(c.(*cached).at) < ttl {
			return c.(*cached).blocks, c.(*cached).at, nil
		}
	}
	bs, err := w.Load(ctx, p)
	if err != nil {
		return nil, now, err
	}
	b.cache.Store(key, &cached{blocks: bs, at: now})
	return bs, now, nil
}

// Invalidate drops the cache of the widget id (every user, every period).
func (b *Builder) Invalidate(id string) {
	b.cache.Range(func(k, _ any) bool {
		if strings.HasPrefix(k.(string), id+"|") {
			b.cache.Delete(k)
		}
		return true
	})
}

// Render is the body of w over p for ctx (its blocks, or its error).
func (b *Builder) Render(ctx *web.EventContext, w *Widget, p Period, refresh bool) h.HTMLComponent {
	msgs := MessagesOf(ctx.Context())
	bs, at, err := b.blocks(ctx, w, p, refresh)
	if err != nil {
		log.Printf("dashboard: widget %s: %v", w.ID, err)
		return v.VAlert(
			h.Div(h.Text(msgs.LoadError)),
			v.VBtn(msgs.Retry).Size(v.SizeSmall).Variant(v.VariantText).Class("mt-1").
				Attr("@click", b.loadEvent(ctx.R, w, p).Query("refresh", "1").Go()),
		).Type("error").Variant(v.VariantTonal).Density(v.DensityCompact).Attr("data-dashboard-error", w.ID)
	}
	rc := &renderCtx{b: b, w: w, ctx: ctx, msgs: msgs, period: p}
	out := h.Div().Attr("data-dashboard-body", w.ID)
	for _, bl := range bs {
		if bl != nil {
			out.AppendChildren(bl.render(rc))
		}
	}
	out.AppendChildren(h.Div(h.Text(strings.ReplaceAll(msgs.UpdatedAt, "{time}", at.Format("15:04")))).
		Class("text-caption text-disabled mt-2 text-end"))
	return out
}

// widgetOf is the request's widget, if visible.
func (b *Builder) widgetOf(r *http.Request) (*Widget, error) {
	w := b.widget(r.FormValue("widget"))
	if w == nil || !w.visible(r) || w.Load == nil {
		return nil, presets.ErrRecordNotFound
	}
	return w, nil
}

func (b *Builder) eventPeriod(ctx *web.EventContext) Period {
	prefs := b.LoadPrefs(userID(ctx.R))
	return PeriodOf(ctx.R, b.Now(), cmp.Or(prefs.Period, PeriodDefault))
}

func (b *Builder) eventLoad(ctx *web.EventContext) (r web.EventResponse, err error) {
	w, err := b.widgetOf(ctx.R)
	if err != nil {
		// not a widget, or not the user's: as if there were none
		http.NotFound(ctx.W, ctx.R)
		return r, nil
	}
	refresh := ctx.R.FormValue("refresh") == "1"
	body := b.Render(ctx, w, b.eventPeriod(ctx), refresh)
	if refresh {
		r.UpdatePortal(PortalName(w.ID), body)
	} else {
		r.Body = body
	}
	return
}

func (b *Builder) eventAction(ctx *web.EventContext) (r web.EventResponse, err error) {
	w, err := b.widgetOf(ctx.R)
	if err != nil {
		http.NotFound(ctx.W, ctx.R)
		return r, nil
	}
	msgs := MessagesOf(ctx.Context())
	f := w.Actions[ctx.R.FormValue("action")]
	if f == nil {
		return r, presets.ErrRecordNotFound
	}
	if aerr := f(ctx, ctx.R.FormValue("arg")); aerr != nil {
		presets.ShowMessage(&r, strings.ReplaceAll(msgs.ActionFailed, "{error}", aerr.Error()), "error")
	} else {
		presets.ShowMessage(&r, msgs.ActionDone)
	}
	b.Invalidate(w.ID)
	r.UpdatePortal(PortalName(w.ID), b.Render(ctx, w, b.eventPeriod(ctx), true))
	return
}

func (b *Builder) eventPref(ctx *web.EventContext) (r web.EventResponse, err error) {
	uid := userID(ctx.R)
	if uid == uuid.Nil {
		return r, presets.ErrRecordNotFound
	}
	prefs := b.LoadPrefs(uid)
	op, arg := ctx.R.FormValue("op"), ctx.R.FormValue("arg")
	switch op {
	case "hide", "show":
		if b.widget(arg) == nil {
			return r, presets.ErrRecordNotFound
		}
		prefs.Hidden = toggle(prefs.Hidden, arg, op == "hide")
	case "collapse", "expand":
		prefs.Collapsed = toggle(prefs.Collapsed, arg, op == "collapse")
	case "up", "down":
		w := b.widget(arg)
		if w == nil {
			return r, presets.ErrRecordNotFound
		}
		ws := b.groupWidgets(ctx.R, w.Group, &prefs, true)
		ids := make([]string, len(ws))
		for i, x := range ws {
			ids[i] = x.ID
		}
		d := 1
		if op == "up" {
			d = -1
		}
		prefs.move(ids, arg, d)
	case "period":
		if _, ok := periodSpans[arg]; !ok {
			return r, presets.ErrRecordNotFound
		}
		prefs.Period = arg
	case "reset":
		if err = b.ResetPrefs(uid); err != nil {
			return
		}
		r.RunScript = "window.location.reload()"
		return
	default:
		return r, presets.ErrRecordNotFound
	}
	if err = b.SavePrefs(uid, prefs); err != nil {
		return
	}
	// the page again as it is (its period, its mode: the URL's)
	r.RunScript = "window.location.reload()"
	return
}
