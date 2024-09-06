package presets

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/perm"
	. "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
)

type ListingBuilder struct {
	mb              *ModelBuilder
	bulkActions     []*BulkActionBuilder
	footerActions   []*FooterActionBuilder
	actions         []*ActionBuilder
	actionsAsMenu   bool
	filterDataFunc  FilterDataFunc
	filterTabsFunc  FilterTabsFunc
	newBtnFunc      ComponentFunc
	pageFunc        web.PageFunc
	cellWrapperFunc vx.CellWrapperFunc
	Searcher        SearchFunc
	Deleter         DeleteFunc
	searchColumns   []string

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

	orderBy           string
	orderableFields   []*OrderableField
	selectableColumns bool
	conditions        []*SQLCondition
	dialogWidth       string
	dialogHeight      string
	dataTableDensity  string
	dataListFormatter func(ctx *web.EventContext) func(r any) any
	FieldsBuilder
	RowMenuFields
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
	return lb
}

func (mb *ModelBuilder) newListing() (lb *ListingBuilder) {
	mb.listing = NewListingBuilder(mb, *mb.NewFieldsBuilder(mb.listFieldBuilders.HasMode(LIST)...))
	mb.listing.DeleteFunc(mb.Deleter)
	mb.listing.SearchFunc(mb.Searcher)

	rmb := mb.listing.RowMenu()
	// rmb.RowMenuItem("Edit").ComponentFunc(editRowMenuItemFunc(mb.Info(), "", url.Values{}))
	rmb.RowMenuItem("Delete").ComponentFunc(NewDeletingMenuItemBuilder(mb.Info()).
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

func (b *ListingBuilder) DataListFormatter() func(ctx *web.EventContext) func(r any) any {
	return b.dataListFormatter
}

func (b *ListingBuilder) SetDataListFormatter(dataListFormatter func(ctx *web.EventContext) func(r any) any) *ListingBuilder {
	b.dataListFormatter = dataListFormatter
	return b
}

func (b *ListingBuilder) GetPageFunc() web.PageFunc {
	if b.pageFunc != nil {
		return b.pageFunc
	}
	return b.defaultPageFunc
}

const (
	bulkPanelOpenParamName   = "bulkOpen"
	actionPanelOpenParamName = "actionOpen"
	DeleteConfirmPortalName  = "deleteConfirm"

	detailingContentPortalName = "detailingContentPortal"
	formPortalName             = "formPortal"

	SelectedEventParamName       = "selectedEvent"
	SelectedEventConfigParamName = "selectedEventConfig"
)

func (b *ListingBuilder) TTitle(r *http.Request) string {
	if b.title != "" {
		return b.title
	}
	return MustGetMessages(r).ListingObjectTitle(i18n.T(r, ModelsI18nModuleKey, b.mb.pluralLabel))
}

func (b *ListingBuilder) defaultPageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	if b.mb.Info().Verifier().Do(PermList).WithReq(ctx.R).IsAllowed() != nil {
		err = perm.PermissionDenied
		return
	}
	title := b.TTitle(ctx.R)
	r.PageTitle = title

	r.Body, err = b.listingComponent(ctx)

	return
}

func (b *ListingBuilder) records(ctx *web.EventContext) (r web.EventResponse, err error) {
	var sr SearchResult

	if sr, err = b.search(ctx); err != nil {
		return
	}

	if ctx.R.FormValue(ParamMustResult) == "true" {
		r.Data = sr.objs
		return
	}

	records := sr.objs

	if b.dataListFormatter != nil {
		var (
			formatter        = b.dataListFormatter(ctx)
			formattedRecords = make([]any, reflect.ValueOf(records).Len())
		)

		reflectutils.ForEach(sr.objs, func(i int, v interface{}) {
			formattedRecords[i] = formatter(v)
		})

		records = formattedRecords
	}

	r.Data = map[string]any{
		"records":    records,
		"totalCount": sr.totalCount,
	}
	return
}

type SearchResult struct {
	objs              any
	totalCount        int
	orderableFieldMap map[string]string
	orderBys          []*ColOrderBy
	Page, PerPage     int64
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
	r.orderBys = GetOrderBysFromQuery(qs)
	// map[FieldName]DBColumn
	r.orderableFieldMap = make(map[string]string)
	for _, v := range b.orderableFields {
		r.orderableFieldMap[v.FieldName] = v.DBColumn
	}
	for _, ob := range r.orderBys {
		dbCol, ok := r.orderableFieldMap[ob.FieldName]
		if !ok {
			continue
		}
		orderBySQL += fmt.Sprintf("%s %s,", dbCol, ob.OrderBy)
	}
	// remove the last ","
	if orderBySQL != "" {
		orderBySQL = orderBySQL[:len(orderBySQL)-1]
	}
	if orderBySQL == "" {
		if b.orderBy != "" {
			orderBySQL = b.orderBy
		} else {
			orderBySQL = fmt.Sprintf("%s DESC", b.mb.primaryField)
		}
	}

	searchParams := &SearchParams{
		KeywordColumns: b.searchColumns,
		Keyword:        qs.Get("keyword"),
		PerPage:        perPage,
		OrderBy:        orderBySQL,
		PageQuery:      qs,
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

	r.objs, r.totalCount, err = b.Searcher(b.mb.NewModelSlice(), searchParams, ctx)
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
	msgr := MustGetMessages(ctx.R)

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
	return VCard(
		VCardTitle(
			h.Text(bulk.NameLabel.label),
		),
		VCardText(
			errComp,
			processSelectedIdsNotice,
			bulk.compFunc(selectedIds, ctx),
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

func (b *ListingBuilder) actionPanel(action *ActionBuilder, ctx *web.EventContext) (r h.HTMLComponent) {
	msgr := MustGetMessages(ctx.R)

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
		Query(ParamListingActionName, action.name).
		MergeQuery(true)
	if isInDialogFromQuery(ctx) {
		onOK.URL(ctx.R.RequestURI)
	}

	var comp h.HTMLComponent
	if action.compFunc != nil {
		comp = action.compFunc("", ctx)
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

func (b *ListingBuilder) openActionDialog(ctx *web.EventContext) (r web.EventResponse, err error) {
	actionName := ctx.R.URL.Query().Get(actionPanelOpenParamName)
	action := getAction(b.actions, actionName)
	if action == nil {
		err = errors.New("cannot find requested action")
		return
	}

	b.mb.p.dialog(
		&r,
		b.actionPanel(action, ctx),
		action.dialogWidth,
	)
	return
}

func (b *ListingBuilder) openBulkActionDialog(ctx *web.EventContext) (r web.EventResponse, err error) {
	msgr := MustGetMessages(ctx.R)
	selected := getSelectedIds(ctx)
	bulkName := ctx.R.URL.Query().Get(bulkPanelOpenParamName)
	bulk := getBulkAction(b.bulkActions, bulkName)

	if bulk == nil {
		err = errors.New("cannot find requested action")
		return
	}

	if len(selected) == 0 {
		ShowMessage(&r, "Please select record", "warning")
		return
	}

	// If selectedIdsProcessorFunc is not nil, process the request in it and skip the confirmation dialog
	var processedSelectedIds []string
	if bulk.selectedIdsProcessorFunc != nil {
		processedSelectedIds, err = bulk.selectedIdsProcessorFunc(selected, ctx)
		if err != nil {
			return
		}
		if len(processedSelectedIds) == 0 {
			if bulk.selectedIdsProcessorNoticeFunc != nil {
				ShowMessage(&r, bulk.selectedIdsProcessorNoticeFunc(selected, processedSelectedIds, selected), "warning")
			} else {
				ShowMessage(&r, msgr.BulkActionNoAvailableRecords, "warning")
			}
			return
		}
	} else {
		processedSelectedIds = selected
	}

	b.mb.p.dialog(
		&r,
		b.bulkPanel(bulk, selected, processedSelectedIds, ctx),
		bulk.dialogWidth,
	)
	return
}

func (b *ListingBuilder) doBulkAction(ctx *web.EventContext) (r web.EventResponse, err error) {
	bulk := getBulkAction(b.bulkActions, ctx.R.FormValue(ParamBulkActionName))
	if bulk == nil {
		panic("bulk required")
	}

	if b.mb.Info().Verifier().SnakeDo(PermBulkActions, bulk.name).WithReq(ctx.R).IsAllowed() != nil {
		ShowMessage(&r, perm.PermissionDenied.Error(), "warning")
		return
	}

	selectedIds := getSelectedIds(ctx)

	var err1 error
	var processedSelectedIds []string
	if bulk.selectedIdsProcessorFunc != nil {
		processedSelectedIds, err1 = bulk.selectedIdsProcessorFunc(selectedIds, ctx)
	} else {
		processedSelectedIds = selectedIds
	}

	if err1 == nil {
		err1 = bulk.updateFunc(processedSelectedIds, ctx)
	}

	if err1 != nil {
		if _, ok := err1.(*web.ValidationErrors); !ok {
			vErr := &web.ValidationErrors{}
			vErr.GlobalError(err1.Error())
			ctx.Flash = vErr
		}
	}

	if ctx.Flash != nil {
		r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
			Name: actions.Dialog.ContentPortalName(),
			Body: b.bulkPanel(bulk, selectedIds, processedSelectedIds, ctx),
		})
		return
	}

	msgr := MustGetMessages(ctx.R)
	ShowMessage(&r, msgr.SuccessfullyUpdated, "")
	if isInDialogFromQuery(ctx) {
		qs := ctx.Queries()
		qs.Del(bulkPanelOpenParamName)
		qs.Del(ParamBulkActionName)
		web.AppendRunScripts(&r,
			closeDialogVarScript,
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

func (b *ListingBuilder) doListingAction(ctx *web.EventContext) (r web.EventResponse, err error) {
	action := getAction(b.actions, ctx.R.FormValue(ParamListingActionName))
	if action == nil {
		panic("action required")
	}

	if b.mb.Info().Verifier().SnakeDo(PermDoListingAction, action.name).WithReq(ctx.R).IsAllowed() != nil {
		ShowMessage(&r, perm.PermissionDenied.Error(), "warning")
		return
	}

	if err := action.updateFunc("", ctx); err != nil {
		if _, ok := err.(*web.ValidationErrors); !ok {
			vErr := &web.ValidationErrors{}
			vErr.GlobalError(err.Error())
			ctx.Flash = vErr
		}
	}

	portalID := ctx.R.FormValue(ParamPortalID)

	if ctx.Flash != nil {
		r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
			Name: actions.Dialog.ContentPortalName() + portalID,
			Body: b.actionPanel(action, ctx),
		})
		return
	}

	msgr := MustGetMessages(ctx.R)
	ShowMessage(&r, msgr.SuccessfullyUpdated, "")

	if isInDialogFromQuery(ctx) {
		qs := ctx.Queries()
		qs.Del(actionPanelOpenParamName)
		qs.Del(ParamListingActionName)
		web.AppendRunScripts(&r,
			closeDialogVarScript,
			web.Plaid().
				URL(ctx.R.RequestURI).
				EventFunc(actions.UpdateListingDialog).
				Queries(qs).
				Query(ParamPortalID, ctx.R.FormValue(ParamPortalID)).
				Go(),
		)
	} else {
		r.PushState = web.Location(url.Values{actionPanelOpenParamName: []string{}}).MergeQuery(true)
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
	for i, tab := range tabsData {
		if tab.ID == "" {
			tab.ID = fmt.Sprintf("tab%d", i)
		}
	}
	value := -1
	activeTabValue := qs.Get(ActiveFilterTabQueryKey)

	for i, td := range tabsData {
		// Find selected tab by active_filter_tab=xx in the url query
		if activeTabValue == td.ID {
			value = i
		}

		tabContent := h.Text(td.Label)
		if td.AdvancedLabel != nil {
			tabContent = td.AdvancedLabel
		}

		totalQuery := url.Values{}
		totalQuery.Set(ActiveFilterTabQueryKey, td.ID)
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
) (btn h.HTMLComponent, displaySortedFields []*FieldBuilder) {
	var (
		_, respath         = path.Split(pageURL.Path)
		displayColumnsName = fmt.Sprintf("%s_display_columns", respath)
		sortedColumnsName  = fmt.Sprintf("%s_sorted_columns", respath)
		originalColumns    []string
		displayColumns     []string
		sortedColumns      []string
	)

	b.fields.EachHavesComponent(func(f *FieldBuilder) bool {
		if b.mb.Info().Verifier().Do(PermList).SnakeOn("f_"+f.name).WithReq(ctx.R).IsAllowed() == nil {
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
		displaySortedFields = b.fields
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
			Label: i18n.PT(ctx.R, ModelsI18nModuleKey, b.mb.label, b.mb.getLabel(b.Field(sc).NameLabel)),
		})
	}

	msgr := MustGetMessages(ctx.R)
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
		VSlot("{ locals }").Init(fmt.Sprintf(`{selectColumnsMenu: false,...%s}`, h.JSONString(selectColumns)))
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
	ft.Clear = msgr.FiltersClear
	ft.Add = msgr.FiltersAdd
	ft.Apply = msgr.FilterApply
	for _, d := range fd {
		d.Translations = vx.FilterIndependentTranslations{
			FilterBy: msgr.FilterBy(d.Label),
		}
	}

	ft.Date.To = msgr.FiltersDateTo

	ft.Number.And = msgr.FiltersNumberAnd
	ft.Number.Equals = msgr.FiltersNumberEquals
	ft.Number.Between = msgr.FiltersNumberBetween
	ft.Number.GreaterThan = msgr.FiltersNumberGreaterThan
	ft.Number.LessThan = msgr.FiltersNumberLessThan

	ft.String.Equals = msgr.FiltersStringEquals
	ft.String.Contains = msgr.FiltersStringContains

	ft.MultipleSelect.In = msgr.FiltersMultipleSelectIn
	ft.MultipleSelect.NotIn = msgr.FiltersMultipleSelectNotIn

	filter := vx.VXFilter(fd).Translations(ft)
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

type ColOrderBy struct {
	FieldName string
	// ASC, DESC
	OrderBy string
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
			OrderBy:   ss[ssl-1],
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
			if oob.OrderBy == "ASC" {
				newOrderBysQueryValue = append(newOrderBysQueryValue, oob.FieldName+"_DESC")
			}
			continue
		}
		newOrderBysQueryValue = append(newOrderBysQueryValue, oob.FieldName+"_"+oob.OrderBy)
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
	lcb := b.listingComponentBuilderCtx(ctx).SetSelection(true)
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
		Respond(&r, web.Scope(body).VSlot("{ form }"))
	return
}

func (b *ListingBuilder) openListingDialog(ctx *web.EventContext) (r web.EventResponse, err error) {
	lcb := b.listingComponentBuilderCtx(ctx)
	ctx.R.Form.Set(ParamOverlay, actions.Dialog.String())

	targetPortal := ctx.R.FormValue(ParamTargetPortal)

	var body h.HTMLComponent
	if body, err = lcb.Build(ctx); err != nil {
		return
	}
	body = web.Scope(body).VSlot("{ form, closer }")

	b.mb.p.DialogPortal(targetPortal).
		SetValidWidth(b.dialogWidth).
		SetValidHeight(b.dialogHeight).
		SetContentPortalName(lcb.portals.Main()).
		SetScrollable(true).
		Respond(&r, body)
	return
}

func (b *ListingBuilder) updateListingDialog(ctx *web.EventContext) (r web.EventResponse, err error) {
	lcb := b.listingComponentBuilderCtx(ctx)
	ctx.R.Form.Set(ParamOverlay, actions.Dialog.String())
	var dataTable, dataTableAdditions h.HTMLComponent

	if dataTable, dataTableAdditions, err = lcb.getTableComponents(ctx); err != nil {
		return
	}

	r.UpdatePortals = append(r.UpdatePortals,
		&web.PortalUpdate{
			Name: lcb.portals.DataTable(),
			Body: dataTable,
		}, &web.PortalUpdate{
			Name: lcb.portals.DataTableAdditions(),
			Body: dataTableAdditions,
		},
	)
	return
}
