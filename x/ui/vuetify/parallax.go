package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VParallaxBuilder struct {
	VTagBuilder[*VParallaxBuilder]
}

func VParallax(children ...h.HTMLComponent) *VParallaxBuilder {
	return VTag(&VParallaxBuilder{}, "v-parallax", children...)
}

func (b *VParallaxBuilder) Scale(v interface{}) (r *VParallaxBuilder) {
	b.Attr(":scale", h.JSONString(v))
	return b
}

func (b *VParallaxBuilder) On(name string, value string) (r *VParallaxBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VParallaxBuilder) Bind(name string, value string) (r *VParallaxBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
