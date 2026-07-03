package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableBuilder struct {
	VTagBuilder[*VDataTableBuilder]
}

func VDataTable(children ...h.HTMLComponent) *VDataTableBuilder {
	return VTag(&VDataTableBuilder{}, "v-data-table", children...)
}

func (b *VDataTableBuilder) Search(v string) (r *VDataTableBuilder) {
	b.Attr("search", v)
	return b
}

func (b *VDataTableBuilder) Tag(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Theme(v string) (r *VDataTableBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDataTableBuilder) Items(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) FilterMode(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":filter-mode", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) NoFilter(v bool) (r *VDataTableBuilder) {
	b.Attr(":no-filter", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) CustomFilter(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":custom-filter", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) CustomKeyFilter(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":custom-key-filter", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) FilterKeys(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":filter-keys", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) HideNoData(v bool) (r *VDataTableBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) SelectStrategy(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Density(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Height(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Width(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ItemValue(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ReturnObject(v bool) (r *VDataTableBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) ValueComparator(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Color(v string) (r *VDataTableBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDataTableBuilder) ModelValue(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) NoDataText(v string) (r *VDataTableBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VDataTableBuilder) Loading(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Mobile(v bool) (r *VDataTableBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) Page(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":page", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ItemSelectable(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":item-selectable", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ShowSelect(v bool) (r *VDataTableBuilder) {
	b.Attr(":show-select", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) SortBy(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":sort-by", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) MultiSort(v bool) (r *VDataTableBuilder) {
	b.Attr(":multi-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) MustSort(v bool) (r *VDataTableBuilder) {
	b.Attr(":must-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) CustomKeySort(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":custom-key-sort", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ItemsPerPage(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":items-per-page", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ExpandOnClick(v bool) (r *VDataTableBuilder) {
	b.Attr(":expand-on-click", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) ShowExpand(v bool) (r *VDataTableBuilder) {
	b.Attr(":show-expand", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) Expanded(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":expanded", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) GroupBy(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":group-by", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) HeaderProps(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":header-props", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) CellProps(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":cell-props", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) DisableSort(v bool) (r *VDataTableBuilder) {
	b.Attr(":disable-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) Headers(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":headers", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) LoadingText(v string) (r *VDataTableBuilder) {
	b.Attr("loading-text", v)
	return b
}

func (b *VDataTableBuilder) MobileBreakpoint(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) RowProps(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":row-props", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) HideDefaultBody(v bool) (r *VDataTableBuilder) {
	b.Attr(":hide-default-body", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) HideDefaultFooter(v bool) (r *VDataTableBuilder) {
	b.Attr(":hide-default-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) HideDefaultHeader(v bool) (r *VDataTableBuilder) {
	b.Attr(":hide-default-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) FixedHeader(v bool) (r *VDataTableBuilder) {
	b.Attr(":fixed-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) SortAscIcon(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":sort-asc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) SortDescIcon(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":sort-desc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) Sticky(v bool) (r *VDataTableBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) FixedFooter(v bool) (r *VDataTableBuilder) {
	b.Attr(":fixed-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) Hover(v bool) (r *VDataTableBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) PrevIcon(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) NextIcon(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) FirstIcon(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":first-icon", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) LastIcon(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":last-icon", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ItemsPerPageText(v string) (r *VDataTableBuilder) {
	b.Attr("items-per-page-text", v)
	return b
}

func (b *VDataTableBuilder) PageText(v string) (r *VDataTableBuilder) {
	b.Attr("page-text", v)
	return b
}

func (b *VDataTableBuilder) FirstPageLabel(v string) (r *VDataTableBuilder) {
	b.Attr("first-page-label", v)
	return b
}

func (b *VDataTableBuilder) PrevPageLabel(v string) (r *VDataTableBuilder) {
	b.Attr("prev-page-label", v)
	return b
}

func (b *VDataTableBuilder) NextPageLabel(v string) (r *VDataTableBuilder) {
	b.Attr("next-page-label", v)
	return b
}

func (b *VDataTableBuilder) LastPageLabel(v string) (r *VDataTableBuilder) {
	b.Attr("last-page-label", v)
	return b
}

func (b *VDataTableBuilder) ItemsPerPageOptions(v interface{}) (r *VDataTableBuilder) {
	b.Attr(":items-per-page-options", h.JSONString(v))
	return b
}

func (b *VDataTableBuilder) ShowCurrentPage(v bool) (r *VDataTableBuilder) {
	b.Attr(":show-current-page", fmt.Sprint(v))
	return b
}

func (b *VDataTableBuilder) On(name string, value string) (r *VDataTableBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableBuilder) Bind(name string, value string) (r *VDataTableBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDataTableBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDataTableBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDataTableBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotDataTableGroup(child ...h.HTMLComponent) {
	b.SetSlot("data-table-group", child...)
}

func (b *VDataTableBuilder) SetScopedSlotDataTableGroup(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("data-table-group", scope, child...)
}

func (b *VDataTableBuilder) SlotDataTableGroup(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotDataTableGroup(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotDataTableGroup(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotDataTableGroup(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("data-table-select", child...)
}

func (b *VDataTableBuilder) SetScopedSlotDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("data-table-select", scope, child...)
}

func (b *VDataTableBuilder) SlotDataTableSelect(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotDataTableSelect(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotItemDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("item.data-table-select", child...)
}

func (b *VDataTableBuilder) SetScopedSlotItemDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item.data-table-select", scope, child...)
}

func (b *VDataTableBuilder) SlotItemDataTableSelect(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotItemDataTableSelect(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotItemDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotItemDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotItemDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("item.data-table-expand", child...)
}

func (b *VDataTableBuilder) SetScopedSlotItemDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item.data-table-expand", scope, child...)
}

func (b *VDataTableBuilder) SlotItemDataTableExpand(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotItemDataTableExpand(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotItemDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotItemDataTableExpand(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotHeaderDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-select", child...)
}

func (b *VDataTableBuilder) SetScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-select", scope, child...)
}

func (b *VDataTableBuilder) SlotHeaderDataTableSelect(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotHeaderDataTableSelect(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotHeaderDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotHeaderDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-expand", child...)
}

func (b *VDataTableBuilder) SetScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-expand", scope, child...)
}

func (b *VDataTableBuilder) SlotHeaderDataTableExpand(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotHeaderDataTableExpand(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotHeaderDataTableExpand(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VDataTableBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VDataTableBuilder) SlotItem(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotLoading(child ...h.HTMLComponent) {
	b.SetSlot("loading", child...)
}

func (b *VDataTableBuilder) SetScopedSlotLoading(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loading", scope, child...)
}

func (b *VDataTableBuilder) SlotLoading(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotLoading(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotLoading(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotLoading(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotGroupHeader(child ...h.HTMLComponent) {
	b.SetSlot("group-header", child...)
}

func (b *VDataTableBuilder) SetScopedSlotGroupHeader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("group-header", scope, child...)
}

func (b *VDataTableBuilder) SlotGroupHeader(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotGroupHeader(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotGroupHeader(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotGroupHeader(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotNoData(child ...h.HTMLComponent) {
	b.SetSlot("no-data", child...)
}

func (b *VDataTableBuilder) SetScopedSlotNoData(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("no-data", scope, child...)
}

func (b *VDataTableBuilder) SlotNoData(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotNoData(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotNoData(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotNoData(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotExpandedRow(child ...h.HTMLComponent) {
	b.SetSlot("expanded-row", child...)
}

func (b *VDataTableBuilder) SetScopedSlotExpandedRow(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("expanded-row", scope, child...)
}

func (b *VDataTableBuilder) SlotExpandedRow(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotExpandedRow(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotExpandedRow(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotExpandedRow(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotHeaders(child ...h.HTMLComponent) {
	b.SetSlot("headers", child...)
}

func (b *VDataTableBuilder) SetScopedSlotHeaders(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("headers", scope, child...)
}

func (b *VDataTableBuilder) SlotHeaders(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotHeaders(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotHeaders(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotHeaders(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VDataTableBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VDataTableBuilder) SlotLoader(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDataTableBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDataTableBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotColgroup(child ...h.HTMLComponent) {
	b.SetSlot("colgroup", child...)
}

func (b *VDataTableBuilder) SetScopedSlotColgroup(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("colgroup", scope, child...)
}

func (b *VDataTableBuilder) SlotColgroup(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotColgroup(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotColgroup(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotColgroup(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotTop(child ...h.HTMLComponent) {
	b.SetSlot("top", child...)
}

func (b *VDataTableBuilder) SetScopedSlotTop(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("top", scope, child...)
}

func (b *VDataTableBuilder) SlotTop(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotTop(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotTop(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotTop(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotBody(child ...h.HTMLComponent) {
	b.SetSlot("body", child...)
}

func (b *VDataTableBuilder) SetScopedSlotBody(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("body", scope, child...)
}

func (b *VDataTableBuilder) SlotBody(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotBody(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotBody(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotBody(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotTbody(child ...h.HTMLComponent) {
	b.SetSlot("tbody", child...)
}

func (b *VDataTableBuilder) SetScopedSlotTbody(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("tbody", scope, child...)
}

func (b *VDataTableBuilder) SlotTbody(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotTbody(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotTbody(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotTbody(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotThead(child ...h.HTMLComponent) {
	b.SetSlot("thead", child...)
}

func (b *VDataTableBuilder) SetScopedSlotThead(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("thead", scope, child...)
}

func (b *VDataTableBuilder) SlotThead(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotThead(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotThead(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotThead(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotTfoot(child ...h.HTMLComponent) {
	b.SetSlot("tfoot", child...)
}

func (b *VDataTableBuilder) SetScopedSlotTfoot(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("tfoot", scope, child...)
}

func (b *VDataTableBuilder) SlotTfoot(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotTfoot(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotTfoot(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotTfoot(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotBottom(child ...h.HTMLComponent) {
	b.SetSlot("bottom", child...)
}

func (b *VDataTableBuilder) SetScopedSlotBottom(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("bottom", scope, child...)
}

func (b *VDataTableBuilder) SlotBottom(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotBottom(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotBottom(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotBottom(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotBodyPrepend(child ...h.HTMLComponent) {
	b.SetSlot("body.prepend", child...)
}

func (b *VDataTableBuilder) SetScopedSlotBodyPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("body.prepend", scope, child...)
}

func (b *VDataTableBuilder) SlotBodyPrepend(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotBodyPrepend(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotBodyPrepend(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotBodyPrepend(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotBodyAppend(child ...h.HTMLComponent) {
	b.SetSlot("body.append", child...)
}

func (b *VDataTableBuilder) SetScopedSlotBodyAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("body.append", scope, child...)
}

func (b *VDataTableBuilder) SlotBodyAppend(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotBodyAppend(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotBodyAppend(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotBodyAppend(scope, child...)
	return b
}

func (b *VDataTableBuilder) SetSlotFooterPrepend(child ...h.HTMLComponent) {
	b.SetSlot("footer.prepend", child...)
}

func (b *VDataTableBuilder) SetScopedSlotFooterPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("footer.prepend", scope, child...)
}

func (b *VDataTableBuilder) SlotFooterPrepend(child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetSlotFooterPrepend(child...)
	return b
}

func (b *VDataTableBuilder) ScopedSlotFooterPrepend(scope string, child ...h.HTMLComponent) (r *VDataTableBuilder) {
	b.SetScopedSlotFooterPrepend(scope, child...)
	return b
}
