package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VConfirmEditBuilder struct {
	VTagBuilder[*VConfirmEditBuilder]
}

func VConfirmEdit(children ...h.HTMLComponent) *VConfirmEditBuilder {
	return VTag(&VConfirmEditBuilder{}, "v-confirm-edit", children...)
}

func (b *VConfirmEditBuilder) Disabled(v interface{}) (r *VConfirmEditBuilder) {
	b.Attr(":disabled", h.JSONString(v))
	return b
}

func (b *VConfirmEditBuilder) Color(v string) (r *VConfirmEditBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VConfirmEditBuilder) ModelValue(v interface{}) (r *VConfirmEditBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VConfirmEditBuilder) CancelText(v string) (r *VConfirmEditBuilder) {
	b.Attr("cancel-text", v)
	return b
}

func (b *VConfirmEditBuilder) OkText(v string) (r *VConfirmEditBuilder) {
	b.Attr("ok-text", v)
	return b
}

func (b *VConfirmEditBuilder) HideActions(v bool) (r *VConfirmEditBuilder) {
	b.Attr(":hide-actions", fmt.Sprint(v))
	return b
}

func (b *VConfirmEditBuilder) On(name string, value string) (r *VConfirmEditBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VConfirmEditBuilder) Bind(name string, value string) (r *VConfirmEditBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VConfirmEditBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VConfirmEditBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VConfirmEditBuilder) Slot(name string, child ...h.HTMLComponent) (r *VConfirmEditBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VConfirmEditBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VConfirmEditBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VConfirmEditBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VConfirmEditBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VConfirmEditBuilder) SlotDefault(child ...h.HTMLComponent) (r *VConfirmEditBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VConfirmEditBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VConfirmEditBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
