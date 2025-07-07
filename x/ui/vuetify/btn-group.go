package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBtnGroupBuilder struct {
	VTagBuilder[*VBtnGroupBuilder]
}

func VBtnGroup(children ...h.HTMLComponent) *VBtnGroupBuilder {
	return VTag(&VBtnGroupBuilder{}, "v-btn-group", children...)
}

func (b *VBtnGroupBuilder) BaseColor(v string) (r *VBtnGroupBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VBtnGroupBuilder) Divided(v bool) (r *VBtnGroupBuilder) {
	b.Attr(":divided", fmt.Sprint(v))
	return b
}

func (b *VBtnGroupBuilder) Border(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) Density(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) Elevation(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) Rounded(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) Tile(v bool) (r *VBtnGroupBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBtnGroupBuilder) Tag(v string) (r *VBtnGroupBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VBtnGroupBuilder) Theme(v string) (r *VBtnGroupBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBtnGroupBuilder) Color(v string) (r *VBtnGroupBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBtnGroupBuilder) Variant(v interface{}) (r *VBtnGroupBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VBtnGroupBuilder) On(name string, value string) (r *VBtnGroupBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBtnGroupBuilder) Bind(name string, value string) (r *VBtnGroupBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
