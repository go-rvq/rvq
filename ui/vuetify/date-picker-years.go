package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VDatePickerYearsBuilder struct {
	VTagBuilder[*VDatePickerYearsBuilder]
}

func VDatePickerYears(children ...h.HTMLComponent) *VDatePickerYearsBuilder {
	return VTag(&VDatePickerYearsBuilder{}, "v-date-picker-years", children...)
}

func (b *VDatePickerYearsBuilder) Color(v string) (r *VDatePickerYearsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDatePickerYearsBuilder) Height(v interface{}) (r *VDatePickerYearsBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VDatePickerYearsBuilder) ModelValue(v int) (r *VDatePickerYearsBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VDatePickerYearsBuilder) Min(v interface{}) (r *VDatePickerYearsBuilder) {
	b.Attr(":min", h.JSONString(v))
	return b
}

func (b *VDatePickerYearsBuilder) Max(v interface{}) (r *VDatePickerYearsBuilder) {
	b.Attr(":max", h.JSONString(v))
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
