package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListItemActionBuilder struct {
	VTagBuilder[*VListItemActionBuilder]
}

func VListItemAction(children ...h.HTMLComponent) *VListItemActionBuilder {
	return VTag(&VListItemActionBuilder{}, "v-list-item-action", children...)
}

func (b *VListItemActionBuilder) Tag(v interface{}) (r *VListItemActionBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VListItemActionBuilder) Start(v bool) (r *VListItemActionBuilder) {
	b.Attr(":start", fmt.Sprint(v))
	return b
}

func (b *VListItemActionBuilder) End(v bool) (r *VListItemActionBuilder) {
	b.Attr(":end", fmt.Sprint(v))
	return b
}

func (b *VListItemActionBuilder) On(name string, value string) (r *VListItemActionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListItemActionBuilder) Bind(name string, value string) (r *VListItemActionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VListItemActionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VListItemActionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VListItemActionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VListItemActionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VListItemActionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VListItemActionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VListItemActionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VListItemActionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VListItemActionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VListItemActionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VListItemActionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VListItemActionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
