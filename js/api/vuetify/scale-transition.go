package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VScaleTransitionBuilder struct {
	VTagBuilder[*VScaleTransitionBuilder]
}

func VScaleTransition(children ...h.HTMLComponent) *VScaleTransitionBuilder {
	return VTag(&VScaleTransitionBuilder{}, "v-scale-transition", children...)
}

func (b *VScaleTransitionBuilder) Disabled(v bool) (r *VScaleTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VScaleTransitionBuilder) Mode(v string) (r *VScaleTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VScaleTransitionBuilder) Group(v bool) (r *VScaleTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VScaleTransitionBuilder) Origin(v string) (r *VScaleTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VScaleTransitionBuilder) HideOnLeave(v bool) (r *VScaleTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VScaleTransitionBuilder) LeaveAbsolute(v bool) (r *VScaleTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VScaleTransitionBuilder) On(name string, value string) (r *VScaleTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VScaleTransitionBuilder) Bind(name string, value string) (r *VScaleTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VScaleTransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VScaleTransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VScaleTransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VScaleTransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VScaleTransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VScaleTransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VScaleTransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VScaleTransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VScaleTransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VScaleTransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VScaleTransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VScaleTransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
