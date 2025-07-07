package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSheetBuilder struct {
	VTagBuilder[*VSheetBuilder]
}

func VSheet(children ...h.HTMLComponent) *VSheetBuilder {
	return VTag(&VSheetBuilder{}, "v-sheet", children...)
}

func (b *VSheetBuilder) Color(v string) (r *VSheetBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSheetBuilder) Border(v interface{}) (r *VSheetBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Height(v interface{}) (r *VSheetBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) MaxHeight(v interface{}) (r *VSheetBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) MaxWidth(v interface{}) (r *VSheetBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) MinHeight(v interface{}) (r *VSheetBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) MinWidth(v interface{}) (r *VSheetBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Width(v interface{}) (r *VSheetBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Elevation(v interface{}) (r *VSheetBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Location(v interface{}) (r *VSheetBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Position(v interface{}) (r *VSheetBuilder) {
	b.Attr(":position", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Rounded(v interface{}) (r *VSheetBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VSheetBuilder) Tile(v bool) (r *VSheetBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VSheetBuilder) Tag(v string) (r *VSheetBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VSheetBuilder) Theme(v string) (r *VSheetBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSheetBuilder) On(name string, value string) (r *VSheetBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSheetBuilder) Bind(name string, value string) (r *VSheetBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
