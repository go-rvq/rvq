package presets

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
)

type ListingBuilder struct {
	mb                 *ModelBuilder
	bulkActions        []*BulkActionBuilder
	footerActions      []*FooterActionBuilder
	actions            []*ActionBuilder
	actionsAsMenu      bool
	filterDataFunc     FilterDataFunc
	filterTabsFunc     FilterTabsFunc
	newBtnFunc         ComponentFunc
	pageFunc           web.PageFunc
	cellWrapperFunc    vx.CellWrapperFunc
	Searcher           SearchFunc
	Deleter            DeleteFunc
	searchColumns      []string
	itemActions        []*ActionBuilder
	prependListButtons []func(ctx *web.EventContext) h.HTMLComponents
	appendListButtons  []func(ctx *web.EventContext) h.HTMLComponents

	// title is the title of the listing page.
	// its default value is "Listing ${modelName}".
	title string

	// perPage is the number of records per page.
	// if request query param "per_page" is set, it will be set to that value.
	// if the final value is less than 0, it will be set to 50.
	// if the final value is greater than 1000, it will be set to 1000.
	perPage int64

	// disablePagination is used to disable pagination, its default value is false.
	// if it is true, the following will happen:
	// 1. the pagination component will not display on listing page.
	// 2. the perPage will actually be ignored.
	// 3. all data will be returned in one page.
	disablePagination bool

	orderBy                string
	orderableFields        []*OrderableField
	selectableColumns      bool
	conditions             []*SQLCondition
	dialogWidth            string
	dialogHeight           string
	dataTableDensity       string
	recordEncoderFactories map[string]RecordEncoderFactory[any]
	configureComponent     func(lcb *ListingComponentBuilder)

	FieldsBuilder
	RowMenuFields

	CreatingRestrictionField[*ListingBuilder]
	EditingRestrictionField[*ListingBuilder]
	DetailingRestrictionField[*ListingBuilder]
	DeletingRestrictionField[*ListingBuilder]
}

func NewListingBuilder(mb *ModelBuilder, fieldsBuilder FieldsBuilder) *ListingBuilder {
	lb := &ListingBuilder{mb: mb, FieldsBuilder: fieldsBuilder}
	lb.RowMenuFields.init(mb)

	strType := reflect.TypeOf("")
	for _, field := range fieldsBuilder.fields {
		if field.structField != nil && field.structField.Type == strType {
			lb.searchColumns = append(lb.searchColumns, field.ColumnName())
		}
	}
	lb.CreatingRestriction = NewRestriction(lb, func(r *Restriction[*ListingBuilder]) {
		r.Insert(mb.CreatingRestriction)
	})
	lb.EditingRestriction = NewObjRestriction(lb, func(r *ObjRestriction[*ListingBuilder]) {
		r.Insert(mb.EditingRestriction)
	})
	lb.DetailingRestriction = NewObjRestriction(lb, func(r *ObjRestriction[*ListingBuilder]) {
		r.Insert(mb.DetailingRestriction)
	})
	lb.DeletingRestriction = NewObjRestriction(lb, func(r *ObjRestriction[*ListingBuilder]) {
		r.Insert(mb.DetailingRestriction)
	})
	return lb
}

func (mb *ModelBuilder) newListing() (lb *ListingBuilder) {
	mb.listing = NewListingBuilder(mb, *mb.NewFieldsBuilder(mb.listFieldBuilders.HasMode(LIST)...))
	mb.listing.DeleteFunc(mb.Deleter)
	mb.listing.SearchFunc(mb.Searcher)

	rmb := mb.listing.RowMenu()
	// rmb.RowMenuItem("Edit").ComponentFunc(editRowMenuItemFunc(mb.Info(), "", url.Values{}))
	rmb.SetRowMenuItem("Delete").ComponentFunc(NewDeletingMenuItemBuilder(mb.Info()).
		SetWrapEvent(func(rctx *RecordMenuItemContext, e *web.VueEventTagBuilder) {
			e.Query(ParamPostChangeCallback, web.CallbackScript(mb.listing.reloadURI(rctx.Ctx)).Encode())
		}).
		Build())
	return
}

func (mb *ModelBuilder) Listing(vs ...string) (r *ListingBuilder) {
	r = mb.listing
	if len(vs) == 0 {
		return
	}

	r.Only(vs...)
	return r
}

func (mb *ModelBuilder) ListingAny(vs ...interface{}) (r *ListingBuilder) {
	var names = make([]string, len(vs))
	for i, v := range vs {
		names[i] = v.(string)
	}
	return mb.Listing(names...)
}

func (b *ListingBuilder) ModelBuilder() *ModelBuilder {
	return b.mb
}

func (b *ListingBuilder) Only(vs ...string) (r *ListingBuilder) {
	r = b
	ivs := make([]interface{}, 0, len(vs))
	for _, v := range vs {
		ivs = append(ivs, v)
	}
	r.FieldsBuilder = *r.FieldsBuilder.Only(ivs...)
	return
}

func (b *ListingBuilder) Except(vs ...string) (r *ListingBuilder) {
	r = b
	r.FieldsBuilder = *r.FieldsBuilder.Except(vs...)
	return
}

func (b *ListingBuilder) PageFunc(pf web.PageFunc) (r *ListingBuilder) {
	b.pageFunc = pf
	return b
}

func (b *ListingBuilder) CellWrapperFunc(cwf vx.CellWrapperFunc) (r *ListingBuilder) {
	b.cellWrapperFunc = cwf
	return b
}

func (b *ListingBuilder) DisablePagination(v bool) (r *ListingBuilder) {
	b.disablePagination = v
	return b
}

func (b *ListingBuilder) SearchFunc(v SearchFunc) (r *ListingBuilder) {
	b.Searcher = v
	return b
}

func (b *ListingBuilder) WrapSearchFunc(w func(in SearchFunc) SearchFunc) (r *ListingBuilder) {
	b.Searcher = w(b.Searcher)
	return b
}

func (b *ListingBuilder) Title(title string) (r *ListingBuilder) {
	b.title = title
	return b
}

func (b *ListingBuilder) SearchColumns(vs ...string) (r *ListingBuilder) {
	b.searchColumns = vs
	return b
}

func (b *ListingBuilder) PerPage(v int64) (r *ListingBuilder) {
	b.perPage = v
	return b
}

func (b *ListingBuilder) OrderBy(v string) (r *ListingBuilder) {
	b.orderBy = v
	return b
}

func (b *ListingBuilder) NewButtonFunc(v ComponentFunc) (r *ListingBuilder) {
	b.newBtnFunc = v
	return b
}

func (b *ListingBuilder) ActionsAsMenu(v bool) (r *ListingBuilder) {
	b.actionsAsMenu = v
	return b
}

type OrderableField struct {
	FieldName string
	DBColumn  string
}

func (b *ListingBuilder) OrderableFields(v []*OrderableField) (r *ListingBuilder) {
	b.orderableFields = v
	return b
}

func (b *ListingBuilder) SelectableColumns(v bool) (r *ListingBuilder) {
	b.selectableColumns = v
	return b
}

func (b *ListingBuilder) Conditions(v []*SQLCondition) (r *ListingBuilder) {
	b.conditions = v
	return b
}

func (b *ListingBuilder) DialogWidth(v string) (r *ListingBuilder) {
	b.dialogWidth = v
	return b
}

func (b *ListingBuilder) DialogHeight(v string) (r *ListingBuilder) {
	b.dialogHeight = v
	return b
}

func (b *ListingBuilder) RecordEncoderFactory(name string) RecordEncoderFactory[any] {
	return b.recordEncoderFactories[name]
}

func (b *ListingBuilder) SetRecordEncoderFactory(name string, encoderFactory RecordEncoderFactory[any]) *ListingBuilder {
	if b.recordEncoderFactories == nil {
		b.recordEncoderFactories = make(map[string]RecordEncoderFactory[any])
	}
	b.recordEncoderFactories[name] = encoderFactory
	return b
}

func (b *ListingBuilder) GetPageFunc() web.PageFunc {
	if b.pageFunc != nil {
		return b.pageFunc
	}
	return b.defaultPageFunc
}

func (b *ListingBuilder) RowMenuOfItems(ctx *web.EventContext) (fs RecordMenuItemFuncs) {
	fs = b.RowMenu().listingItemFuncs(ctx)
	actions := b.itemActions

	for _, action := range b.mb.detailing.actions {
		if action.showInList {
			actions = append(actions, action)
		}
	}

	if len(actions) > 0 {
		fs = append(fs, func(rctx *RecordMenuItemContext) h.HTMLComponent {
			actionsMenus, actionsErrors := BuildMenuItemCompomentsOfActions(rctx.TempPortal, ctx, b.mb, rctx.ID, rctx.Obj, actions...)
			var items h.HTMLComponents

			if len(actionsErrors) > 0 {
				items = append(items, actionsErrors...)
			}

			for _, item := range actionsMenus {
				items = append(items, item)
			}
			return items
		})
	}

	return
}

func (b *ListingBuilder) PrependListAction(f func(ctx *web.EventContext) h.HTMLComponents) *ListingBuilder {
	b.prependListButtons = append(b.prependListButtons, f)
	return b
}

func (b *ListingBuilder) AppendListAction(f func(ctx *web.EventContext) h.HTMLComponents) *ListingBuilder {
	b.appendListButtons = append(b.appendListButtons, f)
	return b
}

func (b *ListingBuilder) ConfigureComponent() func(lcb *ListingComponentBuilder) {
	return b.configureComponent
}

func (b *ListingBuilder) SetConfigureComponent(configureComponent func(lcb *ListingComponentBuilder)) *ListingBuilder {
	b.configureComponent = configureComponent
	return b
}

const (
	bulkPanelOpenParamName  = "bulkOpen"
	DeleteConfirmPortalName = "deleteConfirm"

	detailingContentPortalName = "detailingContentPortal"
	formPortalName             = "formPortal"

	SelectedEventParamName       = "selectedEvent"
	SelectedEventConfigParamName = "selectedEventConfig"
)

func (b *ListingBuilder) TTitle(ctx context.Context) string {
	if b.title != "" {
		return b.title
	}
	return MustGetMessages(ctx).ListingObjectTitle(b.mb.TTitlePlural(ctx))
}

func (b *ListingBuilder) defaultPageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	if b.mb.permissioner.ReqLister(ctx.R).Denied() {
		err = perm.PermissionDenied
		return
	}
	title := b.TTitle(ctx.R.Context())
	r.PageTitle = title

	r.Body, err = b.listingComponent(ctx)
	r.Wrap(func(comp h.HTMLComponent) h.HTMLComponent {
		// return comp
		return b.wrapComp(ctx, comp)
	})
	return
}

func (b *ListingBuilder) records(ctx *web.EventContext) (r web.EventResponse, err error) {
	if b.mb.permissioner.ReqLister(ctx.R).Denied() {
		err = perm.PermissionDenied
		return
	}

	var sr SearchResult

	if sr, err = b.search(ctx); err != nil {
		return
	}

	if ctx.R.FormValue(ParamMustResult) == "true" {
		r.Data = sr.Records
		return
	}

	if encName := ctx.R.FormValue(ParamListingEncoder); encName != "" {
		encFacory, ok := b.recordEncoderFactories[encName]
		if !ok {
			err = fmt.Errorf("%s: %q not found", ParamListingEncoder, encName)
			return
		}

		encodedRecords := make([]any, reflect.ValueOf(sr.Records).Len())
		enc := encFacory(ctx)

		var i int
		reflectutils.ForEach(sr.Records, func(v interface{}) {
			encodedRecords[i] = enc(v)
			i++
		})

		sr.Records = encodedRecords
	}

	r.Data = &sr
	return
}

type SearchResult struct {
	Records           any               `json:"Records,omitempty"`
	TotalCount        int               `json:"TotalCount,omitempty"`
	Page              int64             `json:"Page,omitempty"`
	PerPage           int64             `json:"PerPage,omitempty"`
	OrderBys          []*ColOrderBy     `json:"OrderBys,omitempty"`
	orderableFieldMap map[string]string // map[FieldName]DBColumn
}

func (b *ListingBuilder) search(ctx *web.EventContext) (r SearchResult, err error) {
	var (
		qs            = ctx.R.Form
		perPage int64 = 0
	)

	if !b.disablePagination {
		var requestPerPage int64
		qPerPageStr := qs.Get("per_page")
		qPerPage, _ := strconv.ParseInt(qPerPageStr, 10, 64)
		if qPerPage != 0 {
			setLocalPerPage(ctx, b.mb, qPerPage)
			requestPerPage = qPerPage
		} else if cPerPage := getLocalPerPage(ctx, b.mb); cPerPage != 0 {
			requestPerPage = cPerPage
		}

		perPage = b.perPage
		if requestPerPage != 0 {
			perPage = requestPerPage
		}
		if perPage == 0 {
			perPage = 50
		}
		if perPage > 1000 {
			perPage = 1000
		}
	}

	var orderBySQL string
	r.OrderBys = GetOrderBysFromQuery(qs)
	// map[FieldName]DBColumn
	r.orderableFieldMap = make(map[string]string)
	for _, v := range b.orderableFields {
		r.orderableFieldMap[v.FieldName] = v.DBColumn
	}
	for _, ob := range r.OrderBys {
		dbCol, ok := r.orderableFieldMap[ob.FieldName]
		if !ok {
			continue
		}
		orderBySQL += fmt.Sprintf("%s %s,", dbCol, ob.Asc)
	}
	// remove the last ","
	if orderBySQL != "" {
		orderBySQL = orderBySQL[:len(orderBySQL)-1]
	}
	if orderBySQL == "" {
		if b.orderBy != "" {
			orderBySQL = b.orderBy
		} else if fields := b.mb.Schema().PrimaryFields(); len(fields) > 0 {
			orderBySQL = fmt.Sprintf("%s DESC", strings.Join(fields.QuotedFullDBNames(), ", "))
		}
	}

	searchParams := &SearchParams{
		KeywordColumns: b.searchColumns,
		Keyword:        qs.Get("keyword"),
		PerPage:        perPage,
		OrderBy:        orderBySQL,
		Query:          web.Query(qs),
		SQLConditions:  b.conditions,
	}

	searchParams.Page, _ = strconv.ParseInt(qs.Get("page"), 10, 64)
	if searchParams.Page == 0 {
		searchParams.Page = 1
	}

	r.Page = searchParams.Page
	r.PerPage = perPage

	var fd vx.FilterData
	if b.filterDataFunc != nil {
		fd = b.filterDataFunc(ctx)
		cond, args := fd.SetByQueryString(ctx.R.URL.RawQuery)

		searchParams.SQLConditions = append(searchParams.SQLConditions, &SQLCondition{
			Query: cond,
			Args:  args,
		})
	}

	if b.Searcher == nil || b.mb.CurrentDataOperator() == nil {
		err = errors.New("presets.New().DataOperator(...) required")
		return
	}

	r.Records, r.TotalCount, err = b.Searcher(b.mb.NewModelSlice(), searchParams, ctx)
	return
}

func getSelectedIds(ctx *web.EventContext) (selected []string) {
	selectedValue := ctx.R.URL.Query().Get(ParamSelectedIds)
	if len(selectedValue) > 0 {
		selected = strings.Split(selectedValue, ",")
	}
	return selected
}

func (b *ListingBuilder) bulkPanel(
	bulk *BulkActionBuilder,
	selectedIds []string,
	processedSelectedIds []string,
	ctx *web.EventContext,
) (r h.HTMLComponent) {
	msgr := MustGetMessages(ctx.Context())

	var errComp h.HTMLComponent
	if vErr, ok := ctx.Flash.(*web.ValidationErrors); ok {
		if gErr := vErr.GetGlobalError(); gErr != "" {
			errComp = VAlert(h.Text(gErr)).
				Border("left").
				Type("error").
				Elevation(2)
		}
	}
	var processSelectedIdsNotice h.HTMLComponent
	if len(processedSelectedIds) < len(selectedIds) {
		unactionables := make([]string, 0, len(selectedIds))
		{
			processedSelectedIdsM := make(map[string]struct{})
			for _, v := range processedSelectedIds {
				processedSelectedIdsM[v] = struct{}{}
			}
			for _, v := range selectedIds {
				if _, ok := processedSelectedIdsM[v]; !ok {
					unactionables = append(unactionables, v)
				}
			}
		}

		if len(unactionables) > 0 {
			var noticeText string
			if bulk.selectedIdsProcessorNoticeFunc != nil {
				noticeText = bulk.selectedIdsProcessorNoticeFunc(selectedIds, processedSelectedIds, unactionables)
			} else {
				var idsText string
				if len(unactionables) <= 10 {
					idsText = strings.Join(unactionables, ", ")
				} else {
					idsText = fmt.Sprintf("%s...(+%d)", strings.Join(unactionables[:10], ", "), len(unactionables)-10)
				}
				noticeText = msgr.BulkActionSelectedIdsProcessNotice(idsText)
			}
			processSelectedIdsNotice = VAlert(h.Text(noticeText)).
				Type("warning")
		}
	}

	onOK := web.Plaid().EventFunc(actions.DoBulkAction).
		Query(ParamBulkActionName, bulk.name).
		MergeQuery(true)
	if isInDialogFromQuery(ctx) {
		onOK.URL(ctx.R.RequestURI)
	}

	var err error
	r, err = bulk.compFunc(selectedIds, ctx)

	if err != nil {
		return VCard(
			VCardTitle(
				h.Text(bulk.NameLabel.label),
			),
			VCardText(
				VAlert(h.RawHTML(err.Error())).
					Type(ColorError).
					Variant(VariantTonal).
					Density(DensityCompact),
			),
		)
	}

	return VCard(
		VCardTitle(
			h.Text(bulk.NameLabel.label),
		),
		VCardText(
			errComp,
			processSelectedIdsNotice,
			r,
		),
		VCardActions(
			VSpacer(),
			VBtn(msgr.Cancel).
				Variant(VariantFlat).
				Class("ml-2").
				Attr("@click", closeDialogVarScript),

			VBtn(msgr.OK).
				Color("primary").
				Variant(VariantFlat).
				Theme(ThemeDark).
				Attr("@click", onOK.Go()),
		),
	)
}

func (b *ListingBuilder) openBulkActionDialog(ctx *web.EventContext) (r web.EventResponse, err error) {
	msgr := MustGetMessages(ctx.Context())
	selected := getSelectedIds(ctx)
	bulkName := ctx.R.URL.Query().Get(bulkPanelOpenParamName)
	a := getBulkAction(b.bulkActions, bulkName)

	if a == nil {
		err = errors.New("cannot find requested action")
		return
	}

	if a.Verifier(b.mb.permissioner.ReqList(ctx.R)).Denied() {
		err = perm.PermissionDenied
		return
	}

	if len(selected) == 0 && !a.allowEmpty {
		ShowMessage(&r, msgr.PleaseSelectRecord, "warning")
		return
	}

	// If selectedIdsProcessorFunc is not nil, process the request in it and skip the confirmation dialog
	var processedSelectedIds []string
	if a.selectedIdsProcessorFunc != nil {
		processedSelectedIds, err = a.selectedIdsProcessorFunc(selected, ctx)
		if err != nil {
			return
		}
		if len(processedSelectedIds) == 0 {
			if a.selectedIdsProcessorNoticeFunc != nil {
				ShowMessage(&r, a.selectedIdsProcessorNoticeFunc(selected, processedSelectedIds, selected), "warning")
			} else {
				ShowMessage(&r, msgr.BulkActionNoAvailableRecords, "warning")
			}
			return
		}
	} else {
		processedSelectedIds = selected
	}

	err = a.View(b.mb, processedSelectedIds, ctx, &r)
	return
}

func (b *ListingBuilder) doBulkAction(ctx *web.EventContext) (r web.EventResponse, err error) {
	a := getBulkAction(b.bulkActions, ctx.R.FormValue(ParamBulkActionName))
	if a == nil {
		panic("bulk required")
	}

	if a.Verifier(b.mb.permissioner.ReqList(ctx.R)).Denied() {
		err = perm.PermissionDenied
		return
	}

	selectedIds := getSelectedIds(ctx)

	var err1 error
	var processedSelectedIds []string
	if a.selectedIdsProcessorFunc != nil {
		processedSelectedIds, err1 = a.selectedIdsProcessorFunc(selectedIds, ctx)
	} else {
		processedSelectedIds = selectedIds
	}

	opts := &ListingDoBulkActionOptions{}

	if err1 == nil {
		reset := ctx.WithData(opts)
		err1 = a.Do(processedSelectedIds, ctx, &r)
		reset()
	}

	if err1 != nil {
		if _, ok := err1.(*web.ValidationErrors); !ok {
			vErr := &web.ValidationErrors{}
			vErr.GlobalError(err1.Error())
			ctx.Flash = vErr

			err = a.View(b.mb, processedSelectedIds, ctx, &r)
			return
		}
	}

	if ctx.Flash != nil {
		r.UpdatePortal(actions.Dialog.ContentPortalName(), b.bulkPanel(a, selectedIds, processedSelectedIds, ctx))
		return
	}

	web.AppendRunScripts(&r, "closer.show = false")

	if opts.Flash != nil {
		ctx.Flash = opts.Flash
	} else {
		msgr := MustGetMessages(ctx.Context())
		ctx.Flash = msgr.SuccessfullyUpdated
	}

	if opts.ListReloadDisabled {
		return
	}

	if isInDialogFromQuery(ctx) {
		qs := ctx.Queries()
		qs.Del(bulkPanelOpenParamName)
		qs.Del(ParamBulkActionName)
		web.AppendRunScripts(&r,
			web.Plaid().
				URL(ctx.R.RequestURI).
				EventFunc(actions.UpdateListingDialog).
				Queries(qs).
				Go(),
		)
	} else {
		r.PushState = web.Location(url.Values{bulkPanelOpenParamName: []string{}}).MergeQuery(true)
	}

	return
}

func (b *ListingBuilder) actionPanel(action *ActionBuilder, ctx *web.EventContext) (r h.HTMLComponent) {
	msgr := MustGetMessages(ctx.Context())

	var errComp h.HTMLComponent

	if vErr, ok := ctx.Flash.(*web.ValidationErrors); ok {
		if gErr := vErr.GetGlobalError(); gErr != "" {
			errComp = VAlert(h.Text(gErr)).
				Border("left").
				Type("error").
				Elevation(2)
		}
	}

	onOK := web.Plaid().EventFunc(actions.DoListingAction).
		Query(ParamAction, action.name).
		MergeQuery(true)
	if isInDialogFromQuery(ctx) {
		onOK.URL(ctx.R.RequestURI)
	}

	var comp h.HTMLComponent
	if action.compFunc != nil {
		var err error
		if comp, err = action.compFunc("", ctx); err != nil {
			errComp = VAlert(h.Text(err.Error())).
				Border("left").
				Type("error").
				Elevation(2)
			return
		}
	}

	return VCard(
		VCardTitle(
			h.Text(action.NameLabel.label),
		),
		VCardText(
			errComp,
			comp,
		),
		VCardActions(
			VSpacer(),
			VBtn(msgr.Cancel).
				Variant(VariantFlat).
				Class("ml-2").
				Attr("@click", closeDialogVarScript),

			VBtn(msgr.OK).
				Color("primary").
				Variant(VariantFlat).
				Theme(ThemeDark).
				Attr("@click", onOK.Go()),
		),
	)
}

func (b *ListingBuilder) openItemActionDialog(ctx *web.EventContext) (r web.EventResponse, err error) {
	return b.openActionDialogInternal(b.itemActions, ctx)
}

func (b *ListingBuilder) openActionDialog(ctx *web.EventContext) (r web.EventResponse, err error) {
	return b.openActionDialogInternal(b.actions, ctx)
}

func (b *ListingBuilder) openActionDialogInternal(actionList []*ActionBuilder, ctx *web.EventContext) (r web.EventResponse, err error) {
	actionName := ctx.R.URL.Query().Get(ParamAction)
	action := getAction(actionList, actionName)
	if action == nil {
		err = errors.New("cannot find requested action")
		return
	}

	if b.mb.permissioner.ReqListActioner(ctx.R, action.name).Denied() {
		err = perm.PermissionDenied
		return
	}

	r.RunScript = `console.log("action dialog", presetsListing.uri)`
	err = action.View(b.mb, ctx.R.Form.Get(ParamID), ctx, &r)
	return
}

func (b *ListingBuilder) doListingItemAction(ctx *web.EventContext) (r web.EventResponse, err error) {
	return b.doListingActionInternal(b.itemActions, ctx)
}

func (b *ListingBuilder) doListingAction(ctx *web.EventContext) (r web.EventResponse, err error) {
	return b.doListingActionInternal(b.actions, ctx)
}

func (b *ListingBuilder) doListingActionInternal(actionList []*ActionBuilder, ctx *web.EventContext) (r web.EventResponse, err error) {
	action := getAction(actionList, ctx.R.FormValue(ParamAction))
	if action == nil {
		err = errors.New("cannot find requested action")
		return
	}

	if b.mb.permissioner.ReqListActioner(ctx.R, action.PermName()).Denied() {
		err = perm.PermissionDenied
		return
	}

	var success bool
	if success, err = action.Do(b.mb, "", ctx, &r); err != nil || !success {
		return
	}
	return
}

const ActiveFilterTabQueryKey = "active_filter_tab"

func (b *ListingBuilder) filterTabs(portals *ListingPortals,
	ctx *web.EventContext,
	inDialog bool,
) (r h.HTMLComponent) {
	if b.filterTabsFunc == nil {
		return
	}

	qs := ctx.R.URL.Query()

	tabs := VTabs().
		Class("mb-2").
		ShowArrows(true).
		Color("primary").
		Density(DensityCompact)

	tabsData := b.filterTabsFunc(ctx)
	var defaultTab *FilterTab
	for _, tab := range tabsData {
		if tab.Default {
			if defaultTab != nil {
				return VAlert(h.RawHTML("Many filter tabs with <b>Default</b> flag.")).Color("error")
			}
			defaultTab = tab
		}
	}
	value := -1
	activeTabValue := qs.Get(ActiveFilterTabQueryKey)

	for i, td := range tabsData {
		// Find selected tab by active_filter_tab=xx in the url query
		if activeTabValue != "" && activeTabValue == td.ID {
			value = i
		}

		tabContent := h.Text(td.Label)
		if td.AdvancedLabel != nil {
			tabContent = td.AdvancedLabel
		}

		totalQuery := url.Values{}
		if td.ID != "" {
			totalQuery.Set(ActiveFilterTabQueryKey, td.ID)
		}
		for k, v := range td.Query {
			totalQuery[k] = v
		}

		onclick := web.Plaid().Queries(totalQuery)
		portalID := ctx.R.FormValue(ParamPortalID)
		if inDialog {
			onclick.URL(ctx.R.RequestURI).
				EventFunc(actions.UpdateListingDialog).
				ValidQuery(ParamPortalID, portalID)
		} else {
			onclick.PushState(true)
		}
		tabs.AppendChild(
			VTab(tabContent).
				Attr("@click", onclick.Go()),
		)
	}

	if value == -1 {
	loop:
		for i, td := range tabsData {
			if td.ID == "" {
				if len(td.Query) > 0 {
					for k := range td.Query {
						if ctx.R.FormValue(k) != td.Query.Get(k) {
							continue loop
						}
					}
					value = i
					break loop
				}
			}
		}
	}

	return tabs.ModelValue(value)
}

type selectColumns struct {
	DisplayColumns []string       `json:"displayColumns,omitempty"`
	SortedColumns  []sortedColumn `json:"sortedColumns,omitempty"`
}
type sortedColumn struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

func (b *ListingBuilder) selectColumnsBtn(
	pageURL *url.URL,
	ctx *web.EventContext,
	inDialog bool,
) (btn h.HTMLComponent, displaySortedFields FieldBuilders) {
	var (
		_, respath         = path.Split(pageURL.Path)
		displayColumnsName = fmt.Sprintf("%s_display_columns", respath)
		sortedColumnsName  = fmt.Sprintf("%s_sorted_columns", respath)
		originalColumns    []string
		displayColumns     []string
		sortedColumns      []string
	)

	b.fields.EachHavesComponent(func(f *FieldBuilder) bool {
		if b.mb.permissioner.ReqLister(ctx.R).SnakeOn(FieldPerm(f.name)).Denied() {
			originalColumns = append(originalColumns, f.name)
		}
		return true
	})

	// get the columns setting from url params or cookie data
	if urldata := pageURL.Query().Get(displayColumnsName); urldata != "" {
		if urlColumns := strings.Split(urldata, ","); len(urlColumns) > 0 {
			displayColumns = urlColumns
		}
	}

	if urldata := pageURL.Query().Get(sortedColumnsName); urldata != "" {
		if urlColumns := strings.Split(urldata, ","); len(urlColumns) > 0 {
			sortedColumns = urlColumns
		}
	}

	// get the columns setting from  cookie data
	if len(displayColumns) == 0 {
		cookiedata, err := ctx.R.Cookie(displayColumnsName)
		if err == nil {
			if cookieColumns := strings.Split(cookiedata.Value, ","); len(cookieColumns) > 0 {
				displayColumns = cookieColumns
			}
		}
	}

	if len(sortedColumns) == 0 {
		cookiedata, err := ctx.R.Cookie(sortedColumnsName)
		if err == nil {
			if cookieColumns := strings.Split(cookiedata.Value, ","); len(cookieColumns) > 0 {
				sortedColumns = cookieColumns
			}
		}
	}

	// check if listing fileds is changed. if yes, use the original columns
	var originalFiledsChanged bool

	if len(sortedColumns) > 0 && len(originalColumns) != len(sortedColumns) {
		originalFiledsChanged = true
	}

	if len(sortedColumns) > 0 && !originalFiledsChanged {
		for _, sortedColumn := range sortedColumns {
			var find bool
			for _, originalColumn := range originalColumns {
				if sortedColumn == originalColumn {
					find = true
					break
				}
			}
			if !find {
				originalFiledsChanged = true
				break
			}
		}
	}

	if len(displayColumns) > 0 && !originalFiledsChanged {
		for _, displayColumn := range displayColumns {
			var find bool
			for _, originalColumn := range originalColumns {
				if displayColumn == originalColumn {
					find = true
					break
				}
			}
			if !find {
				originalFiledsChanged = true
				break
			}
		}
	}

	// save display columns setting on cookie
	if !originalFiledsChanged && len(displayColumns) > 0 {
		http.SetCookie(ctx.W, &http.Cookie{
			Name:  displayColumnsName,
			Value: strings.Join(displayColumns, ","),
		})
	}

	// save sorted columns setting on cookie
	if !originalFiledsChanged && len(sortedColumns) > 0 {
		http.SetCookie(ctx.W, &http.Cookie{
			Name:  sortedColumnsName,
			Value: strings.Join(sortedColumns, ","),
		})
	}

	// set the data for displaySortedFields on data table
	if originalFiledsChanged || (len(sortedColumns) == 0 && len(displayColumns) == 0) {
		displaySortedFields = b.fields.FieldsFromLayout(b.CurrentLayout(), FieldRenderable())
	}

	if originalFiledsChanged || len(displayColumns) == 0 {
		displayColumns = originalColumns
	}

	if originalFiledsChanged || len(sortedColumns) == 0 {
		sortedColumns = originalColumns
	}

	if len(displaySortedFields) == 0 {
		for _, sortedColumn := range sortedColumns {
			for _, displayColumn := range displayColumns {
				if sortedColumn == displayColumn {
					displaySortedFields = append(displaySortedFields, b.Field(sortedColumn))
					break
				}
			}
		}
	}

	// set the data for selected columns on toolbar
	selectColumns := selectColumns{
		DisplayColumns: displayColumns,
	}

	for _, sc := range sortedColumns {
		selectColumns.SortedColumns = append(selectColumns.SortedColumns, sortedColumn{
			Name:  sc,
			Label: b.FieldsBuilder.GetField(sc).ContextLabel(b.mb.Info(), ctx.Context()),
		})
	}

	msgr := MustGetMessages(ctx.Context())
	onOK := web.Plaid().
		Query(displayColumnsName, web.Var("locals.displayColumns")).
		Query(sortedColumnsName, web.Var("locals.sortedColumns.map(column => column.name )")).
		MergeQuery(true)
	if inDialog {
		onOK.URL(ctx.R.RequestURI).
			EventFunc(actions.UpdateListingDialog).
			Query(ParamPortalID, ctx.R.FormValue(ParamPortalID))
	}
	// add the HTML component of columns setting into toolbar
	btn = web.Scope(VMenu(
		web.Slot(
			VBtn("").Icon("mdi-cog").Attr("v-bind", "props").Variant(VariantText).Size(SizeSmall),
		).Name("activator").Scope("{ props }"),
		VList(
			h.Tag("vx-draggable").Attr("item-key", "name").Attr("v-model", "locals.sortedColumns", "handle", ".handle", "animation", "300").Children(
				h.Template(
					VListItem(
						VListItemTitle(
							VSwitch().Density(DensityCompact).Attr("v-model", "locals.displayColumns", ":value",
								"element.name",
								":label", "element.label").Color("primary").Class(" mt-2 "),
							VIcon("mdi-reorder-vertical").Class("handle cursor-grab mt-4"),
						).Class("d-flex justify-space-between "),
						VDivider(),
					),
				).Attr("#item", " { element } "),
			),
			VListItem(
				VBtn(msgr.Cancel).Elevation(0).Attr("@click", `locals.selectColumnsMenu = false`),
				VBtn(msgr.OK).Elevation(0).Color("primary").Attr("@click", `locals.selectColumnsMenu = false ; `+onOK.Go()),
			).Class("d-flex justify-space-between"),
		).Density(DensityCompact),
	).CloseOnContentClick(false).Width(240).
		Attr("v-model", "locals.selectColumnsMenu")).
		Slot("{ locals }").LocalsInit(fmt.Sprintf(`{selectColumnsMenu: false,...%s}`, h.JSONString(selectColumns)))
	return
}

func (b *ListingBuilder) filterBar(
	ctx *web.EventContext,
	msgr *Messages,
	fd vx.FilterData,
	inDialog bool,
) (filterBar h.HTMLComponent) {
	if fd == nil {
		return nil
	}
	noVisiableItem := true
	for _, d := range fd {
		if !d.Invisible {
			noVisiableItem = false
			break
		}
	}
	if noVisiableItem {
		return nil
	}

	ft := vx.FilterTranslations{}
	ft.ClearAll = msgr.FiltersClear
	ft.Clear = msgr.Clear
	ft.Add = msgr.FiltersAdd
	ft.Apply = msgr.FilterApply

	for _, d := range fd {
		d.Translations = vx.FilterIndependentTranslations{
			FilterBy: msgr.FilterBy(d.Label),
		}
	}

	ft.To = msgr.FiltersTo
	ft.Month.Year = msgr.Year
	ft.Month.Month = msgr.Month
	ft.Month.MonthNames = msgr.MonthNames

	ft.Number.And = msgr.FiltersNumberAnd
	ft.Number.Equals = msgr.FiltersNumberEquals
	ft.Number.Between = msgr.FiltersNumberBetween
	ft.Number.GreaterThan = msgr.FiltersNumberGreaterThan
	ft.Number.LessThan = msgr.FiltersNumberLessThan

	ft.String.Equals = msgr.FiltersStringEquals
	ft.String.Contains = msgr.FiltersStringContains

	ft.MultipleSelect.In = msgr.FiltersMultipleSelectIn
	ft.MultipleSelect.NotIn = msgr.FiltersMultipleSelectNotIn

	filter := vx.VXFilter(fd).Translations(ft).Attr("v-model:visibility", "filterBarVisible.value")
	if inDialog {
		filter.UpdateModelValue(web.Plaid().
			URL(ctx.R.RequestURI).
			StringQuery(web.Var("$event.encodedFilterData")).
			Query("page", 1).
			ClearMergeQuery(web.Var("$event.filterKeys")).
			EventFunc(actions.UpdateListingDialog).
			Query(ParamPortalID, ctx.R.FormValue(ParamPortalID)).
			Go())
	}
	return filter
}

func getLocalPerPage(
	ctx *web.EventContext,
	mb *ModelBuilder,
) int64 {
	// c is the cookie value of a serials of per page value, split by "$".
	// each value is split by "#".
	// the first part is the uri name of the model builder.
	// the second part is the per page value.
	c, err := ctx.R.Cookie("_perPage")
	if err != nil {
		return 0
	}
	vals := strings.Split(c.Value, "$")
	for _, v := range vals {
		vvs := strings.Split(v, "#")
		if len(vvs) != 2 {
			continue
		}
		if vvs[0] == mb.uriName {
			r, _ := strconv.ParseInt(vvs[1], 10, 64)
			return r
		}
	}

	return 0
}

// setLocalPerPage set the per page value to cookie.
// v is the per page value to set.
func setLocalPerPage(
	ctx *web.EventContext,
	mb *ModelBuilder,
	v int64,
) {
	var oldVals []string
	{
		c, err := ctx.R.Cookie("_perPage")
		if err == nil {
			oldVals = strings.Split(c.Value, "$")
		}
	}
	newVals := []string{fmt.Sprintf("%s#%d", mb.uriName, v)}
	for _, v := range oldVals {
		vvs := strings.Split(v, "#")
		if len(vvs) != 2 {
			continue
		}
		if vvs[0] == mb.uriName {
			continue
		}
		newVals = append(newVals, v)
	}
	http.SetCookie(ctx.W, &http.Cookie{
		Name:  "_perPage",
		Value: strings.Join(newVals, "$"),
	})
}

type OrderMode bool

const (
	OrderDesc OrderMode = false
	OrderAsc  OrderMode = true
)

func (m OrderMode) String() string {
	if m {
		return "ASC"
	}
	return "DESC"
}

type ColOrderBy struct {
	FieldName string
	// ASC, DESC
	Asc OrderMode
}

// GetOrderBysFromQuery gets order bys from query string.
func GetOrderBysFromQuery(query url.Values) []*ColOrderBy {
	r := make([]*ColOrderBy, 0)
	// qs is like "field1_ASC,field2_DESC"
	qs := strings.Split(query.Get("order_by"), ",")
	for _, q := range qs {
		ss := strings.Split(q, "_")
		ssl := len(ss)
		if ssl == 1 {
			continue
		}
		if ss[ssl-1] != "ASC" && ss[ssl-1] != "DESC" {
			continue
		}
		r = append(r, &ColOrderBy{
			FieldName: strings.Join(ss[:ssl-1], "_"),
			Asc:       ss[ssl-1] == "ASC",
		})
	}

	return r
}

func newQueryWithFieldToggleOrderBy(query url.Values, fieldName string) url.Values {
	oldOrderBys := GetOrderBysFromQuery(query)
	var newOrderBysQueryValue []string
	existed := false
	for _, oob := range oldOrderBys {
		if oob.FieldName == fieldName {
			existed = true
			if oob.Asc {
				newOrderBysQueryValue = append(newOrderBysQueryValue, oob.FieldName+"_DESC")
			}
			continue
		}
		newOrderBysQueryValue = append(newOrderBysQueryValue, oob.FieldName+"_"+oob.Asc.String())
	}
	if !existed {
		newOrderBysQueryValue = append(newOrderBysQueryValue, fieldName+"_ASC")
	}

	newQuery := make(url.Values)
	for k, v := range query {
		newQuery[k] = v
	}
	newQuery.Set("order_by", strings.Join(newOrderBysQueryValue, ","))
	return newQuery
}

func (b *ListingBuilder) openListingDialogForSelection(ctx *web.EventContext) (r web.EventResponse, err error) {
	lcb := b.ListingComponentBuilderCtx(ctx).SetSelection(true)
	ctx.R.Form.Set(ParamOverlay, actions.Dialog.String())
	var body h.HTMLComponent
	if body, err = lcb.Build(ctx); err != nil {
		return
	}
	targetPortal := ctx.R.FormValue(ParamTargetPortal)
	b.mb.p.DialogPortal(targetPortal).
		SetValidWidth(b.dialogWidth).
		SetValidHeight(b.dialogHeight).
		SetContentPortalName(lcb.portals.Main()).
		Respond(ctx, &r, web.Scope(body).Slot("{ form }"))
	return
}

func (b *ListingBuilder) openListingDialog(ctx *web.EventContext) (r web.EventResponse, err error) {
	lcb := b.ListingComponentBuilderCtx(ctx)
	ctx.R.Form.Set(ParamOverlay, actions.Dialog.String())

	targetPortal := ctx.R.FormValue(ParamTargetPortal)

	var body h.HTMLComponent
	if body, err = lcb.Build(ctx); err != nil {
		return
	}

	b.mb.p.DialogPortal(targetPortal).
		SetValidWidth(b.dialogWidth).
		SetValidHeight(b.dialogHeight).
		SetContentPortalName(lcb.portals.Main()).
		SetScrollable(true).
		Respond(ctx, &r, b.wrapComp(ctx, body))
	return
}

func (b *ListingBuilder) wrapComp(ctx *web.EventContext, comp h.HTMLComponent) h.HTMLComponent {
	plaid := web.POST().NoCache()
	if IsInDialog(ctx) {
		plaid.ParseURL(ctx.R.URL.Path + "?" + ctx.R.Form.Encode()).
			EventFunc(actions.UpdateListingDialog)
	} else {
		plaid.EventFunc("__reload__")
	}
	return vue.UserComponent(web.Scope(comp).Slot("{ form }")).
		ScopeVar("parentPresetsListing", "presetsListing").
		ScopeVar("presetsListing", "{}").
		Setup(fmt.Sprintf(`({scope}) => {
	scope.presetsListing.parent = scope.parentPresetsListing;
	const presetsListing = scope.presetsListing;
	scope.presetsListing.uri = %q
	scope.presetsListing.loader = %s
}`,
			ctx.R.URL.Path,
			plaid.String(),
		))
}

func (b *ListingBuilder) updateListingDialog(ctx *web.EventContext) (r web.EventResponse, err error) {
	lcb := b.ListingComponentBuilderCtx(ctx)
	ctx.R.Form.Set(ParamOverlay, actions.Dialog.String())
	var dataTable, dataTableAdditions h.HTMLComponent

	if dataTable, dataTableAdditions, err = lcb.GetTableComponents(ctx); err != nil {
		return
	}

	r.UpdatePortal(
		lcb.portals.DataTable(),
		dataTable,
	).UpdatePortal(
		lcb.portals.DataTableAdditions(),
		dataTableAdditions,
	)
	r.RunScript = fmt.Sprintf(`presetsListing.loader.parseUrl(%q)`, ctx.R.RequestURI)
	return
}
