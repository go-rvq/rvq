package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VScrollYtransitionBuilder struct {
	VTagBuilder[*VScrollYtransitionBuilder]
}

func VScrollYtransition(children ...h.HTMLComponent) *VScrollYtransitionBuilder {
	return VTag(&VScrollYtransitionBuilder{}, "v-scroll-ytransition", children...)
}

func (b *VScrollYtransitionBuilder) Disabled(v bool) (r *VScrollYtransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VScrollYtransitionBuilder) Mode(v string) (r *VScrollYtransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VScrollYtransitionBuilder) Group(v bool) (r *VScrollYtransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VScrollYtransitionBuilder) Origin(v string) (r *VScrollYtransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VScrollYtransitionBuilder) HideOnLeave(v bool) (r *VScrollYtransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VScrollYtransitionBuilder) LeaveAbsolute(v bool) (r *VScrollYtransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VScrollYtransitionBuilder) On(name string, value string) (r *VScrollYtransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VScrollYtransitionBuilder) Bind(name string, value string) (r *VScrollYtransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VScrollYtransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VScrollYtransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VScrollYtransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VScrollYtransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VScrollYtransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VScrollYtransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VScrollYtransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VScrollYtransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VScrollYtransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VScrollYtransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VScrollYtransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VScrollYtransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
