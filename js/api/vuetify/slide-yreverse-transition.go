package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSlideYreverseTransitionBuilder struct {
	VTagBuilder[*VSlideYreverseTransitionBuilder]
}

func VSlideYreverseTransition(children ...h.HTMLComponent) *VSlideYreverseTransitionBuilder {
	return VTag(&VSlideYreverseTransitionBuilder{}, "v-slide-yreverse-transition", children...)
}

func (b *VSlideYreverseTransitionBuilder) Mode(v string) (r *VSlideYreverseTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VSlideYreverseTransitionBuilder) Disabled(v bool) (r *VSlideYreverseTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideYreverseTransitionBuilder) Origin(v string) (r *VSlideYreverseTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VSlideYreverseTransitionBuilder) Group(v bool) (r *VSlideYreverseTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VSlideYreverseTransitionBuilder) HideOnLeave(v bool) (r *VSlideYreverseTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VSlideYreverseTransitionBuilder) LeaveAbsolute(v bool) (r *VSlideYreverseTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VSlideYreverseTransitionBuilder) On(name string, value string) (r *VSlideYreverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideYreverseTransitionBuilder) Bind(name string, value string) (r *VSlideYreverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSlideYreverseTransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSlideYreverseTransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSlideYreverseTransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSlideYreverseTransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSlideYreverseTransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSlideYreverseTransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSlideYreverseTransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSlideYreverseTransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSlideYreverseTransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSlideYreverseTransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSlideYreverseTransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSlideYreverseTransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
