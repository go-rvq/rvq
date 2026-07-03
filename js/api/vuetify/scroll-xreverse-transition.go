package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VScrollXreverseTransitionBuilder struct {
	VTagBuilder[*VScrollXreverseTransitionBuilder]
}

func VScrollXreverseTransition(children ...h.HTMLComponent) *VScrollXreverseTransitionBuilder {
	return VTag(&VScrollXreverseTransitionBuilder{}, "v-scroll-xreverse-transition", children...)
}

func (b *VScrollXreverseTransitionBuilder) Disabled(v bool) (r *VScrollXreverseTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VScrollXreverseTransitionBuilder) Mode(v string) (r *VScrollXreverseTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VScrollXreverseTransitionBuilder) Group(v bool) (r *VScrollXreverseTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VScrollXreverseTransitionBuilder) Origin(v string) (r *VScrollXreverseTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VScrollXreverseTransitionBuilder) HideOnLeave(v bool) (r *VScrollXreverseTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VScrollXreverseTransitionBuilder) LeaveAbsolute(v bool) (r *VScrollXreverseTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VScrollXreverseTransitionBuilder) On(name string, value string) (r *VScrollXreverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VScrollXreverseTransitionBuilder) Bind(name string, value string) (r *VScrollXreverseTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VScrollXreverseTransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VScrollXreverseTransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VScrollXreverseTransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VScrollXreverseTransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VScrollXreverseTransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VScrollXreverseTransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VScrollXreverseTransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VScrollXreverseTransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VScrollXreverseTransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VScrollXreverseTransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VScrollXreverseTransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VScrollXreverseTransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
