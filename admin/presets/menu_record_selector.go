package presets

import (
	"context"
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
// —, the resources of the one chosen under it (its children, Items, and
// pages of its detail, Pages, with the record in their path); then the item
// that lists the records.
//
// It is on for a model of no parent, a child of a singleton, or a child of
// a model that has it too (MenuRecordSelector of both): the records of its
// parents are then known — the ones chosen above it.
//
// The record chosen is the one of the page (…/<model>/{id}/…), else the one
// kept in the browser for its parents (a cookie), else the only one, else
// the default (DefaultConditions).

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
	// Conditions narrow the records that may be chosen (the active ones).
	Conditions SQLConditions
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

// menuRecordCookie is the cookie that keeps, in the browser, the record of
// mb chosen: "<parents>|<record>".
func (mb *ModelBuilder) menuRecordCookie() string {
	return "rvq_menu_" + strings.ReplaceAll(mb.id, "/", "_")
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

// menuRecordChosen is the record of mb chosen, of the parents: the one of
// the page, else the browser's, else the only one, else the default; nil,
// none.
func (mb *ModelBuilder) menuRecordChosen(ctx *web.EventContext, parents IDSlice) (any, ID) {
	ok := func(obj any) (any, ID) {
		if obj == nil {
			return nil, ID{}
		}
		mid, readable := mb.menuRecordReadable(ctx.R, parents, obj)
		if !readable {
			return nil, ID{}
		}
		return obj, mid
	}
	// the page's
	if rest, found := strings.CutPrefix(ctx.R.URL.Path, mb.Info().ListingHref(parents...)+"/"); found {
		id, _, _ := strings.Cut(rest, "/")
		if obj := mb.menuRecordFind(ctx, parents, id); obj != nil {
			return ok(obj)
		}
	}
	// the browser's
	if ck, err := ctx.R.Cookie(mb.menuRecordCookie()); err == nil {
		if p, id, found := strings.Cut(ck.Value, "|"); found && p == parents.String() {
			if obj := mb.menuRecordFind(ctx, parents, id); obj != nil {
				return ok(obj)
			}
		}
	}
	// the only one
	if objs, total, err := mb.menuRecordSearch(ctx, parents, "", nil, 1, 2); err == nil && total == 1 && len(objs) == 1 {
		if obj, mid := ok(objs[0]); obj != nil {
			return obj, mid
		}
	}
	// the default
	if c := mb.MenuRecordSelectorConfig(); len(c.DefaultConditions) > 0 {
		if objs, _, err := mb.menuRecordSearch(ctx, parents, "", c.DefaultConditions, 1, 1); err == nil && len(objs) > 0 {
			return ok(objs[0])
		}
	}
	return nil, ID{}
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
	key := mb.uriName + ":" + parents.String()
	title := mb.TTitle(ctx.Context())
	selector := &MenuNode{
		Title: title,
		Value: "g:menu-record:" + key,
		Props: map[string]any{"prependIcon": icon},
	}
	chip := &MenuChip{Title: msgs.MenuRecordSelect, Icon: "mdi-swap-horizontal", Color: "grey",
		Hint:  fmt.Sprintf(msgs.MenuRecordChange, strings.ToLower(title)),
		Event: actions.OpenActionDialog,
		URL:   listHref,
		Query: url.Values{
			ParamAction:  {ActionChooseMenuRecord},
			ParamOverlay: {string(actions.Dialog)},
		}.Encode()}
	selector.SetChip(chip)
	list := &MenuNode{Title: fmt.Sprintf(msgs.MenuRecordList, strings.ToLower(mb.TTitlePlural(ctx.Context()))),
		Value: listHref, Props: map[string]any{"href": listHref, "prependIcon": "mdi-format-list-bulleted"}}
	n = &MenuNode{
		Title:    mb.TTitlePlural(ctx.Context()),
		Value:    "g:menu-records:" + key,
		Props:    map[string]any{"prependIcon": icon},
		Children: []*MenuNode{selector, list},
	}

	if obj, mid := mb.menuRecordChosen(ctx, parents); obj != nil {
		chip.Title, chip.Color = mb.RecordTitle(obj, ctx), "primary"
		opened = append(opened, n.Value, selector.Value)
		under := append(append(IDSlice{}, parents...), mid)
		for _, child := range mb.menuRecordItems() {
			if child.notInMenu || child.permissioner.Lister(ctx.R, under...).Denied() {
				continue
			}
			cn, ca, co := child.MenuNodeOf(ctx, under)
			selector.Children = append(selector.Children, cn)
			opened = append(opened, co...)
			if ca != "" {
				active = ca
			}
		}
		for _, p := range mb.MenuRecordSelectorConfig().Pages {
			page := mb.Detailing().PagesRegistrator().GetHttpPage(p)
			// the record readable (chosen); the page's own permission, when
			// it has one
			if page == nil || (page.GetBaseVerifier() != nil && page.Verifier(ctx.R).Denied()) {
				continue
			}
			href := mb.Info().DetailingHref(mid.String(), parents...) + p
			pi := page.GetMenuIcon()
			if pi == "" {
				pi = "mdi-circle-small"
			}
			selector.Children = append(selector.Children, &MenuNode{Title: page.TTitle(ctx.Context()), Value: href,
				Props: map[string]any{"href": href, "prependIcon": pi}})
			if isActive(href) {
				active = href
			}
		}
	}
	if active == "" && isActive(listHref) {
		active = listHref
	}
	if active != "" {
		opened = append(opened, n.Value)
	}
	return
}

// ActionChooseMenuRecord is the listing action of a model of a record chosen
// in the menu that chooses it (its dialog).
const ActionChooseMenuRecord = "choose_menu_record"

const (
	menuRecordListEvent  = "presets_menuRecord_list"
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
	if obj, mid := mb.menuRecordChosen(&rctx, parents); obj != nil {
		chosenID = mid.String()
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

// configureMenuRecordChoose installs the action ActionChooseMenuRecord on
// the model (its listing's), and the event of its list.
func (mb *ModelBuilder) configureMenuRecordChoose() {
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
			path := mb.p.GetURIPrefix()
			if path == "" {
				path = "/"
			}
			http.SetCookie(ctx.W, &http.Cookie{Name: mb.menuRecordCookie(), Value: parents.String() + "|" + mid.String(),
				Path: path, MaxAge: 365 * 24 * 3600, HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: ctx.R.TLS != nil})
			// the first of its resources, else its page
			under := append(append(IDSlice{}, parents...), mid)
			ctx.Resp.RedirectURL = mb.Info().DetailingHref(mid.String(), parents...)
			for _, child := range mb.menuRecordItems() {
				if !child.permissioner.Lister(ctx.R, under...).Denied() {
					ctx.Resp.RedirectURL = child.Info().ListingHref(under...)
					break
				}
			}
			return nil
		})
}
