package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableRowsBuilder struct {
	VTagBuilder[*VDataTableRowsBuilder]
}

func VDataTableRows(children ...h.HTMLComponent) *VDataTableRowsBuilder {
	return VTag(&VDataTableRowsBuilder{}, "v-data-table-rows", children...)
}

func (b *VDataTableRowsBuilder) Items(v interface{}) (r *VDataTableRowsBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VDataTableRowsBuilder) HideNoData(v bool) (r *VDataTableRowsBuilder) {
	b.Attr(":hide-no-data", fmt.Sprint(v))
	return b
}

func (b *VDataTableRowsBuilder) NoDataText(v string) (r *VDataTableRowsBuilder) {
	b.Attr("no-data-text", v)
	return b
}

func (b *VDataTableRowsBuilder) Loading(v interface{}) (r *VDataTableRowsBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDataTableRowsBuilder) Mobile(v bool) (r *VDataTableRowsBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableRowsBuilder) CellProps(v interface{}) (r *VDataTableRowsBuilder) {
	b.Attr(":cell-props", h.JSONString(v))
	return b
}

func (b *VDataTableRowsBuilder) LoadingText(v string) (r *VDataTableRowsBuilder) {
	b.Attr("loading-text", v)
	return b
}

func (b *VDataTableRowsBuilder) MobileBreakpoint(v interface{}) (r *VDataTableRowsBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableRowsBuilder) RowProps(v interface{}) (r *VDataTableRowsBuilder) {
	b.Attr(":row-props", h.JSONString(v))
	return b
}

func (b *VDataTableRowsBuilder) On(name string, value string) (r *VDataTableRowsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableRowsBuilder) Bind(name string, value string) (r *VDataTableRowsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDataTableRowsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDataTableRowsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDataTableRowsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotDataTableGroup(child ...h.HTMLComponent) {
	b.SetSlot("data-table-group", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotDataTableGroup(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("data-table-group", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotDataTableGroup(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotDataTableGroup(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotDataTableGroup(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotDataTableGroup(scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("data-table-select", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("data-table-select", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotDataTableSelect(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotDataTableSelect(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotItemDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("item.data-table-select", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotItemDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item.data-table-select", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotItemDataTableSelect(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotItemDataTableSelect(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotItemDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotItemDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotItemDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("item.data-table-expand", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotItemDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item.data-table-expand", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotItemDataTableExpand(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotItemDataTableExpand(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotItemDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotItemDataTableExpand(scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotHeaderDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-select", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-select", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotHeaderDataTableSelect(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotHeaderDataTableSelect(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotHeaderDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotHeaderDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-expand", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-expand", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotHeaderDataTableExpand(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotHeaderDataTableExpand(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotHeaderDataTableExpand(scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotItem(child ...h.HTMLComponent) {
	b.SetSlot("item", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotItem(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotItem(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotItem(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotItem(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotItem(scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotLoading(child ...h.HTMLComponent) {
	b.SetSlot("loading", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotLoading(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loading", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotLoading(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotLoading(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotLoading(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotLoading(scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotGroupHeader(child ...h.HTMLComponent) {
	b.SetSlot("group-header", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotGroupHeader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("group-header", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotGroupHeader(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotGroupHeader(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotGroupHeader(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotGroupHeader(scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotNoData(child ...h.HTMLComponent) {
	b.SetSlot("no-data", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotNoData(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("no-data", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotNoData(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotNoData(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotNoData(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotNoData(scope, child...)
	return b
}

func (b *VDataTableRowsBuilder) SetSlotExpandedRow(child ...h.HTMLComponent) {
	b.SetSlot("expanded-row", child...)
}

func (b *VDataTableRowsBuilder) SetScopedSlotExpandedRow(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("expanded-row", scope, child...)
}

func (b *VDataTableRowsBuilder) SlotExpandedRow(child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetSlotExpandedRow(child...)
	return b
}

func (b *VDataTableRowsBuilder) ScopedSlotExpandedRow(scope string, child ...h.HTMLComponent) (r *VDataTableRowsBuilder) {
	b.SetScopedSlotExpandedRow(scope, child...)
	return b
}
