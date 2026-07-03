package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperWindowItemBuilder struct {
	VTagBuilder[*VStepperWindowItemBuilder]
}

func VStepperWindowItem(children ...h.HTMLComponent) *VStepperWindowItemBuilder {
	return VTag(&VStepperWindowItemBuilder{}, "v-stepper-window-item", children...)
}

func (b *VStepperWindowItemBuilder) Value(v interface{}) (r *VStepperWindowItemBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VStepperWindowItemBuilder) Eager(v bool) (r *VStepperWindowItemBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VStepperWindowItemBuilder) Disabled(v bool) (r *VStepperWindowItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperWindowItemBuilder) Transition(v interface{}) (r *VStepperWindowItemBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VStepperWindowItemBuilder) ReverseTransition(v interface{}) (r *VStepperWindowItemBuilder) {
	b.Attr(":reverse-transition", h.JSONString(v))
	return b
}

func (b *VStepperWindowItemBuilder) SelectedClass(v string) (r *VStepperWindowItemBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperWindowItemBuilder) On(name string, value string) (r *VStepperWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperWindowItemBuilder) Bind(name string, value string) (r *VStepperWindowItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VStepperWindowItemBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VStepperWindowItemBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VStepperWindowItemBuilder) Slot(name string, child ...h.HTMLComponent) (r *VStepperWindowItemBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VStepperWindowItemBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VStepperWindowItemBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VStepperWindowItemBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VStepperWindowItemBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VStepperWindowItemBuilder) SlotDefault(child ...h.HTMLComponent) (r *VStepperWindowItemBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VStepperWindowItemBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VStepperWindowItemBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
