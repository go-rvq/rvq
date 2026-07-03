package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VWindowItemBuilder struct {
	VTagBuilder[*VWindowItemBuilder]
}

func VWindowItem(children ...h.HTMLComponent) *VWindowItemBuilder {
	return VTag(&VWindowItemBuilder{}, "v-window-item", children...)
}

func (b *VWindowItemBuilder) Value(v interface{}) (r *VWindowItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VWindowItemBuilder) Eager(v bool) (r *VWindowItemBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VWindowItemBuilder) Disabled(v bool) (r *VWindowItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VWindowItemBuilder) Transition(v interface{}) (r *VWindowItemBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VWindowItemBuilder) ReverseTransition(v interface{}) (r *VWindowItemBuilder) {
	b.Attr(":reverse-transition", h.JSONString(v))
	return b
}

func (b *VWindowItemBuilder) SelectedClass(v string) (r *VWindowItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VWindowItemBuilder) On(name string, value string) (r *VWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VWindowItemBuilder) Bind(name string, value string) (r *VWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VWindowItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VWindowItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VWindowItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VWindowItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VWindowItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VWindowItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VWindowItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VWindowItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VWindowItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VWindowItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VWindowItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VWindowItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
