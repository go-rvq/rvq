package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VSkeletonLoaderBuilder struct {
	VTagBuilder[*VSkeletonLoaderBuilder]
}

func VSkeletonLoader(children ...h.HTMLComponent) *VSkeletonLoaderBuilder {
	return VTag(&VSkeletonLoaderBuilder{}, "v-skeleton-loader", children...)
}

func (b *VSkeletonLoaderBuilder) Boilerplate(v bool) (r *VSkeletonLoaderBuilder) {
	b.Attr(":boilerplate", fmt.Sprint(v))
	return b
}

func (b *VSkeletonLoaderBuilder) Color(v string) (r *VSkeletonLoaderBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VSkeletonLoaderBuilder) Loading(v bool) (r *VSkeletonLoaderBuilder) {
	b.Attr(":loading", fmt.Sprint(v))
	return b
}

func (b *VSkeletonLoaderBuilder) LoadingText(v string) (r *VSkeletonLoaderBuilder) {
	b.Attr("loading-text", v)
	return b
}

func (b *VSkeletonLoaderBuilder) Type(v interface{}) (r *VSkeletonLoaderBuilder) {
	b.Attr(":type", h.JSONString(v))
	return b
}

func (b *VSkeletonLoaderBuilder) Height(v interface{}) (r *VSkeletonLoaderBuilder) {
	b.Attr(":height", h.JSONString(v))
	return b
}

func (b *VSkeletonLoaderBuilder) MaxHeight(v interface{}) (r *VSkeletonLoaderBuilder) {
	b.Attr(":max-height", h.JSONString(v))
	return b
}

func (b *VSkeletonLoaderBuilder) MaxWidth(v interface{}) (r *VSkeletonLoaderBuilder) {
	b.Attr(":max-width", h.JSONString(v))
	return b
}

func (b *VSkeletonLoaderBuilder) MinHeight(v interface{}) (r *VSkeletonLoaderBuilder) {
	b.Attr(":min-height", h.JSONString(v))
	return b
}

func (b *VSkeletonLoaderBuilder) MinWidth(v interface{}) (r *VSkeletonLoaderBuilder) {
	b.Attr(":min-width", h.JSONString(v))
	return b
}

func (b *VSkeletonLoaderBuilder) Width(v interface{}) (r *VSkeletonLoaderBuilder) {
	b.Attr(":width", h.JSONString(v))
	return b
}

func (b *VSkeletonLoaderBuilder) Elevation(v interface{}) (r *VSkeletonLoaderBuilder) {
	b.Attr(":elevation", h.JSONString(v))
	return b
}

func (b *VSkeletonLoaderBuilder) Theme(v string) (r *VSkeletonLoaderBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VSkeletonLoaderBuilder) On(name string, value string) (r *VSkeletonLoaderBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VSkeletonLoaderBuilder) Bind(name string, value string) (r *VSkeletonLoaderBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
