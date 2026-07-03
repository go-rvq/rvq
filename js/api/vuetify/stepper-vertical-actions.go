package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperVerticalActionsBuilder struct {
	VTagBuilder[*VStepperVerticalActionsBuilder]
}

func VStepperVerticalActions(children ...h.HTMLComponent) *VStepperVerticalActionsBuilder {
	return VTag(&VStepperVerticalActionsBuilder{}, "v-stepper-vertical-actions", children...)
}

func (b *VStepperVerticalActionsBuilder) Disabled(v interface{}) (r *VStepperVerticalActionsBuilder) {
	b.Attr(":disabled", h.JSONString(v))
	return b
}

func (b *VStepperVerticalActionsBuilder) Color(v string) (r *VStepperVerticalActionsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VStepperVerticalActionsBuilder) PrevText(v string) (r *VStepperVerticalActionsBuilder) {
	b.Attr("prev-text", v)
	return b
}

func (b *VStepperVerticalActionsBuilder) NextText(v string) (r *VStepperVerticalActionsBuilder) {
	b.Attr("next-text", v)
	return b
}

func (b *VStepperVerticalActionsBuilder) On(name string, value string) (r *VStepperVerticalActionsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperVerticalActionsBuilder) Bind(name string, value string) (r *VStepperVerticalActionsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VStepperVerticalActionsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VStepperVerticalActionsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VStepperVerticalActionsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VStepperVerticalActionsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VStepperVerticalActionsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VStepperVerticalActionsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VStepperVerticalActionsBuilder) SetSlotPrev(child ...h.HTMLComponent) {
	b.SetSlot("prev", child...)
}

func (b *VStepperVerticalActionsBuilder) SetScopedSlotPrev(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prev", scope, child...)
}

func (b *VStepperVerticalActionsBuilder) SlotPrev(child ...h.HTMLComponent) (r *VStepperVerticalActionsBuilder) {
	b.SetSlotPrev(child...)
	return b
}

func (b *VStepperVerticalActionsBuilder) ScopedSlotPrev(scope string, child ...h.HTMLComponent) (r *VStepperVerticalActionsBuilder) {
	b.SetScopedSlotPrev(scope, child...)
	return b
}

func (b *VStepperVerticalActionsBuilder) SetSlotNext(child ...h.HTMLComponent) {
	b.SetSlot("next", child...)
}

func (b *VStepperVerticalActionsBuilder) SetScopedSlotNext(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("next", scope, child...)
}

func (b *VStepperVerticalActionsBuilder) SlotNext(child ...h.HTMLComponent) (r *VStepperVerticalActionsBuilder) {
	b.SetSlotNext(child...)
	return b
}

func (b *VStepperVerticalActionsBuilder) ScopedSlotNext(scope string, child ...h.HTMLComponent) (r *VStepperVerticalActionsBuilder) {
	b.SetScopedSlotNext(scope, child...)
	return b
}
