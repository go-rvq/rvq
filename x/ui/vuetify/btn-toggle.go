package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBtnToggleBuilder struct {
	VTagBuilder[*VBtnToggleBuilder]
}

func VBtnToggle(children ...h.HTMLComponent) *VBtnToggleBuilder {
	return VTag(&VBtnToggleBuilder{}, "v-btn-toggle", children...)
}

func (b *VBtnToggleBuilder) BaseColor(v string) (r *VBtnToggleBuilder) {
	b.Attr("base-color", v)
	return b
}

func (b *VBtnToggleBuilder) Divided(v bool) (r *VBtnToggleBuilder) {
	b.Attr(":divided", fmt.Sprint(v))
	return b
}

func (b *VBtnToggleBuilder) Border(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Density(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Elevation(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Rounded(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Tile(v bool) (r *VBtnToggleBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBtnToggleBuilder) Tag(v string) (r *VBtnToggleBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VBtnToggleBuilder) Theme(v string) (r *VBtnToggleBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VBtnToggleBuilder) Color(v string) (r *VBtnToggleBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBtnToggleBuilder) Variant(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) ModelValue(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":model-value", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) Multiple(v bool) (r *VBtnToggleBuilder) {
	b.Attr(":multiple", fmt.Sprint(v))
	return b
}

func (b *VBtnToggleBuilder) Max(v int) (r *VBtnToggleBuilder) {
	b.Attr(":max", fmt.Sprint(v))
	return b
}

func (b *VBtnToggleBuilder) SelectedClass(v string) (r *VBtnToggleBuilder) {
	b.Attr("selected-class", v)
	return b
}

func (b *VBtnToggleBuilder) Disabled(v bool) (r *VBtnToggleBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VBtnToggleBuilder) Mandatory(v interface{}) (r *VBtnToggleBuilder) {
	b.Attr(":mandatory", h.JSONString(v))
	return b
}

func (b *VBtnToggleBuilder) On(name string, value string) (r *VBtnToggleBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBtnToggleBuilder) Bind(name string, value string) (r *VBtnToggleBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
