package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VExpandTransitionBuilder struct {
	VTagBuilder[*VExpandTransitionBuilder]
}

func VExpandTransition(children ...h.HTMLComponent) *VExpandTransitionBuilder {
	return VTag(&VExpandTransitionBuilder{}, "v-expand-transition", children...)
}

func (b *VExpandTransitionBuilder) Mode(v interface{}) (r *VExpandTransitionBuilder) {
	b.Attr(":mode", h.JSONString(v))
	return b
}

func (b *VExpandTransitionBuilder) Disabled(v bool) (r *VExpandTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VExpandTransitionBuilder) Group(v bool) (r *VExpandTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VExpandTransitionBuilder) On(name string, value string) (r *VExpandTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VExpandTransitionBuilder) Bind(name string, value string) (r *VExpandTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VExpandTransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VExpandTransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VExpandTransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VExpandTransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VExpandTransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VExpandTransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VExpandTransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VExpandTransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VExpandTransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VExpandTransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VExpandTransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VExpandTransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
