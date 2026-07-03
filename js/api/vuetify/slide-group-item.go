package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSlideGroupItemBuilder struct {
	VTagBuilder[*VSlideGroupItemBuilder]
}

func VSlideGroupItem(children ...h.HTMLComponent) *VSlideGroupItemBuilder {
	return VTag(&VSlideGroupItemBuilder{}, "v-slide-group-item", children...)
}

func (b *VSlideGroupItemBuilder) Value(v interface{}) (r *VSlideGroupItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VSlideGroupItemBuilder) Disabled(v bool) (r *VSlideGroupItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideGroupItemBuilder) SelectedClass(v string) (r *VSlideGroupItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VSlideGroupItemBuilder) On(name string, value string) (r *VSlideGroupItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideGroupItemBuilder) Bind(name string, value string) (r *VSlideGroupItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSlideGroupItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSlideGroupItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSlideGroupItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSlideGroupItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSlideGroupItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSlideGroupItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSlideGroupItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSlideGroupItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSlideGroupItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSlideGroupItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSlideGroupItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSlideGroupItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
