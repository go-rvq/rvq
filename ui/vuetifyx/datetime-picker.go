package vuetifyx

import (
	"fmt"

	v "github.com/qor5/x/v3/ui/vuetify"

	h "github.com/theplant/htmlgo"
)

type VXDateTimePickerBuilder struct {
	v.VTagBuilder[*VXDateTimePickerBuilder]
}

func VXDateTimePicker() (r *VXDateTimePickerBuilder) {
	return v.VTag(&VXDateTimePickerBuilder{}, "vx-datetimepicker")
}

type DatePickerProps struct{}

type TimePickerProps struct {
	Format     string `json:"format"` // 可用的选项是 ampm 和 24hr
	Scrollable bool   `json:"scrollable"`
	UseSeconds bool   `json:"use-seconds"`
	NoTitle    bool   `json:"no-title"`
}

func (b *VXDateTimePickerBuilder) Value(v string) (r *VXDateTimePickerBuilder) {
	b.Attr(":value", h.JSONString(v))
	return b
}

func (b *VXDateTimePickerBuilder) Label(v string) (r *VXDateTimePickerBuilder) {
	b.Attr(":label", h.JSONString(v))
	return b
}

func (b *VXDateTimePickerBuilder) DialogWidth(v int) (r *VXDateTimePickerBuilder) {
	b.Attr(":dialogWidth", h.JSONString(v))
	return b
}

func (b *VXDateTimePickerBuilder) DateFormat(v string) (r *VXDateTimePickerBuilder) {
	b.Attr(":dateFormat", h.JSONString(v))
	return b
}

func (b *VXDateTimePickerBuilder) TimeFormat(v string) (r *VXDateTimePickerBuilder) {
	b.Attr(":timeFormat", h.JSONString(v))
	return b
}

func (b *VXDateTimePickerBuilder) ClearText(v string) (r *VXDateTimePickerBuilder) {
	b.Attr(":clearText", h.JSONString(v))
	return b
}

func (b *VXDateTimePickerBuilder) OkText(v string) (r *VXDateTimePickerBuilder) {
	b.Attr(":okText", h.JSONString(v))
	return b
}

func (b *VXDateTimePickerBuilder) Disabled(v bool) (r *VXDateTimePickerBuilder) {
	b.Attr(":disabled", h.JSONString(v))
	return b
}

func (b *VXDateTimePickerBuilder) DatePickerProps(v DatePickerProps) (r *VXDateTimePickerBuilder) {
	b.Attr(":datePickerProps", h.JSONString(v))
	return b
}

func (b *VXDateTimePickerBuilder) TimePickerProps(v TimePickerProps) (r *VXDateTimePickerBuilder) {
	b.Attr(":timePickerProps", h.JSONString(v))
	return b
}

func (b *VXDateTimePickerBuilder) HideDetails(v bool) (r *VXDateTimePickerBuilder) {
	b.Attr(":hide-details", fmt.Sprint(v))
	return b
}
