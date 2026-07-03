package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VHoverBuilder struct {
	VTagBuilder[*VHoverBuilder]
}

func VHover(children ...h.HTMLComponent) *VHoverBuilder {
	return VTag(&VHoverBuilder{}, "v-hover", children...)
}

func (b *VHoverBuilder) Disabled(v bool) (r *VHoverBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VHoverBuilder) ModelValue(v bool) (r *VHoverBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VHoverBuilder) CloseDelay(v interface{}) (r *VHoverBuilder) {
	b.Attr(":close-delay", h.JSONString(v))
	return b
}

func (b *VHoverBuilder) OpenDelay(v interface{}) (r *VHoverBuilder) {
	b.Attr(":open-delay", h.JSONString(v))
	return b
}

func (b *VHoverBuilder) On(name string, value string) (r *VHoverBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VHoverBuilder) Bind(name string, value string) (r *VHoverBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VHoverBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VHoverBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VHoverBuilder) Slot(name string, child ...h.HTMLComponent) (r *VHoverBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VHoverBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VHoverBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VHoverBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VHoverBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VHoverBuilder) SlotDefault(child ...h.HTMLComponent) (r *VHoverBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VHoverBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VHoverBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
