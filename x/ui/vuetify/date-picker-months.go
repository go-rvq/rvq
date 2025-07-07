package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VDatePickerMonthsBuilder struct {
	VTagBuilder[*VDatePickerMonthsBuilder]
}

func VDatePickerMonths(children ...h.HTMLComponent) *VDatePickerMonthsBuilder {
	return VTag(&VDatePickerMonthsBuilder{}, "v-date-picker-months", children...)
}

func (b *VDatePickerMonthsBuilder) Color(v string) (r *VDatePickerMonthsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDatePickerMonthsBuilder) Height(v interface{}) (r *VDatePickerMonthsBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthsBuilder) ModelValue(v int) (r *VDatePickerMonthsBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VDatePickerMonthsBuilder) Year(v int) (r *VDatePickerMonthsBuilder) {
	b.Attr(":year", fmt.Sprint(v))
	return b
}

func (b *VDatePickerMonthsBuilder) Min(v interface{}) (r *VDatePickerMonthsBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VDatePickerMonthsBuilder) Max(v interface{}) (r *VDatePickerMonthsBuilder) {
	b.Attr(":max", h.JSONString(v))
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
