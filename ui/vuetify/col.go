package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VColBuilder struct {
	VTagBuilder[*VColBuilder]
}

func VCol(children ...h.HTMLComponent) *VColBuilder {
	return VTag(&VColBuilder{}, "v-col", children...)
}

func (b *VColBuilder) Cols(v interface{}) (r *VColBuilder) {
	b.Attr(":cols", h.JSONString(v))
	return b
}

func (b *VColBuilder) Sm(v interface{}) (r *VColBuilder) {
	b.Attr(":sm", h.JSONString(v))
	return b
}

func (b *VColBuilder) Md(v interface{}) (r *VColBuilder) {
	b.Attr(":md", h.JSONString(v))
	return b
}

func (b *VColBuilder) Lg(v interface{}) (r *VColBuilder) {
	b.Attr(":lg", h.JSONString(v))
	return b
}

func (b *VColBuilder) Xl(v interface{}) (r *VColBuilder) {
	b.Attr(":xl", h.JSONString(v))
	return b
}

func (b *VColBuilder) Xxl(v interface{}) (r *VColBuilder) {
	b.Attr(":xxl", h.JSONString(v))
	return b
}

func (b *VColBuilder) Offset(v interface{}) (r *VColBuilder) {
	b.Attr(":offset", h.JSONString(v))
	return b
}

func (b *VColBuilder) OffsetSm(v interface{}) (r *VColBuilder) {
	b.Attr(":offset-sm", h.JSONString(v))
	return b
}

func (b *VColBuilder) OffsetMd(v interface{}) (r *VColBuilder) {
	b.Attr(":offset-md", h.JSONString(v))
	return b
}

func (b *VColBuilder) OffsetLg(v interface{}) (r *VColBuilder) {
	b.Attr(":offset-lg", h.JSONString(v))
	return b
}

func (b *VColBuilder) OffsetXl(v interface{}) (r *VColBuilder) {
	b.Attr(":offset-xl", h.JSONString(v))
	return b
}

func (b *VColBuilder) OffsetXxl(v interface{}) (r *VColBuilder) {
	b.Attr(":offset-xxl", h.JSONString(v))
	return b
}

func (b *VColBuilder) Order(v interface{}) (r *VColBuilder) {
	b.Attr(":order", h.JSONString(v))
	return b
}

func (b *VColBuilder) OrderSm(v interface{}) (r *VColBuilder) {
	b.Attr(":order-sm", h.JSONString(v))
	return b
}

func (b *VColBuilder) OrderMd(v interface{}) (r *VColBuilder) {
	b.Attr(":order-md", h.JSONString(v))
	return b
}

func (b *VColBuilder) OrderLg(v interface{}) (r *VColBuilder) {
	b.Attr(":order-lg", h.JSONString(v))
	return b
}

func (b *VColBuilder) OrderXl(v interface{}) (r *VColBuilder) {
	b.Attr(":order-xl", h.JSONString(v))
	return b
}

func (b *VColBuilder) OrderXxl(v interface{}) (r *VColBuilder) {
	b.Attr(":order-xxl", h.JSONString(v))
	return b
}

func (b *VColBuilder) AlignSelf(v interface{}) (r *VColBuilder) {
	b.Attr(":align-self", h.JSONString(v))
	return b
}

func (b *VColBuilder) Tag(v string) (r *VColBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VColBuilder) On(name string, value string) (r *VColBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VColBuilder) Bind(name string, value string) (r *VColBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
