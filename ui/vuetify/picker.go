package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VPickerBuilder struct {
	VTagBuilder[*VPickerBuilder]
}

func VPicker(children ...h.HTMLComponent) *VPickerBuilder {
	return VTag(&VPickerBuilder{}, "v-picker", children...)
}

func (b *VPickerBuilder) Title(v string) (r *VPickerBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VPickerBuilder) BgColor(v string) (r *VPickerBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VPickerBuilder) Landscape(v bool) (r *VPickerBuilder) {
	b.Attr(":landscape", fmt.Sprint(v))
	return b
}

func (b *VPickerBuilder) HideHeader(v bool) (r *VPickerBuilder) {
	b.Attr(":hide-header", fmt.Sprint(v))
	return b
}

func (b *VPickerBuilder) Color(v string) (r *VPickerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VPickerBuilder) Border(v interface{}) (r *VPickerBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Height(v interface{}) (r *VPickerBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) MaxHeight(v interface{}) (r *VPickerBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) MaxWidth(v interface{}) (r *VPickerBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) MinHeight(v interface{}) (r *VPickerBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) MinWidth(v interface{}) (r *VPickerBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Width(v interface{}) (r *VPickerBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Elevation(v interface{}) (r *VPickerBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Location(v interface{}) (r *VPickerBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Position(v interface{}) (r *VPickerBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Rounded(v interface{}) (r *VPickerBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VPickerBuilder) Tile(v bool) (r *VPickerBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VPickerBuilder) Tag(v string) (r *VPickerBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VPickerBuilder) Theme(v string) (r *VPickerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VPickerBuilder) On(name string, value string) (r *VPickerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VPickerBuilder) Bind(name string, value string) (r *VPickerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
