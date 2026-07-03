package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VPickerTitleBuilder struct {
	VTagBuilder[*VPickerTitleBuilder]
}

func VPickerTitle(children ...h.HTMLComponent) *VPickerTitleBuilder {
	return VTag(&VPickerTitleBuilder{}, "v-picker-title", children...)
}

func (b *VPickerTitleBuilder) Tag(v string) (r *VPickerTitleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VPickerTitleBuilder) On(name string, value string) (r *VPickerTitleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VPickerTitleBuilder) Bind(name string, value string) (r *VPickerTitleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VPickerTitleBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VPickerTitleBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VPickerTitleBuilder) Slot(name string, child ...h.HTMLComponent) (r *VPickerTitleBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VPickerTitleBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VPickerTitleBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VPickerTitleBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VPickerTitleBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VPickerTitleBuilder) SlotDefault(child ...h.HTMLComponent) (r *VPickerTitleBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VPickerTitleBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VPickerTitleBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}
