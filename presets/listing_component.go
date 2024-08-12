package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	. "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type ListingComponentBuilder struct {
	b         *ListingBuilder
	portals   *ListingPortals
	selection bool
}

func NewListingComponentBuilder(b *ListingBuilder, portals *ListingPortals) *ListingComponentBuilder {
	return &ListingComponentBuilder{b: b, portals: portals}
}

func (b *ListingBuilder) listingComponentBuilder(
	portals *ListingPortals,
) *ListingComponentBuilder {
	return NewListingComponentBuilder(b, portals)
}

func (b *ListingBuilder) listingComponentBuilderCtx(
	ctx *web.EventContext,
) *ListingComponentBuilder {
	ctx.R.Form.Set(ParamPortalID, GetOrNewPortalID(ctx.R))
	return b.listingComponentBuilder(b.Portals(ctx.R.FormValue(ParamPortalID)))
}

func (b *ListingBuilder) listingComponent(
	ctx *web.EventContext,
) (h.HTMLComponent, error) {
	return b.listingComponentBuilderCtx(ctx).Build(ctx)
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

func (lcb *ListingComponentBuilder) Build(ctx *web.EventContext) (comp h.HTMLComponent, err error) {
	b := lcb.b
	inDialog := IsInDialog(ctx)
	msgr := MustGetMessages(ctx.R)
	portalID := GetPortalID(ctx.R)

	var tabsAndActionsBar h.HTMLComponent
	{
		filterTabs := b.filterTabs(lcb.portals, ctx, inDialog)

		actionsComponent := lcb.actionsComponent(msgr, ctx, inDialog)
		// if v := ; v != nil {
		//	actionsComponent = append(actionsComponent, v)
		// }
		// || len(actionsComponent) > 0
		if filterTabs != nil {
			tabsAndActionsBar = filterTabs
		}
		ctx.WithContextValue(CtxActionsComponent, actionsComponent)
	}

	var filterBar h.HTMLComponent
	if b.filterDataFunc != nil {
		fd := b.filterDataFunc(ctx)
		fd.SetByQueryString(ctx.R.URL.RawQuery)
		filterBar = b.filterBar(ctx, msgr, fd, inDialog)
	}
	searchBoxDefault := VResponsive(
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
					VBtn("").
						Theme("dark").
						// Size(SizeSmall).
						Attr("@click", web.Plaid().
							PushState(true).
							Go()).
						Icon(true).
						Density("comfortable").
						Children(VIcon("mdi-reload")),
				).Attr("style", "flex-grow: 0;padding-left:0"),
			),
		).VSlot("{ locals }").Init(`{isFocus: false}`),
	)

	var (
		dataTable          h.HTMLComponent
		dataTableAdditions h.HTMLComponent
	)

	if dataTable, dataTableAdditions, err = lcb.getTableComponents(ctx); err != nil {
		return
	}

	var footerCardAction h.HTMLComponent

	if len(b.footerActions) > 0 {
		var footerActions []h.HTMLComponent
		footerActions = append(footerActions, VSpacer())
		for _, action := range b.footerActions {
			footerActions = append(footerActions, action.buttonCompFunc(ctx))
		}
		footerCardAction = VCardActions(footerActions...)
	}

	cb := &ContentComponentBuilder{
		Overlay: &ContentComponentBuilderOverlay{
			Mode: actions.Dialog,
		},
		TopRightActions: GetActionsComponent(ctx),
		Scope:           web.Scope().VSlot("{ locals, closer, form }").Init(`{currEditingListItemID: ""}`),
	}

	if tabsAndActionsBar != nil {
		cb.Tabs = append(cb.Tabs, tabsAndActionsBar)
	}

	if inDialog {
		cb.Title = b.title
		if cb.Title == "" {
			cb.Title = msgr.ListingObjectTitle(i18n.T(ctx.R, ModelsI18nModuleKey, b.mb.pluralLabel))
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
									MergeQuery(true).
									Query(ParamPortalID, portalID).
									EventFunc(actions.UpdateListingDialog).
									Go()).
								Attr("@click:clear", web.Plaid().
									URL(ctx.R.RequestURI).
									Query("keyword", "").
									MergeQuery(true).
									Query(ParamPortalID, portalID).
									EventFunc(actions.UpdateListingDialog).
									Go()).
								Class("ma-0 pa-0"),
						),
						VCol(
							VBtn("").
								Theme("dark").
								// Size(SizeSmall).
								Attr("@click", web.Plaid().
									URL(ctx.R.RequestURI).
									MergeQuery(true).
									Query(ParamPortalID, portalID).
									EventFunc(actions.UpdateListingDialog).
									Go()).
								Icon(true).
								Density("comfortable").
								Children(VIcon("mdi-reload")),
						).Attr("style", "flex-grow: 0;padding-left:0"),
					),
				).VSlot("{ locals }").Init(`{isFocus: false}`),
			).Width(100)
		}

		cb.Body = VCard(
			VToolbar(
				searchBoxDefault,
				filterBar,
			).Flat(true).Color("surface").AutoHeight(true).Class("pa-2"),
			VCardText(
				web.Portal().Name(lcb.portals.Temp()),
				web.Portal(dataTable).Name(lcb.portals.DataTable()),
				web.Portal(dataTableAdditions).Name(lcb.portals.DataTableAdditions()),
			).Class("pa-2"),
		)

		return cb.BuildOverlay(), nil
	}
	return web.Scope(
		VLayout(
			VMain(
				tabsAndActionsBar,
				VCard(
					VToolbar(
						searchBoxDefault,
						filterBar,
					).Flat(true).Color("surface").AutoHeight(true).Class("pa-2"),
					VCardText(
						web.Portal().Name(lcb.portals.Temp()),
						web.Portal(dataTable).Name(lcb.portals.DataTable()),
						web.Portal(dataTableAdditions).Name(lcb.portals.DataTableAdditions()),
					).Class("pa-2"),
					footerCardAction,
				),
			),
		),
	).VSlot("{ locals }").Init(`{currEditingListItemID: ""}`), nil
}
