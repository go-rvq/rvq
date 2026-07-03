package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VListItemMediaBuilder struct {
	VTagBuilder[*VListItemMediaBuilder]
}

func VListItemMedia(children ...h.HTMLComponent) *VListItemMediaBuilder {
	return VTag(&VListItemMediaBuilder{}, "v-list-item-media", children...)
}

func (b *VListItemMediaBuilder) Tag(v interface{}) (r *VListItemMediaBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VListItemMediaBuilder) Start(v bool) (r *VListItemMediaBuilder) {
	b.Attr(":start", fmt.Sprint(v))
	return b
}

func (b *VListItemMediaBuilder) End(v bool) (r *VListItemMediaBuilder) {
	b.Attr(":end", fmt.Sprint(v))
	return b
}

func (b *VListItemMediaBuilder) On(name string, value string) (r *VListItemMediaBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VListItemMediaBuilder) Bind(name string, value string) (r *VListItemMediaBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VListItemMediaBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VListItemMediaBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VListItemMediaBuilder) Slot(name string, child ...h.HTMLComponent) (r *VListItemMediaBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VListItemMediaBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VListItemMediaBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VListItemMediaBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VListItemMediaBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VListItemMediaBuilder) SlotDefault(child ...h.HTMLComponent) (r *VListItemMediaBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VListItemMediaBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VListItemMediaBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
