package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VScrollYreverseTransitionBuilder struct {
	VTagBuilder[*VScrollYreverseTransitionBuilder]
}

func VScrollYreverseTransition(children ...h.HTMLComponent) *VScrollYreverseTransitionBuilder {
	return VTag(&VScrollYreverseTransitionBuilder{}, "v-scroll-yreverse-transition", children...)
}

func (b *VScrollYreverseTransitionBuilder) Disabled(v bool) (r *VScrollYreverseTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VScrollYreverseTransitionBuilder) Mode(v string) (r *VScrollYreverseTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VScrollYreverseTransitionBuilder) Group(v bool) (r *VScrollYreverseTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VScrollYreverseTransitionBuilder) Origin(v string) (r *VScrollYreverseTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VScrollYreverseTransitionBuilder) HideOnLeave(v bool) (r *VScrollYreverseTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VScrollYreverseTransitionBuilder) LeaveAbsolute(v bool) (r *VScrollYreverseTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VScrollYreverseTransitionBuilder) On(name string, value string) (r *VScrollYreverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VScrollYreverseTransitionBuilder) Bind(name string, value string) (r *VScrollYreverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VScrollYreverseTransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VScrollYreverseTransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VScrollYreverseTransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VScrollYreverseTransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VScrollYreverseTransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VScrollYreverseTransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VScrollYreverseTransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VScrollYreverseTransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VScrollYreverseTransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VScrollYreverseTransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VScrollYreverseTransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VScrollYreverseTransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
