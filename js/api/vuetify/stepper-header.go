package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VStepperHeaderBuilder struct {
	VTagBuilder[*VStepperHeaderBuilder]
}

func VStepperHeader(children ...h.HTMLComponent) *VStepperHeaderBuilder {
	return VTag(&VStepperHeaderBuilder{}, "v-stepper-header", children...)
}

func (b *VStepperHeaderBuilder) Tag(v string) (r *VStepperHeaderBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VStepperHeaderBuilder) On(name string, value string) (r *VStepperHeaderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VStepperHeaderBuilder) Bind(name string, value string) (r *VStepperHeaderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VStepperHeaderBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VStepperHeaderBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VStepperHeaderBuilder) Slot(name string, child ...h.HTMLComponent) (r *VStepperHeaderBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VStepperHeaderBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VStepperHeaderBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VStepperHeaderBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VStepperHeaderBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VStepperHeaderBuilder) SlotDefault(child ...h.HTMLComponent) (r *VStepperHeaderBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VStepperHeaderBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VStepperHeaderBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
