package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VIconBuilder struct {
	VTagBuilder[*VIconBuilder]
}

func VIcon(name string) *VIconBuilder {
	return VTag(&VIconBuilder{}, "v-icon").Icon(name)
}

func (b *VIconBuilder) Color(v string) (r *VIconBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VIconBuilder) Disabled(v bool) (r *VIconBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VIconBuilder) Start(v bool) (r *VIconBuilder) {
	b.Attr(":start", fmt.Sprint(v))
	return b
}

func (b *VIconBuilder) End(v bool) (r *VIconBuilder) {
	b.Attr(":end", fmt.Sprint(v))
	return b
}

func (b *VIconBuilder) Icon(v interface{}) (r *VIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VIconBuilder) Size(v interface{}) (r *VIconBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VIconBuilder) Tag(v string) (r *VIconBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VIconBuilder) Theme(v string) (r *VIconBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VIconBuilder) On(name string, value string) (r *VIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VIconBuilder) Bind(name string, value string) (r *VIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
