package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VLazyBuilder struct {
	VTagBuilder[*VLazyBuilder]
}

func VLazy(children ...h.HTMLComponent) *VLazyBuilder {
	return VTag(&VLazyBuilder{}, "v-lazy", children...)
}

func (b *VLazyBuilder) ModelValue(v bool) (r *VLazyBuilder) {
	b.Attr(":model-value", fmt.Sprint(v))
	return b
}

func (b *VLazyBuilder) Options(v interface{}) (r *VLazyBuilder) {
	b.Attr(":options", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) Height(v interface{}) (r *VLazyBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) MaxHeight(v interface{}) (r *VLazyBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) MaxWidth(v interface{}) (r *VLazyBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) MinHeight(v interface{}) (r *VLazyBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) MinWidth(v interface{}) (r *VLazyBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) Width(v interface{}) (r *VLazyBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) Tag(v string) (r *VLazyBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VLazyBuilder) Transition(v interface{}) (r *VLazyBuilder) {
	b.Attr(":transition", h.JSONString(v))
	return b
}

func (b *VLazyBuilder) On(name string, value string) (r *VLazyBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLazyBuilder) Bind(name string, value string) (r *VLazyBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
