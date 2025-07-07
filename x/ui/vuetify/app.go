package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VAppBuilder struct {
	VTagBuilder[*VAppBuilder]
}

func VApp(children ...h.HTMLComponent) *VAppBuilder {
	return VTag(&VAppBuilder{}, "v-app", children...)
}

func (b *VAppBuilder) FullHeight(v bool) (r *VAppBuilder) {
	b.Attr(":full-height", fmt.Sprint(v))
	return b
}

func (b *VAppBuilder) Overlaps(v interface{}) (r *VAppBuilder) {
	b.Attr(":overlaps", h.JSONString(v))
	return b
}

func (b *VAppBuilder) Theme(v string) (r *VAppBuilder) {
	b.Attr("theme", v)
	return b
}

func (b *VAppBuilder) On(name string, value string) (r *VAppBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VAppBuilder) Bind(name string, value string) (r *VAppBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
