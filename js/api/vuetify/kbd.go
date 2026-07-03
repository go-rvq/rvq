package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VKbdBuilder struct {
	VTagBuilder[*VKbdBuilder]
}

func VKbd(children ...h.HTMLComponent) *VKbdBuilder {
	return VTag(&VKbdBuilder{}, "v-kbd", children...)
}

func (b *VKbdBuilder) Tag(v string) (r *VKbdBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VKbdBuilder) On(name string, value string) (r *VKbdBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VKbdBuilder) Bind(name string, value string) (r *VKbdBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VKbdBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VKbdBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VKbdBuilder) Slot(name string, child ...h.HTMLComponent) (r *VKbdBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VKbdBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VKbdBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VKbdBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VKbdBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VKbdBuilder) SlotDefault(child ...h.HTMLComponent) (r *VKbdBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VKbdBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VKbdBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
