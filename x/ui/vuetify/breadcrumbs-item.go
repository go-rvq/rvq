package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VBreadcrumbsItemBuilder struct {
	VTagBuilder[*VBreadcrumbsItemBuilder]
}

func VBreadcrumbsItem(children ...h.HTMLComponent) *VBreadcrumbsItemBuilder {
	return VTag(&VBreadcrumbsItemBuilder{}, "v-breadcrumbs-item", children...)
}

func (b *VBreadcrumbsItemBuilder) Active(v bool) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":active", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) ActiveClass(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("active-class", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) ActiveColor(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("active-color", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) Color(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) Disabled(v bool) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":disabled", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) Title(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("title", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) Href(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("href", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) Replace(v bool) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":replace", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) Exact(v bool) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":exact", fmt.Sprint(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) To(v interface{}) (r *VBreadcrumbsItemBuilder) {
	b.Attr(":to", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsItemBuilder) Tag(v string) (r *VBreadcrumbsItemBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VBreadcrumbsItemBuilder) On(name string, value string) (r *VBreadcrumbsItemBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBreadcrumbsItemBuilder) Bind(name string, value string) (r *VBreadcrumbsItemBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
