package vuetify

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VComponentIconBuilder struct {
	VTagBuilder[*VComponentIconBuilder]
}

func VComponentIcon(children ...h.HTMLComponent) *VComponentIconBuilder {
	return VTag(&VComponentIconBuilder{}, "v-component-icon", children...)
}

func (b *VComponentIconBuilder) Icon(v interface{}) (r *VComponentIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VComponentIconBuilder) Tag(v string) (r *VComponentIconBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VComponentIconBuilder) On(name string, value string) (r *VComponentIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VComponentIconBuilder) Bind(name string, value string) (r *VComponentIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
