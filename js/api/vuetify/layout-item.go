package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VLayoutItemBuilder struct {
	VTagBuilder[*VLayoutItemBuilder]
}

func VLayoutItem(children ...h.HTMLComponent) *VLayoutItemBuilder {
	return VTag(&VLayoutItemBuilder{}, "v-layout-item", children...)
}

func (b *VLayoutItemBuilder) Name(v string) (r *VLayoutItemBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VLayoutItemBuilder) Size(v interface{}) (r *VLayoutItemBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VLayoutItemBuilder) Absolute(v bool) (r *VLayoutItemBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VLayoutItemBuilder) ModelValue(v bool) (r *VLayoutItemBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VLayoutItemBuilder) Order(v interface{}) (r *VLayoutItemBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VLayoutItemBuilder) Position(v interface{}) (r *VLayoutItemBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VLayoutItemBuilder) On(name string, value string) (r *VLayoutItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLayoutItemBuilder) Bind(name string, value string) (r *VLayoutItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VLayoutItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VLayoutItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VLayoutItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VLayoutItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VLayoutItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VLayoutItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VLayoutItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VLayoutItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VLayoutItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VLayoutItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VLayoutItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VLayoutItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
