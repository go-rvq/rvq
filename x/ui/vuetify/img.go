package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VImgBuilder struct {
	VTagBuilder[*VImgBuilder]
}

func VImg(children ...h.HTMLComponent) *VImgBuilder {
	return VTag(&VImgBuilder{}, "v-img", children...)
}

func (b *VImgBuilder) Alt(v string) (r *VImgBuilder) {
	b.Attr("alt", v)
	return b
}

func (b *VImgBuilder) Cover(v bool) (r *VImgBuilder) {
	b.Attr(":cover", fmt.Sprint(v))
	return b
}

func (b *VImgBuilder) Color(v string) (r *VImgBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VImgBuilder) Draggable(v interface{}) (r *VImgBuilder) {
	b.Attr(":draggable", h.JSONString(v))
	return b
}

func (b *VImgBuilder) Eager(v bool) (r *VImgBuilder) {
	b.Attr(":eager", fmt.Sprint(v))
	return b
}

func (b *VImgBuilder) Gradient(v string) (r *VImgBuilder) {
	b.Attr("gradient", v)
	return b
}

func (b *VImgBuilder) LazySrc(v string) (r *VImgBuilder) {
	b.Attr("lazy-src", v)
	return b
}

func (b *VImgBuilder) Options(v interface{}) (r *VImgBuilder) {
	b.Attr(":options", h.JSONString(v))
	return b
}

func (b *VImgBuilder) Sizes(v string) (r *VImgBuilder) {
	b.Attr("sizes", v)
	return b
}

func (b *VImgBuilder) Src(v interface{}) (r *VImgBuilder) {
	b.Attr(":src", h.JSONString(v))
	return b
}

func (b *VImgBuilder) Srcset(v string) (r *VImgBuilder) {
	b.Attr("srcset", v)
	return b
}

func (b *VImgBuilder) Position(v string) (r *VImgBuilder) {
	b.Attr("position", v)
	return b
}

func (b *VImgBuilder) AspectRatio(v interface{}) (r *VImgBuilder) {
	b.Attr(":aspect-ratio", h.JSONString(v))
	return b
}

func (b *VImgBuilder) ContentClass(v interface{}) (r *VImgBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VImgBuilder) Inline(v bool) (r *VImgBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VImgBuilder) Height(v interface{}) (r *VImgBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VImgBuilder) MaxHeight(v interface{}) (r *VImgBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VImgBuilder) MaxWidth(v interface{}) (r *VImgBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VImgBuilder) MinHeight(v interface{}) (r *VImgBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VImgBuilder) MinWidth(v interface{}) (r *VImgBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VImgBuilder) Width(v interface{}) (r *VImgBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VImgBuilder) Rounded(v interface{}) (r *VImgBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VImgBuilder) Tile(v bool) (r *VImgBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VImgBuilder) Transition(v interface{}) (r *VImgBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VImgBuilder) Crossorigin(v interface{}) (r *VImgBuilder) {
	b.Attr(":crossorigin", h.JSONString(v))
	return b
}

func (b *VImgBuilder) Referrerpolicy(v interface{}) (r *VImgBuilder) {
	b.Attr(":referrerpolicy", h.JSONString(v))
	return b
}

func (b *VImgBuilder) On(name string, value string) (r *VImgBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VImgBuilder) Bind(name string, value string) (r *VImgBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
