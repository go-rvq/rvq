package presets

import (
	"fmt"
	"net/url"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	. "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
)

func (b *ListingBuilder) DataTableDensity() string {
	return b.dataTableDensity
}

func (b *ListingBuilder) SetDataTableDensity(dataTableDensity string) *ListingBuilder {
	b.dataTableDensity = dataTableDensity
	return b
}

func (b *ListingBuilder) cellComponentFunc(obj interface{}, f *FieldBuilder) vx.CellComponentFunc {
	return func(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent {
		return f.compFunc(obj, b.mb.GetComponentFuncField(obj, f), ctx)
	}
}

func (lcb *ListingComponentBuilder) getTableComponents(ctx *web.EventContext) (
	dataTable h.HTMLComponent,
	datatableAdditions h.HTMLComponent,
	err error,
) {
	var (
		b    = lcb.b
		msgr = MustGetMessages(ctx.R)
		qs   = ctx.R.URL.Query()

		overlayMode = actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
		inDialog    = overlayMode.IsDialog()
		sr          SearchResult
	)

	if sr, err = lcb.b.search(ctx); err != nil {
		return
	}

	haveCheckboxes := !lcb.selection && len(b.bulkActions) > 0

	tempPortal := lcb.portals.Temp()

	cellWrapperFunc := func(cell h.MutableAttrHTMLComponent, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent {
		return cell
	}

	if b.cellWrapperFunc != nil {
		cellWrapperFunc = b.cellWrapperFunc
	} else if lcb.selection {
		selectedEventName := ctx.R.URL.Query().Get(SelectedEventParamName)
		selectedEventConfig := ctx.R.URL.Query().Get(SelectedEventConfigParamName)
		portalID := ctx.R.FormValue(ParamPortalID)
		cellWrapperFunc = func(cell h.MutableAttrHTMLComponent, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent {
			onclick := web.Plaid().
				EventFunc(selectedEventName).
				Query(ParamSelectedID, id).
				Query(SelectedEventConfigParamName, selectedEventConfig)
			cell.SetAttr("@click.self",
				onclick.Go()+`;vars.presetsListingDialog`+portalID+` = false`)
			return cell
		}
	} else if !CallOkHandled(b.mb.editionDisabled, ctx) {
		reloadCb := b.reloadCallback(ctx).Encode()

		cellWrapperFunc = func(cell h.MutableAttrHTMLComponent, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent {
			onclick := web.Plaid().
				Query(ParamID, id).
				URL(ctx.R.RequestURI).
				Query(ParamPostChangeCallback, reloadCb)

			if b.mb.hasDetailing {
				onclick.EventFunc(actions.Detailing)
			} else {
				onclick.EventFunc(actions.Edit)
			}

			onclick.Query(ParamOverlay, overlayMode.Up())

			if overlayMode.Overlayed() {
				onclick.Query(ParamTargetPortal, tempPortal)
			}

			cell.SetAttr("@click.self",
				onclick.Go()+fmt.Sprintf(`; locals.currEditingListItemID="%s-%s"`, dataTableID, id))
			cell.SetAttr("@click.middle",
				fmt.Sprintf("(e) => e.view.window.open(%q, '_blank')", b.mb.Info().DetailingHrefCtx(ctx, id)))
			return cell
		}
	}

	displayFields := b.fields
	var selectColumnsBtn h.HTMLComponent

	if b.selectableColumns {
		selectColumnsBtn, displayFields = b.selectColumnsBtn(ctx.R.URL, ctx, inDialog)
	} else {
		displayFields = displayFields.Renderable()
	}

	recordPortals := make(map[int][]string)

	sDataTable := vx.DataTable(sr.objs).
		SetDensity(b.dataTableDensity).
		CellWrapperFunc(cellWrapperFunc).
		HeadCellWrapperFunc(func(cell h.MutableAttrHTMLComponent, field string, title string, ctx *web.EventContext) h.HTMLComponent {
			if _, ok := sr.orderableFieldMap[field]; ok {
				var orderBy string
				var orderByIdx int
				for i, ob := range sr.orderBys {
					if ob.FieldName == field {
						orderBy = ob.OrderBy
						orderByIdx = i + 1
						break
					}
				}
				th := h.Th("").Style("cursor: pointer; white-space: nowrap;").
					Children(
						h.Span(title).
							Style("text-decoration: underline;"),
						h.If(orderBy == "ASC",
							VIcon("arrow_drop_up").Size(SizeSmall),
							h.Span(fmt.Sprint(orderByIdx)),
						).ElseIf(orderBy == "DESC",
							VIcon("arrow_drop_down").Size(SizeSmall),
							h.Span(fmt.Sprint(orderByIdx)),
						).Else(
							// take up place
							h.Span("").Style("visibility: hidden;").Children(
								VIcon("arrow_drop_down").Size(SizeSmall),
								h.Span(fmt.Sprint(orderByIdx)),
							),
						),
					)
				qs.Del(web.ExecuteEventPatam)
				newQuery := newQueryWithFieldToggleOrderBy(qs, field)
				onclick := web.Plaid().
					Queries(newQuery)
				if inDialog {
					onclick.URL(ctx.R.RequestURI).
						EventFunc(actions.UpdateListingDialog)
				} else {
					onclick.PushState(true)
				}
				th.Attr("@click", onclick.Go())

				cell = th
			}

			return cell
		}).
		RowWrapperFunc(func(row h.MutableAttrHTMLComponent, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent {
			row.SetAttr(":class", fmt.Sprintf(`{"vx-list-item--active primary--text": vars.presetsRightDrawer && locals.currEditingListItemID==="%s-%s"}`, dataTableID, id))

			return row
		}).
		RowMenuItemFuncs(b.RowMenu().listingItemFuncs(ctx).ToRowMenuItemFuncs(tempPortal, func(rctx *RecordMenuItemContext, name string) string {
			portalName := tempPortal + "--" + rctx.ID
			recordPortals[rctx.RecordIndex] = append(recordPortals[rctx.RecordIndex], portalName)
			return portalName
		})...).
		RowStarter(func(rowIndex int, id string, obj interface{}, dataTableID string, ctx *web.EventContext) vx.RowEndFunc {
			return func(row *h.HTMLTagBuilder) {
				for _, portal := range recordPortals[rowIndex] {
					row.AppendChildren(web.Portal().Name(portal))
				}
			}
		}).
		Selectable(haveCheckboxes).
		SelectionParamName(ParamSelectedIds).
		SelectedCountLabel(msgr.ListingSelectedCountNotice).
		SelectableColumnsBtn(selectColumnsBtn).
		ClearSelectionLabel(msgr.ListingClearSelection)

	targetPortal := ctx.R.FormValue(ParamTargetPortal)

	if inDialog {
		sDataTable.OnSelectAllFunc(func(idsOfPage []string, ctx *web.EventContext) string {
			return web.Plaid().
				URL(ctx.R.RequestURI).
				EventFunc(actions.UpdateListingDialog).
				ValidQuery(ParamTargetPortal, targetPortal).
				Query(ParamSelectedIds,
					web.Var(fmt.Sprintf(`{value: %s, add: $event, remove: !$event}`, h.JSONString(idsOfPage))),
				).
				MergeQuery(true).
				Go()
		})
		sDataTable.OnSelectFunc(func(id string, ctx *web.EventContext) string {
			return web.Plaid().
				URL(ctx.R.RequestURI).
				EventFunc(actions.UpdateListingDialog).
				ValidQuery(ParamTargetPortal, targetPortal).
				Query(ParamSelectedIds,
					web.Var(fmt.Sprintf(`{value: %s, add: $event, remove: !$event}`, h.JSONString(id))),
				).
				MergeQuery(true).
				Go()
		})
		sDataTable.OnClearSelectionFunc(func(ctx *web.EventContext) string {
			return web.Plaid().
				URL(ctx.R.RequestURI).
				EventFunc(actions.UpdateListingDialog).
				MergeQuery(true).
				ValidQuery(ParamTargetPortal, targetPortal).
				Query(ParamSelectedIds, "").
				Go()
		})
	}

	dataTable = sDataTable
	mode := FieldModeStack{LIST}

	for _, f := range displayFields {
		fctx := &FieldContext{
			EventContext: ctx,
			Field:        f,
			Mode:         mode,
			Name:         f.name,
			Label:        i18n.PT(ctx.R, ModelsI18nModuleKey, b.mb.label, b.mb.getLabel(f.NameLabel)),
		}

		if f.IsEnabled(fctx) {
			if b.mb.Info().Verifier().Do(PermList).SnakeOn("f_"+f.name).WithReq(ctx.R).IsAllowed() != nil {
				continue
			}
			f = b.GetFieldOrDefault(f.name) // fill in empty compFunc and setter func with default
			dataTable.(*vx.DataTableBuilder).Column(f.name).
				Title(fctx.Label).
				CellComponentFunc(b.cellComponentFunc(nil, f))
		}
	}

	if b.disablePagination {
		// if disable pagination, we don't need to add
		// the pagination component and the no-record message to page.
		return
	}
	if sr.totalCount > 0 {
		tpb := vx.VXTablePagination().
			Total(int64(sr.totalCount)).
			CurrPage(sr.Page).
			PerPage(sr.PerPage).
			CustomPerPages([]int64{b.perPage}).
			PerPageText(msgr.PaginationRowsPerPage)

		if inDialog {
			tpb.OnSelectPerPage(web.Plaid().
				URL(ctx.R.RequestURI).
				Query("per_page", web.Var("[$event]")).
				MergeQuery(true).
				ValidQuery(ParamTargetPortal, targetPortal).
				EventFunc(actions.UpdateListingDialog).
				Go())
			tpb.OnPrevPage(web.Plaid().
				URL(ctx.R.RequestURI).
				Query("page", sr.Page-1).
				MergeQuery(true).
				ValidQuery(ParamTargetPortal, targetPortal).
				EventFunc(actions.UpdateListingDialog).
				Go())
			tpb.OnNextPage(web.Plaid().
				URL(ctx.R.RequestURI).
				Query("page", sr.Page+1).
				MergeQuery(true).
				ValidQuery(ParamTargetPortal, targetPortal).
				EventFunc(actions.UpdateListingDialog).
				Go())
		}

		datatableAdditions = h.Div(tpb).Class("mt-2")
	} else {
		datatableAdditions = h.Div(h.Text(msgr.ListingNoRecordToShow)).Class("mt-10 text-center grey--text text--darken-2")
	}

	return
}

func (b *ListingBuilder) reloadList(ctx *web.EventContext) (r web.EventResponse, err error) {
	var (
		lcb                = b.listingComponentBuilderCtx(ctx)
		dataTable          h.HTMLComponent
		dataTableAdditions h.HTMLComponent
	)

	if dataTable, dataTableAdditions, err = lcb.getTableComponents(ctx); err != nil {
		return
	}

	r.UpdatePortals = append(r.UpdatePortals,
		&web.PortalUpdate{
			Name: lcb.portals.DataTable(),
			Body: dataTable,
		},
		&web.PortalUpdate{
			Name: lcb.portals.DataTableAdditions(),
			Body: dataTableAdditions,
		},
	)

	return
}

func (lcb *ListingComponentBuilder) actionsComponent(
	msgr *Messages,
	ctx *web.EventContext,
	inDialog bool,
) h.HTMLComponent {
	var (
		b          = lcb.b
		actionBtns h.HTMLComponents
	)

	// Render bulk actions
	for _, ba := range b.bulkActions {
		if b.mb.Info().Verifier().SnakeDo(PermBulkActions, ba.name).WithReq(ctx.R).IsAllowed() != nil {
			continue
		}

		var btn h.HTMLComponent
		if ba.buttonCompFunc != nil {
			btn = ba.buttonCompFunc(ctx)
		} else {
			buttonColor := ba.buttonColor
			if buttonColor == "" {
				buttonColor = ColorSecondary
			}
			onclick := web.Plaid().EventFunc(actions.OpenBulkActionDialog).
				Queries(url.Values{bulkPanelOpenParamName: []string{ba.name}}).
				MergeQuery(true)
			if inDialog {
				onclick.URL(ctx.R.RequestURI).
					Query(ParamOverlay, actions.Dialog)
			}
			btn = VBtn(b.mb.getLabel(ba.NameLabel)).
				Color(buttonColor).
				Variant(VariantFlat).
				// Size(SizeSmall).
				Class("ml-2").
				Attr("@click", onclick.Go())
		}

		actionBtns = append(actionBtns, btn)
	}

	// Render actions
	for _, ba := range b.actions {
		if b.mb.Info().Verifier().SnakeDo(PermActions, ba.name).WithReq(ctx.R).IsAllowed() != nil {
			continue
		}

		var btn h.HTMLComponent
		if ba.buttonCompFunc != nil {
			btn = ba.buttonCompFunc(ctx)
		} else {
			buttonColor := ba.buttonColor
			if buttonColor == "" {
				buttonColor = ColorPrimary
			}

			onclick := web.Plaid().EventFunc(actions.OpenActionDialog).
				Queries(url.Values{actionPanelOpenParamName: []string{ba.name}}).
				MergeQuery(true)
			if inDialog {
				onclick.URL(ctx.R.RequestURI).
					Query(ParamOverlay, actions.Dialog)
			}
			btn = VBtn(b.mb.getLabel(ba.NameLabel)).
				Color(buttonColor).
				Variant(VariantFlat).
				// Size(SizeSmall).
				Class("ml-2").
				Attr("@click", onclick.Go())
		}

		actionBtns = append(actionBtns, btn)
	}

	// if len(actionBtns) == 0 {
	//	return nil
	// }

	if b.actionsAsMenu {
		var listItems []h.HTMLComponent
		for _, btn := range actionBtns {
			listItems = append(listItems, VListItem(btn))
		}
		return h.Components(VMenu(
			web.Slot(
				VBtn("Actions").
					Attr("v-bind", "props").
					Attr("v-on", "on").
					Size(SizeSmall),
			).Name("activator").Scope("{ on, props }"),
			VList(listItems...),
		).OpenOnHover(true))
	}

	if b.newBtnFunc != nil {
		if btn := b.newBtnFunc(ctx); btn != nil {
			actionBtns = append(actionBtns, b.newBtnFunc(ctx))
		}
	} else if b.mb.Info().CanCreate(ctx.R) {
		mode := OverlayMode(ctx)

		onclick := web.Plaid().
			EventFunc(actions.New).URL(ctx.R.RequestURI).
			Query(ParamTargetPortal, lcb.portals.Temp()).
			Query(ParamOverlay, mode.Up().String()).
			Query(ParamPostChangeCallback, b.reloadCallback(ctx).Encode())

		actionBtns = append(actionBtns, VBtn("").
			Color("primary").
			Variant(VariantFlat).
			Theme("dark").Class("ml-2").
			// Size(SizeSmall).
			Attr("@click", onclick.Go()).
			Icon(true).
			Density("comfortable").
			Children(VIcon("mdi-plus")))

	}
	return actionBtns
}

func (b *ListingBuilder) reloadURI(ctx *web.EventContext) string {
	mode := OverlayMode(ctx)
	portalID := GetPortalID(ctx.R)
	return web.Plaid().
		URL(ctx.R.RequestURI).
		EventFunc(actions.ReloadList).
		ValidQuery(ParamOverlay, mode.String()).
		Query(ParamPortalID, portalID).
		StringQuery(ctx.Queries().Encode()).
		Go()
}

func (b *ListingBuilder) reloadCallback(ctx *web.EventContext) (cb *web.Callback) {
	cb = new(web.Callback)
	cb.AddScript(b.reloadURI(ctx))
	return
}
