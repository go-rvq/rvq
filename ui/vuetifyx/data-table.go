package vuetifyx

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
)

type (
	CellComponentFunc    func(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent
	CellWrapperFunc      func(cell h.MutableAttrHTMLComponent, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent
	HeadCellWrapperFunc  func(cell h.MutableAttrHTMLComponent, field string, title string, ctx *web.EventContext) h.HTMLComponent
	RowWrapperFunc       func(row h.MutableAttrHTMLComponent, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent
	RowMenuItemFunc      func(recordIndex int, obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent
	RowComponentFunc     func(obj interface{}, ctx *web.EventContext) h.HTMLComponent
	OnSelectFunc         func(id string, ctx *web.EventContext) string
	OnSelectAllFunc      func(idsOfPage []string, ctx *web.EventContext) string
	OnClearSelectionFunc func(ctx *web.EventContext) string
	RowStartFunc         func(rowIndex int, id string, obj interface{}, dataTableID string, ctx *web.EventContext) RowEndFunc
	RowEndFunc           func(row *h.HTMLTagBuilder)
)

type DataTableBuilder struct {
	data               interface{}
	selectable         bool
	withoutHeaders     bool
	selectionParamName string
	cellWrapper        CellWrapperFunc
	headCellWrapper    HeadCellWrapperFunc
	rowWrapper         RowWrapperFunc
	rowStart           RowStartFunc
	rowMenuItemFuncs   []RowMenuItemFunc
	rowExpandFunc      RowComponentFunc
	columns            []*DataTableColumnBuilder
	loadMoreCount      int
	loadMoreLabel      string
	loadMoreURL        string
	// e.g. {count} records are selected.
	selectedCountLabel   string
	clearSelectionLabel  string
	onClearSelectionFunc OnClearSelectionFunc
	tfootChildren        []h.HTMLComponent
	selectableColumnsBtn h.HTMLComponent
	onSelectFunc         OnSelectFunc
	onSelectAllFunc      OnSelectAllFunc
	density              string
}

func DataTable(data interface{}) (r *DataTableBuilder) {
	r = &DataTableBuilder{
		data:                data,
		selectionParamName:  "selected_ids",
		selectedCountLabel:  "{count} records are selected.",
		clearSelectionLabel: "clear selection",
	}
	return
}

func (b *DataTableBuilder) LoadMoreAt(count int, label string) (r *DataTableBuilder) {
	b.loadMoreCount = count
	b.loadMoreLabel = label
	return b
}

func (b *DataTableBuilder) LoadMoreURL(url string) (r *DataTableBuilder) {
	b.loadMoreURL = url
	return b
}

func (b *DataTableBuilder) Tfoot(children ...h.HTMLComponent) (r *DataTableBuilder) {
	b.tfootChildren = children
	return b
}

func (b *DataTableBuilder) Selectable(v bool) (r *DataTableBuilder) {
	b.selectable = v
	return b
}

func (b *DataTableBuilder) Data(v interface{}) (r *DataTableBuilder) {
	b.data = v
	return b
}

func (b *DataTableBuilder) SelectionParamName(v string) (r *DataTableBuilder) {
	b.selectionParamName = v
	return b
}

func (b *DataTableBuilder) WithoutHeader(v bool) (r *DataTableBuilder) {
	b.withoutHeaders = v
	return b
}

func (b *DataTableBuilder) CellWrapperFunc(v CellWrapperFunc) (r *DataTableBuilder) {
	b.cellWrapper = v
	return b
}

func (b *DataTableBuilder) HeadCellWrapperFunc(v HeadCellWrapperFunc) (r *DataTableBuilder) {
	b.headCellWrapper = v
	return b
}

func (b *DataTableBuilder) RowWrapperFunc(v RowWrapperFunc) (r *DataTableBuilder) {
	b.rowWrapper = v
	return b
}

func (b *DataTableBuilder) RowStarter(v RowStartFunc) *DataTableBuilder {
	b.rowStart = v
	return b
}

func (b *DataTableBuilder) RowMenuItemFuncs(vs ...RowMenuItemFunc) (r *DataTableBuilder) {
	b.rowMenuItemFuncs = vs
	return b
}

func (b *DataTableBuilder) RowMenuItemFunc(v RowMenuItemFunc) (r *DataTableBuilder) {
	b.rowMenuItemFuncs = append(b.rowMenuItemFuncs, v)
	return b
}

func (b *DataTableBuilder) RowExpandFunc(v RowComponentFunc) (r *DataTableBuilder) {
	b.rowExpandFunc = v
	return b
}

func (b *DataTableBuilder) SelectedCountLabel(v string) (r *DataTableBuilder) {
	b.selectedCountLabel = v
	return b
}

func (b *DataTableBuilder) ClearSelectionLabel(v string) (r *DataTableBuilder) {
	b.clearSelectionLabel = v
	return b
}

func (b *DataTableBuilder) OnClearSelectionFunc(v OnClearSelectionFunc) (r *DataTableBuilder) {
	b.onClearSelectionFunc = v
	return b
}

func (b *DataTableBuilder) SelectableColumnsBtn(v h.HTMLComponent) (r *DataTableBuilder) {
	b.selectableColumnsBtn = v
	return b
}

func (b *DataTableBuilder) OnSelectAllFunc(v OnSelectAllFunc) (r *DataTableBuilder) {
	b.onSelectAllFunc = v
	return b
}

func (b *DataTableBuilder) OnSelectFunc(v OnSelectFunc) (r *DataTableBuilder) {
	b.onSelectFunc = v
	return b
}

func (b *DataTableBuilder) Density() string {
	return b.density
}

func (b *DataTableBuilder) SetDensity(density string) *DataTableBuilder {
	b.density = density
	return b
}

type primarySlugger interface {
	PrimarySlug() string
}

func (b *DataTableBuilder) MarshalHTML(c context.Context) (r []byte, err error) {
	ctx := web.MustGetEventContext(c)

	selected := getSelectedIds(ctx, b.selectionParamName)

	loadMoreVarName := "loadmore"
	expandVarName := "expand"
	selectedCountVarName := "selected_count"

	initContextLocalsMap := map[string]interface{}{
		selectedCountVarName: len(selected),
	}
	// map[obj_id]{rowMenus}
	objRowMenusMap := make(map[string][]h.HTMLComponent)
	var i int
	reflectutils.ForEach(b.data, func(obj interface{}) {
		id := ObjectID(obj)
		var opMenuItems []h.HTMLComponent
		for _, f := range b.rowMenuItemFuncs {
			item := f(i, obj, id, ctx)
			if item == nil {
				continue
			}
			opMenuItems = append(opMenuItems, item)
		}
		if len(opMenuItems) > 0 {
			objRowMenusMap[id] = opMenuItems
		}
		i++
	})

	hasRowMenuCol := len(objRowMenusMap) > 0 || b.selectableColumnsBtn != nil

	var rows []h.HTMLComponent
	var idsOfPage []string

	inPlaceLoadMore := b.loadMoreCount > 0 && b.loadMoreURL == ""

	hasExpand := b.rowExpandFunc != nil

	i = 0
	tdCount := 0
	haveMoreRecord := false
	reflectutils.ForEach(b.data, func(obj interface{}) {
		id := ObjectID(obj)
		idsOfPage = append(idsOfPage, id)

		inputValue := ""
		if slices.Contains(selected, id) {
			inputValue = id
		}
		var (
			row *h.HTMLTagBuilder
			tds []h.HTMLComponent
			end RowEndFunc
		)

		if b.rowStart != nil {
			end = b.rowStart(i, id, obj, "", ctx)
		}

		if hasExpand {
			initContextLocalsMap[fmt.Sprintf("%s_%d", expandVarName, i)] = false
			localsExpandVarName := fmt.Sprintf("locals.%s_%d", expandVarName, i)
			tds = append(tds, h.Td(
				v.VIcon("").
					Attr(":icon", fmt.Sprintf(`%s?"mdi-chevron-up-circle":"mdi-chevron-down"`, localsExpandVarName)).
					Attr(":class", fmt.Sprintf(`{"v-data-table__expand-icon--active": locals.%s_%d, "v-data-table__expand-icon": true}`, expandVarName, i)).
					On("click", fmt.Sprintf("locals.%s_%d = !locals.%s_%d", expandVarName, i, expandVarName, i)),
			).Class("pr-0").Style("width: 40px;"))
		}

		if b.selectable {
			onChange := web.Plaid().
				PushState(true).
				MergeQuery(true).
				Query(b.selectionParamName,
					web.Var(fmt.Sprintf(`{value: %s, add: $event, remove: !$event}`, h.JSONString(id))),
				).RunPushState()
			if b.onSelectFunc != nil {
				onChange = b.onSelectFunc(id, ctx)
			}

			tds = append(tds, h.Td(
				web.Scope(
					v.VCheckbox().
						Density(v.DensityCompact).
						Class("mt-0").
						TrueValue(id).
						FalseValue("").
						HideDetails(true).
						Attr("v-model", "itemLocals.inputValue").
						Attr("@update:model-value", onChange+";locals.selected_count+=($event?1:-1);"),
				).VSlot("{ locals: itemLocals }").Init(fmt.Sprintf(`{ inputValue :"%v"} `, inputValue)),
			).Class("pr-0"))
		}

		for _, f := range b.columns {
			tds = append(tds, f.cellComponentFunc(obj, f.name, ctx))
		}

		var bindTds []h.HTMLComponent
		for _, td := range tds {
			std, ok := td.(h.MutableAttrHTMLComponent)
			if !ok {
				bindTds = append(bindTds, td)
				continue
			}

			var tdWrapped h.HTMLComponent = std
			if b.cellWrapper != nil {
				tdWrapped = b.cellWrapper(std, id, obj, "", ctx)
			}

			bindTds = append(bindTds, tdWrapped)
		}

		if hasRowMenuCol {
			var td h.HTMLComponent
			rowMenus, ok := objRowMenusMap[id]
			if ok {
				td = h.Td(
					v.VMenu(
						web.Slot(
							v.VBtn("").Children(
								v.VIcon("mdi-dots-horizontal"),
							).Attr("v-bind", "props").Variant("text").Size("small"),
						).Name("activator").Scope("{ props }"),

						v.VList(
							rowMenus...,
						).Elevation(8),
					),
				).Style("width: 64px;").Class("pl-0")
			} else {
				td = h.Td().Style("width: 64px;").Class("pl-0")
			}
			bindTds = append(bindTds, td)
		}

		tdCount = len(bindTds)
		row = h.Tr(bindTds...)
		if b.loadMoreCount > 0 && i >= b.loadMoreCount {
			if len(b.loadMoreURL) > 0 {
				return
			} else {
				row.Attr("v-if", "locals.loadmore")
			}
			haveMoreRecord = true
		}

		if end != nil {
			end(row)
		}

		if b.rowWrapper != nil {
			rows = append(rows, b.rowWrapper(row, id, obj, "", ctx))
		} else {
			rows = append(rows, row)
		}

		if hasExpand {
			rows = append(rows,
				h.Tr(
					h.Td(
						v.VExpandTransition(
							h.Div(
								b.rowExpandFunc(obj, ctx),
								v.VDivider(),
							).Attr("v-if", fmt.Sprintf("locals.%s_%d", expandVarName, i)).
								Class("bg-grey-lighten-5"), // bg-grey-lighten-5 | grey lighten-5
						),
					).Attr("colspan", fmt.Sprint(tdCount)).Class("pa-0").Style("height: auto; border-bottom: none"),
				).Class("v-data-table__expand-row"),
			)
		}
		i++
	})

	var thead h.HTMLComponent

	if !b.withoutHeaders {

		var heads []h.HTMLComponent

		if hasExpand {
			heads = append(heads, h.Th(" "))
		}

		if b.selectable {
			allInputValue := ""
			idsOfPageComma := strings.Join(idsOfPage, ",")
			if allSelected(selected, idsOfPage) {
				allInputValue = idsOfPageComma
			}

			onChange := web.Plaid().
				PushState(true).
				MergeQuery(true).
				Query(b.selectionParamName,
					web.Var(fmt.Sprintf(`{value: %s, add: $event, remove: !$event}`,
						h.JSONString(idsOfPage))),
				).Go()
			if b.onSelectAllFunc != nil {
				onChange = b.onSelectAllFunc(idsOfPage, ctx)
			}
			heads = append(heads, h.Th("").Children(
				web.Scope(
					v.VCheckbox().
						Density(v.DensityCompact).
						Class("mt-0").
						TrueValue(idsOfPageComma).
						HideDetails(true).
						Attr("v-model", "itemLocals.allInputValue").
						Attr("@update:model-value", onChange),
				).VSlot("{ locals: itemLocals }").Init(fmt.Sprintf(`{ allInputValue :"%v"} `, allInputValue)),
			).Style("width: 48px;").Class("pr-0"))
		}

		for _, f := range b.columns {
			var head h.HTMLComponent
			th := h.Th(f.title)
			head = th
			if b.headCellWrapper != nil {
				head = b.headCellWrapper(
					(h.MutableAttrHTMLComponent)(th),
					f.name,
					f.title,
					ctx,
				)
			}
			heads = append(heads, head)
		}

		if hasRowMenuCol {
			heads = append(heads, h.Th("").Children(b.selectableColumnsBtn).Style("width: 64px;").Class("pl-0")) // Edit, Delete menu
		}
		thead = h.Thead(
			h.Tr(heads...),
		).Class("bg-grey-lighten-5")
	}

	var tfoot h.HTMLComponent
	if len(b.tfootChildren) > 0 {
		tfoot = h.Tfoot(b.tfootChildren...)
	}
	if b.loadMoreCount > 0 && haveMoreRecord {
		var btn h.HTMLComponent

		if inPlaceLoadMore {
			btn = v.VBtn(b.loadMoreLabel).
				Variant("text").
				Size("small").
				Class("mt-2").
				On("click", "locals.loadmore = !locals.loadmore")
		} else {
			btn = v.VBtn(b.loadMoreLabel).
				Variant("text").
				Size("small").
				Class("mt-2").
				Href(b.loadMoreURL)
		}

		tfoot = h.Tfoot(
			h.Tr(
				h.Td(
					h.If(!hasExpand, v.VDivider()),
					btn,
				).Class("text-center pa-0").Attr("colspan", fmt.Sprint(tdCount)),
			),
		).Attr("v-if", "!locals.loadmore")
	}

	var selectedCountNotice h.HTMLComponents
	onClearSelection := web.Plaid().
		MergeQuery(true).
		Query(b.selectionParamName, "").
		PushState(true).
		Go()
	if b.onClearSelectionFunc != nil {
		onClearSelection = b.onClearSelectionFunc(ctx)
	}
	{
		ss := strings.Split(b.selectedCountLabel, "{count}")
		if len(ss) == 1 {
			selectedCountNotice = append(selectedCountNotice, h.Text(ss[0]))
		} else {
			selectedCountNotice = append(selectedCountNotice,
				h.Text(ss[0]),
				h.Strong("{{locals.selected_count}}"),
				h.Text(ss[1]),
			)
		}
	}
	vtable := v.VTable(
		thead,
		h.Tbody(rows...),
		tfoot,
	)

	if b.density != "" {
		vtable.Density(b.density)
	}

	table := web.Scope(
		h.Div(
			selectedCountNotice,
			v.VBtn(b.clearSelectionLabel).
				Variant("plain").
				Size("small").
				On("click", onClearSelection),
		).
			Class("bg-grey-lighten-3 text-center pt-2 pb-2").
			Attr("v-show", "locals.selected_count > 0"),
		vtable,
	).VSlot("{ locals }").Init(fmt.Sprintf(` { selected_count : %v , loadmore : false }`, len(selected)))

	if inPlaceLoadMore {
		initContextLocalsMap[loadMoreVarName] = false
	}

	// if len(initContextLocalsMap) > 0 {
	//	table.AppendChildren(web.ObjectAssignTag("vars", initContextLocalsMap))
	// }

	return table.MarshalHTML(c)
}

func ObjectID(obj interface{}) (id string) {
	if slugger, ok := obj.(primarySlugger); ok {
		id = slugger.PrimarySlug()
	} else {
		defer func() {
			if r := recover(); r != nil {
				id = ""
				return
			}
		}()
		id = fmt.Sprint(reflectutils.MustGet(obj, "ID"))
	}
	if id == "0" {
		id = ""
	}
	return id
}

func getSelectedIds(ctx *web.EventContext, selectedParamName string) (selected []string) {
	selectedValue := ctx.R.URL.Query().Get(selectedParamName)
	if len(selectedValue) > 0 {
		selected = strings.Split(selectedValue, ",")
	}
	return selected
}

func allSelected(selectedInURL []string, pageSelected []string) bool {
	for _, ps := range pageSelected {
		if !slices.Contains(selectedInURL, ps) {
			return false
		}
	}
	return true
}

func (b *DataTableBuilder) Column(name string) (r *DataTableColumnBuilder) {
	r = &DataTableColumnBuilder{}
	for _, c := range b.columns {
		if c.name == name {
			return c
		}
	}

	r.Name(name).CellComponentFunc(defaultCellComponentFunc)
	b.columns = append(b.columns, r)
	return
}

func defaultCellComponentFunc(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent {
	return h.Td(h.Text(fmt.Sprint(reflectutils.MustGet(obj, fieldName))))
}

type DataTableColumnBuilder struct {
	name              string
	title             string
	cellComponentFunc CellComponentFunc
}

func (b *DataTableColumnBuilder) Name(v string) (r *DataTableColumnBuilder) {
	b.name = v
	return b
}

func (b *DataTableColumnBuilder) Title(v string) (r *DataTableColumnBuilder) {
	b.title = v
	return b
}

func (b *DataTableColumnBuilder) CellComponentFunc(v CellComponentFunc) (r *DataTableColumnBuilder) {
	b.cellComponentFunc = v
	return b
}
