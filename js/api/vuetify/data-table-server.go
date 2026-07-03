package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableServerBuilder struct {
	VTagBuilder[*VDataTableServerBuilder]
}

func VDataTableServer(children ...h.HTMLComponent) *VDataTableServerBuilder {
	return VTag(&VDataTableServerBuilder{}, "v-data-table-server", children...)
}

func (b *VDataTableServerBuilder) Search(v string) (r *VDataTableServerBuilder) {
	b.Attr("search", v)
	return b
}

func (b *VDataTableServerBuilder) Tag(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Theme(v string) (r *VDataTableServerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDataTableServerBuilder) Items(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) HideNoData(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) SelectStrategy(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Density(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Height(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Width(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ItemValue(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ReturnObject(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) ValueComparator(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Color(v string) (r *VDataTableServerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDataTableServerBuilder) ModelValue(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) NoDataText(v string) (r *VDataTableServerBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VDataTableServerBuilder) Loading(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Mobile(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) Page(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":page", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ItemSelectable(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":item-selectable", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ShowSelect(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":show-select", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) SortBy(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":sort-by", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) MultiSort(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":multi-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) MustSort(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":must-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) CustomKeySort(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":custom-key-sort", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ItemsPerPage(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":items-per-page", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ExpandOnClick(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":expand-on-click", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) ShowExpand(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":show-expand", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) Expanded(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":expanded", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) GroupBy(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":group-by", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) HeaderProps(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":header-props", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) CellProps(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":cell-props", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ItemsLength(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":items-length", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) DisableSort(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":disable-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) Headers(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":headers", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) LoadingText(v string) (r *VDataTableServerBuilder) {
	b.Attr("loading-text", v)
	return b
}

func (b *VDataTableServerBuilder) MobileBreakpoint(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) RowProps(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":row-props", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) HideDefaultBody(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":hide-default-body", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) HideDefaultFooter(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":hide-default-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) HideDefaultHeader(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":hide-default-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) FixedHeader(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":fixed-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) SortAscIcon(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":sort-asc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) SortDescIcon(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":sort-desc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) Sticky(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) FixedFooter(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":fixed-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) Hover(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) PrevIcon(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":prev-icon", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) NextIcon(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":next-icon", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) FirstIcon(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":first-icon", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) LastIcon(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":last-icon", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ItemsPerPageText(v string) (r *VDataTableServerBuilder) {
	b.Attr("items-per-page-text", v)
	return b
}

func (b *VDataTableServerBuilder) PageText(v string) (r *VDataTableServerBuilder) {
	b.Attr("page-text", v)
	return b
}

func (b *VDataTableServerBuilder) FirstPageLabel(v string) (r *VDataTableServerBuilder) {
	b.Attr("first-page-label", v)
	return b
}

func (b *VDataTableServerBuilder) PrevPageLabel(v string) (r *VDataTableServerBuilder) {
	b.Attr("prev-page-label", v)
	return b
}

func (b *VDataTableServerBuilder) NextPageLabel(v string) (r *VDataTableServerBuilder) {
	b.Attr("next-page-label", v)
	return b
}

func (b *VDataTableServerBuilder) LastPageLabel(v string) (r *VDataTableServerBuilder) {
	b.Attr("last-page-label", v)
	return b
}

func (b *VDataTableServerBuilder) ItemsPerPageOptions(v interface{}) (r *VDataTableServerBuilder) {
	b.Attr(":items-per-page-options", h.JSONString(v))
	return b
}

func (b *VDataTableServerBuilder) ShowCurrentPage(v bool) (r *VDataTableServerBuilder) {
	b.Attr(":show-current-page", fmt.Sprint(v))
	return b
}

func (b *VDataTableServerBuilder) On(name string, value string) (r *VDataTableServerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableServerBuilder) Bind(name string, value string) (r *VDataTableServerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDataTableServerBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDataTableServerBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDataTableServerBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotDataTableGroup(child ...h.HTMLComponent) {
	b.SetSlot("data-table-group", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotDataTableGroup(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("data-table-group", scope, child...)
}

func (b *VDataTableServerBuilder) SlotDataTableGroup(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotDataTableGroup(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotDataTableGroup(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotDataTableGroup(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("data-table-select", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("data-table-select", scope, child...)
}

func (b *VDataTableServerBuilder) SlotDataTableSelect(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotDataTableSelect(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotItemDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("item.data-table-select", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotItemDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item.data-table-select", scope, child...)
}

func (b *VDataTableServerBuilder) SlotItemDataTableSelect(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotItemDataTableSelect(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotItemDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotItemDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotItemDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("item.data-table-expand", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotItemDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item.data-table-expand", scope, child...)
}

func (b *VDataTableServerBuilder) SlotItemDataTableExpand(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotItemDataTableExpand(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotItemDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotItemDataTableExpand(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotHeaderDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-select", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-select", scope, child...)
}

func (b *VDataTableServerBuilder) SlotHeaderDataTableSelect(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotHeaderDataTableSelect(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotHeaderDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotHeaderDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-expand", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-expand", scope, child...)
}

func (b *VDataTableServerBuilder) SlotHeaderDataTableExpand(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotHeaderDataTableExpand(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotHeaderDataTableExpand(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VDataTableServerBuilder) SlotItem(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotLoading(child ...h.HTMLComponent) {
	b.SetSlot("loading", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotLoading(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loading", scope, child...)
}

func (b *VDataTableServerBuilder) SlotLoading(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotLoading(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotLoading(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotLoading(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotGroupHeader(child ...h.HTMLComponent) {
	b.SetSlot("group-header", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotGroupHeader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("group-header", scope, child...)
}

func (b *VDataTableServerBuilder) SlotGroupHeader(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotGroupHeader(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotGroupHeader(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotGroupHeader(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotNoData(child ...h.HTMLComponent) {
	b.SetSlot("no-data", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotNoData(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("no-data", scope, child...)
}

func (b *VDataTableServerBuilder) SlotNoData(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotNoData(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotNoData(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotNoData(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotExpandedRow(child ...h.HTMLComponent) {
	b.SetSlot("expanded-row", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotExpandedRow(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("expanded-row", scope, child...)
}

func (b *VDataTableServerBuilder) SlotExpandedRow(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotExpandedRow(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotExpandedRow(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotExpandedRow(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotHeaders(child ...h.HTMLComponent) {
	b.SetSlot("headers", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotHeaders(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("headers", scope, child...)
}

func (b *VDataTableServerBuilder) SlotHeaders(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotHeaders(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotHeaders(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotHeaders(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VDataTableServerBuilder) SlotLoader(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDataTableServerBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotColgroup(child ...h.HTMLComponent) {
	b.SetSlot("colgroup", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotColgroup(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("colgroup", scope, child...)
}

func (b *VDataTableServerBuilder) SlotColgroup(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotColgroup(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotColgroup(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotColgroup(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotTop(child ...h.HTMLComponent) {
	b.SetSlot("top", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotTop(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("top", scope, child...)
}

func (b *VDataTableServerBuilder) SlotTop(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotTop(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotTop(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotTop(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotBody(child ...h.HTMLComponent) {
	b.SetSlot("body", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotBody(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("body", scope, child...)
}

func (b *VDataTableServerBuilder) SlotBody(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotBody(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotBody(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotBody(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotTbody(child ...h.HTMLComponent) {
	b.SetSlot("tbody", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotTbody(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("tbody", scope, child...)
}

func (b *VDataTableServerBuilder) SlotTbody(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotTbody(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotTbody(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotTbody(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotThead(child ...h.HTMLComponent) {
	b.SetSlot("thead", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotThead(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("thead", scope, child...)
}

func (b *VDataTableServerBuilder) SlotThead(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotThead(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotThead(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotThead(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotTfoot(child ...h.HTMLComponent) {
	b.SetSlot("tfoot", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotTfoot(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("tfoot", scope, child...)
}

func (b *VDataTableServerBuilder) SlotTfoot(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotTfoot(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotTfoot(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotTfoot(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotBottom(child ...h.HTMLComponent) {
	b.SetSlot("bottom", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotBottom(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("bottom", scope, child...)
}

func (b *VDataTableServerBuilder) SlotBottom(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotBottom(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotBottom(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotBottom(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotBodyPrepend(child ...h.HTMLComponent) {
	b.SetSlot("body.prepend", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotBodyPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("body.prepend", scope, child...)
}

func (b *VDataTableServerBuilder) SlotBodyPrepend(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotBodyPrepend(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotBodyPrepend(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotBodyPrepend(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotBodyAppend(child ...h.HTMLComponent) {
	b.SetSlot("body.append", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotBodyAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("body.append", scope, child...)
}

func (b *VDataTableServerBuilder) SlotBodyAppend(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotBodyAppend(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotBodyAppend(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotBodyAppend(scope, child...)
	return b
}

func (b *VDataTableServerBuilder) SetSlotFooterPrepend(child ...h.HTMLComponent) {
	b.SetSlot("footer.prepend", child...)
}

func (b *VDataTableServerBuilder) SetScopedSlotFooterPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("footer.prepend", scope, child...)
}

func (b *VDataTableServerBuilder) SlotFooterPrepend(child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetSlotFooterPrepend(child...)
	return b
}

func (b *VDataTableServerBuilder) ScopedSlotFooterPrepend(scope string, child ...h.HTMLComponent) (r *VDataTableServerBuilder) {
	b.SetScopedSlotFooterPrepend(scope, child...)
	return b
}
