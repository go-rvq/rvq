package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTabsWindowItemBuilder struct {
	VTagBuilder[*VTabsWindowItemBuilder]
}

func VTabsWindowItem(children ...h.HTMLComponent) *VTabsWindowItemBuilder {
	return VTag(&VTabsWindowItemBuilder{}, "v-tabs-window-item", children...)
}

func (b *VTabsWindowItemBuilder) Value(v interface{}) (r *VTabsWindowItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VTabsWindowItemBuilder) Disabled(v bool) (r *VTabsWindowItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTabsWindowItemBuilder) SelectedClass(v string) (r *VTabsWindowItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VTabsWindowItemBuilder) Transition(v interface{}) (r *VTabsWindowItemBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VTabsWindowItemBuilder) Eager(v bool) (r *VTabsWindowItemBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VTabsWindowItemBuilder) ReverseTransition(v interface{}) (r *VTabsWindowItemBuilder) {
	b.Attr(":reverse-transition", h.JSONString(v))
	return b
}

func (b *VTabsWindowItemBuilder) On(name string, value string) (r *VTabsWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTabsWindowItemBuilder) Bind(name string, value string) (r *VTabsWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTabsWindowItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTabsWindowItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTabsWindowItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTabsWindowItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTabsWindowItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTabsWindowItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTabsWindowItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VTabsWindowItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VTabsWindowItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VTabsWindowItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VTabsWindowItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VTabsWindowItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
