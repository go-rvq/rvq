package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCounterBuilder struct {
	VTagBuilder[*VCounterBuilder]
}

func VCounter(children ...h.HTMLComponent) *VCounterBuilder {
	return VTag(&VCounterBuilder{}, "v-counter", children...)
}

func (b *VCounterBuilder) Value(v interface{}) (r *VCounterBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VCounterBuilder) Disabled(v bool) (r *VCounterBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VCounterBuilder) Transition(v interface{}) (r *VCounterBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VCounterBuilder) Active(v bool) (r *VCounterBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VCounterBuilder) Max(v interface{}) (r *VCounterBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VCounterBuilder) On(name string, value string) (r *VCounterBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCounterBuilder) Bind(name string, value string) (r *VCounterBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VCounterBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VCounterBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VCounterBuilder) Slot(name string, child ...h.HTMLComponent) (r *VCounterBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VCounterBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VCounterBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VCounterBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VCounterBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VCounterBuilder) SlotDefault(child ...h.HTMLComponent) (r *VCounterBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VCounterBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VCounterBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
