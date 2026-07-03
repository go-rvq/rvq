package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VItemBuilder struct {
	VTagBuilder[*VItemBuilder]
}

func VItem(children ...h.HTMLComponent) *VItemBuilder {
	return VTag(&VItemBuilder{}, "v-item", children...)
}

func (b *VItemBuilder) Value(v interface{}) (r *VItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VItemBuilder) Disabled(v bool) (r *VItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VItemBuilder) SelectedClass(v string) (r *VItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VItemBuilder) On(name string, value string) (r *VItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VItemBuilder) Bind(name string, value string) (r *VItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
