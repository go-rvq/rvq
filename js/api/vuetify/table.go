package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTableBuilder struct {
	VTagBuilder[*VTableBuilder]
}

func VTable(children ...h.HTMLComponent) *VTableBuilder {
	return VTag(&VTableBuilder{}, "v-table", children...)
}

func (b *VTableBuilder) Tag(v interface{}) (r *VTableBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VTableBuilder) Theme(v string) (r *VTableBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTableBuilder) Density(v interface{}) (r *VTableBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VTableBuilder) Height(v interface{}) (r *VTableBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VTableBuilder) FixedHeader(v bool) (r *VTableBuilder) {
	b.Attr(":fixed-header", fmt.Sprint(v))
	return b
}

func (b *VTableBuilder) FixedFooter(v bool) (r *VTableBuilder) {
	b.Attr(":fixed-footer", fmt.Sprint(v))
	return b
}

func (b *VTableBuilder) Hover(v bool) (r *VTableBuilder) {
	b.Attr(":hover", fmt.Sprint(v))
	return b
}

func (b *VTableBuilder) On(name string, value string) (r *VTableBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTableBuilder) Bind(name string, value string) (r *VTableBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTableBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTableBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTableBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTableBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTableBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTableBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTableBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VTableBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VTableBuilder) SlotDefault(child ...h.HTMLComponent) (r *VTableBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VTableBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VTableBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VTableBuilder) SetSlotTop(child ...h.HTMLComponent) {
	b.SetSlot("top", child...)
}

func (b *VTableBuilder) SetScopedSlotTop(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("top", scope, child...)
}

func (b *VTableBuilder) SlotTop(child ...h.HTMLComponent) (r *VTableBuilder) {
	b.SetSlotTop(child...)
	return b
}

func (b *VTableBuilder) ScopedSlotTop(scope string, child ...h.HTMLComponent) (r *VTableBuilder) {
	b.SetScopedSlotTop(scope, child...)
	return b
}

func (b *VTableBuilder) SetSlotBottom(child ...h.HTMLComponent) {
	b.SetSlot("bottom", child...)
}

func (b *VTableBuilder) SetScopedSlotBottom(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("bottom", scope, child...)
}

func (b *VTableBuilder) SlotBottom(child ...h.HTMLComponent) (r *VTableBuilder) {
	b.SetSlotBottom(child...)
	return b
}

func (b *VTableBuilder) ScopedSlotBottom(scope string, child ...h.HTMLComponent) (r *VTableBuilder) {
	b.SetScopedSlotBottom(scope, child...)
	return b
}

func (b *VTableBuilder) SetSlotWrapper(child ...h.HTMLComponent) {
	b.SetSlot("wrapper", child...)
}

func (b *VTableBuilder) SetScopedSlotWrapper(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("wrapper", scope, child...)
}

func (b *VTableBuilder) SlotWrapper(child ...h.HTMLComponent) (r *VTableBuilder) {
	b.SetSlotWrapper(child...)
	return b
}

func (b *VTableBuilder) ScopedSlotWrapper(scope string, child ...h.HTMLComponent) (r *VTableBuilder) {
	b.SetScopedSlotWrapper(scope, child...)
	return b
}
