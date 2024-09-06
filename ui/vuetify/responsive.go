package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VResponsiveBuilder struct {
	VTagBuilder[*VResponsiveBuilder]
}

func VResponsive(children ...h.HTMLComponent) *VResponsiveBuilder {
	return VTag(&VResponsiveBuilder{}, "v-responsive", children...)
}

func (b *VResponsiveBuilder) AspectRatio(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":aspect-ratio", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) ContentClass(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":content-class", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) Inline(v bool) (r *VResponsiveBuilder) {
	b.Attr(":inline", fmt.Sprint(v))
	return b
}

func (b *VResponsiveBuilder) Height(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) MaxHeight(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) MaxWidth(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) MinHeight(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) MinWidth(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) Width(v interface{}) (r *VResponsiveBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VResponsiveBuilder) On(name string, value string) (r *VResponsiveBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VResponsiveBuilder) Bind(name string, value string) (r *VResponsiveBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
