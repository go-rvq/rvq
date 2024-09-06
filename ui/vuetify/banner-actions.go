package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VBannerActionsBuilder struct {
	VTagBuilder[*VBannerActionsBuilder]
}

func VBannerActions(children ...h.HTMLComponent) *VBannerActionsBuilder {
	return VTag(&VBannerActionsBuilder{}, "v-banner-actions", children...)
}

func (b *VBannerActionsBuilder) Color(v string) (r *VBannerActionsBuilder) {
	b.Attr("color", v)
	return b
}

func (b *VBannerActionsBuilder) Density(v string) (r *VBannerActionsBuilder) {
	b.Attr("density", v)
	return b
}

func (b *VBannerActionsBuilder) On(name string, value string) (r *VBannerActionsBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBannerActionsBuilder) Bind(name string, value string) (r *VBannerActionsBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
