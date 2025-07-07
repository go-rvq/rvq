package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VRowBuilder struct {
	VTagBuilder[*VRowBuilder]
}

func VRow(children ...h.HTMLComponent) *VRowBuilder {
	return VTag(&VRowBuilder{}, "v-row", children...)
}

func (b *VRowBuilder) Dense(v bool) (r *VRowBuilder) {
	b.Attr(":dense", fmt.Sprint(v))
	return b
}

func (b *VRowBuilder) NoGutters(v bool) (r *VRowBuilder) {
	b.Attr(":no-gutters", fmt.Sprint(v))
	return b
}

func (b *VRowBuilder) Align(v interface{}) (r *VRowBuilder) {
	b.Attr(":align", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignSm(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-sm", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignMd(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-md", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignLg(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-lg", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignXl(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-xl", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignXxl(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-xxl", h.JSONString(v))
	return b
}

func (b *VRowBuilder) JustifySm(v interface{}) (r *VRowBuilder) {
	b.Attr(":justify-sm", h.JSONString(v))
	return b
}

func (b *VRowBuilder) JustifyMd(v interface{}) (r *VRowBuilder) {
	b.Attr(":justify-md", h.JSONString(v))
	return b
}

func (b *VRowBuilder) JustifyLg(v interface{}) (r *VRowBuilder) {
	b.Attr(":justify-lg", h.JSONString(v))
	return b
}

func (b *VRowBuilder) JustifyXl(v interface{}) (r *VRowBuilder) {
	b.Attr(":justify-xl", h.JSONString(v))
	return b
}

func (b *VRowBuilder) JustifyXxl(v interface{}) (r *VRowBuilder) {
	b.Attr(":justify-xxl", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignContentSm(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-content-sm", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignContentMd(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-content-md", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignContentLg(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-content-lg", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignContentXl(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-content-xl", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignContentXxl(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-content-xxl", h.JSONString(v))
	return b
}

func (b *VRowBuilder) Justify(v interface{}) (r *VRowBuilder) {
	b.Attr(":justify", h.JSONString(v))
	return b
}

func (b *VRowBuilder) AlignContent(v interface{}) (r *VRowBuilder) {
	b.Attr(":align-content", h.JSONString(v))
	return b
}

func (b *VRowBuilder) Tag(v string) (r *VRowBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VRowBuilder) On(name string, value string) (r *VRowBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VRowBuilder) Bind(name string, value string) (r *VRowBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
