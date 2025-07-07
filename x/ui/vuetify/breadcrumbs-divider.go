package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBreadcrumbsDividerBuilder struct {
	VTagBuilder[*VBreadcrumbsDividerBuilder]
}

func VBreadcrumbsDivider(children ...h.HTMLComponent) *VBreadcrumbsDividerBuilder {
	return VTag(&VBreadcrumbsDividerBuilder{}, "v-breadcrumbs-divider", children...)
}

func (b *VBreadcrumbsDividerBuilder) Divider(v interface{}) (r *VBreadcrumbsDividerBuilder) {
	b.Attr(":divider", h.JSONString(v))
	return b
}

func (b *VBreadcrumbsDividerBuilder) On(name string, value string) (r *VBreadcrumbsDividerBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBreadcrumbsDividerBuilder) Bind(name string, value string) (r *VBreadcrumbsDividerBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
