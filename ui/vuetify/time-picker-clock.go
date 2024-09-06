package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VTimePickerClockBuilder struct {
	VTagBuilder[*VTimePickerClockBuilder]
}

func VTimePickerClock(children ...h.HTMLComponent) *VTimePickerClockBuilder {
	return VTag(&VTimePickerClockBuilder{}, "v-time-picker-clock", children...)
}

func (b *VTimePickerClockBuilder) Ampm(v bool) (r *VTimePickerClockBuilder) {
	b.Attr(":ampm", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Color(v string) (r *VTimePickerClockBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTimePickerClockBuilder) Disabled(v bool) (r *VTimePickerClockBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) DisplayedValue(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":displayed-value", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) Double(v bool) (r *VTimePickerClockBuilder) {
	b.Attr(":double", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Format(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":format", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) Max(v int) (r *VTimePickerClockBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Min(v int) (r *VTimePickerClockBuilder) {
	b.Attr(":min", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Scrollable(v bool) (r *VTimePickerClockBuilder) {
	b.Attr(":scrollable", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Readonly(v bool) (r *VTimePickerClockBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Rotate(v int) (r *VTimePickerClockBuilder) {
	b.Attr(":rotate", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) Step(v int) (r *VTimePickerClockBuilder) {
	b.Attr(":step", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) ModelValue(v int) (r *VTimePickerClockBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VTimePickerClockBuilder) AllowedValues(v interface{}) (r *VTimePickerClockBuilder) {
	b.Attr(":allowed-values", h.JSONString(v))
	return b
}

func (b *VTimePickerClockBuilder) On(name string, value string) (r *VTimePickerClockBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTimePickerClockBuilder) Bind(name string, value string) (r *VTimePickerClockBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
