package userdocs

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// The tree of the permissions as a page shows it (PermissionsTreeView): a
// tree table — the name in the tree; its kind; what of it a user's stack
// allows (Active, when asked); its description; its permissions, one a
// line —, a search, a filter of the kinds, the groups of the menu or not,
// the views (by verb, by kind) as tabs, a button to open or close every row.
//
// The tree is never in the page: the one of the view chosen comes by an
// event (PermissionsTreeEventFunc) when it is shown, and is kept in the
// browser, out of Vue's reactivity — the same array every render, so the
// table is not made again by a change of the page's state.

// PermissionsTreeEvent is the event of the documentation that answers the
// tree of the permissions (PermissionsTreeEventFunc).
const PermissionsTreeEvent = "userdocs_PermissionsTree"

// PermTreeViewOptions are where the tree of a PermissionsTreeView comes
// from and what it shows.
type PermTreeViewOptions struct {
	// URL, Event: the event that answers the tree (PermissionsTreeEventFunc),
	// at the address of the model that registered it; Query, what it is
	// asked with besides the view (the user's id)
	URL, Event string
	Query      map[string]string
	// Active: the column of what the user's stack allows
	Active bool
}

// PermKindLabels are the labels of the kinds of the rows (PermKinds), in the
// language of ctx.
func PermKindLabels(ctx context.Context) map[string]string {
	m := GetMessages(ctx)
	return map[string]string{
		PermKindGroup:     m.PermKindGroup,
		PermKindModel:     m.PermKindModel,
		PermKindSingleton: m.PermKindSingleton,
		PermKindVerb:      m.PermKindVerb,
		PermKindAction:    m.PermKindAction,
		PermKindField:     m.PermKindField,
		PermKindInline:    m.PermKindInline,
		PermKindSection:   m.PermKindSection,
		PermKindPage:      m.PermKindPage,
		PermKindCheck:     m.PermKindCheck,
	}
}

// PermissionsTreeEventFunc answers the tree of the permissions of the admin
// p (r.Data): of the view asked (the form's view, "verb" or "kind"), with
// the groups of the menu or not (groups, "1"). allows, when not nil, gives
// what a stack allows, for the column Active (SetActive): the user's, of
// the request.
func PermissionsTreeEventFunc(p *presets.Builder, allows func(ctx *web.EventContext) (func(PermCheck) bool, error)) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		opts := PermTreeOptions{View: PermTreeByVerb, Groups: ctx.R.FormValue("groups") == "1"}
		if PermTreeView(ctx.R.FormValue("view")) == PermTreeByKind {
			opts.View = PermTreeByKind
		}
		items := PermissionsTree(ctx.Context(), p, opts)
		if allows != nil {
			a, err := allows(ctx)
			if err != nil {
				return r, err
			}
			SetActive(items, a)
		}
		r.Data = items
		return
	}
}

// PermissionsTreeView is the tree of the permissions as a tree table, its
// tree from the event of o.
func PermissionsTreeView(ctx context.Context, o PermTreeViewOptions) h.HTMLComponent {
	m := GetMessages(ctx)
	labels := PermKindLabels(ctx)
	var kinds []map[string]string
	for _, k := range PermKinds {
		kinds = append(kinds, map[string]string{"title": labels[k], "value": k})
	}
	query := ""
	keys := make([]string, 0, len(o.Query))
	for k := range o.Query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		query += ".query(" + strconv.Quote(k) + ", " + strconv.Quote(o.Query[k]) + ")"
	}

	headers := []map[string]any{
		{"title": m.PermTreeKind, "key": "kind", "sortable": false, "nowrap": true},
		{"title": m.PermTreePath, "key": "path", "sortable": false},
	}
	state := func(value, icon, color, title string) h.HTMLComponent {
		return v.VIcon(icon).Color(color).Size(v.SizeSmall).
			Attr("v-if", fmt.Sprintf("item.active === %q", value)).Attr("title", title)
	}
	slots := []h.HTMLComponent{
		h.Template(h.Div(h.Code("").Attr("v-text", "r")).Class("text-caption").
			Attr("v-for", "(r, i) in item.resources || []").Attr(":key", "i")).
			Attr("v-slot:item.resources", "{ item }"),
		h.Template(h.Span("").Attr("v-text", fmt.Sprintf("(%s)[item.kind] || ''", h.JSONString(labels)))).
			Attr("v-slot:item.kind", "{ item }"),
		// the path: its labels ([id, label] each)
		h.Template(h.Span("").Attr("v-text", fmt.Sprintf("(item.path || []).map((p) => p[1]).join(%q)", PathSep))).
			Attr("v-slot:item.path", "{ item }"),
	}
	if o.Active {
		headers = append(headers, map[string]any{"title": m.PermTreeActive, "key": "active", "sortable": false, "align": "center", "width": "1%"})
		slots = append(slots, h.Template(
			state(PermActiveAll, "mdi-check-circle", "success", m.PermTreeActiveAll),
			state(PermActiveSome, "mdi-circle-half-full", "warning", m.PermTreeActiveSome),
			state(PermActiveNone, "mdi-close-circle-outline", "grey", m.PermTreeActiveNone),
		).Attr("v-slot:item.active", "{ item }"))
	}
	headers = append(headers,
		map[string]any{"title": m.PermTreeDescription, "key": "description", "sortable": false},
		map[string]any{"title": m.PermTreePermission, "key": "resources", "sortable": false})

	table := vx.VXTreeDataTable(slots...).ItemsVar("st.items").Headers(headers).
		// a tree a table: made again for another
		Attr(":key", "st.shown").
		Attr("v-if", "st.items").
		Attr("control-title", m.PermTreeName).
		Attr(":open-all", "st.open").
		// the search typed, once the typing stops (st.q)
		Attr(":search", "st.q || ''").
		Attr(":search-keys", `["title", "path", "description", "resources"]`).
		// the kinds chosen: their rows and their ancestors
		Attr("filter-key", "kind").
		Attr(":filter-values", "st.kinds").
		Attr(":settings", "false").
		Attr(":guides", "true").
		Attr(":table-props", `{density: "compact"}`)

	body := h.Div(
		h.Div(
			v.VTextField().Attr("v-model", "st.search").
				Label(m.PermTreeSearch).PrependInnerIcon("mdi-magnify").
				Clearable(true).Density(v.DensityCompact).HideDetails(true).
				Variant(v.FieldVariantUnderlined).Attr("data-permissions-tree-search", true),
			v.VSelect().Items(kinds).Attr("v-model", "st.kinds").
				Label(m.PermTreeKinds).Multiple(true).Chips(true).ClosableChips(true).
				Clearable(true).Density(v.DensityCompact).HideDetails(true).
				Variant(v.FieldVariantUnderlined).Style("max-width: 420px").
				Attr("data-permissions-tree-kinds", true),
			v.VSwitch().Label(m.PermTreeGroups).Attr("v-model", "st.groups").
				Density(v.DensityCompact).Color("primary").HideDetails(true).Class("flex-grow-0 text-no-wrap").
				Attr("data-permissions-tree-groups", true),
			v.VBtn("").Variant(v.VariantTonal).Size(v.SizeSmall).Color("primary").
				Attr(":prepend-icon", "st.open ? 'mdi-collapse-all-outline' : 'mdi-expand-all-outline'").
				Attr("@click", "st.open = !st.open").
				Attr("data-permissions-tree-toggle", true).
				Children(h.Span("").Attr("v-text", fmt.Sprintf("st.open ? %q : %q", m.PermTreeCollapseAll, m.PermTreeExpandAll))),
		).Class("d-flex align-center ga-3 mb-2"),
		// the views, as tabs
		v.VTabs(
			v.VTab(h.Text(m.PermTreeByVerb)).Value(string(PermTreeByVerb)),
			v.VTab(h.Text(m.PermTreeByKind)).Value(string(PermTreeByKind)),
		).Attr("v-model", "st.view").Color("primary").Density(v.DensityCompact).Class("mb-2").
			Attr("data-permissions-tree-views", true),
		v.VProgressLinear().Indeterminate(true).Color("primary").Attr("v-if", "st.loading"),
		v.VAlert(h.Text(m.PermTreeLoadError)).Type("error").Density(v.DensityCompact).Attr("v-if", "st.error"),
		table,
	).Attr("data-permissions-full-tree", true)

	// the tree of the view shown: asked once, kept out of the reactivity
	// (markRaw) — the same array every render
	setup := `({$scope, inject, watch, debounce, Vue}) => {
	const plaid = inject('plaid'), st = $scope.value.st, trees = {};
	// the search: the table's once the typing stops — every row is gone
	// through by it
	watch(() => st.search, debounce((q) => { st.q = q }, 300));
	const key = () => st.view + '/' + (st.groups ? 1 : 0);
	const load = () => {
		const k = key();
		if (trees[k]) { st.items = trees[k]; st.shown = k; return }
		st.loading = true; st.error = false;
		plaid().url(` + strconv.Quote(o.URL) + `).eventFunc(` + strconv.Quote(o.Event) + `)` + query + `.
			query("view", st.view).query("groups", st.groups ? "1" : "").go().
			then((r) => {
				trees[k] = Vue.markRaw((r && r.data) || []);
				if (k === key()) { st.items = trees[k]; st.shown = k }
			}).
			catch(() => { st.error = true }).
			finally(() => { st.loading = false });
	};
	watch(() => [st.view, st.groups], load);
	load();
}`
	return vue.UserComponent(body).
		Scope("st", map[string]any{"view": string(PermTreeByVerb), "groups": false, "kinds": []string{},
			"search": "", "q": "", "open": nil, "items": nil, "shown": "", "loading": false, "error": false}).
		Setup(setup)
}
