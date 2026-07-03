package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDataTableHeadersBuilder struct {
	VTagBuilder[*VDataTableHeadersBuilder]
}

func VDataTableHeaders(children ...h.HTMLComponent) *VDataTableHeadersBuilder {
	return VTag(&VDataTableHeadersBuilder{}, "v-data-table-headers", children...)
}

func (b *VDataTableHeadersBuilder) Color(v string) (r *VDataTableHeadersBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDataTableHeadersBuilder) Loading(v interface{}) (r *VDataTableHeadersBuilder) {
	b.Attr(":loading", h.JSONString(v))
	return b
}

func (b *VDataTableHeadersBuilder) Mobile(v bool) (r *VDataTableHeadersBuilder) {
	b.Attr(":mobile", fmt.Sprint(v))
	return b
}

func (b *VDataTableHeadersBuilder) MultiSort(v bool) (r *VDataTableHeadersBuilder) {
	b.Attr(":multi-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableHeadersBuilder) HeaderProps(v interface{}) (r *VDataTableHeadersBuilder) {
	b.Attr(":header-props", h.JSONString(v))
	return b
}

func (b *VDataTableHeadersBuilder) DisableSort(v bool) (r *VDataTableHeadersBuilder) {
	b.Attr(":disable-sort", fmt.Sprint(v))
	return b
}

func (b *VDataTableHeadersBuilder) MobileBreakpoint(v interface{}) (r *VDataTableHeadersBuilder) {
	b.Attr(":mobile-breakpoint", h.JSONString(v))
	return b
}

func (b *VDataTableHeadersBuilder) FixedHeader(v bool) (r *VDataTableHeadersBuilder) {
	b.Attr(":fixed-header", fmt.Sprint(v))
	return b
}

func (b *VDataTableHeadersBuilder) SortAscIcon(v interface{}) (r *VDataTableHeadersBuilder) {
	b.Attr(":sort-asc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableHeadersBuilder) SortDescIcon(v interface{}) (r *VDataTableHeadersBuilder) {
	b.Attr(":sort-desc-icon", h.JSONString(v))
	return b
}

func (b *VDataTableHeadersBuilder) Sticky(v bool) (r *VDataTableHeadersBuilder) {
	b.Attr(":sticky", fmt.Sprint(v))
	return b
}

func (b *VDataTableHeadersBuilder) On(name string, value string) (r *VDataTableHeadersBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDataTableHeadersBuilder) Bind(name string, value string) (r *VDataTableHeadersBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDataTableHeadersBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDataTableHeadersBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDataTableHeadersBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDataTableHeadersBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDataTableHeadersBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDataTableHeadersBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDataTableHeadersBuilder) SetSlotHeaders(child ...h.HTMLComponent) {
	b.SetSlot("headers", child...)
}

func (b *VDataTableHeadersBuilder) SetScopedSlotHeaders(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("headers", scope, child...)
}

func (b *VDataTableHeadersBuilder) SlotHeaders(child ...h.HTMLComponent) (r *VDataTableHeadersBuilder) {
	b.SetSlotHeaders(child...)
	return b
}

func (b *VDataTableHeadersBuilder) ScopedSlotHeaders(scope string, child ...h.HTMLComponent) (r *VDataTableHeadersBuilder) {
	b.SetScopedSlotHeaders(scope, child...)
	return b
}

func (b *VDataTableHeadersBuilder) SetSlotLoader(child ...h.HTMLComponent) {
	b.SetSlot("loader", child...)
}

func (b *VDataTableHeadersBuilder) SetScopedSlotLoader(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("loader", scope, child...)
}

func (b *VDataTableHeadersBuilder) SlotLoader(child ...h.HTMLComponent) (r *VDataTableHeadersBuilder) {
	b.SetSlotLoader(child...)
	return b
}

func (b *VDataTableHeadersBuilder) ScopedSlotLoader(scope string, child ...h.HTMLComponent) (r *VDataTableHeadersBuilder) {
	b.SetScopedSlotLoader(scope, child...)
	return b
}

func (b *VDataTableHeadersBuilder) SetSlotHeaderDataTableSelect(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-select", child...)
}

func (b *VDataTableHeadersBuilder) SetScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-select", scope, child...)
}

func (b *VDataTableHeadersBuilder) SlotHeaderDataTableSelect(child ...h.HTMLComponent) (r *VDataTableHeadersBuilder) {
	b.SetSlotHeaderDataTableSelect(child...)
	return b
}

func (b *VDataTableHeadersBuilder) ScopedSlotHeaderDataTableSelect(scope string, child ...h.HTMLComponent) (r *VDataTableHeadersBuilder) {
	b.SetScopedSlotHeaderDataTableSelect(scope, child...)
	return b
}

func (b *VDataTableHeadersBuilder) SetSlotHeaderDataTableExpand(child ...h.HTMLComponent) {
	b.SetSlot("header.data-table-expand", child...)
}

func (b *VDataTableHeadersBuilder) SetScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("header.data-table-expand", scope, child...)
}

func (b *VDataTableHeadersBuilder) SlotHeaderDataTableExpand(child ...h.HTMLComponent) (r *VDataTableHeadersBuilder) {
	b.SetSlotHeaderDataTableExpand(child...)
	return b
}

func (b *VDataTableHeadersBuilder) ScopedSlotHeaderDataTableExpand(scope string, child ...h.HTMLComponent) (r *VDataTableHeadersBuilder) {
	b.SetScopedSlotHeaderDataTableExpand(scope, child...)
	return b
}
