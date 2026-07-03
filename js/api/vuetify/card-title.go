package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCardTitleBuilder struct {
	VTagBuilder[*VCardTitleBuilder]
}

func VCardTitle(children ...h.HTMLComponent) *VCardTitleBuilder {
	return VTag(&VCardTitleBuilder{}, "v-card-title", children...)
}

func (b *VCardTitleBuilder) Tag(v string) (r *VCardTitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VCardTitleBuilder) On(name string, value string) (r *VCardTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardTitleBuilder) Bind(name string, value string) (r *VCardTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCardTitleBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCardTitleBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCardTitleBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCardTitleBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCardTitleBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCardTitleBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCardTitleBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCardTitleBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCardTitleBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCardTitleBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCardTitleBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCardTitleBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
