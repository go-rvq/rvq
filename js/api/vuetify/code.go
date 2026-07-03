package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCodeBuilder struct {
	VTagBuilder[*VCodeBuilder]
}

func VCode(children ...h.HTMLComponent) *VCodeBuilder {
	return VTag(&VCodeBuilder{}, "v-code", children...)
}

func (b *VCodeBuilder) Tag(v string) (r *VCodeBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VCodeBuilder) On(name string, value string) (r *VCodeBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCodeBuilder) Bind(name string, value string) (r *VCodeBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCodeBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCodeBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCodeBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCodeBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCodeBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCodeBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCodeBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCodeBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCodeBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCodeBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCodeBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCodeBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
