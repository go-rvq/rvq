package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBreadcrumbsBuilder struct {
	VTagBuilder[*VBreadcrumbsBuilder]
}

func VBreadcrumbs(children ...h.HTMLComponent) *VBreadcrumbsBuilder {
	return VTag(&VBreadcrumbsBuilder{}, "v-breadcrumbs", children...)
}

func (b *VBreadcrumbsBuilder) Divider(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("divider", v)
	return b
}

func (b *VBreadcrumbsBuilder) ActiveClass(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VBreadcrumbsBuilder) ActiveColor(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VBreadcrumbsBuilder) BgColor(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("bg-color", v)
	return b
}

func (b *VBreadcrumbsBuilder) Color(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBreadcrumbsBuilder) Disabled(v bool) (r *VBreadcrumbsBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsBuilder) Icon(v interface{}) (r *VBreadcrumbsBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsBuilder) Items(v interface{}) (r *VBreadcrumbsBuilder) {
	b.Attr(":items", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsBuilder) Density(v interface{}) (r *VBreadcrumbsBuilder) {
	b.Attr(":density", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsBuilder) Rounded(v interface{}) (r *VBreadcrumbsBuilder) {
	b.Attr(":rounded", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsBuilder) Tile(v bool) (r *VBreadcrumbsBuilder) {
	b.Attr(":tile", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsBuilder) Tag(v string) (r *VBreadcrumbsBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VBreadcrumbsBuilder) On(name string, value string) (r *VBreadcrumbsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBreadcrumbsBuilder) Bind(name string, value string) (r *VBreadcrumbsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
