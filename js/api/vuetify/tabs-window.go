package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTabsWindowBuilder struct {
	VTagBuilder[*VTabsWindowBuilder]
}

func VTabsWindow(children ...h.HTMLComponent) *VTabsWindowBuilder {
	return VTag(&VTabsWindowBuilder{}, "v-tabs-window", children...)
}

func (b *VTabsWindowBuilder) ModelValue(v interface{}) (r *VTabsWindowBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTabsWindowBuilder) Reverse(v bool) (r *VTabsWindowBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VTabsWindowBuilder) Tag(v interface{}) (r *VTabsWindowBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VTabsWindowBuilder) Theme(v string) (r *VTabsWindowBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTabsWindowBuilder) Disabled(v bool) (r *VTabsWindowBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTabsWindowBuilder) SelectedClass(v string) (r *VTabsWindowBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VTabsWindowBuilder) Direction(v interface{}) (r *VTabsWindowBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VTabsWindowBuilder) On(name string, value string) (r *VTabsWindowBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTabsWindowBuilder) Bind(name string, value string) (r *VTabsWindowBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VTabsWindowBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VTabsWindowBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VTabsWindowBuilder) Slot(name string, child ...h.HTMLComponent) (r *VTabsWindowBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VTabsWindowBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VTabsWindowBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VTabsWindowBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VTabsWindowBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VTabsWindowBuilder) SlotDefault(child ...h.HTMLComponent) (r *VTabsWindowBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VTabsWindowBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VTabsWindowBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
