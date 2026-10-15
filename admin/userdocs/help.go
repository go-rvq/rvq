package userdocs

import (
	"net/url"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// HelpEvent is the event of the documentation page that answers the help of
// an address (?uri=): the document of the part of the admin it is (NodeOfURI)
// — the content of the help bar.
const HelpEvent = "userdocs_Help"

// HelpBarStorageKey is where the browser keeps the help bar: open or not,
// and its width (…​.width).
const HelpBarStorageKey = "rvq.presets.helpBar"

// helpBarRailWidth is the width of the help bar collapsed: its icon.
const helpBarRailWidth = 56

// NodeOfURI is the node of the documentation of the address uri — of a page
// of the admin, or of the event that filled a drawer or a dialog (its path
// and query, presets.HelpURI) —, among the nodes of tree (what the request
// sees): the deepest that there is of the model, the form, the action or the
// page it is; "" when none (the address of no part).
func (b *Builder) NodeOfURI(tree []*Node, uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	for _, c := range b.uriNodes(u) {
		if c != "" && findNode(tree, c) != nil {
			return c
		}
	}
	return ""
}

// uriNodes are the nodes the address u may be of, the most precise first.
func (b *Builder) uriNodes(u *url.URL) (nodes []string) {
	p := "/" + strings.Trim(u.Path, "/")
	// a document of the documentation
	if b.mb != nil {
		base := strings.TrimSuffix(b.mb.Info().ListingHref(), "/")
		if rest, ok := strings.CutPrefix(p, base+"/"); ok {
			if n, err := url.PathUnescape(rest); err == nil {
				nodes = append(nodes, n)
			}
			return
		}
	}

	segs := splitPath(p)
	mb, rest := b.modelOfPath(segs)
	if mb == nil {
		// a page of the menu
		prefix := strings.Trim(b.p.GetURIPrefix(), "/")
		page := strings.Trim(p, "/")
		if prefix != "" {
			page = strings.Trim(strings.TrimPrefix(page, prefix), "/")
		}
		if page != "" {
			nodes = append(nodes, "pages/"+page)
		}
		return
	}

	node := b.modelNode(mb)
	q := u.Query()
	switch {
	case q.Get(presets.ParamAction) != "":
		nodes = append(nodes, node+"/actions/"+q.Get(presets.ParamAction))
	case q.Get(presets.ParamBulkActionName) != "":
		nodes = append(nodes, node+"/actions/"+q.Get(presets.ParamBulkActionName))
	}
	switch q.Get(web.EventFuncIDName) {
	case actions.New, actions.NewForm, actions.Create:
		nodes = append(nodes, node+"/forms/"+FormNew)
	case actions.Edit, actions.EditForm, actions.Update:
		if q.Get(presets.ParamID) == "" && !mb.GetSingleton() {
			nodes = append(nodes, node+"/forms/"+FormNew)
		} else {
			nodes = append(nodes, node+"/forms/"+FormEdit)
		}
	case actions.Detailing, actions.DetailingContent, actions.ReloadDetail:
		nodes = append(nodes, node+"/forms/"+FormDetail)
	}

	// the path under the model's: its forms and pages
	if mb.GetSingleton() {
		switch {
		case len(rest) == 0 && mb.HasDetailing():
			nodes = append(nodes, node+"/forms/"+FormDetail)
		case len(rest) == 0, len(rest) == 1 && rest[0] == "edit":
			nodes = append(nodes, node+"/forms/"+FormEdit)
		default:
			nodes = append(nodes, node+"/pages/"+pageKey(strings.Join(rest, "/")))
		}
	} else {
		switch {
		case len(rest) == 0:
		case len(rest) == 1 && rest[0] == "new":
			nodes = append(nodes, node+"/forms/"+FormNew)
		case len(rest) == 1:
			// a record — or a page of the listing
			nodes = append(nodes, node+"/pages/"+pageKey(rest[0]), node+"/forms/"+FormDetail)
		case len(rest) == 2 && rest[1] == "edit":
			nodes = append(nodes, node+"/forms/"+FormEdit)
		default:
			// a page of a record, or of the listing
			nodes = append(nodes, node+"/pages/"+pageKey(strings.Join(rest[1:], "/")),
				node+"/pages/"+pageKey(strings.Join(rest, "/")))
		}
	}
	return append(nodes, node)
}

// modelOfPath is the model whose address segs is under — the one of the
// longest (a nested model's, under its parent's record) —, and what of segs
// is under it.
func (b *Builder) modelOfPath(segs []string) (best *presets.ModelBuilder, rest []string) {
	n := -1
	var walk func(mbs []*presets.ModelBuilder)
	seen := map[*presets.ModelBuilder]bool{}
	walk = func(mbs []*presets.ModelBuilder) {
		for _, mb := range mbs {
			if seen[mb] {
				continue
			}
			seen[mb] = true
			pat := splitPath(mb.Info().ListingHref())
			if len(pat) > n && len(pat) <= len(segs) && matchSegments(pat, segs[:len(pat)]) {
				best, n, rest = mb, len(pat), segs[len(pat):]
			}
			walk(mb.Children())
		}
	}
	walk(b.p.Models())
	return
}

// matchSegments says the segments of an address are those of a pattern —
// "{…}" any one.
func matchSegments(pattern, segs []string) bool {
	for i, p := range pattern {
		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			continue
		}
		if p != segs[i] {
			return false
		}
	}
	return true
}

// splitPath are the segments of the path p.
func splitPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

// nodePath are the nodes from the root of tree to id: its breadcrumbs.
func nodePath(tree []*Node, id string) []*Node {
	for _, n := range tree {
		if n.ID == id {
			return []*Node{n}
		}
		if p := nodePath(n.Children, id); p != nil {
			return append([]*Node{n}, p...)
		}
	}
	return nil
}

// helpEvent answers the help of ?uri=: the document of its node
// (NodeOfURI), the first when none.
func (b *Builder) helpEvent(ctx *web.EventContext) (r web.EventResponse, err error) {
	tree := b.Tree(ctx)
	node := b.NodeOfURI(tree, ctx.R.FormValue("uri"))
	if node == "" && len(tree) > 0 {
		node = tree[0].ID
	}
	r.Body = b.helpContent(ctx, tree, node)
	// the title of its tab in the help bar
	var title string
	if n := findNode(tree, node); n != nil {
		title = n.Title
	}
	r.Data = map[string]string{"node": node, "title": title}
	return
}

// helpContent is the document of node in the help bar: where it is (its
// breadcrumbs), the document, a link to it in the documentation page.
func (b *Builder) helpContent(ctx *web.EventContext, tree []*Node, node string) h.HTMLComponent {
	msgs := GetMessages(ctx.Context())
	var crumbs []string
	for _, n := range nodePath(tree, node) {
		crumbs = append(crumbs, n.Title)
	}
	href := b.DocHref(ctx, node)
	return h.Div(
		h.Tag("component").Attr(":is", "'style'").Children(h.RawHTML(docStyle)),
		h.If(len(crumbs) > 0, h.Div(h.Text(strings.Join(crumbs, " › "))).
			Class("text-caption text-medium-emphasis mb-2 rvq-help-bar__crumbs")),
		b.document(ctx, tree, node),
		v.VDivider().Class("my-3"),
		h.A(h.Text(msgs.HelpOpenDocs)).Href(href).Target("_blank").Class("text-caption"),
	).Class("user-docs rvq-help-bar__doc").
		Attr("data-help-node", node).
		Attr("data-help-href", href).
		Attr("@click", b.helpOpenScript())
}

// helpOpenScript is the click of a document in the help bar: a link to
// another document opens it there (vars.helpBar.nav: from the tab of the
// focus, in a new tab; from another, in that one); another link goes as
// links go.
func (b *Builder) helpOpenScript() string {
	base := strconv.Quote(b.mb.Info().ListingHref())
	return `(e) => {
	const a = e.target && e.target.closest && e.target.closest("a[href]");
	if (!a || a.target === "_blank" || e.ctrlKey || e.metaKey || e.shiftKey || e.button) return;
	const w = a.ownerDocument.defaultView;
	const u = new w.URL(a.getAttribute("href"), w.location.href);
	const base = ` + base + `;
	if (!u.pathname.startsWith(base + "/")) return;
	e.preventDefault();
	vars.helpBar.nav(u.pathname);
}`
}

// HelpBar is the help bar of the admin (presets.Builder.HelpBarFunc): a
// panel at its right edge, out of the page — the page narrows beside it —,
// collapsed to its icon. Open, its tabs: the first, the help of what is in
// focus — the dialog open, else the drawer or the page last clicked — by the
// address that filled it (presets.HelpURIAttr; the page's, location), again
// as it changes; the others, documents opened from it (its links, the help
// buttons of the fields: vars.helpBar.show), each with its own back and
// forward, closed by their X. A tab's content is a vx-portal whose content
// (v-model) is the body of the help event, asked by plaid. None for who may
// not read the documentation, and on its own page.
func (b *Builder) HelpBar(ctx *web.EventContext) h.HTMLComponent {
	if b.mb == nil || !b.canList(b.mb, ctx) {
		return nil
	}
	base := strings.TrimSuffix(b.mb.Info().ListingHref(), "/")
	if p := strings.TrimSuffix(ctx.R.URL.Path, "/"); p == base || strings.HasPrefix(p, base+"/") {
		return nil
	}
	msgs := GetMessages(ctx.Context())

	kinds := h.JSONString(map[string]string{
		"page": msgs.HelpFocusPage, "drawer": msgs.HelpFocusDrawer, "dialog": msgs.HelpFocusDialog,
	})

	btn := func(icon, title, click string) *v.VBtnBuilder {
		return v.VBtn("").Icon(icon).Size(v.SizeSmall).Variant(v.VariantText).
			Density(v.DensityComfortable).Attr("title", title).Attr("aria-label", title).
			Attr("@click", click)
	}
	const cur = "vars.helpBar.tabs[vars.helpBar.active]"

	header := h.Div(
		btn("mdi-arrow-left", msgs.HelpBack, "vars.helpBar.go(-1)").
			Attr(":disabled", "!("+cur+" && "+cur+".pos > 0)"),
		btn("mdi-arrow-right", msgs.HelpForward, "vars.helpBar.go(1)").
			Attr(":disabled", "!("+cur+" && "+cur+".pos < "+cur+".history.length - 1)").
			Class("mr-1"),
		v.VIcon("mdi-help-circle-outline").Color("primary").Class("mr-2"),
		h.Span(msgs.HelpBar).Class("text-subtitle-1 font-weight-bold"),
		v.VSpacer(),
		btn("mdi-open-in-new", msgs.HelpOpenDocs, "vars.helpBar.openDocs()"),
		btn("mdi-refresh", msgs.HelpReload, "vars.helpBar.reload()"),
		btn("mdi-chevron-double-right", msgs.HelpCollapse, "vars.helpBar.toggle(false)"),
	).Class("d-flex align-center px-2 py-1")

	tabs := v.VTabs(
		// the tab of the focus: always the first, never closed
		v.VTab(
			v.VIcon("mdi-crosshairs-gps").Size(v.SizeSmall).Class("mr-1"),
			h.Text(msgs.HelpFocusTab),
			v.VChip(h.Text("{{ vars.helpBar.kinds[vars.helpBar.tabs[0].kind] || '' }}")).
				Size(v.SizeXSmall).Class("ml-2").Attr("v-if", "vars.helpBar.kinds[vars.helpBar.tabs[0].kind]"),
		).Attr(":value", "0").Class("text-none"),
		h.Template(
			v.VTab(
				h.Span("{{ t.title || '…' }}").Class("text-truncate").Style("max-width: 160px"),
				v.VBtn("").Icon("mdi-close").Size(v.SizeXSmall).Variant(v.VariantText).
					Density(v.DensityCompact).Class("ml-1").
					Attr("title", msgs.HelpCloseTab).Attr("aria-label", msgs.HelpCloseTab).
					Attr("@click.stop", "vars.helpBar.close(i + 1)"),
			).Attr(":value", "i + 1").Attr(":key", "t.id").Attr(":title", "t.title").Class("text-none"),
		).Attr("v-for", "(t, i) in vars.helpBar.tabs.slice(1)"),
	).Attr("v-model", "vars.helpBar.active").Density(v.DensityCompact).ShowArrows(true).
		Class("border-b rvq-help-bar__tabs")

	open := h.Div(
		header,
		tabs,
		v.VProgressLinear().Indeterminate(true).Height(2).Color("primary").
			Attr("v-if", cur+" && "+cur+".loading"),
		v.VAlert(h.Text(msgs.HelpLoadError)).Type("warning").Variant(v.VariantTonal).
			Density(v.DensityCompact).Class("ma-3").Attr("v-if", cur+" && "+cur+".error"),
		// each tab its document: the body of the help event (vars.helpBar
		// load), its templates with the page's vars (their links: nav)
		h.Div(
			vx.VXPortal().
				Tag("div").
				Class("pa-4").
				Model("t.content").
				Data("{vars}"),
		).Attr("v-for", "(t, i) in vars.helpBar.tabs").Attr(":key", "t.id").
			Attr("v-show", "vars.helpBar.active === i").
			Class("rvq-help-bar__content").
			Attr(":class", "{'rvq-help-bar__tab--active': vars.helpBar.active === i}"),
	).Class("d-flex flex-column h-100").Attr("v-if", "vars.helpBar.open")

	rail := h.Div(
		v.VBtn("").Icon("mdi-help-circle-outline").Variant(v.VariantText).Color("primary").
			Attr("title", msgs.HelpBar+" (F1)").Attr("aria-label", msgs.HelpBar).
			Attr("@click", "vars.helpBar.toggle(true)"),
		h.Div(h.Text(msgs.HelpBar)).Class("text-caption text-medium-emphasis"),
	).Class("d-flex flex-column align-center pt-3").Attr("v-else", true)

	drawer := vx.VXNavigationDrawer(open, rail).
		Location(v.LocationRight).
		Width("400").
		Resizable(HelpBarStorageKey+".width").
		MinWidth(320).
		MaxWidth(900).
		Attr("v-model", "vars.helpBar.visible").
		Class("rvq-help-bar").
		Attr("role", "complementary").
		Attr("aria-label", msgs.HelpBar).
		// at the right edge, the first there (order): the page's drawers
		// open beside it; on a narrow screen, open, over the page
		ContainerProps(`{order: -1, floating: true, class: "border-s rvq-help-bar",
			permanent: !(vars.navMobile && vars.helpBar.open),
			temporary: !!(vars.navMobile && vars.helpBar.open),
			rail: !vars.helpBar.open, railWidth: ` + strconv.Itoa(helpBarRailWidth) + `}`)

	return h.Tag("user-component").
		Attr(":assign", `[[vars, {helpBar: vars.helpBar || {visible: true, open: false, active: 0, seq: 0,
			tabs: [{id: 0, src: "", kind: "", title: "", content: "", history: [], pos: -1, loading: false, error: false}]}}]]`).
		Attr(":setup", "[("+helpBarSetupScript(base, kinds)+")]").
		Children(
			h.Tag("component").Attr(":is", "'style'").Children(h.RawHTML(helpBarStyle)),
			drawer,
		)
}

// helpBarStyle: the dialogs open beside the help bar, not over it — their
// layer and scrim end where it begins (--rvq-help-bar-width) —; its edge to
// drag only open; its content scrolls.
const helpBarStyle = `
.v-overlay.v-dialog, .v-overlay.v-dialog > .v-overlay__scrim { right: var(--rvq-help-bar-width, 0px); }
.rvq-help-bar.v-navigation-drawer--rail .vx-navigation-drawer__resizer { display: none; }
.rvq-help-bar .rvq-help-bar__content { flex: 1 1 auto; overflow-y: auto; }
.rvq-help-bar .rvq-help-bar__tabs { flex: 0 0 auto; }
`

// helpBarSetupScript runs once a page: the help bar's state (open, kept in
// the browser) and its acts; and, while it is open, the address of what is
// in focus, every little while — a new one the document of the first tab
// anew.
func helpBarSetupScript(base, kinds string) string {
	return `({window}) => {
	const w = window, d = w.document, hb = vars.helpBar;
	const key = ` + strconv.Quote(HelpBarStorageKey) + `;
	hb.kinds = ` + kinds + `;
	if (!hb.restored) {
		hb.restored = true;
		try { hb.open = !!JSON.parse(w.localStorage.getItem(key) || "{}").open } catch (e) {}
	}
	const save = () => { try { w.localStorage.setItem(key, JSON.stringify({open: !!hb.open})) } catch (e) {} };
	const page = () => w.location.pathname + w.location.search;
	const at = (kind, uri) => JSON.stringify({kind, uri});
	// a tab's document: the body of the help event for its address (src) —
	// the answer of the last asked only —, its title the document's
	const load = (t) => {
		const n = t.asked = (t.asked || 0) + 1, f = JSON.parse(t.src);
		t.loading = true; t.error = false;
		plaid().vars(vars).url(` + strconv.Quote(base) + `).eventFunc(` + strconv.Quote(HelpEvent) + `).
			query("uri", f.uri).loadPortalBody(true).go().
			then((r) => { if (n === t.asked) { t.content = (r && r.body) || ""; t.title = (r && r.data && r.data.title) || t.title; t.loading = false } }).
			catch(() => { if (n === t.asked) { t.loading = false; t.error = true } });
	};
	const cur = () => hb.tabs[hb.active] || hb.tabs[0];
	// shown: inside no drawer or dialog closed
	const shown = (el) => !!el && el.isConnected &&
		!el.closest(".v-navigation-drawer:not(.v-navigation-drawer--active)") &&
		!el.closest(".v-overlay:not(.v-overlay--active)");
	// the outermost address of el's (its drawer's, its dialog's)
	const outer = (el) => { let m = null; for (let x = el; x; x = x.parentElement) { if (x.hasAttribute && x.hasAttribute("data-help-uri")) m = x } return m };
	const focus = () => {
		const last = shown(w.__rvqHelpBarLast) ? w.__rvqHelpBarLast : null;
		const lastAt = last && outer(last);
		const dialogs = [...d.querySelectorAll(".v-overlay.v-overlay--active.v-dialog")];
		if (dialogs.length) {
			const top = dialogs[dialogs.length - 1];
			const m = top.hasAttribute("data-help-uri") ? top : top.querySelector("[data-help-uri]");
			return {kind: "dialog", uri: m ? m.getAttribute("data-help-uri") : (lastAt ? lastAt.getAttribute("data-help-uri") : page())};
		}
		if (lastAt) return {kind: "drawer", uri: lastAt.getAttribute("data-help-uri")};
		return {kind: "page", uri: page()};
	};
	const width = () => {
		const el = d.querySelector(".v-navigation-drawer.rvq-help-bar");
		d.documentElement.style.setProperty("--rvq-help-bar-width", (el && !vars.navMobile ? el.offsetWidth : 0) + "px");
	};
	// the first tab follows the focus
	hb.tick = () => {
		width();
		if (!hb.open) return;
		const t = hb.tabs[0], f = focus(), s = at(f.kind, f.uri);
		if (s !== t.src) { t.kind = f.kind; t.src = s; load(t) }
	};
	hb.toggle = (v) => { hb.open = v === undefined ? !hb.open : v; save(); w.setTimeout(hb.tick, 0) };
	// show: the document of uri in a new tab, the one shown — the tab that
	// shows it already, when there is one
	hb.show = (uri) => {
		const src = at("doc", uri);
		const i = hb.tabs.findIndex((t, j) => j > 0 && t.src === src);
		if (i > 0) { hb.active = i; return }
		hb.tabs.push({id: ++hb.seq, src, kind: "doc", title: "", content: "", history: [src], pos: 0, loading: false, error: false});
		hb.active = hb.tabs.length - 1;
		load(hb.tabs[hb.active]);
	};
	// nav: a link of a document — from the tab of the focus, a new tab; from
	// another, in it (its back and forward)
	hb.nav = (uri) => {
		if (hb.active === 0) return hb.show(uri);
		const t = cur(), src = at("doc", uri);
		if (t.history[t.pos] === src) return;
		t.history = t.history.slice(0, t.pos + 1).concat([src]).slice(-50);
		t.pos = t.history.length - 1; t.src = src; load(t);
	};
	hb.go = (dir) => {
		const t = cur(), i = t.pos + dir;
		if (hb.active === 0 || i < 0 || i >= t.history.length) return;
		t.pos = i; t.src = t.history[i]; load(t);
	};
	hb.close = (i) => {
		if (i <= 0 || i >= hb.tabs.length) return;
		hb.tabs.splice(i, 1);
		if (hb.active >= i) hb.active = Math.max(0, hb.active - 1);
	};
	hb.reload = () => { const t = cur(); if (t.src) load(t) };
	hb.openDocs = () => {
		const el = d.querySelector(".rvq-help-bar .rvq-help-bar__tab--active [data-help-href]");
		if (el) w.open(el.getAttribute("data-help-href"), "_blank");
	};
	if (!w.__rvqHelpBarListening) {
		w.__rvqHelpBarListening = true;
		const remember = (e) => { const t = e.target; if (t && t.closest && !t.closest(".rvq-help-bar")) w.__rvqHelpBarLast = t };
		d.addEventListener("pointerdown", remember, true);
		d.addEventListener("focusin", remember, true);
		w.addEventListener("keydown", (e) => { if (e.key === "F1") { e.preventDefault(); w.__rvqHelpBarToggle && w.__rvqHelpBarToggle() } });
	}
	w.__rvqHelpBarToggle = () => hb.toggle();
	if (w.__rvqHelpBarTimer) w.clearInterval(w.__rvqHelpBarTimer);
	w.__rvqHelpBarTimer = w.setInterval(() => hb.tick(), 400);
	hb.tick();
}`
}
