package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VToolbarBuilder struct {
	VTagBuilder[*VToolbarBuilder]
}

func VToolbar(children ...h.HTMLComponent) *VToolbarBuilder {
	return VTag(&VToolbarBuilder{}, "v-toolbar", children...)
}

func (b *VToolbarBuilder) Image(v string) (r *VToolbarBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VToolbarBuilder) Title(v string) (r *VToolbarBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VToolbarBuilder) Flat(v bool) (r *VToolbarBuilder) {
	b.Attr(":flat", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) Absolute(v bool) (r *VToolbarBuilder) {
	b.Attr(":absolute", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) Collapse(v bool) (r *VToolbarBuilder) {
	b.Attr(":collapse", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) Color(v string) (r *VToolbarBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VToolbarBuilder) Density(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Extended(v bool) (r *VToolbarBuilder) {
	b.Attr(":extended", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) ExtensionHeight(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":extension-height", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Floating(v bool) (r *VToolbarBuilder) {
	b.Attr(":floating", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) Height(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Border(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":border", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Elevation(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Rounded(v interface{}) (r *VToolbarBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VToolbarBuilder) Tile(v bool) (r *VToolbarBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VToolbarBuilder) Tag(v string) (r *VToolbarBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VToolbarBuilder) Theme(v string) (r *VToolbarBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VToolbarBuilder) On(name string, value string) (r *VToolbarBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VToolbarBuilder) Bind(name string, value string) (r *VToolbarBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
