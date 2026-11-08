package presets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// A record chosen in the side menu (MenuRecordSelector): a model whose
// records have resources of their own — a project, its movements and its
// payments. Its item in the menu is a group: first the selector — a chip
// with the record chosen, or "Select", that opens the dialog of its records
// —, the resources of the one chosen under it (its children, Items, in
// their menu groups, and pages of its detail, Pages, with the record in their
// path); then the item that lists the records.
//
// It is on for a model of no parent, a child of a singleton, or a child of
// a model that has it too (MenuRecordSelector of both): the records of its
// parents are then known — the ones chosen above it.
//
// The record chosen is the one of the page (…/<model>/{id}/…), else the one
// chosen in this tab (vars.menuRecords, carried by every request of the page
// in MenuRecordsHeader; none when it was cleared), else the only one, else
// the default of the request (DefaultID: the user's), else the default
// (DefaultConditions). Nothing is kept in the browser: a page
// opened anew (F5) says it by its URI. A change in the dialog, or the clear,
// keeps it in the tab and draws its item again, the page staying. The chip is the title of its item (no label of the
// model): the record chosen — an icon at the right of the item clears it —,
// or "Select".

// MenuRecordSelectorConfig says what of the record chosen in the menu is
// shown and which records may be chosen.
type MenuRecordSelectorConfig struct {
	enabled bool
	// Items are the URI names of the children shown under the record (its
	// resources), in their order; none, every child.
	Items []string
	// Pages are the paths of pages of the record's detail shown after them
	// ("/relatorios").
	Pages []string
	// PageGroups are groups of pages of the record's detail, after Pages:
	// a group of the menu each (Relatórios: its pages under it).
	PageGroups []MenuRecordPageGroup
	// Conditions narrow the records that may be chosen (the active ones).
	Conditions SQLConditions
	// DefaultID is the id of the record chosen when none was, for the
	// request — the user's own default —, before DefaultConditions; "" or
	// one that may not be chosen, none.
	DefaultID func(ctx *web.EventContext) string
	// DefaultConditions say the record chosen when none was: the first that
	// meets them (the default one); none, none.
	DefaultConditions SQLConditions
	// KeywordColumns are searched in the dialog; none, the listing's.
	KeywordColumns []string
	// OrderBy orders the records of the dialog; empty, the listing's.
	OrderBy string
	// External: the item of the model — and its choice — is whoever owns
	// it's (the organization's, at the top of the menu); its children have
	// their selectors still. Set before MenuRecordSelector.
	External bool
}

// MenuRecordPageGroup is a group of pages of the detail of the record chosen
// in the menu (MenuRecordSelectorConfig.PageGroups): its name (of its node),
// title and icon, the paths of its pages.
type MenuRecordPageGroup struct {
	Name  string
	Title func(ctx context.Context) string
	Icon  string
	Pages []string
}

// MenuRecordSelector turns on (or off) the record chosen in the side menu
// for the model (see MenuRecordSelectorConfig).
func (mb *ModelBuilder) MenuRecordSelector(v bool) *ModelBuilder {
	c := mb.MenuRecordSelectorConfig()
	if v && !c.enabled && !c.External {
		mb.configureMenuRecordChoose()
	}
	c.enabled = v
	return mb
}

// GetMenuRecordSelector says whether the model asked for the record chosen
// in the menu.
func (mb *ModelBuilder) GetMenuRecordSelector() bool {
	return mb.menuRecordSelector != nil && mb.menuRecordSelector.enabled
}

// MenuRecordSelectorConfig is the config of the record chosen in the menu.
func (mb *ModelBuilder) MenuRecordSelectorConfig() *MenuRecordSelectorConfig {
	if mb.menuRecordSelector == nil {
		mb.menuRecordSelector = &MenuRecordSelectorConfig{}
	}
	return mb.menuRecordSelector
}

// HasMenuRecordSelector says the model has its record chosen in the menu:
// it asked for it, and its parent — none, a singleton, or one of a record
// chosen in the menu — lets it.
func (mb *ModelBuilder) HasMenuRecordSelector() bool {
	if !mb.GetMenuRecordSelector() {
		return false
	}
	p := mb.parent
	return p == nil || p.singleton || p.HasMenuRecordSelector()
}

// withParents is ctx with parents as the ids of the parents of its request:
// what the reads of a child are scoped by.
func withParents(ctx *web.EventContext, parents IDSlice) *web.EventContext {
	c := *ctx
	c.R = ctx.R.WithContext(context.WithValue(ctx.R.Context(), ParentsModelIDKey, parents))
	return &c
}

// MenuRecordsHeader is the header of the requests of the page that carries
// the records chosen in the menu in this tab (vars.menuRecords, a JSON
// object): the key of the item of each selector (menuRecordKey) and the
// record as its path of permission (RecordPath), "" cleared.
const MenuRecordsHeader = "X-Rvq-Menu-Records"

// menuRecordsSetup is the setup of the side menu (CreateMenus: scope.menu,
// its locals) for the records chosen in it in this tab: vars.menuRecords,
// the key of the item of each selector and its record (its path of
// permission, "" cleared), carried by every request (MenuRecordsHeader); the
// ones the page's URI says, kept in it; vars.menuRecordSet(key, value, url),
// a record kept and its item drawn again (menuRecordNodeEvent at url).
var menuRecordsSetup = `({scope, inject}) => {
	const vars = inject('vars')
	const plaid = inject('plaid')
	const menu = scope.menu
	if (!vars.menuRecords) vars.menuRecords = {}
	if (!vars.__requestHeaders) vars.__requestHeaders = {}
	vars.__requestHeaders['` + MenuRecordsHeader + `'] = () => {
		const m = vars.menuRecords
		return m && Object.keys(m).length ? JSON.stringify(m) : ''
	}
	const walk = (ns, f) => (ns || []).forEach((n) => { f(n); walk(n.children, f) })
	walk(menu.items, (n) => {
		const p = n.props
		if (p && p.recordKey && p.recordFromURI) vars.menuRecords[p.recordKey] = p.recordPath
	})
	vars.menuRecordSet = (key, value, url) => {
		vars.menuRecords[key] = value
		return plaid().eventFunc('` + menuRecordNodeEvent + `').url(url).go().then((r) => {
			const d = r && r.data
			if (!d || !d.node) return
			const replace = (ns) => {
				for (let i = 0; i < (ns || []).length; i++) {
					if (ns[i].value === key) { ns.splice(i, 1, d.node); return true }
					if (replace(ns[i].children)) return true
				}
				return false
			}
			replace(menu.items)
			;(d.opened || []).forEach((v) => { if (!menu.opened.includes(v)) menu.opened.push(v) })
		})
	}
}`

// menuRecordKey is the unique id of the item of the selector of mb, of the
// parents: the key of its record in vars.menuRecords.
func (mb *ModelBuilder) menuRecordKey(parents IDSlice) string {
	return "g:menu-record:" + mb.uriName + ":" + parents.String()
}

// menuRecordValue is the record mid of mb, of the parents, as
// vars.menuRecords keeps it: its path of permission (RecordPath); its id
// when there is no generator of them.
func (mb *ModelBuilder) menuRecordValue(parents IDSlice, mid ID) string {
	if p := mb.RecordPath(parents, mid); p != "" {
		return p
	}
	return mid.String()
}

// menuRecordsOf are the records chosen in the menu the request carries
// (MenuRecordsHeader); nil, none.
func menuRecordsOf(r *http.Request) map[string]string {
	v := r.Header.Get(MenuRecordsHeader)
	if v == "" {
		return nil
	}
	var m map[string]string
	if json.Unmarshal([]byte(v), &m) != nil {
		return nil
	}
	return m
}

// menuRecordOfValue is the record of mb, of the parents, that v names (a
// value of vars.menuRecords): its path of permission, read back and the
// same as the one of the record found (of mb, of these parents); else its
// id. nil when it names none.
func (mb *ModelBuilder) menuRecordOfValue(ctx *web.EventContext, parents IDSlice, v string) any {
	id := v
	if strings.HasSuffix(v, ":") {
		t, err := mb.p.ParseRecordPath(v)
		if err != nil || len(t.Models) == 0 || t.Listing || t.Page != nil || t.Action != nil {
			return nil
		}
		last := t.Models[len(t.Models)-1]
		if last.Model != mb || mb.RecordPath(parents, last.ID) != v {
			return nil
		}
		id = last.ID.String()
	}
	return mb.menuRecordFind(ctx, parents, id)
}

// menuRecordSearch searches the records of mb that may be chosen, of the
// parents: keyword, more conditions, the page (from 1) of perPage.
func (mb *ModelBuilder) menuRecordSearch(ctx *web.EventContext, parents IDSlice, keyword string, conds SQLConditions,
	page, perPage int64) (objs []any, total int, err error) {
	c := mb.MenuRecordSelectorConfig()
	cols := c.KeywordColumns
	if len(cols) == 0 && mb.listing != nil {
		cols = mb.listing.searchColumns
	}
	order := c.OrderBy
	if order == "" && mb.listing != nil {
		order = mb.listing.orderBy
	}
	params := &SearchParams{
		KeywordColumns: cols,
		Keyword:        keyword,
		SQLConditions:  append(append(SQLConditions{}, c.Conditions...), conds...),
		Page:           page,
		PerPage:        perPage,
		OrderBy:        order,
		Context:        ctx.R.Context(),
	}
	r, total, err := mb.Searcher(mb.NewModelSlice(), params, withParents(ctx, parents))
	if err != nil {
		return
	}
	v := reflect.Indirect(reflect.ValueOf(r))
	if v.Kind() != reflect.Slice {
		return
	}
	for i := 0; i < v.Len(); i++ {
		e := v.Index(i)
		if e.Kind() != reflect.Ptr {
			e = e.Addr()
		}
		objs = append(objs, e.Interface())
	}
	return
}

// menuRecordReadable says the user may read the record obj of the parents;
// its id.
func (mb *ModelBuilder) menuRecordReadable(r *http.Request, parents IDSlice, obj any) (ID, bool) {
	mid, err := mb.RecordID(obj)
	if err != nil {
		return mid, false
	}
	return mid, !mb.permissioner.Reader(r, mid, parents...).Denied()
}

// menuRecordFind is the record id of the parents, when the user may read
// it; nil, none.
func (mb *ModelBuilder) menuRecordFind(ctx *web.EventContext, parents IDSlice, id string) any {
	if id == "" {
		return nil
	}
	mid, err := mb.ParseRecordID(id)
	if err != nil || mid.IsZero() {
		return nil
	}
	obj := mb.NewModel()
	if mb.Fetcher(obj, mid, withParents(ctx, parents)) != nil {
		return nil
	}
	if _, ok := mb.menuRecordReadable(ctx.R, parents, obj); !ok {
		return nil
	}
	return obj
}

// menuRecordChoice is the record of mb chosen in the menu: the record and
// its id (nil, none), and whether it is the one of the page (its URI).
type menuRecordChoice struct {
	obj     any
	id      ID
	fromURI bool
}

// menuRecordChosen is the record of mb chosen, of the parents: the one of
// the page (its URI), else the one of this tab (MenuRecordsHeader: none,
// cleared), else the only one, else the default; nil, none.
func (mb *ModelBuilder) menuRecordChosen(ctx *web.EventContext, parents IDSlice) (c menuRecordChoice) {
	ok := func(obj any) menuRecordChoice {
		if obj == nil {
			return menuRecordChoice{}
		}
		mid, readable := mb.menuRecordReadable(ctx.R, parents, obj)
		if !readable {
			return menuRecordChoice{}
		}
		return menuRecordChoice{obj: obj, id: mid}
	}
	// the page's
	if rest, found := strings.CutPrefix(ctx.R.URL.Path, mb.Info().ListingHref(parents...)+"/"); found {
		id, _, _ := strings.Cut(rest, "/")
		if obj := mb.menuRecordFind(ctx, parents, id); obj != nil {
			if c = ok(obj); c.obj != nil {
				c.fromURI = true
				return
			}
		}
	}
	// the tab's — "", cleared
	if v, found := menuRecordsOf(ctx.R)[mb.menuRecordKey(parents)]; found {
		if v == "" {
			return menuRecordChoice{}
		}
		if c = ok(mb.menuRecordOfValue(ctx, parents, v)); c.obj != nil {
			return
		}
	}
	// the only one
	if objs, total, err := mb.menuRecordSearch(ctx, parents, "", nil, 1, 2); err == nil && total == 1 && len(objs) == 1 {
		if c = ok(objs[0]); c.obj != nil {
			return
		}
	}
	// the default of the request (the user's)
	if f := mb.MenuRecordSelectorConfig().DefaultID; f != nil {
		if obj := mb.menuRecordFind(ctx, parents, f(ctx)); obj != nil {
			if c = ok(obj); c.obj != nil {
				return
			}
		}
	}
	// the default
	if cfg := mb.MenuRecordSelectorConfig(); len(cfg.DefaultConditions) > 0 {
		if objs, _, err := mb.menuRecordSearch(ctx, parents, "", cfg.DefaultConditions, 1, 1); err == nil && len(objs) > 0 {
			return ok(objs[0])
		}
	}
	return menuRecordChoice{}
}

// menuRecordItems are the children shown under the record chosen.
func (mb *ModelBuilder) menuRecordItems() (out []*ModelBuilder) {
	c := mb.MenuRecordSelectorConfig()
	if len(c.Items) == 0 {
		return mb.children
	}
	for _, name := range c.Items {
		for _, child := range mb.children {
			if child.uriName == name {
				out = append(out, child)
			}
		}
	}
	return
}

// MenuNodeOf is the item of mb in the side menu, of the records of its
// parents (their ids): its listing; with its record chosen in the menu
// (HasMenuRecordSelector), the group of the selector — the resources of the
// record chosen under it — and of its listing. active is the value of the
// item of the page, opened the values of the items to open.
func (mb *ModelBuilder) MenuNodeOf(ctx *web.EventContext, parents IDSlice) (n *MenuNode, active string, opened []string) {
	path := ctx.R.URL.Path
	isActive := func(href string) bool { return path == href || strings.HasPrefix(path, href+"/") }
	icon := mb.menuIcon
	if icon == "" {
		icon = defaultMenuIcon(mb.label)
	}
	listHref := mb.Info().ListingHref(parents...)
	if !mb.HasMenuRecordSelector() {
		n = &MenuNode{Title: mb.TTitlePlural(ctx.Context()), Value: listHref,
			Props: map[string]any{"href": listHref, "prependIcon": icon}}
		if isActive(listHref) {
			active = listHref
		}
		return
	}

	msgs := MustGetMessages(ctx.Context())
	selector, active, opened := mb.menuRecordSelectorNode(ctx, parents, mb.menuRecordChosen(ctx, parents))
	list := &MenuNode{Title: fmt.Sprintf(msgs.MenuRecordList, strings.ToLower(mb.TTitlePlural(ctx.Context()))),
		Value: listHref, Props: map[string]any{"href": listHref, "prependIcon": "mdi-format-list-bulleted"}}
	n = &MenuNode{
		Title:    mb.TTitlePlural(ctx.Context()),
		Value:    "g:menu-records:" + mb.uriName + ":" + parents.String(),
		Props:    map[string]any{"prependIcon": icon},
		Children: []*MenuNode{selector, list},
	}
	if len(selector.Children) > 0 {
		opened = append(opened, n.Value)
	}
	if active == "" && isActive(listHref) {
		active = listHref
	}
	if active != "" {
		opened = append(opened, n.Value)
	}
	return
}

// menuRecordSelectorNode is the item of the selector of mb, of the parents,
// with the record ch chosen: its chip — the record, or "Select" —, and under
// it, chosen, the resources of the record (in their menu groups), the pages
// of its detail and its groups of pages. Its props carry what the menu keeps
// of it in the tab (vars.menuRecords): recordKey, recordPath (""
// none) and recordFromURI.
func (mb *ModelBuilder) menuRecordSelectorNode(ctx *web.EventContext, parents IDSlice, ch menuRecordChoice) (selector *MenuNode, active string, opened []string) {
	path := ctx.R.URL.Path
	isActive := func(href string) bool { return path == href || strings.HasPrefix(path, href+"/") }
	icon := mb.menuIcon
	if icon == "" {
		icon = defaultMenuIcon(mb.label)
	}
	listHref := mb.Info().ListingHref(parents...)
	msgs := MustGetMessages(ctx.Context())
	key := mb.menuRecordKey(parents)
	title := mb.TTitle(ctx.Context())
	selector = &MenuNode{
		Title: title,
		Value: key,
		Props: map[string]any{"prependIcon": icon, "recordKey": key, "recordPath": "", "recordFromURI": ch.fromURI},
	}
	// the chip the title of the item (no label of the model), chosen or not
	chip := &MenuChip{Title: msgs.MenuRecordSelect, Icon: "mdi-swap-horizontal", Color: "grey", Inline: true,
		Hint:  fmt.Sprintf(msgs.MenuRecordChange, strings.ToLower(title)),
		Event: actions.OpenActionDialog,
		URL:   listHref,
		Query: url.Values{
			ParamAction:  {ActionChooseMenuRecord},
			ParamOverlay: {string(actions.Dialog)},
		}.Encode()}
	selector.SetChip(chip)
	if ch.obj == nil {
		return
	}
	obj, mid := ch.obj, ch.id
	selector.Props["recordPath"] = mb.menuRecordValue(parents, mid)
	// its record in the chip, the clear at the right of the item
	chip.Title, chip.Color = mb.RecordTitle(obj, ctx), "primary"
	selector.Title = chip.Title
	chip.Clear = &MenuChipAction{Icon: "mdi-close", Hint: msgs.MenuRecordClear,
		Event: menuRecordClearEvent, URL: listHref}
	opened = append(opened, selector.Value)
	under := append(append(IDSlice{}, parents...), mid)
	// its resources, in their menu groups (Financeiro, People…) in the
	// order they come
	groups := map[*MenuGroupBuilder]*MenuNode{}
	for _, child := range mb.menuRecordItems() {
		if child.notInMenu || child.permissioner.Lister(ctx.R, under...).Denied() {
			continue
		}
		cn, ca, co := child.MenuNodeOf(ctx, under)
		opened = append(opened, co...)
		g := child.MenuGroup()
		if g == nil {
			selector.Children = append(selector.Children, cn)
		} else {
			gn := groups[g]
			if gn == nil {
				gi := g.GetIcon()
				if gi == "" {
					gi = defaultMenuIcon(g.Name())
				}
				gn = &MenuNode{Title: g.TTitle(ctx.Context()), Value: selector.Value + ":" + g.Name(),
					Props: map[string]any{"prependIcon": gi}}
				groups[g] = gn
				selector.Children = append(selector.Children, gn)
			}
			gn.Children = append(gn.Children, cn)
			if ca != "" {
				opened = append(opened, gn.Value)
			}
		}
		if ca != "" {
			active = ca
		}
	}
	// a page of its detail: its node, nil when not to be shown
	pageNode := func(p string) *MenuNode {
		page := mb.Detailing().PagesRegistrator().GetHttpPage(p)
		// the record readable (chosen); the page's own permission, when
		// it has one
		if page == nil || (page.GetBaseVerifier() != nil && page.Verifier(ctx.R).Denied()) {
			return nil
		}
		href := mb.Info().DetailingHref(mid.String(), parents...) + p
		pi := page.GetMenuIcon()
		if pi == "" {
			pi = "mdi-circle-small"
		}
		if isActive(href) {
			active = href
		}
		return &MenuNode{Title: page.TTitle(ctx.Context()), Value: href,
			Props: map[string]any{"href": href, "prependIcon": pi}}
	}
	for _, p := range mb.MenuRecordSelectorConfig().Pages {
		if pn := pageNode(p); pn != nil {
			selector.Children = append(selector.Children, pn)
		}
	}
	for _, pg := range mb.MenuRecordSelectorConfig().PageGroups {
		gi := pg.Icon
		if gi == "" {
			gi = defaultMenuIcon(pg.Name)
		}
		gt := pg.Name
		if pg.Title != nil {
			gt = pg.Title(ctx.Context())
		}
		gn := &MenuNode{Title: gt, Value: selector.Value + ":pages:" + pg.Name,
			Props: map[string]any{"prependIcon": gi}}
		was := active
		for _, p := range pg.Pages {
			if pn := pageNode(p); pn != nil {
				gn.Children = append(gn.Children, pn)
			}
		}
		if len(gn.Children) == 0 {
			continue
		}
		if active != was {
			opened = append(opened, gn.Value)
		}
		selector.Children = append(selector.Children, gn)
	}
	return
}

// ActionChooseMenuRecord is the listing action of a model of a record chosen
// in the menu that chooses it (its dialog).
const ActionChooseMenuRecord = "choose_menu_record"

const (
	menuRecordListEvent  = "presets_menuRecord_list"
	menuRecordClearEvent = "presets_menuRecord_clear"
	// menuRecordNodeEvent is the item of the selector again, of the record
	// the tab chose (MenuRecordsHeader): a change of it in the dialog, its
	// clear
	menuRecordNodeEvent  = "presets_menuRecord_node"
	menuRecordListPortal = "presetsMenuRecordList"
	paramMenuRecordQ     = "menu_record_q"
	paramMenuRecordPage  = "menu_record_page"
	paramMenuRecordID    = "menu_record_id"
	menuRecordPageSize   = 10
)

// refererRequest is r at the page the dialog was opened from (its Referer):
// the record of that page is the one chosen.
func refererRequest(r *http.Request) *http.Request {
	ref, err := url.Parse(r.Referer())
	if err != nil || ref.Path == "" {
		return r
	}
	rr := r.Clone(r.Context())
	rr.URL = ref
	return rr
}

// menuRecordList is the list of the dialog: the records of the parents of
// the request with q, the page page of them, the one chosen marked — a click
// on one does the action (ActionChooseMenuRecord) with it.
func (mb *ModelBuilder) menuRecordList(ctx *web.EventContext, q string, page int64) h.HTMLComponent {
	msgs := MustGetMessages(ctx.Context())
	parents := ParentsModelID(ctx.R)
	if page < 1 {
		page = 1
	}
	objs, total, err := mb.menuRecordSearch(ctx, parents, strings.TrimSpace(q), nil, page, menuRecordPageSize)
	if err != nil {
		return VAlert(h.Text(err.Error())).Type("error")
	}
	rctx := *ctx
	rctx.R = refererRequest(ctx.R)
	var chosenID string
	if ch := mb.menuRecordChosen(&rctx, parents); ch.obj != nil {
		chosenID = ch.id.String()
	}
	href := mb.Info().ListingHref(parents...)
	var items h.HTMLComponents
	for _, obj := range objs {
		mid, ok := mb.menuRecordReadable(ctx.R, parents, obj)
		if !ok {
			continue
		}
		item := VListItem(VListItemTitle(h.Text(mb.RecordTitle(obj, ctx)))).
			PrependIcon(mb.menuIcon).
			Attr("@click", web.Plaid().URL(href).EventFunc(actions.DoListingAction).
				Query(ParamAction, ActionChooseMenuRecord).
				Query(paramMenuRecordID, mid.String()).Go())
		if mid.String() == chosenID {
			item.Active(true).Color("primary").AppendIcon("mdi-check")
		}
		items = append(items, item)
	}
	var list h.HTMLComponent = VList(items...).Density(DensityCompact)
	if len(items) == 0 {
		list = h.Div(h.Text(msgs.MenuRecordNone)).Class("text-medium-emphasis pa-4")
	}
	var pager h.HTMLComponent
	if pages := (int64(total) + menuRecordPageSize - 1) / menuRecordPageSize; pages > 1 {
		pager = VPagination().Length(int(pages)).ModelValue(int(page)).Density(DensityCompact).TotalVisible(7).
			Attr("@update:model-value", web.Plaid().URL(href).EventFunc(menuRecordListEvent).
				Query(paramMenuRecordQ, q).Query(paramMenuRecordPage, web.Var("$event")).Go())
	}
	return h.Div(list, pager).Attr("data-menu-record-list", mb.uriName)
}

// menuRecordSetScript is the script that keeps in the tab the record v (a
// value of vars.menuRecords, "" none) of the selector of mb, of the
// parents, and draws its item again (vars.menuRecordSet, CreateMenus).
func (mb *ModelBuilder) menuRecordSetScript(parents IDSlice, v string) string {
	return fmt.Sprintf("vars.menuRecordSet && vars.menuRecordSet(%s, %s, %s)",
		h.JSONString(mb.menuRecordKey(parents)), h.JSONString(v), h.JSONString(mb.Info().ListingHref(parents...)))
}

// menuRecordSamePage is, of the page path of a record of mb (of the
// parents), the same page of the record mid: its detail, a page of it
// (/relatorios); a page of a resource of it — of a record of the resource,
// or below it, the other record's —, the listing of the resource; "" when
// path is no page of a record of mb.
func (mb *ModelBuilder) menuRecordSamePage(path string, parents IDSlice, mid ID) string {
	listHref := mb.Info().ListingHref(parents...)
	rest, found := strings.CutPrefix(path, listHref+"/")
	if !found || rest == "" {
		return ""
	}
	seg, tail, _ := strings.Cut(rest, "/")
	to := mb.Info().DetailingHref(mid.String(), parents...)
	if tail == "" {
		return to
	}
	if old, err := mb.ParseRecordID(seg); err == nil {
		was := append(append(IDSlice{}, parents...), old)
		under := append(append(IDSlice{}, parents...), mid)
		for _, child := range mb.children {
			if h := child.Info().ListingHref(was...); path == h || strings.HasPrefix(path, h+"/") {
				return child.Info().ListingHref(under...)
			}
		}
	}
	return to + "/" + tail
}

// menuRecordUnderScript is the script that leaves the selectors of the
// children of the record mid (of the parents) with none chosen in the tab:
// a record chosen anew, the ones under it are chosen again.
func (mb *ModelBuilder) menuRecordUnderScript(parents IDSlice, mid ID) (script string) {
	under := append(append(IDSlice{}, parents...), mid)
	for _, child := range mb.children {
		if child.HasMenuRecordSelector() {
			script += fmt.Sprintf("vars.menuRecords && (vars.menuRecords[%s] = '');", h.JSONString(child.menuRecordKey(under)))
		}
	}
	return
}

// configureMenuRecordChoose installs the action ActionChooseMenuRecord on
// the model (its listing's), the event of its list, and the one that
// clears the choice.
func (mb *ModelBuilder) configureMenuRecordChoose() {
	mb.RegisterEventFunc(menuRecordClearEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		parents := ParentsModelID(ctx.R)
		// none, in the tab; on a page of the record (its resources), the
		// listing of the model
		script := mb.menuRecordSetScript(parents, "")
		listHref := mb.Info().ListingHref(parents...)
		if ref := refererRequest(ctx.R); strings.HasPrefix(ref.URL.Path, listHref+"/") {
			script += fmt.Sprintf(";plaid().vars(vars).pushStateURL(%s).go()", h.JSONString(listHref))
		}
		web.AppendRunScripts(&r, script)
		return
	})
	mb.RegisterEventFunc(menuRecordNodeEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		parents := ParentsModelID(ctx.R)
		// the record of the tab (the header), not of the page
		var ch menuRecordChoice
		if v, found := menuRecordsOf(ctx.R)[mb.menuRecordKey(parents)]; found && v != "" {
			if obj := mb.menuRecordOfValue(ctx, parents, v); obj != nil {
				if mid, ok := mb.menuRecordReadable(ctx.R, parents, obj); ok {
					ch = menuRecordChoice{obj: obj, id: mid}
				}
			}
		}
		node, _, opened := mb.menuRecordSelectorNode(ctx, parents, ch)
		r.Data = map[string]any{"node": node, "opened": opened}
		return
	})
	mb.RegisterEventFunc(menuRecordListEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		page, _ := strconv.ParseInt(ctx.R.FormValue(paramMenuRecordPage), 10, 64)
		r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
			Name: menuRecordListPortal,
			Body: mb.menuRecordList(ctx, ctx.R.FormValue(paramMenuRecordQ), page),
		})
		return
	})
	mb.Listing().Action(ActionChooseMenuRecord).
		Icon("mdi-swap-horizontal").
		SetI18nLabel(func(ctx context.Context) string {
			return fmt.Sprintf(MustGetMessages(ctx).MenuRecordChange, strings.ToLower(mb.TTitle(ctx)))
		}).
		SetDialogWidth("520").
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			msgs := MustGetMessages(ctx.Context())
			href := mb.Info().ListingHrefCtx(ctx)
			search := VTextField().
				Placeholder(msgs.MenuRecordSearch).
				PrependInnerIcon("mdi-magnify").
				Clearable(true).
				Density(DensityCompact).
				Variant(FieldVariantOutlined).
				HideDetails(true).
				Autofocus(true).
				Class("mb-2").
				Attr("@update:model-value", web.Plaid().URL(href).EventFunc(menuRecordListEvent).
					Query(paramMenuRecordQ, web.Var("$event || ''")).Query(paramMenuRecordPage, 1).Go()).
				Attr(":debounce", 300)
			return h.Div(search, web.Portal(mb.menuRecordList(ctx, "", 1)).Name(menuRecordListPortal)), nil
		}).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			msgs := MustGetMessages(ctx.Context())
			parents := ParentsModelID(ctx.R)
			obj := mb.menuRecordFind(ctx, parents, ctx.R.FormValue(paramMenuRecordID))
			if obj == nil {
				return errors.New(msgs.MenuRecordRequired)
			}
			mid, _ := mb.RecordID(obj)
			// kept in the tab — none chosen under it —, its item drawn again;
			// a page of the record before: the same of the record chosen
			script := mb.menuRecordUnderScript(parents, mid) + mb.menuRecordSetScript(parents, mb.menuRecordValue(parents, mid))
			if to := mb.menuRecordSamePage(refererRequest(ctx.R).URL.Path, parents, mid); to != "" {
				script += fmt.Sprintf(";plaid().vars(vars).pushStateURL(%s).go()", h.JSONString(to))
			}
			web.AppendRunScripts(ctx.Resp, script)
			return nil
		})
}
