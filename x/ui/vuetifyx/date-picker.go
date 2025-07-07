package vuetifyx

import (
	"fmt"

	v "github.com/go-rvq/rvq/x/ui/vuetify"

	h "github.com/theplant/htmlgo"
)

type VXDatePickerBuilder struct {
	v.VTagBuilder[*VXDatePickerBuilder]
}

func VXDatePicker() (r *VXDatePickerBuilder) {
	return v.VTag(&VXDatePickerBuilder{}, "vx-datepicker")
}

func (b *VXDatePickerBuilder) Value(v string) (r *VXDatePickerBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VXDatePickerBuilder) Label(v string) (r *VXDatePickerBuilder) {
	b.Attr(":label", h.JSONString(v))
	return b
}

func (b *VXDatePickerBuilder) DialogWidth(v int) (r *VXDatePickerBuilder) {
	b.Attr(":dialogWidth", h.JSONString(v))
	return b
}

func (b *VXDatePickerBuilder) DateFormat(v string) (r *VXDatePickerBuilder) {
	b.Attr(":dateFormat", h.JSONString(v))
	return b
}

func (b *VXDatePickerBuilder) TimeFormat(v string) (r *VXDatePickerBuilder) {
	b.Attr(":timeFormat", h.JSONString(v))
	return b
}

func (b *VXDatePickerBuilder) ClearText(v string) (r *VXDatePickerBuilder) {
	b.Attr(":clearText", h.JSONString(v))
	return b
}

func (b *VXDatePickerBuilder) OkText(v string) (r *VXDatePickerBuilder) {
	b.Attr(":okText", h.JSONString(v))
	return b
}

func (b *VXDatePickerBuilder) Disabled(v bool) (r *VXDatePickerBuilder) {
	b.Attr(":disabled", h.JSONString(v))
	return b
}

func (b *VXDatePickerBuilder) DatePickerProps(v DatePickerProps) (r *VXDatePickerBuilder) {
	b.Attr(":datePickerProps", h.JSONString(v))
	return b
}

func (b *VXDatePickerBuilder) TimePickerProps(v TimePickerProps) (r *VXDatePickerBuilder) {
	b.Attr(":timePickerProps", h.JSONString(v))
	return b
}

func (b *VXDatePickerBuilder) HideDetails(v bool) (r *VXDatePickerBuilder) {
	b.Attr(":hide-details", fmt.Sprint(v))
	return b
}
