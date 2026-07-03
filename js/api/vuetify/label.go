package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VLabelBuilder struct {
	VTagBuilder[*VLabelBuilder]
}

func VLabel(children ...h.HTMLComponent) *VLabelBuilder {
	return VTag(&VLabelBuilder{}, "v-label", children...)
}

func (b *VLabelBuilder) Theme(v string) (r *VLabelBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VLabelBuilder) Text(v string) (r *VLabelBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VLabelBuilder) On(name string, value string) (r *VLabelBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLabelBuilder) Bind(name string, value string) (r *VLabelBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VLabelBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VLabelBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VLabelBuilder) Slot(name string, child ...h.HTMLComponent) (r *VLabelBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VLabelBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VLabelBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VLabelBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VLabelBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VLabelBuilder) SlotDefault(child ...h.HTMLComponent) (r *VLabelBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VLabelBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VLabelBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
