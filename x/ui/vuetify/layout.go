package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VLayoutBuilder struct {
	VTagBuilder[*VLayoutBuilder]
}

func VLayout(children ...h.HTMLComponent) *VLayoutBuilder {
	return VTag(&VLayoutBuilder{}, "v-layout", children...)
}

func (b *VLayoutBuilder) Height(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) MaxHeight(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) MaxWidth(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) MinHeight(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) MinWidth(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) Width(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) FullHeight(v bool) (r *VLayoutBuilder) {
	b.Attr(":full-height", fmt.Sprint(v))
	return b
}

func (b *VLayoutBuilder) Overlaps(v interface{}) (r *VLayoutBuilder) {
	b.Attr(":overlaps", h.JSONString(v))
	return b
}

func (b *VLayoutBuilder) On(name string, value string) (r *VLayoutBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VLayoutBuilder) Bind(name string, value string) (r *VLayoutBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
