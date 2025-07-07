package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VCodeBuilder struct {
	VTagBuilder[*VCodeBuilder]
}

func VCode(children ...h.HTMLComponent) *VCodeBuilder {
	return VTag(&VCodeBuilder{}, "v-code", children...)
}

func (b *VCodeBuilder) Tag(v string) (r *VCodeBuilder) {
	b.Attr("tag", v)
	return b
}

func (b *VCodeBuilder) On(name string, value string) (r *VCodeBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VCodeBuilder) Bind(name string, value string) (r *VCodeBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
