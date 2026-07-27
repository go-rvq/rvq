package presets

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/datafield"
	"github.com/go-rvq/rvq/web/vue"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

type ListingTableBuilder func(lcb *ListingComponentBuilder, ctx *web.EventContext, sr *SearchResult, overlayMode actions.OverlayMode) (
	comp h.HTMLComponent,
	err error,
)

type ListingPreBuild func(lcb *ListingComponentBuilder, ctx *web.EventContext) (err error)

type ListingFilterComponentsBuilder func(lcb *ListingComponentBuilder, ctx *web.EventContext) h.HTMLComponent

type ListingComponentBuilder struct {
	b                     *ListingBuilder
	portals               *ListingPortals
	selection             bool
	configureComponent    func(cb *ContentComponentBuilder)
	tableBuilder          ListingTableBuilder
	preBuild              ListingPreBuild
	componentWrap         func(ctx *web.EventContext, comp h.HTMLComponent) h.HTMLComponent
	headWrap              func(cell h.HTMLComponent, rowspan, colspan int, field string, dataTableID string, ctx *web.EventContext) h.HTMLComponent
	cellWrap              func(cell h.HTMLComponent, field string, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent
	filterComponentsBuild ListingFilterComponentsBuilder
	SearchbarDisabled     bool
	FilterDisabled        bool
	datafield.DataField[*ListingComponentBuilder]
}

func (lcb *ListingComponentBuilder) FilterComponentsBuild() ListingFilterComponentsBuilder {
	return lcb.filterComponentsBuild
}

func (lcb *ListingComponentBuilder) SetFilterComponentsBuild(filterComponentsBuild ListingFilterComponentsBuilder) *ListingComponentBuilder {
	lcb.filterComponentsBuild = filterComponentsBuild
	return lcb
}

func (lcb *ListingComponentBuilder) ComponentWrap(f func(ctx *web.EventContext, comp h.HTMLComponent) h.HTMLComponent) *ListingComponentBuilder {
	lcb.componentWrap = f
	return lcb
}

func (lcb *ListingComponentBuilder) CellWrap(f func(cell h.HTMLComponent, field string, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent) *ListingComponentBuilder {
	lcb.cellWrap = f
	return lcb
}

func (lcb *ListingComponentBuilder) HeadWrap(f func(cell h.HTMLComponent, rowspan, colspan int, field string, dataTableID string, ctx *web.EventContext) h.HTMLComponent) *ListingComponentBuilder {
	lcb.headWrap = f
	return lcb
}

func NewListingComponentBuilder(b *ListingBuilder, portals *ListingPortals) *ListingComponentBuilder {
	return datafield.New(&ListingComponentBuilder{b: b, portals: portals})
}

func (b *ListingBuilder) listingComponentBuilder(
	portals *ListingPortals,
) *ListingComponentBuilder {
	lcb := NewListingComponentBuilder(b, portals)
	if b.configureComponent != nil {
		b.configureComponent(lcb)
	}
	return lcb
}

func (b *ListingBuilder) ListingComponentBuilderCtx(
	ctx *web.EventContext,
) *ListingComponentBuilder {
	ctx.R.Form.Set(ParamPortalID, GetOrNewPortalID(ctx.R))
	return b.listingComponentBuilder(b.Portals(ctx.R.FormValue(ParamPortalID)))
}

func (b *ListingBuilder) listingComponent(
	ctx *web.EventContext,
) (h.HTMLComponent, error) {
	return b.ListingComponentBuilderCtx(ctx).
		Build(ctx)
}

func (lcb *ListingComponentBuilder) Portals() *ListingPortals {
	return lcb.portals
}

func (lcb *ListingComponentBuilder) SetPortals(portals *ListingPortals) {
	lcb.portals = portals
}

func (lcb *ListingComponentBuilder) Selection() bool {
	return lcb.selection
}

func (lcb *ListingComponentBuilder) SetSelection(selection bool) *ListingComponentBuilder {
	lcb.selection = selection
	return lcb
}

func (lcb *ListingComponentBuilder) PreBuild() ListingPreBuild {
	return lcb.preBuild
}

func (lcb *ListingComponentBuilder) SetPreBuild(preBuild ListingPreBuild) *ListingComponentBuilder {
	lcb.preBuild = preBuild
	return lcb
}

func (lcb *ListingComponentBuilder) ConfigureComponent() func(cb *ContentComponentBuilder) {
	return lcb.configureComponent
}

func (lcb *ListingComponentBuilder) SetConfigureComponent(configureComponent func(cb *ContentComponentBuilder)) *ListingComponentBuilder {
	lcb.configureComponent = configureComponent
	return lcb
}

func (lcb *ListingComponentBuilder) TableBuilder() ListingTableBuilder {
	return lcb.tableBuilder
}

func (lcb *ListingComponentBuilder) SetTableBuilder(tableBuilder ListingTableBuilder) *ListingComponentBuilder {
	lcb.tableBuilder = tableBuilder
	return lcb
}

func (lcb *ListingComponentBuilder) Build(ctx *web.EventContext) (comp h.HTMLComponent, err error) {
	b := lcb.b
	inDialog := IsInDialog(ctx)
	msgr := MustGetMessages(ctx.Context())
	portalID := GetPortalID(ctx.R)

	// publish the per-item hosts before anything renders a row, so the rows open
	// the listing's shared overlays instead of carrying their own plaid.
	itemHosts := lcb.itemFormHosts(ctx)

	filterTabs := b.filterTabs(lcb.portals, ctx, inDialog)

	actionsComponent := lcb.actionsComponent(msgr, ctx, inDialog)
	// if v := ; v != nil {
	//	actionsComponent = append(actionsComponent, v)
	// }
	// || len(actionsComponent) > 0

	if !inDialog {
		WithActionsComponent(ctx, actionsComponent)
	}

	var filterBar h.HTMLComponent
	if !lcb.FilterDisabled && b.filterDataFunc != nil {
		fd := b.filterDataFunc(ctx)
		fd.SetByQueryString(ctx.R.URL.RawQuery)
		filterBar = b.filterBar(ctx, msgr, fd, inDialog)
	}

	var searchBoxDefault h.HTMLComponent
	if !lcb.SearchbarDisabled && (b.mb.layoutConfig == nil || !b.mb.layoutConfig.SearchBoxInvisible) {
		searchBoxDefault = VResponsive(
			web.Scope(
				VRow(
					VSpacer(),
					VCol(
						VTextField(
							web.Slot(VIcon("mdi-magnify")).Name("append-inner"),
						).Density(DensityCompact).
							Variant(FieldVariantOutlined).
							Label(msgr.Search).
							Flat(true).
							Clearable(true).
							HideDetails(true).
							SingleLine(true).
							ModelValue(ctx.R.URL.Query().Get("keyword")).
							Attr("@keyup.enter", web.Plaid().
								ClearMergeQuery("page").
								Query("keyword", web.Var("[$event.target.value]")).
								MergeQuery(true).
								PushState(true).
								Go()).
							Attr("@click:clear", web.Plaid().
								Query("keyword", "").
								PushState(true).
								Go()),
					),
					VCol(
						VLayout(
							VBtn("").
								// Size(SizeSmall).
								Attr("@click", web.Plaid().
									PushState(true).
									Go()).
								Icon(true).
								Variant(VariantFlat).
								Density(DensityCompact).
								Children(VIcon("mdi-reload")),
							h.If(filterBar != nil,
								VBtn("").
									Attr("@click", "filterBarVisible.value = !filterBarVisible.value").
									Attr(":color", `filterBarVisible.value ? "primary": ""`).
									Icon(true).
									Variant(VariantFlat).
									Density(DensityCompact).
									Children(VIcon("mdi-filter")),
							),
						),
					).Class("ps-0"),
					VSpacer(),
				),
			).Slot("{ locals }").LocalsInit(`{isFocus: false}`),
		)
	}

	var (
		dataTable          h.HTMLComponent
		dataTableAdditions h.HTMLComponent
	)

	if lcb.preBuild != nil {
		if err = lcb.preBuild(lcb, ctx); err != nil {
			return
		}
	}

	if dataTable, dataTableAdditions, err = lcb.GetTableComponents(ctx); err != nil {
		return
	}

	var footerActions h.HTMLComponents

	if len(b.footerActions) > 0 {
		footerActions = append(footerActions, VSpacer())
		for _, action := range b.footerActions {
			footerActions = append(footerActions, action.buttonCompFunc(ctx))
		}
	}

	cb := &ContentComponentBuilder{
		Context: ctx,
		Overlay: &ContentComponentBuilderOverlay{
			Mode: actions.Dialog,
		},
		TopRightActions: h.HTMLComponents{actionsComponent},
		Scope:           web.Scope().Slot("{ locals, closer, form }").LocalsInit(`{currEditingListItemID: ""}`),
	}

	if filterTabs != nil {
		cb.PreBody = append(cb.PreBody, filterTabs)
	}

	if inDialog {
		cb.Title = b.title
		if cb.Title == "" {
			cb.Title = msgr.ListingObjectTitle(b.mb.TTitlePlural(ctx.Context()))
		}

		if len(footerActions) > 0 {
			cb.BottomActions = append(cb.BottomActions, VCardActions(footerActions...))
		}

		if !lcb.SearchbarDisabled && (b.mb.layoutConfig == nil || !b.mb.layoutConfig.SearchBoxInvisible) {
			searchBoxDefault = VResponsive(
				web.Scope(
					VRow(
						VCol(
							VTextField(
								web.Slot(VIcon("mdi-magnify")).Name("append-inner"),
							).Density(DensityCompact).
								Variant(FieldVariantOutlined).
								Label(msgr.Search).
								Flat(true).
								Clearable(true).
								HideDetails(true).
								SingleLine(true).
								ModelValue(ctx.R.URL.Query().Get("keyword")).
								Attr("@keyup.enter", web.Plaid().
									URL(ctx.R.RequestURI).
									Query("keyword", web.Var("[$event.target.value]")).
									Scope(vue.Var("{presetsListing: presetsListing}")).
									MergeQuery(true).
									Query(ParamPortalID, portalID).
									EventFunc(actions.UpdateListingDialog).
									Go()).
								Attr("@click:clear", web.Plaid().
									URL(ctx.R.RequestURI).
									Query("keyword", "").
									Scope(vue.Var("{presetsListing: presetsListing}")).
									MergeQuery(true).
									Query(ParamPortalID, portalID).
									EventFunc(actions.UpdateListingDialog).
									Go()).
								Class("ma-0 pa-0"),
						),
						VCol(
							VBtn("").
								// Size(SizeSmall).
								Attr("@click", web.Plaid().
									URL(ctx.R.RequestURI).
									MergeQuery(true).
									Query(ParamPortalID, portalID).
									EventFunc(actions.UpdateListingDialog).
									Go()).
								Icon(true).
								Variant(VariantFlat).
								Density(DensityCompact).
								Children(VIcon("mdi-reload")),
						).Attr("style", "flex-grow: 0;padding-left:0"),
					),
				).Slot("{ locals }").LocalsInit(`{isFocus: false}`),
			).Width(100)
		}

		if searchBoxDefault != nil || filterBar != nil {
			cb.TopBar = h.HTMLComponents{
				VDivider(),
				VToolbar(
					searchBoxDefault,
					filterBar,
				).Flat(true).Color("surface").AutoHeight(true).Class("pa-2"),
			}
		}

		cb.BottomBar = VCardActions(web.Portal(dataTableAdditions).Name(lcb.portals.DataTableAdditions()))

		cb.Body = h.HTMLComponents{
			web.Portal().Name(lcb.portals.Temp()),
			web.Portal(dataTable).Name(lcb.portals.DataTable()),
		}

		if lcb.configureComponent != nil {
			lcb.configureComponent(cb)
		}

		comp := cb.BuildOverlay()

		if lcb.componentWrap != nil {
			comp = lcb.componentWrap(ctx, comp)
		}

		comp = vue.UserComponent(comp).ScopeVar("filterBarVisible", "{value: false}")

		return lcb.hostForms(ctx, comp, itemHosts), nil
	}

	if lcb.configureComponent != nil {
		lcb.configureComponent(cb)
	}

	var preContent h.HTMLComponent
	if lcb.filterComponentsBuild != nil {
		preContent = lcb.filterComponentsBuild(lcb, ctx)
	}

	cb.PreBody = append(cb.PreBody,
		preContent,
	)

	if searchBoxDefault != nil {
		cb.PreBody = append(cb.PreBody, h.Div(searchBoxDefault).Class("mb-2"))
	}

	if filterBar != nil {
		cb.PreBody = append(cb.PreBody,
			h.Div(
				VDivider(),
				VContainer(VRow(filterBar.(*vx.VXFilterBuilder).Attr("@data", `data => filterBarVisible.value = true`))).Style("background-color:#fafafa"),
				VDivider(),
				h.Div().Class("mb-2"),
			).Attr(":style", `filterBarVisible.value ? "" : "display:none"`),
		)
	}

	cb.PreBody = append(cb.PreBody, VDivider().Class("mb-2"))

	cb.Body = h.HTMLComponents{
		web.Portal().Name(lcb.portals.Temp()),
		web.Portal(dataTable).Name(lcb.portals.DataTable()),
		web.Portal(dataTableAdditions).Name(lcb.portals.DataTableAdditions()),
	}

	cb.BottomActions = append(cb.BottomActions, footerActions...)

	comp = cb.BuildPage()

	if lcb.componentWrap != nil {
		comp = lcb.componentWrap(ctx, comp)
	}

	comp = vue.UserComponent(web.Scope(
		comp,
	).Slot("{ locals }").LocalsInit(`{currEditingListItemID: ""}`),
	).ScopeVar("filterBarVisible", "{value: false}")

	return lcb.hostForms(ctx, comp, itemHosts), nil
}

// itemFormHosts builds the listing's per-item hosts (edit and detail) and
// publishes them on the context, so every row rendered afterwards can open them
// with `<scope>.id = "<id>"; <scope>.show = true` instead of carrying its own
// overlay plaid. There is ONE host for the whole listing, keyed by the id var.
func (lcb *ListingComponentBuilder) itemFormHosts(ctx *web.EventContext) *ItemFormHosts {
	var (
		b       = lcb.b
		overlay = OverlayMode(ctx).Up().String()
		reload  = b.reloadURI(ctx)
		hosts   = &ItemFormHosts{}
		// this listing's own address, and the addresses of the PAGES its
		// overlays stand for. Which record it is, is only known when a row opens
		// the host, so the address is a function of the closer's `id`.
		listingHref = b.mb.Info().ListingHrefCtx(ctx)
		itemURL     = func(suffix string) string {
			s := fmt.Sprintf(`(closer) => %s + "/" + closer.id`, h.JSONString(listingHref))
			if suffix != "" {
				s += " + " + h.JSONString(suffix)
			}
			return s
		}
		host = func(scope, event string) *FormHostBuilder {
			portal := ctx.UID()
			// whatever this host opens — a detail, an edit form — a successful save
			// inside it changes a row of THIS listing, so the listing reloads.
			hb := FormHost(scope, portal, nil).Var("id", "null").OnSave(reload)
			// the id comes from the host's own state, whose reference depends on
			// where it lives (`vars` on a page, a slot variable elsewhere).
			hb.load = web.Plaid().
				URL(listingHref).
				EventFunc(event).
				Query(ParamID, web.Var(hb.ScopeVarExpr("id"))).
				Query(ParamTargetPortal, portal).
				Query(ParamOverlay, overlay)
			return hb
		}
	)

	if !b.mb.editingDisabled {
		hosts.Edit = host(ListingItemEditScope, actions.Edit).URLExpr(itemURL("/edit"))
	}
	if b.mb.hasDetailing && !b.mb.detailingDisabled {
		hosts.Detail = host(ListingItemDetailScope, actions.Detailing).URLExpr(itemURL(""))
	}

	newPortal := ctx.UID()
	hosts.New = FormHost(ListingNewScope, newPortal,
		web.Plaid().
			URL(ctx.R.RequestURI).
			EventFunc(actions.New).
			Query(ParamTargetPortal, newPortal).
			Query(ParamOverlay, overlay)).
		OnSave(reload)

	if !b.mb.editingDisabled {
		hosts.New.URL(listingHref + "/new")
	}

	WithItemFormHosts(ctx, hosts)
	return hosts
}

// hostForms wraps the listing in the hosts its ROWS open: the per-item edit and
// detail overlays. Each host owns its overlay's closer, so turning its variable
// off destroys the overlay and turning it on loads a fresh one — one request.
//
// The create host joins them only when the New button stays inline (in a
// dialog); on a page that button is rendered by the layout into the app bar, so
// the host travels WITH the button instead (see actionsComponent).
//
// They are declared by a SINGLE component: each host renders a
// `<template v-slot>`, and one nested in another is not compiled by the runtime
// template parser — the listing inside would render inert.
func (lcb *ListingComponentBuilder) hostForms(ctx *web.EventContext, comp h.HTMLComponent, hosts *ItemFormHosts) h.HTMLComponent {
	all := []*FormHostBuilder{hosts.Detail, hosts.Edit}
	if IsInDialog(ctx) {
		all = append(all, hosts.New)
	}
	return FormHosts(h.HTMLComponents{comp}, all...)
}
