package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VDatePickerYearsBuilder struct {
	VTagBuilder[*VDatePickerYearsBuilder]
}

func VDatePickerYears(children ...h.HTMLComponent) *VDatePickerYearsBuilder {
	return VTag(&VDatePickerYearsBuilder{}, "v-date-picker-years", children...)
}

func (b *VDatePickerYearsBuilder) Height(v interface{}) (r *VDatePickerYearsBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDatePickerYearsBuilder) Color(v string) (r *VDatePickerYearsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDatePickerYearsBuilder) ModelValue(v interface{}) (r *VDatePickerYearsBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VDatePickerYearsBuilder) Max(v interface{}) (r *VDatePickerYearsBuilder) {
	b.Attr(":max", h.JSONString(v))
	return b
}

func (b *VDatePickerYearsBuilder) Min(v interface{}) (r *VDatePickerYearsBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VDatePickerYearsBuilder) On(name string, value string) (r *VDatePickerYearsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDatePickerYearsBuilder) Bind(name string, value string) (r *VDatePickerYearsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}

func (b *VDatePickerYearsBuilder) SetSlot(name string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, true).Attr())
}

func (b *VDatePickerYearsBuilder) SetScopedSlot(name string, scope string, child ...h.HTMLComponent) {
	b.AppendChild(h.Template(child...).Attr("v-slot:"+name, scope).Attr())
}

func (b *VDatePickerYearsBuilder) Slot(name string, child ...h.HTMLComponent) (r *VDatePickerYearsBuilder) {
	b.SetSlot(name, child...)
	return b
}

func (b *VDatePickerYearsBuilder) ScopedSlot(name string, scope string, child ...h.HTMLComponent) (r *VDatePickerYearsBuilder) {
	b.SetScopedSlot(name, scope, child...)
	return b
}

func (b *VDatePickerYearsBuilder) SetSlotYear(child ...h.HTMLComponent) {
	b.SetSlot("year", child...)
}

func (b *VDatePickerYearsBuilder) SetScopedSlotYear(scope string, child ...h.HTMLComponent) {
	b.SetScopedSlot("year", scope, child...)
}

func (b *VDatePickerYearsBuilder) SlotYear(child ...h.HTMLComponent) (r *VDatePickerYearsBuilder) {
	b.SetSlotYear(child...)
	return b
}

func (b *VDatePickerYearsBuilder) ScopedSlotYear(scope string, child ...h.HTMLComponent) (r *VDatePickerYearsBuilder) {
	b.SetScopedSlotYear(scope, child...)
	return b
}
