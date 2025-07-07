package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VMainBuilder struct {
	VTagBuilder[*VMainBuilder]
}

func VMain(children ...h.HTMLComponent) *VMainBuilder {
	return VTag(&VMainBuilder{}, "v-main", children...)
}

func (b *VMainBuilder) Scrollable(v bool) (r *VMainBuilder) {
	b.Attr(":scrollable", fmt.Sprint(v))
	return b
}

func (b *VMainBuilder) Height(v interface{}) (r *VMainBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VMainBuilder) MaxHeight(v interface{}) (r *VMainBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VMainBuilder) MaxWidth(v interface{}) (r *VMainBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VMainBuilder) MinHeight(v interface{}) (r *VMainBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VMainBuilder) MinWidth(v interface{}) (r *VMainBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VMainBuilder) Width(v interface{}) (r *VMainBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VMainBuilder) Tag(v string) (r *VMainBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VMainBuilder) On(name string, value string) (r *VMainBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VMainBuilder) Bind(name string, value string) (r *VMainBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
