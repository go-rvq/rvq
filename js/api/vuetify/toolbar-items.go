package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VToolbarItemsBuilder struct {
	VTagBuilder[*VToolbarItemsBuilder]
}

func VToolbarItems(children ...h.HTMLComponent) *VToolbarItemsBuilder {
	return VTag(&VToolbarItemsBuilder{}, "v-toolbar-items", children...)
}

func (b *VToolbarItemsBuilder) Color(v string) (r *VToolbarItemsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VToolbarItemsBuilder) Variant(v interface{}) (r *VToolbarItemsBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VToolbarItemsBuilder) On(name string, value string) (r *VToolbarItemsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VToolbarItemsBuilder) Bind(name string, value string) (r *VToolbarItemsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VToolbarItemsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VToolbarItemsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VToolbarItemsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VToolbarItemsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VToolbarItemsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VToolbarItemsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VToolbarItemsBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VToolbarItemsBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VToolbarItemsBuilder) SlotDefault(child ...h.HTMLComponent) (r *VToolbarItemsBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VToolbarItemsBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VToolbarItemsBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
