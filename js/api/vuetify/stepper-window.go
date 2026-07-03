package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperWindowBuilder struct {
	VTagBuilder[*VStepperWindowBuilder]
}

func VStepperWindow(children ...h.HTMLComponent) *VStepperWindowBuilder {
	return VTag(&VStepperWindowBuilder{}, "v-stepper-window", children...)
}

func (b *VStepperWindowBuilder) ModelValue(v interface{}) (r *VStepperWindowBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VStepperWindowBuilder) Reverse(v bool) (r *VStepperWindowBuilder) {
	b.Attr(":reverse", fmt.Sprint(v))
	return b
}

func (b *VStepperWindowBuilder) Tag(v interface{}) (r *VStepperWindowBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VStepperWindowBuilder) Theme(v string) (r *VStepperWindowBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VStepperWindowBuilder) Disabled(v bool) (r *VStepperWindowBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VStepperWindowBuilder) SelectedClass(v string) (r *VStepperWindowBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VStepperWindowBuilder) Direction(v interface{}) (r *VStepperWindowBuilder) {
	b.Attr(":direction", h.JSONString(v))
	return b
}

func (b *VStepperWindowBuilder) On(name string, value string) (r *VStepperWindowBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperWindowBuilder) Bind(name string, value string) (r *VStepperWindowBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VStepperWindowBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VStepperWindowBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VStepperWindowBuilder) Slot(name string, child ...h.HTMLComponent) (r *VStepperWindowBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VStepperWindowBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VStepperWindowBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VStepperWindowBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VStepperWindowBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VStepperWindowBuilder) SlotDefault(child ...h.HTMLComponent) (r *VStepperWindowBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VStepperWindowBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VStepperWindowBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
