package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDialogTransitionBuilder struct {
	VTagBuilder[*VDialogTransitionBuilder]
}

func VDialogTransition(children ...h.HTMLComponent) *VDialogTransitionBuilder {
	return VTag(&VDialogTransitionBuilder{}, "v-dialog-transition", children...)
}

func (b *VDialogTransitionBuilder) Target(v interface{}) (r *VDialogTransitionBuilder) {
	b.Attr(":target", h.JSONString(v))
	return b
}

func (b *VDialogTransitionBuilder) On(name string, value string) (r *VDialogTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDialogTransitionBuilder) Bind(name string, value string) (r *VDialogTransitionBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDialogTransitionBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDialogTransitionBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDialogTransitionBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDialogTransitionBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDialogTransitionBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDialogTransitionBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDialogTransitionBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDialogTransitionBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDialogTransitionBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDialogTransitionBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDialogTransitionBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDialogTransitionBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
