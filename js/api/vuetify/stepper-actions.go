package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperActionsBuilder struct {
	VTagBuilder[*VStepperActionsBuilder]
}

func VStepperActions(children ...h.HTMLComponent) *VStepperActionsBuilder {
	return VTag(&VStepperActionsBuilder{}, "v-stepper-actions", children...)
}

func (b *VStepperActionsBuilder) Disabled(v interface{}) (r *VStepperActionsBuilder) {
	b.Attr(":disabled", h.JSONString(v))
	return b
}

func (b *VStepperActionsBuilder) Color(v string) (r *VStepperActionsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperActionsBuilder) PrevText(v string) (r *VStepperActionsBuilder) {
	b.Attr("prev-text", v)
	return b
}

func (b *VStepperActionsBuilder) NextText(v string) (r *VStepperActionsBuilder) {
	b.Attr("next-text", v)
	return b
}

func (b *VStepperActionsBuilder) On(name string, value string) (r *VStepperActionsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperActionsBuilder) Bind(name string, value string) (r *VStepperActionsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VStepperActionsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VStepperActionsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VStepperActionsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VStepperActionsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VStepperActionsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VStepperActionsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VStepperActionsBuilder) SetSlotPrev(child ...h.HTMLComponent) {
	b.SetSlot("prev", child...)
}

func (b *VStepperActionsBuilder) SetScopedSlotPrev(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prev", scope, child...)
}

func (b *VStepperActionsBuilder) SlotPrev(child ...h.HTMLComponent) (r *VStepperActionsBuilder) {
	b.SetSlotPrev(child...)
	return b
}

func (b *VStepperActionsBuilder) ScopedSlotPrev(scope string, child ...h.HTMLComponent) (r *VStepperActionsBuilder) {
	b.SetScopedSlotPrev(scope, child...)
	return b
}

func (b *VStepperActionsBuilder) SetSlotNext(child ...h.HTMLComponent) {
	b.SetSlot("next", child...)
}

func (b *VStepperActionsBuilder) SetScopedSlotNext(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("next", scope, child...)
}

func (b *VStepperActionsBuilder) SlotNext(child ...h.HTMLComponent) (r *VStepperActionsBuilder) {
	b.SetSlotNext(child...)
	return b
}

func (b *VStepperActionsBuilder) ScopedSlotNext(scope string, child ...h.HTMLComponent) (r *VStepperActionsBuilder) {
	b.SetScopedSlotNext(scope, child...)
	return b
}
