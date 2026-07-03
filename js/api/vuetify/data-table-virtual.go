package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableVirtualBuilder struct {
	VTagBuilder[*VDataTableVirtualBuilder]
}

func VDataTableVirtual(children ...h.HTMLComponent) *VDataTableVirtualBuilder {
	return VTag(&VDataTableVirtualBuilder{}, "v-data-table-virtual", children...)
}

func (b *VDataTableVirtualBuilder) Search(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("search", v)
	return b
}

func (b *VDataTableVirtualBuilder) Tag(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Theme(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDataTableVirtualBuilder) Items(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) FilterMode(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":filter-mode", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) NoFilter(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":no-filter", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) CustomFilter(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":custom-filter", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) CustomKeyFilter(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":custom-key-filter", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) FilterKeys(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":filter-keys", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) HideNoData(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) SelectStrategy(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":select-strategy", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Density(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Height(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Width(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) ItemValue(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":item-value", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) ReturnObject(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":return-object", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) ValueComparator(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":value-comparator", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Color(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDataTableVirtualBuilder) ModelValue(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) NoDataText(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VDataTableVirtualBuilder) Loading(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Mobile(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) ItemSelectable(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":item-selectable", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) ShowSelect(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":show-select", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) SortBy(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":sort-by", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) MultiSort(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":multi-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) MustSort(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":must-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) CustomKeySort(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":custom-key-sort", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) ExpandOnClick(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":expand-on-click", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) ShowExpand(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":show-expand", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) Expanded(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":expanded", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) GroupBy(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":group-by", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) HeaderProps(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":header-props", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) CellProps(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":cell-props", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) DisableSort(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":disable-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) Headers(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":headers", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) LoadingText(v string) (r *VDataTableVirtualBuilder) {
	b.Attr("loading-text", v)
	return b
}

func (b *VDataTableVirtualBuilder) MobileBreakpoint(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) RowProps(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":row-props", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) HideDefaultBody(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":hide-default-body", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) HideDefaultHeader(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":hide-default-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) FixedHeader(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":fixed-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) SortAscIcon(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":sort-asc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) SortDescIcon(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":sort-desc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) Sticky(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) FixedFooter(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":fixed-footer", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) Hover(v bool) (r *VDataTableVirtualBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VDataTableVirtualBuilder) ItemHeight(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":item-height", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) ItemKey(v interface{}) (r *VDataTableVirtualBuilder) {
	b.Attr(":item-key", h.JSONString(v))
	return b
}

func (b *VDataTableVirtualBuilder) On(name string, value string) (r *VDataTableVirtualBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableVirtualBuilder) Bind(name string, value string) (r *VDataTableVirtualBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDataTableVirtualBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDataTableVirtualBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotDataTableGroup(child ...h.HTMLComponent) {
	b.SetSlot("data-table-group", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotDataTableGroup(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("data-table-group", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotDataTableGroup(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotDataTableGroup(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotDataTableGroup(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotDataTableGroup(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("data-table-select", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("data-table-select", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotDataTableSelect(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotDataTableSelect(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotItemDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("item.data-table-select", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotItemDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item.data-table-select", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotItemDataTableSelect(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotItemDataTableSelect(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotItemDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotItemDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotItemDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("item.data-table-expand", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotItemDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item.data-table-expand", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotItemDataTableExpand(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotItemDataTableExpand(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotItemDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotItemDataTableExpand(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotHeaderDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-select", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-select", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotHeaderDataTableSelect(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotHeaderDataTableSelect(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotHeaderDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotHeaderDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-expand", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-expand", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotHeaderDataTableExpand(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotHeaderDataTableExpand(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotHeaderDataTableExpand(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotItem(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotLoading(child ...h.HTMLComponent) {
	b.SetSlot("loading", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotLoading(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loading", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotLoading(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotLoading(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotLoading(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotLoading(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotGroupHeader(child ...h.HTMLComponent) {
	b.SetSlot("group-header", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotGroupHeader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("group-header", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotGroupHeader(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotGroupHeader(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotGroupHeader(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotGroupHeader(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotNoData(child ...h.HTMLComponent) {
	b.SetSlot("no-data", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotNoData(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("no-data", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotNoData(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotNoData(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotNoData(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotNoData(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotExpandedRow(child ...h.HTMLComponent) {
	b.SetSlot("expanded-row", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotExpandedRow(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("expanded-row", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotExpandedRow(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotExpandedRow(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotExpandedRow(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotExpandedRow(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotHeaders(child ...h.HTMLComponent) {
	b.SetSlot("headers", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotHeaders(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("headers", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotHeaders(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotHeaders(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotHeaders(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotHeaders(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotLoader(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotColgroup(child ...h.HTMLComponent) {
	b.SetSlot("colgroup", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotColgroup(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("colgroup", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotColgroup(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotColgroup(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotColgroup(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotColgroup(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotTop(child ...h.HTMLComponent) {
	b.SetSlot("top", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotTop(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("top", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotTop(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotTop(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotTop(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotTop(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotTbody(child ...h.HTMLComponent) {
	b.SetSlot("tbody", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotTbody(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("tbody", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotTbody(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotTbody(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotTbody(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotTbody(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotThead(child ...h.HTMLComponent) {
	b.SetSlot("thead", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotThead(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("thead", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotThead(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotThead(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotThead(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotThead(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotTfoot(child ...h.HTMLComponent) {
	b.SetSlot("tfoot", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotTfoot(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("tfoot", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotTfoot(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotTfoot(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotTfoot(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotTfoot(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotBottom(child ...h.HTMLComponent) {
	b.SetSlot("bottom", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotBottom(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("bottom", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotBottom(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotBottom(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotBottom(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotBottom(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotBodyPrepend(child ...h.HTMLComponent) {
	b.SetSlot("body.prepend", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotBodyPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("body.prepend", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotBodyPrepend(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotBodyPrepend(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotBodyPrepend(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotBodyPrepend(scope, child...)
	return b
}

func (b *VDataTableVirtualBuilder) SetSlotBodyAppend(child ...h.HTMLComponent) {
	b.SetSlot("body.append", child...)
}

func (b *VDataTableVirtualBuilder) SetScopedSlotBodyAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("body.append", scope, child...)
}

func (b *VDataTableVirtualBuilder) SlotBodyAppend(child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetSlotBodyAppend(child...)
	return b
}

func (b *VDataTableVirtualBuilder) ScopedSlotBodyAppend(scope string, child ...h.HTMLComponent) (r *VDataTableVirtualBuilder) {
	b.SetScopedSlotBodyAppend(scope, child...)
	return b
}
