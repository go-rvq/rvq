package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VFormBuilder struct {
	VTagBuilder[*VFormBuilder]
}

func VForm(children ...h.HTMLComponent) *VFormBuilder {
	return VTag(&VFormBuilder{}, "v-form", children...)
}

func (b *VFormBuilder) Disabled(v bool) (r *VFormBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VFormBuilder) ModelValue(v bool) (r *VFormBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VFormBuilder) Readonly(v bool) (r *VFormBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VFormBuilder) ValidateOn(v interface{}) (r *VFormBuilder) {
	b.Attr(":validate-on", h.JSONString(v))
	return b
}

func (b *VFormBuilder) FastFail(v bool) (r *VFormBuilder) {
	b.Attr(":fast-fail", fmt.Sprint(v))
	return b
}

func (b *VFormBuilder) On(name string, value string) (r *VFormBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VFormBuilder) Bind(name string, value string) (r *VFormBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VFormBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VFormBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VFormBuilder) Slot(name string, child ...h.HTMLComponent) (r *VFormBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VFormBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VFormBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VFormBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VFormBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VFormBuilder) SlotDefault(child ...h.HTMLComponent) (r *VFormBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VFormBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VFormBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
