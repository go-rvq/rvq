package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VSystemBarBuilder struct {
	VTagBuilder[*VSystemBarBuilder]
}

func VSystemBar(children ...h.HTMLComponent) *VSystemBarBuilder {
	return VTag(&VSystemBarBuilder{}, "v-system-bar", children...)
}

func (b *VSystemBarBuilder) Color(v string) (r *VSystemBarBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSystemBarBuilder) Height(v interface{}) (r *VSystemBarBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VSystemBarBuilder) Window(v bool) (r *VSystemBarBuilder) {
	b.Attr(":window", fmt.Sprint(v))
	return b
}

func (b *VSystemBarBuilder) Elevation(v interface{}) (r *VSystemBarBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VSystemBarBuilder) Name(v string) (r *VSystemBarBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VSystemBarBuilder) Order(v interface{}) (r *VSystemBarBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VSystemBarBuilder) Absolute(v bool) (r *VSystemBarBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VSystemBarBuilder) Rounded(v interface{}) (r *VSystemBarBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VSystemBarBuilder) Tile(v bool) (r *VSystemBarBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VSystemBarBuilder) Tag(v string) (r *VSystemBarBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VSystemBarBuilder) Theme(v string) (r *VSystemBarBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSystemBarBuilder) On(name string, value string) (r *VSystemBarBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSystemBarBuilder) Bind(name string, value string) (r *VSystemBarBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
