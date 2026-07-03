package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDatePickerMonthsBuilder struct {
	VTagBuilder[*VDatePickerMonthsBuilder]
}

func VDatePickerMonths(children ...h.HTMLComponent) *VDatePickerMonthsBuilder {
	return VTag(&VDatePickerMonthsBuilder{}, "v-date-picker-months", children...)
}

func (b *VDatePickerMonthsBuilder) Height(v interface{}) (r *VDatePickerMonthsBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthsBuilder) Color(v string) (r *VDatePickerMonthsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDatePickerMonthsBuilder) ModelValue(v interface{}) (r *VDatePickerMonthsBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthsBuilder) Max(v interface{}) (r *VDatePickerMonthsBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthsBuilder) Min(v interface{}) (r *VDatePickerMonthsBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthsBuilder) Year(v interface{}) (r *VDatePickerMonthsBuilder) {
	b.Attr(":year", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthsBuilder) On(name string, value string) (r *VDatePickerMonthsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDatePickerMonthsBuilder) Bind(name string, value string) (r *VDatePickerMonthsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDatePickerMonthsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDatePickerMonthsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDatePickerMonthsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDatePickerMonthsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDatePickerMonthsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDatePickerMonthsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDatePickerMonthsBuilder) SetSlotMonth(child ...h.HTMLComponent) {
	b.SetSlot("month", child...)
}

func (b *VDatePickerMonthsBuilder) SetScopedSlotMonth(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("month", scope, child...)
}

func (b *VDatePickerMonthsBuilder) SlotMonth(child ...h.HTMLComponent) (r *VDatePickerMonthsBuilder) {
	b.SetSlotMonth(child...)
	return b
}

func (b *VDatePickerMonthsBuilder) ScopedSlotMonth(scope string, child ...h.HTMLComponent) (r *VDatePickerMonthsBuilder) {
	b.SetScopedSlotMonth(scope, child...)
	return b
}
