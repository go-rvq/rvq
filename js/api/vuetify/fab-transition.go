package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFabTransitionBuilder struct {
	VTagBuilder[*VFabTransitionBuilder]
}

func VFabTransition(children ...h.HTMLComponent) *VFabTransitionBuilder {
	return VTag(&VFabTransitionBuilder{}, "v-fab-transition", children...)
}

func (b *VFabTransitionBuilder) Mode(v string) (r *VFabTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VFabTransitionBuilder) Disabled(v bool) (r *VFabTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFabTransitionBuilder) Origin(v string) (r *VFabTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VFabTransitionBuilder) Group(v bool) (r *VFabTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VFabTransitionBuilder) HideOnLeave(v bool) (r *VFabTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VFabTransitionBuilder) LeaveAbsolute(v bool) (r *VFabTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VFabTransitionBuilder) On(name string, value string) (r *VFabTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFabTransitionBuilder) Bind(name string, value string) (r *VFabTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VFabTransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VFabTransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VFabTransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VFabTransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VFabTransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VFabTransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VFabTransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VFabTransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VFabTransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VFabTransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VFabTransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VFabTransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
