package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAppBarBuilder struct {
	VTagBuilder[*VAppBarBuilder]
}

func VAppBar(children ...h.HTMLComponent) *VAppBarBuilder {
	return VTag(&VAppBarBuilder{}, "v-app-bar", children...)
}

func (b *VAppBarBuilder) Image(v string) (r *VAppBarBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VAppBarBuilder) Title(v string) (r *VAppBarBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VAppBarBuilder) Flat(v bool) (r *VAppBarBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Collapse(v bool) (r *VAppBarBuilder) {
	b.Attr(":collapse", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) ModelValue(v bool) (r *VAppBarBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Location(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":location", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Absolute(v bool) (r *VAppBarBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Color(v string) (r *VAppBarBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VAppBarBuilder) Density(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Extended(v bool) (r *VAppBarBuilder) {
	b.Attr(":extended", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) ExtensionHeight(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":extension-height", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Floating(v bool) (r *VAppBarBuilder) {
	b.Attr(":floating", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Height(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Border(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Elevation(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Rounded(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) Tile(v bool) (r *VAppBarBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VAppBarBuilder) Tag(v string) (r *VAppBarBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VAppBarBuilder) Theme(v string) (r *VAppBarBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VAppBarBuilder) Name(v string) (r *VAppBarBuilder) {
	b.Attr("name", v)
	return b
}

func (b *VAppBarBuilder) Order(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) ScrollTarget(v string) (r *VAppBarBuilder) {
	b.Attr("scroll-target", v)
	return b
}

func (b *VAppBarBuilder) ScrollThreshold(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":scroll-threshold", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) ScrollBehavior(v interface{}) (r *VAppBarBuilder) {
	b.Attr(":scroll-behavior", h.JSONString(v))
	return b
}

func (b *VAppBarBuilder) On(name string, value string) (r *VAppBarBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAppBarBuilder) Bind(name string, value string) (r *VAppBarBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
