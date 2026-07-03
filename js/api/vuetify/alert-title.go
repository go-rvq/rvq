package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAlertTitleBuilder struct {
	VTagBuilder[*VAlertTitleBuilder]
}

func VAlertTitle(children ...h.HTMLComponent) *VAlertTitleBuilder {
	return VTag(&VAlertTitleBuilder{}, "v-alert-title", children...)
}

func (b *VAlertTitleBuilder) Tag(v string) (r *VAlertTitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VAlertTitleBuilder) On(name string, value string) (r *VAlertTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAlertTitleBuilder) Bind(name string, value string) (r *VAlertTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VAlertTitleBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VAlertTitleBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VAlertTitleBuilder) Slot(name string, child ...h.HTMLComponent) (r *VAlertTitleBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VAlertTitleBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VAlertTitleBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VAlertTitleBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VAlertTitleBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VAlertTitleBuilder) SlotDefault(child ...h.HTMLComponent) (r *VAlertTitleBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VAlertTitleBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VAlertTitleBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
