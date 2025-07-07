package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VTimePickerBuilder struct {
	VTagBuilder[*VTimePickerBuilder]
}

func VTimePicker(children ...h.HTMLComponent) *VTimePickerBuilder {
	return VTag(&VTimePickerBuilder{}, "v-time-picker", children...)
}

func (b *VTimePickerBuilder) Title(v string) (r *VTimePickerBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VTimePickerBuilder) AmpmInTitle(v bool) (r *VTimePickerBuilder) {
	b.Attr(":ampm-in-title", fmt.Sprint(v))
	return b
}

func (b *VTimePickerBuilder) Disabled(v bool) (r *VTimePickerBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VTimePickerBuilder) Format(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":format", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) Max(v string) (r *VTimePickerBuilder) {
	b.Attr("max", v)
	return b
}

func (b *VTimePickerBuilder) Min(v string) (r *VTimePickerBuilder) {
	b.Attr("min", v)
	return b
}

func (b *VTimePickerBuilder) Readonly(v bool) (r *VTimePickerBuilder) {
	b.Attr(":readonly", fmt.Sprint(v))
	return b
}

func (b *VTimePickerBuilder) Scrollable(v bool) (r *VTimePickerBuilder) {
	b.Attr(":scrollable", fmt.Sprint(v))
	return b
}

func (b *VTimePickerBuilder) UseSeconds(v bool) (r *VTimePickerBuilder) {
	b.Attr(":use-seconds", fmt.Sprint(v))
	return b
}

func (b *VTimePickerBuilder) BgColor(v string) (r *VTimePickerBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VTimePickerBuilder) HideHeader(v bool) (r *VTimePickerBuilder) {
	b.Attr(":hide-header", fmt.Sprint(v))
	return b
}

func (b *VTimePickerBuilder) Color(v string) (r *VTimePickerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VTimePickerBuilder) Border(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) Height(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) MaxHeight(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) MaxWidth(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) MinHeight(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) MinWidth(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) Width(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) Elevation(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) Location(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) Position(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) Rounded(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) Tile(v bool) (r *VTimePickerBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VTimePickerBuilder) Tag(v string) (r *VTimePickerBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VTimePickerBuilder) Theme(v string) (r *VTimePickerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VTimePickerBuilder) AllowedHours(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":allowed-hours", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) AllowedMinutes(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":allowed-minutes", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) AllowedSeconds(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":allowed-seconds", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) ModelValue(v interface{}) (r *VTimePickerBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VTimePickerBuilder) On(name string, value string) (r *VTimePickerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VTimePickerBuilder) Bind(name string, value string) (r *VTimePickerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
