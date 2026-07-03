package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VToolbarTitleBuilder struct {
	VTagBuilder[*VToolbarTitleBuilder]
}

func VToolbarTitle(children ...h.HTMLComponent) *VToolbarTitleBuilder {
	return VTag(&VToolbarTitleBuilder{}, "v-toolbar-title", children...)
}

func (b *VToolbarTitleBuilder) Tag(v interface{}) (r *VToolbarTitleBuilder) {
	b.Attr(":tag", h.JSONString(v))
	return b
}

func (b *VToolbarTitleBuilder) Text(v string) (r *VToolbarTitleBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VToolbarTitleBuilder) On(name string, value string) (r *VToolbarTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VToolbarTitleBuilder) Bind(name string, value string) (r *VToolbarTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VToolbarTitleBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VToolbarTitleBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VToolbarTitleBuilder) Slot(name string, child ...h.HTMLComponent) (r *VToolbarTitleBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VToolbarTitleBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VToolbarTitleBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VToolbarTitleBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VToolbarTitleBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VToolbarTitleBuilder) SlotDefault(child ...h.HTMLComponent) (r *VToolbarTitleBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VToolbarTitleBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VToolbarTitleBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VToolbarTitleBuilder) SetSlotText(child ...h.HTMLComponent) {
	b.SetSlot("text", child...)
}

func (b *VToolbarTitleBuilder) SetScopedSlotText(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("text", scope, child...)
}

func (b *VToolbarTitleBuilder) SlotText(child ...h.HTMLComponent) (r *VToolbarTitleBuilder) {
	b.SetSlotText(child...)
	return b
}

func (b *VToolbarTitleBuilder) ScopedSlotText(scope string, child ...h.HTMLComponent) (r *VToolbarTitleBuilder) {
	b.SetScopedSlotText(scope, child...)
	return b
}
