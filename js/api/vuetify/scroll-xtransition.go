package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VScrollXtransitionBuilder struct {
	VTagBuilder[*VScrollXtransitionBuilder]
}

func VScrollXtransition(children ...h.HTMLComponent) *VScrollXtransitionBuilder {
	return VTag(&VScrollXtransitionBuilder{}, "v-scroll-xtransition", children...)
}

func (b *VScrollXtransitionBuilder) Disabled(v bool) (r *VScrollXtransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VScrollXtransitionBuilder) Mode(v string) (r *VScrollXtransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VScrollXtransitionBuilder) Group(v bool) (r *VScrollXtransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VScrollXtransitionBuilder) Origin(v string) (r *VScrollXtransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VScrollXtransitionBuilder) HideOnLeave(v bool) (r *VScrollXtransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VScrollXtransitionBuilder) LeaveAbsolute(v bool) (r *VScrollXtransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VScrollXtransitionBuilder) On(name string, value string) (r *VScrollXtransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VScrollXtransitionBuilder) Bind(name string, value string) (r *VScrollXtransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VScrollXtransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VScrollXtransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VScrollXtransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VScrollXtransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VScrollXtransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VScrollXtransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VScrollXtransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VScrollXtransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VScrollXtransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VScrollXtransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VScrollXtransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VScrollXtransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
