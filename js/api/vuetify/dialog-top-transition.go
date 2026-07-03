package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDialogTopTransitionBuilder struct {
	VTagBuilder[*VDialogTopTransitionBuilder]
}

func VDialogTopTransition(children ...h.HTMLComponent) *VDialogTopTransitionBuilder {
	return VTag(&VDialogTopTransitionBuilder{}, "v-dialog-top-transition", children...)
}

func (b *VDialogTopTransitionBuilder) Mode(v string) (r *VDialogTopTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VDialogTopTransitionBuilder) Disabled(v bool) (r *VDialogTopTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VDialogTopTransitionBuilder) Origin(v string) (r *VDialogTopTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VDialogTopTransitionBuilder) Group(v bool) (r *VDialogTopTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VDialogTopTransitionBuilder) HideOnLeave(v bool) (r *VDialogTopTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VDialogTopTransitionBuilder) LeaveAbsolute(v bool) (r *VDialogTopTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VDialogTopTransitionBuilder) On(name string, value string) (r *VDialogTopTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDialogTopTransitionBuilder) Bind(name string, value string) (r *VDialogTopTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDialogTopTransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDialogTopTransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDialogTopTransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDialogTopTransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDialogTopTransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDialogTopTransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDialogTopTransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDialogTopTransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDialogTopTransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDialogTopTransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDialogTopTransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDialogTopTransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
