package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSpacerBuilder struct {
	VTagBuilder[*VSpacerBuilder]
}

func VSpacer(children ...h.HTMLComponent) *VSpacerBuilder {
	return VTag(&VSpacerBuilder{}, "v-spacer", children...)
}

func (b *VSpacerBuilder) Tag(v string) (r *VSpacerBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VSpacerBuilder) On(name string, value string) (r *VSpacerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSpacerBuilder) Bind(name string, value string) (r *VSpacerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSpacerBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSpacerBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSpacerBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSpacerBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSpacerBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSpacerBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSpacerBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSpacerBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSpacerBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSpacerBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSpacerBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSpacerBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
