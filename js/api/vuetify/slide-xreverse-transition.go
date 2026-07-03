package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSlideXreverseTransitionBuilder struct {
	VTagBuilder[*VSlideXreverseTransitionBuilder]
}

func VSlideXreverseTransition(children ...h.HTMLComponent) *VSlideXreverseTransitionBuilder {
	return VTag(&VSlideXreverseTransitionBuilder{}, "v-slide-xreverse-transition", children...)
}

func (b *VSlideXreverseTransitionBuilder) Mode(v string) (r *VSlideXreverseTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VSlideXreverseTransitionBuilder) Disabled(v bool) (r *VSlideXreverseTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideXreverseTransitionBuilder) Origin(v string) (r *VSlideXreverseTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VSlideXreverseTransitionBuilder) Group(v bool) (r *VSlideXreverseTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VSlideXreverseTransitionBuilder) HideOnLeave(v bool) (r *VSlideXreverseTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VSlideXreverseTransitionBuilder) LeaveAbsolute(v bool) (r *VSlideXreverseTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VSlideXreverseTransitionBuilder) On(name string, value string) (r *VSlideXreverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideXreverseTransitionBuilder) Bind(name string, value string) (r *VSlideXreverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSlideXreverseTransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSlideXreverseTransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSlideXreverseTransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSlideXreverseTransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSlideXreverseTransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSlideXreverseTransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSlideXreverseTransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSlideXreverseTransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSlideXreverseTransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSlideXreverseTransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSlideXreverseTransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSlideXreverseTransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
