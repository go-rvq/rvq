package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSlideYtransitionBuilder struct {
	VTagBuilder[*VSlideYtransitionBuilder]
}

func VSlideYtransition(children ...h.HTMLComponent) *VSlideYtransitionBuilder {
	return VTag(&VSlideYtransitionBuilder{}, "v-slide-ytransition", children...)
}

func (b *VSlideYtransitionBuilder) Mode(v string) (r *VSlideYtransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VSlideYtransitionBuilder) Disabled(v bool) (r *VSlideYtransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VSlideYtransitionBuilder) Origin(v string) (r *VSlideYtransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VSlideYtransitionBuilder) Group(v bool) (r *VSlideYtransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VSlideYtransitionBuilder) HideOnLeave(v bool) (r *VSlideYtransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VSlideYtransitionBuilder) LeaveAbsolute(v bool) (r *VSlideYtransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VSlideYtransitionBuilder) On(name string, value string) (r *VSlideYtransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSlideYtransitionBuilder) Bind(name string, value string) (r *VSlideYtransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VSlideYtransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VSlideYtransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VSlideYtransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VSlideYtransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VSlideYtransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VSlideYtransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VSlideYtransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VSlideYtransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VSlideYtransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VSlideYtransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VSlideYtransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VSlideYtransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
