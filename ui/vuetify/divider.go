package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VDividerBuilder struct {
	VTagBuilder[*VDividerBuilder]
}

func VDivider(children ...h.HTMLComponent) *VDividerBuilder {
	return VTag(&VDividerBuilder{}, "v-divider", children...)
}

func (b *VDividerBuilder) Length(v interface{}) (r *VDividerBuilder) {
	b.Attr(":length", h.JSONString(v))
	return b
}

func (b *VDividerBuilder) Color(v string) (r *VDividerBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VDividerBuilder) Inset(v bool) (r *VDividerBuilder) {
	b.Attr(":inset", fmt.Sprint(v))
	return b
}

func (b *VDividerBuilder) Opacity(v interface{}) (r *VDividerBuilder) {
	b.Attr(":opacity", h.JSONString(v))
	return b
}

func (b *VDividerBuilder) Thickness(v interface{}) (r *VDividerBuilder) {
	b.Attr(":thickness", h.JSONString(v))
	return b
}

func (b *VDividerBuilder) Vertical(v bool) (r *VDividerBuilder) {
	b.Attr(":vertical", fmt.Sprint(v))
	return b
}

func (b *VDividerBuilder) Theme(v string) (r *VDividerBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VDividerBuilder) On(name string, value string) (r *VDividerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VDividerBuilder) Bind(name string, value string) (r *VDividerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
