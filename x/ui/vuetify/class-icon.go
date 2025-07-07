package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VClassIconBuilder struct {
	VTagBuilder[*VClassIconBuilder]
}

func VClassIcon(children ...h.HTMLComponent) *VClassIconBuilder {
	return VTag(&VClassIconBuilder{}, "v-class-icon", children...)
}

func (b *VClassIconBuilder) Icon(v interface{}) (r *VClassIconBuilder) {
	b.Attr(":icon", h.JSONString(v))
	return b
}

func (b *VClassIconBuilder) Tag(v string) (r *VClassIconBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VClassIconBuilder) On(name string, value string) (r *VClassIconBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VClassIconBuilder) Bind(name string, value string) (r *VClassIconBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
