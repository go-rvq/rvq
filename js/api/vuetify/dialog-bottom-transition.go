package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDialogBottomTransitionBuilder struct {
	VTagBuilder[*VDialogBottomTransitionBuilder]
}

func VDialogBottomTransition(children ...h.HTMLComponent) *VDialogBottomTransitionBuilder {
	return VTag(&VDialogBottomTransitionBuilder{}, "v-dialog-bottom-transition", children...)
}

func (b *VDialogBottomTransitionBuilder) Mode(v string) (r *VDialogBottomTransitionBuilder) {
	b.Attr("mode", v)
	return b
}

func (b *VDialogBottomTransitionBuilder) Disabled(v bool) (r *VDialogBottomTransitionBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VDialogBottomTransitionBuilder) Origin(v string) (r *VDialogBottomTransitionBuilder) {
	b.Attr("origin", v)
	return b
}

func (b *VDialogBottomTransitionBuilder) Group(v bool) (r *VDialogBottomTransitionBuilder) {
	b.Attr(":group", fmt.Sprint(v))
	return b
}

func (b *VDialogBottomTransitionBuilder) HideOnLeave(v bool) (r *VDialogBottomTransitionBuilder) {
	b.Attr(":hide-on-leave", fmt.Sprint(v))
	return b
}

func (b *VDialogBottomTransitionBuilder) LeaveAbsolute(v bool) (r *VDialogBottomTransitionBuilder) {
	b.Attr(":leave-absolute", fmt.Sprint(v))
	return b
}

func (b *VDialogBottomTransitionBuilder) On(name string, value string) (r *VDialogBottomTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDialogBottomTransitionBuilder) Bind(name string, value string) (r *VDialogBottomTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDialogBottomTransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDialogBottomTransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDialogBottomTransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDialogBottomTransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDialogBottomTransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDialogBottomTransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDialogBottomTransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDialogBottomTransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDialogBottomTransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDialogBottomTransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDialogBottomTransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDialogBottomTransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
