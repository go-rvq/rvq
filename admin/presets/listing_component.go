package presets

import (
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
	filterComponentsBuild ListingFilterComponentsBuilder
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
	if b.filterDataFunc != nil {
		fd := b.filterDataFunc(ctx)
		fd.SetByQueryString(ctx.R.URL.RawQuery)
		filterBar = b.filterBar(ctx, msgr, fd, inDialog)
	}
	var searchBoxDefault h.HTMLComponent
	if b.mb.layoutConfig == nil || !b.mb.layoutConfig.SearchBoxInvisible {
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

		if b.mb.layoutConfig == nil || !b.mb.layoutConfig.SearchBoxInvisible {
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

		cb.TopBar = h.HTMLComponents{
			VDivider(),
			VToolbar(
				searchBoxDefault,
				filterBar,
			).Flat(true).Color("surface").AutoHeight(true).Class("pa-2"),
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

		return comp, nil
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
		h.Div(searchBoxDefault).Class("mb-2"),
	)

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

	return vue.UserComponent(web.Scope(
		comp,
	).Slot("{ locals }").LocalsInit(`{currEditingListItemID: ""}`),
	).ScopeVar("filterBarVisible", "{value: false}"), nil
}
