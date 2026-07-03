package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCardActionsBuilder struct {
	VTagBuilder[*VCardActionsBuilder]
}

func VCardActions(children ...h.HTMLComponent) *VCardActionsBuilder {
	return VTag(&VCardActionsBuilder{}, "v-card-actions", children...)
}

func (b *VCardActionsBuilder) On(name string, value string) (r *VCardActionsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCardActionsBuilder) Bind(name string, value string) (r *VCardActionsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCardActionsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCardActionsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCardActionsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCardActionsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCardActionsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCardActionsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCardActionsBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCardActionsBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCardActionsBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCardActionsBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCardActionsBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCardActionsBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
