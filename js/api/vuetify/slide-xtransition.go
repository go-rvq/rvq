package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSlideXtransitionBuilder struct {
	VTagBuilder[*VSlideXtransitionBuilder]
}

func VSlideXtransition(children ...h.HTMLComponent) *VSlideXtransitionBuilder {
	return VTag(&VSlideXtransitionBuilder{}, "v-slide-xtransition", children...)
}

func (b *VSlideXtransitionBuilder) Mode(v string) (r *VSlideXtransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VSlideXtransitionBuilder) Disabled(v bool) (r *VSlideXtransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideXtransitionBuilder) Origin(v string) (r *VSlideXtransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VSlideXtransitionBuilder) Group(v bool) (r *VSlideXtransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VSlideXtransitionBuilder) HideOnLeave(v bool) (r *VSlideXtransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VSlideXtransitionBuilder) LeaveAbsolute(v bool) (r *VSlideXtransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VSlideXtransitionBuilder) On(name string, value string) (r *VSlideXtransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideXtransitionBuilder) Bind(name string, value string) (r *VSlideXtransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSlideXtransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSlideXtransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSlideXtransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSlideXtransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSlideXtransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSlideXtransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSlideXtransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSlideXtransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSlideXtransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSlideXtransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSlideXtransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSlideXtransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
