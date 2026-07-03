package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFieldLabelBuilder struct {
	VTagBuilder[*VFieldLabelBuilder]
}

func VFieldLabel(children ...h.HTMLComponent) *VFieldLabelBuilder {
	return VTag(&VFieldLabelBuilder{}, "v-field-label", children...)
}

func (b *VFieldLabelBuilder) Floating(v bool) (r *VFieldLabelBuilder) {
	b.Attr(":floating", fmt.Sprint(v))
	return b
}

func (b *VFieldLabelBuilder) On(name string, value string) (r *VFieldLabelBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFieldLabelBuilder) Bind(name string, value string) (r *VFieldLabelBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VFieldLabelBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VFieldLabelBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VFieldLabelBuilder) Slot(name string, child ...h.HTMLComponent) (r *VFieldLabelBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VFieldLabelBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VFieldLabelBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VFieldLabelBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VFieldLabelBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VFieldLabelBuilder) SlotDefault(child ...h.HTMLComponent) (r *VFieldLabelBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VFieldLabelBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VFieldLabelBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
