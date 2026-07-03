package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCardTextBuilder struct {
	VTagBuilder[*VCardTextBuilder]
}

func VCardText(children ...h.HTMLComponent) *VCardTextBuilder {
	return VTag(&VCardTextBuilder{}, "v-card-text", children...)
}

func (b *VCardTextBuilder) Tag(v interface{}) (r *VCardTextBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VCardTextBuilder) Opacity(v interface{}) (r *VCardTextBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VCardTextBuilder) On(name string, value string) (r *VCardTextBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardTextBuilder) Bind(name string, value string) (r *VCardTextBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCardTextBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCardTextBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCardTextBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCardTextBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCardTextBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCardTextBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCardTextBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCardTextBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCardTextBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCardTextBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCardTextBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCardTextBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
