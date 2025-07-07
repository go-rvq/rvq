package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAvatarBuilder struct {
	VTagBuilder[*VAvatarBuilder]
}

func VAvatar(children ...h.HTMLComponent) *VAvatarBuilder {
	return VTag(&VAvatarBuilder{}, "v-avatar", children...)
}

func (b *VAvatarBuilder) Start(v bool) (r *VAvatarBuilder) {
	b.Attr(":start", fmt.Sprint(v))
	return b
}

func (b *VAvatarBuilder) End(v bool) (r *VAvatarBuilder) {
	b.Attr(":end", fmt.Sprint(v))
	return b
}

func (b *VAvatarBuilder) Icon(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) Image(v string) (r *VAvatarBuilder) {
	b.Attr("image", v)
	return b
}

func (b *VAvatarBuilder) Text(v string) (r *VAvatarBuilder) {
	b.Attr("text", v)
	return b
}

func (b *VAvatarBuilder) Density(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) Rounded(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) Tile(v bool) (r *VAvatarBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VAvatarBuilder) Size(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":size", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) Tag(v string) (r *VAvatarBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VAvatarBuilder) Theme(v string) (r *VAvatarBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VAvatarBuilder) Color(v string) (r *VAvatarBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VAvatarBuilder) Variant(v interface{}) (r *VAvatarBuilder) {
	b.Attr(":variant", h.JSONString(v))
	return b
}

func (b *VAvatarBuilder) On(name string, value string) (r *VAvatarBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAvatarBuilder) Bind(name string, value string) (r *VAvatarBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
