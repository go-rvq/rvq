package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableRowBuilder struct {
	VTagBuilder[*VDataTableRowBuilder]
}

func VDataTableRow(children ...h.HTMLComponent) *VDataTableRowBuilder {
	return VTag(&VDataTableRowBuilder{}, "v-data-table-row", children...)
}

func (b *VDataTableRowBuilder) Item(v interface{}) (r *VDataTableRowBuilder) {
	b.Attr(":item", h.JSONString(v))
	return b
}

func (b *VDataTableRowBuilder) Mobile(v bool) (r *VDataTableRowBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableRowBuilder) CellProps(v interface{}) (r *VDataTableRowBuilder) {
	b.Attr(":cell-props", h.JSONString(v))
	return b
}

func (b *VDataTableRowBuilder) MobileBreakpoint(v interface{}) (r *VDataTableRowBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableRowBuilder) Index(v interface{}) (r *VDataTableRowBuilder) {
	b.Attr(":index", h.JSONString(v))
	return b
}

func (b *VDataTableRowBuilder) On(name string, value string) (r *VDataTableRowBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableRowBuilder) Bind(name string, value string) (r *VDataTableRowBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDataTableRowBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDataTableRowBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDataTableRowBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDataTableRowBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDataTableRowBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDataTableRowBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDataTableRowBuilder) SetSlotItemDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("item.data-table-select", child...)
}

func (b *VDataTableRowBuilder) SetScopedSlotItemDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item.data-table-select", scope, child...)
}

func (b *VDataTableRowBuilder) SlotItemDataTableSelect(child ...h.HTMLComponent) (r *VDataTableRowBuilder) {
	b.SetSlotItemDataTableSelect(child...)
	return b
}

func (b *VDataTableRowBuilder) ScopedSlotItemDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableRowBuilder) {
	b.SetScopedSlotItemDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableRowBuilder) SetSlotItemDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("item.data-table-expand", child...)
}

func (b *VDataTableRowBuilder) SetScopedSlotItemDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("item.data-table-expand", scope, child...)
}

func (b *VDataTableRowBuilder) SlotItemDataTableExpand(child ...h.HTMLComponent) (r *VDataTableRowBuilder) {
	b.SetSlotItemDataTableExpand(child...)
	return b
}

func (b *VDataTableRowBuilder) ScopedSlotItemDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableRowBuilder) {
	b.SetScopedSlotItemDataTableExpand(scope, child...)
	return b
}

func (b *VDataTableRowBuilder) SetSlotHeaderDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-select", child...)
}

func (b *VDataTableRowBuilder) SetScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-select", scope, child...)
}

func (b *VDataTableRowBuilder) SlotHeaderDataTableSelect(child ...h.HTMLComponent) (r *VDataTableRowBuilder) {
	b.SetSlotHeaderDataTableSelect(child...)
	return b
}

func (b *VDataTableRowBuilder) ScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableRowBuilder) {
	b.SetScopedSlotHeaderDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableRowBuilder) SetSlotHeaderDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-expand", child...)
}

func (b *VDataTableRowBuilder) SetScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-expand", scope, child...)
}

func (b *VDataTableRowBuilder) SlotHeaderDataTableExpand(child ...h.HTMLComponent) (r *VDataTableRowBuilder) {
	b.SetSlotHeaderDataTableExpand(child...)
	return b
}

func (b *VDataTableRowBuilder) ScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableRowBuilder) {
	b.SetScopedSlotHeaderDataTableExpand(scope, child...)
	return b
}
