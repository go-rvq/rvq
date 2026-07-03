package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VMessagesBuilder struct {
	VTagBuilder[*VMessagesBuilder]
}

func VMessages(children ...h.HTMLComponent) *VMessagesBuilder {
	return VTag(&VMessagesBuilder{}, "v-messages", children...)
}

func (b *VMessagesBuilder) Color(v string) (r *VMessagesBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VMessagesBuilder) Transition(v interface{}) (r *VMessagesBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VMessagesBuilder) Messages(v interface{}) (r *VMessagesBuilder) {
	b.Attr(":messages", h.JSONString(v))
	return b
}

func (b *VMessagesBuilder) Active(v bool) (r *VMessagesBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VMessagesBuilder) On(name string, value string) (r *VMessagesBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VMessagesBuilder) Bind(name string, value string) (r *VMessagesBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VMessagesBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VMessagesBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VMessagesBuilder) Slot(name string, child ...h.HTMLComponent) (r *VMessagesBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VMessagesBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VMessagesBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VMessagesBuilder) SetSlotMessage(child ...h.HTMLComponent) {
	b.SetSlot("message", child...)
}

func (b *VMessagesBuilder) SetScopedSlotMessage(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("message", scope, child...)
}

func (b *VMessagesBuilder) SlotMessage(child ...h.HTMLComponent) (r *VMessagesBuilder) {
	b.SetSlotMessage(child...)
	return b
}

func (b *VMessagesBuilder) ScopedSlotMessage(scope string, child ...h.HTMLComponent) (r *VMessagesBuilder) {
	b.SetScopedSlotMessage(scope, child...)
	return b
}
