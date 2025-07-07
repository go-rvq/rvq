package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VBannerTextBuilder struct {
	VTagBuilder[*VBannerTextBuilder]
}

func VBannerText(children ...h.HTMLComponent) *VBannerTextBuilder {
	return VTag(&VBannerTextBuilder{}, "v-banner-text", children...)
}

func (b *VBannerTextBuilder) Tag(v string) (r *VBannerTextBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VBannerTextBuilder) On(name string, value string) (r *VBannerTextBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VBannerTextBuilder) Bind(name string, value string) (r *VBannerTextBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
