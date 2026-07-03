package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListItemTitleBuilder struct {
	VTagBuilder[*VListItemTitleBuilder]
}

func VListItemTitle(children ...h.HTMLComponent) *VListItemTitleBuilder {
	return VTag(&VListItemTitleBuilder{}, "v-list-item-title", children...)
}

func (b *VListItemTitleBuilder) Tag(v string) (r *VListItemTitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VListItemTitleBuilder) On(name string, value string) (r *VListItemTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListItemTitleBuilder) Bind(name string, value string) (r *VListItemTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VListItemTitleBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VListItemTitleBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VListItemTitleBuilder) Slot(name string, child ...h.HTMLComponent) (r *VListItemTitleBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VListItemTitleBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VListItemTitleBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VListItemTitleBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VListItemTitleBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VListItemTitleBuilder) SlotDefault(child ...h.HTMLComponent) (r *VListItemTitleBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VListItemTitleBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VListItemTitleBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
