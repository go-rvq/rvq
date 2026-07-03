package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDatePickerHeaderBuilder struct {
	VTagBuilder[*VDatePickerHeaderBuilder]
}

func VDatePickerHeader(children ...h.HTMLComponent) *VDatePickerHeaderBuilder {
	return VTag(&VDatePickerHeaderBuilder{}, "v-date-picker-header", children...)
}

func (b *VDatePickerHeaderBuilder) Header(v string) (r *VDatePickerHeaderBuilder) {
	b.Attr("header", v)
	return b
}

func (b *VDatePickerHeaderBuilder) Color(v string) (r *VDatePickerHeaderBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDatePickerHeaderBuilder) Transition(v string) (r *VDatePickerHeaderBuilder) {
	b.Attr("transition", v)
	return b
}

func (b *VDatePickerHeaderBuilder) AppendIcon(v interface{}) (r *VDatePickerHeaderBuilder) {
	b.Attr(":append-icon", h.JSONString(v))
	return b
}

func (b *VDatePickerHeaderBuilder) On(name string, value string) (r *VDatePickerHeaderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDatePickerHeaderBuilder) Bind(name string, value string) (r *VDatePickerHeaderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDatePickerHeaderBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDatePickerHeaderBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDatePickerHeaderBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDatePickerHeaderBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDatePickerHeaderBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDatePickerHeaderBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDatePickerHeaderBuilder) SetSlotPrepend(child ...h.HTMLComponent) {
	b.SetSlot("prepend", child...)
}

func (b *VDatePickerHeaderBuilder) SetScopedSlotPrepend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("prepend", scope, child...)
}

func (b *VDatePickerHeaderBuilder) SlotPrepend(child ...h.HTMLComponent) (r *VDatePickerHeaderBuilder) {
	b.SetSlotPrepend(child...)
	return b
}

func (b *VDatePickerHeaderBuilder) ScopedSlotPrepend(scope string, child ...h.HTMLComponent) (r *VDatePickerHeaderBuilder) {
	b.SetScopedSlotPrepend(scope, child...)
	return b
}

func (b *VDatePickerHeaderBuilder) SetSlotDefault(child ...h.HTMLComponent) {
	b.SetSlot("default", child...)
}

func (b *VDatePickerHeaderBuilder) SetScopedSlotDefault(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("default", scope, child...)
}

func (b *VDatePickerHeaderBuilder) SlotDefault(child ...h.HTMLComponent) (r *VDatePickerHeaderBuilder) {
	b.SetSlotDefault(child...)
	return b
}

func (b *VDatePickerHeaderBuilder) ScopedSlotDefault(scope string, child ...h.HTMLComponent) (r *VDatePickerHeaderBuilder) {
	b.SetScopedSlotDefault(scope, child...)
	return b
}

func (b *VDatePickerHeaderBuilder) SetSlotAppend(child ...h.HTMLComponent) {
	b.SetSlot("append", child...)
}

func (b *VDatePickerHeaderBuilder) SetScopedSlotAppend(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("append", scope, child...)
}

func (b *VDatePickerHeaderBuilder) SlotAppend(child ...h.HTMLComponent) (r *VDatePickerHeaderBuilder) {
	b.SetSlotAppend(child...)
	return b
}

func (b *VDatePickerHeaderBuilder) ScopedSlotAppend(scope string, child ...h.HTMLComponent) (r *VDatePickerHeaderBuilder) {
	b.SetScopedSlotAppend(scope, child...)
	return b
}
