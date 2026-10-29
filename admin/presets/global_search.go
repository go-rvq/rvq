package presets

import (
	"reflect"
	"sort"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/account"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// The search of the top bar (Ctrl K): a palette over the page — the items of
// the menu whose names match (in the browser), the records of the models the
// user may list whose words match (GlobalSearchEvent: each model's own search
// of its listing, its keyword columns, a few of each). ↑ ↓ choose, Enter
// opens, Esc closes. The links are of the account of the page (/admin/u/N).

// GlobalSearchEvent answers the records matching q, by model.
const GlobalSearchEvent = "presets_GlobalSearch"

// globalSearchPerModel is how many records of a model the search shows.
const globalSearchPerModel = 5

type globalSearchItem struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Icon  string `json:"icon,omitempty"`
}

type globalSearchGroup struct {
	Title string             `json:"title"`
	Icon  string             `json:"icon,omitempty"`
	Items []globalSearchItem `json:"items"`
}

// GlobalSearch says whether the search of the top bar looks in the records of
// the model (by default it does, when the model is in the menu).
func (mb *ModelBuilder) GlobalSearch(v bool) *ModelBuilder {
	mb.notInGlobalSearch = !v
	return mb
}

// GlobalSearchParents makes the search of the top bar look in the records of
// a child model too (it does not, by default): parents are the ids of the
// parents of a record, root first — as ListingHref takes them —; its link
// and whether the user may list it are under them. The model's listing must
// find its records with no parent in the request (all of them).
func (mb *ModelBuilder) GlobalSearchParents(parents func(obj any) IDSlice) *ModelBuilder {
	mb.globalSearchParents = parents
	return mb
}

// globalSearch is GlobalSearchEvent.
func (b *Builder) globalSearch(ctx *web.EventContext) (r web.EventResponse, err error) {
	q := strings.TrimSpace(ctx.R.FormValue("q"))
	groups := []*globalSearchGroup{}
	r.Data = groups
	if len([]rune(q)) < 2 {
		return
	}
	n := account.IndexOf(ctx.R)
	// in the order of the menu
	nodes, _, _ := b.menuNodes(ctx)
	order := map[string]int{}
	for i, e := range menuSearchEntries(nodes, "") {
		if _, ok := order[e["href"]]; !ok {
			order[e["href"]] = i
		}
	}
	pos := map[*globalSearchGroup]int{}
	// the models and their children, a child after its parent (in its place
	// in the menu)
	var look func(mb *ModelBuilder, p int)
	look = func(mb *ModelBuilder, p int) {
		if g := b.globalSearchModel(ctx, mb, q, n); g != nil {
			pos[g] = p
			groups = append(groups, g)
		}
		for _, c := range mb.children {
			look(c, p)
		}
	}
	for _, mb := range b.models {
		p, ok := order[mb.Info().ListingHref()]
		if !ok {
			p = len(order)
		}
		look(mb, p)
	}
	sort.SliceStable(groups, func(i, j int) bool { return pos[groups[i]] < pos[groups[j]] })
	r.Data = groups
	return
}

// globalSearchModel is the group of the records of mb matching q (nil: none,
// or a model the search does not look in — out of the menu, a child of
// another that does not say its records' parents (GlobalSearchParents), a
// singleton, one the user may not list, one with no words).
func (b *Builder) globalSearchModel(ctx *web.EventContext, mb *ModelBuilder, q string, n int) *globalSearchGroup {
	l := mb.listing
	child := mb.parent != nil
	if child {
		// a child: when it says its records' parents
		if mb.globalSearchParents == nil {
			return nil
		}
	} else if mb.notInMenu || mb.permissioner.ReqLister(ctx.R).Denied() {
		return nil
	}
	if mb.notInGlobalSearch || mb.singleton || l == nil || l.Searcher == nil ||
		len(l.searchColumns) == 0 || mb.CurrentDataOperator() == nil {
		return nil
	}
	params := &SearchParams{
		KeywordColumns: l.searchColumns,
		Keyword:        q,
		PerPage:        globalSearchPerModel,
		Page:           1,
		SQLConditions:  l.conditions,
	}
	_, params.OrderBy = l.BuildOrderBySQL()
	recs, _, err := l.Searcher(mb.NewModelSlice(), params, ctx)
	if err != nil || recs == nil {
		return nil
	}
	v := reflect.Indirect(reflect.ValueOf(recs))
	if v.Kind() != reflect.Slice || v.Len() == 0 {
		return nil
	}
	icon := mb.menuIcon
	if icon == "" {
		icon = defaultMenuIcon(mb.label)
	}
	g := &globalSearchGroup{Title: mb.TPageLabel(ctx.Context()), Icon: icon}
	for i := 0; i < v.Len(); i++ {
		obj := v.Index(i).Interface()
		id, err := mb.RecordID(obj)
		if err != nil {
			continue
		}
		var parents IDSlice
		if child {
			// under its parents: listed there by the user
			if parents = mb.globalSearchParents(obj); len(parents) == 0 || mb.permissioner.Lister(ctx.R, parents...).Denied() {
				continue
			}
		}
		href := mb.Info().EditingHref(id, parents...)
		if mb.hasDetailing {
			href = mb.Info().DetailingHref(id, parents...)
		}
		g.Items = append(g.Items, globalSearchItem{
			Title: mb.RecordTitle(obj, ctx),
			Href:  account.URL(b.prefix, href, n),
		})
	}
	if len(g.Items) == 0 {
		return nil
	}
	return g
}

// menuSearchEntries are the items of the menu the search finds by name: the
// pages, their group the way to them.
func menuSearchEntries(nodes []*MenuNode, trail string) (out []map[string]string) {
	for _, n := range nodes {
		href, _ := n.Props["href"].(string)
		icon, _ := n.Props["prependIcon"].(string)
		if href != "" {
			out = append(out, map[string]string{"title": n.Title, "sub": trail, "href": href, "icon": icon})
		}
		if len(n.Children) > 0 {
			t := n.Title
			if trail != "" {
				t = trail + " › " + n.Title
			}
			out = append(out, menuSearchEntries(n.Children, t)...)
		}
	}
	return
}

// globalSearchBar is the search of the top bar: the field that opens it (an
// icon on a phone), its palette, Ctrl K.
func (b *Builder) globalSearchBar(ctx *web.EventContext, m *Messages) h.HTMLComponent {
	nodes, _, _ := b.menuNodes(ctx)
	menu := menuSearchEntries(nodes, "")
	url := b.URIPrefixCtx(ctx.R.Context())

	field := VBtn("").
		Variant(VariantTonal).
		Class("rvq-global-search__field text-none d-none d-md-flex").
		Attr("@click", "vars.gs.show()").
		Attr("aria-label", m.SearchAdmin).
		Children(
			VIcon("mdi-magnify").Size(SizeSmall).Class("mr-2"),
			h.Span(m.SearchAdmin).Class("text-medium-emphasis flex-grow-1 text-left"),
			h.Tag("kbd").Children(h.Text("Ctrl K")).Class("rvq-global-search__kbd ml-4"),
		)
	icon := VBtn("").Icon("mdi-magnify").Variant(VariantText).
		Class("d-md-none").
		Attr("title", m.SearchAdmin).Attr("aria-label", m.SearchAdmin).
		Attr("@click", "vars.gs.show()")

	results := VList(
		h.Template(
			VListSubheader(h.Text("{{ it.header }}")).Attr("v-if", "it.header"),
			VListItem().
				Attr("v-else", true).
				Attr(":href", "it.href").
				Attr(":title", "it.title").
				Attr(":subtitle", "it.sub").
				Attr(":prepend-icon", "it.icon || 'mdi-file-outline'").
				Attr(":active", "i === vars.gs.sel").
				Attr(":data-gs-index", "i").
				Attr("@mouseenter", "vars.gs.sel = i"),
		).Attr("v-for", "(it, i) in vars.gs.items()").Attr(":key", "i"),
	).Density(DensityCompact).Class("rvq-global-search__results")

	palette := VDialog(
		VCard(
			VTextField().
				Attr("v-model", "vars.gs.q").
				Attr("@update:model-value", "vars.gs.typed()").
				Attr("@keydown", "vars.gs.key($event)").
				Placeholder(m.SearchAdmin).
				PrependInnerIcon("mdi-magnify").
				Variant(FieldVariantSolo).
				Flat(true).
				HideDetails(true).
				Autofocus(true).
				Attr(":loading", "vars.gs.loading").
				Class("rvq-global-search__input"),
			VDivider(),
			h.Div(
				results,
				h.Div(h.Text(m.SearchNoResults)).
					Attr("v-if", "vars.gs.q.trim().length >= 2 && !vars.gs.loading && vars.gs.items().length === 0").
					Class("text-center text-medium-emphasis py-8"),
				h.Div(h.Text(m.SearchHint)).
					Attr("v-if", "vars.gs.q.trim().length < 2").
					Class("text-center text-medium-emphasis py-6 text-body-2"),
			).Style("max-height: 60vh; overflow-y: auto"),
		),
	).
		Attr("v-model", "vars.gs.open").
		MaxWidth(640).
		Location("top").
		Class("rvq-global-search").
		Attr("content-class", "mt-12")

	return h.Tag("user-component").
		Attr(":assign", "[[vars, {gs: vars.gs || {open: false, q: '', groups: [], sel: 0, loading: false, seq: 0}}]]").
		Attr(":setup", "[("+globalSearchScript(menu, m.SearchMenu, url)+")]").
		Children(field, icon, palette)
}

// globalSearchScript runs once a page: the acts of the search (vars.gs) and
// Ctrl K.
func globalSearchScript(menu []map[string]string, menuTitle, url string) string {
	return `({window}) => {
	const w = window, gs = vars.gs;
	gs.menu = ` + h.JSONString(menu) + `;
	gs.menuTitle = ` + h.JSONString(menuTitle) + `;
	const norm = (s) => String(s || "").normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLowerCase();
	gs.show = () => { gs.open = true; gs.sel = 0 };
	// the items: the menu's that match, then the records, by model
	gs.items = () => {
		const q = norm(gs.q.trim());
		if (q.length < 2) return [];
		const out = [];
		const menu = gs.menu.filter((e) => norm(e.title).includes(q) || norm(e.sub).includes(q)).slice(0, 8);
		if (menu.length) { out.push({header: gs.menuTitle}); out.push(...menu) }
		for (const g of gs.groups) {
			out.push({header: g.title});
			for (const it of g.items) out.push({title: it.title, href: it.href, icon: g.icon, sub: g.title});
		}
		return out;
	};
	const search = () => {
		const q = gs.q.trim(), n = ++gs.seq;
		if (q.length < 2) { gs.groups = []; gs.loading = false; return }
		gs.loading = true;
		plaid().vars(vars).url(` + h.JSONString(url) + `).eventFunc(` + h.JSONString(GlobalSearchEvent) + `).query("q", q).go().
			then((r) => { if (n === gs.seq) { gs.groups = (r && r.data) || []; gs.loading = false; gs.first() } }).
			catch(() => { if (n === gs.seq) gs.loading = false });
	};
	gs.typed = () => { gs.first(); w.clearTimeout(gs.timer); gs.timer = w.setTimeout(search, 250) };
	// the first item that is not a title
	gs.first = () => { const its = gs.items(); gs.sel = its.findIndex((it) => !it.header) };
	gs.key = (e) => {
		const its = gs.items();
		const step = (d) => {
			let i = gs.sel;
			for (let k = 0; k < its.length; k++) {
				i = (i + d + its.length) % its.length;
				if (!its[i].header) { gs.sel = i; break }
			}
			const el = w.document.querySelector('[data-gs-index="' + gs.sel + '"]');
			el && el.scrollIntoView({block: "nearest"});
		};
		if (e.key === "ArrowDown") { e.preventDefault(); step(1) }
		else if (e.key === "ArrowUp") { e.preventDefault(); step(-1) }
		else if (e.key === "Enter") {
			const it = its[gs.sel];
			if (it && it.href) { e.preventDefault(); gs.open = false; w.location.href = w.__rvqAccount ? w.__rvqAccount.fix(it.href) : it.href }
		}
	};
	if (!w.__rvqGlobalSearchKey) {
		w.__rvqGlobalSearchKey = true;
		w.addEventListener("keydown", (e) => {
			if ((e.ctrlKey || e.metaKey) && !e.altKey && (e.key === "k" || e.key === "K")) {
				e.preventDefault();
				vars.gs.show();
			}
		});
	}
}`
}
